/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.model.lifecycle;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.DeleteBucketLifecycleRequest;
import com.obs.services.model.GetBucketLifecycleRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.LifecycleConfiguration;
import com.obs.services.model.SetBucketLifecycleRequest;
import com.obs.services.model.StorageClassEnum;
import com.obs.test.TestTools;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

import java.io.IOException;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

public class AbstractBucketClientLifecycleUnitTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";

    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";

    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";

    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;

    public static String bucketNameForTest = "test-bucket-lifecycle";

    @BeforeClass
    public static void setMockServer() {
        // 启动 MockServer
        mockServer = ClientAndServer.startClientAndServer();
        System.setProperty(PROXY_HOST_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_PROPERTY_NAME, "" + mockServer.getLocalPort());
        System.setProperty(PROXY_HOST_S_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_S_PROPERTY_NAME, "" + mockServer.getLocalPort());
    }

    @AfterClass
    public static void clearEnv() {
        // 关闭 MockServer
        mockServer.close();
        System.clearProperty(PROXY_HOST_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_PROPERTY_NAME);
        System.clearProperty(PROXY_HOST_S_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_S_PROPERTY_NAME);
    }

    /*
    测试场景分析：
        getBucketLifecycle(final GetBucketLifecycleRequest request)：
            1.正常场景：
                1）.调用getBucketLifecycle,ruleId和ruleIdMarker不为空且符合1-255的长度范围，返回码200
            2.边界值：
                1）.ruleId不为空长度为256，抛出异常
                2）.ruleId为空字符串，抛出异常
                3）.ruleIdMarker不为空长度为256，抛出异常
                4）.ruleIdMarker为空字符串，抛出异常
            3.异常场景：
                1）.request.getBucketName()为空，抛出bucketName is null异常
    */
    @Test
    public void should_return_200_when_getBucketLifecycle_with_valid_ruleId_success() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        String validLifeCycleResponseXml = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>"
            + "<LifecycleConfiguration xmlns=\"http://obs.cn-north-4.myhuaweicloud.com/doc/2015-06-30/\">"
            + "    <Rule>" + "        <ID>id</ID>" + "        <Filter>" + "          <And>"
            + "             <Prefix>prefix</Prefix>" + "             <Tag><Key>key1</Key><Value>value1</Value></Tag>"
            + "             <Tag><Key>key2</Key><Value>value2</Value></Tag>" + "          </And>" + "        </Filter>"
            + "        <Status>status</Status> " + "        <Expiration> " + "            <Date>date</Date>"
            + "        </Expiration>" + "        <NoncurrentVersionExpiration>"
            + "            <NoncurrentDays>days</NoncurrentDays>" + "        </NoncurrentVersionExpiration>"
            + "        <Transition>" + "         <Date>date</Date>" + "         <StorageClass>WARM</StorageClass>"
            + "        </Transition>" + "        <Transition>" + "         <Date>date</Date>"
            + "         <StorageClass>COLD</StorageClass>" + "        </Transition>"
            + "        <NoncurrentVersionTransition>" + "         <NoncurrentDays>30</NoncurrentDays>"
            + "         <StorageClass>WARM</StorageClass>" + "        </NoncurrentVersionTransition>"
            + "        <NoncurrentVersionTransition>" + "         <NoncurrentDays>60</NoncurrentDays>"
            + "         <StorageClass>COLD</StorageClass>" + "        </NoncurrentVersionTransition>" + "    </Rule>"
            + "</LifecycleConfiguration>";
        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("GET").withPath("").withQueryStringParameter("lifecycle"))
            .respond(response().withStatusCode(responseCodeForTest)
                .withHeader("Content-Type", "application/xml;charset=utf-8")
                .withHeader("Content-Length", String.valueOf(validLifeCycleResponseXml.getBytes().length))
                .withBody(validLifeCycleResponseXml));
        String validRuleId = java.util.stream.IntStream.range(0, 255)
            .mapToObj(i -> String.valueOf((char) (new java.util.Random().nextInt(26) + 'a')))
            .reduce("", String::concat);
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest, validRuleId, validRuleId);
        LifecycleConfiguration response = obsClient.getBucketLifecycle(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    @Test
    public void should_throw_exception_when_getBucketLifecycle_bucketName_is_null() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.getBucketLifecycle(new GetBucketLifecycleRequest(null));
    }

    @Test
    public void should_throw_exception_when_getBucketLifecycle_ruleId_length_is_256() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        String inValidRuleId = java.util.stream.IntStream.range(0, 256)
            .mapToObj(i -> String.valueOf((char) (new java.util.Random().nextInt(26) + 'a')))
            .reduce("", String::concat);
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest, inValidRuleId);
        obsClient.getBucketLifecycle(request);
    }

    @Test
    public void should_throw_exception_when_getBucketLifecycle_ruleId_length_is_0() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest, "");
        obsClient.getBucketLifecycle(request);
    }

    @Test
    public void should_throw_exception_when_getBucketLifecycle_ruleIdMarKer_length_is_256() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        String inValidRuleId = java.util.stream.IntStream.range(0, 256)
            .mapToObj(i -> String.valueOf((char) (new java.util.Random().nextInt(26) + 'a')))
            .reduce("", String::concat);
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest);
        request.setRuleId(inValidRuleId);
        obsClient.getBucketLifecycle(request);
    }

    @Test
    public void should_throw_exception_when_getBucketLifecycle_ruleIdMarKer_length_is_0() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest);
        request.setRuleId("");
        obsClient.getBucketLifecycle(request);
    }

    /*
    测试场景分析：
    deleteBucketLifecycle(final DeleteBucketLifecycleRequest request))：
        1.正常场景：
            1）.调用deleteBucketLifecycle,ruleId不为空且符合1-255的长度范围，返回码200
        2.边界值：
            1）.ruleId不为空长度为256，抛出异常
            2）.ruleId为空字符串，抛出异常
        3.异常场景：
            1）.request.getBucketName()为空，抛出bucketName is null异常
    */
    @Test
    public void should_return_200_when_deleteBucketLifecycle_with_valid_ruleId_success() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("DELETE").withPath("").withQueryStringParameter("lifecycle"))
            .respond(response().withStatusCode(responseCodeForTest));
        String validRuleId = java.util.stream.IntStream.range(0, 255)
            .mapToObj(i -> String.valueOf((char) (new java.util.Random().nextInt(26) + 'a')))
            .reduce("", String::concat);
        DeleteBucketLifecycleRequest request = new DeleteBucketLifecycleRequest(bucketNameForTest, validRuleId);
        HeaderResponse response = obsClient.deleteBucketLifecycle(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    @Test
    public void should_throw_exception_when_deleteBucketLifecycle_bucketName_is_null() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.deleteBucketLifecycle(new DeleteBucketLifecycleRequest(null));
    }

    @Test
    public void should_throw_exception_when_deleteBucketLifecycle_ruleId_length_is_256() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        String inValidRuleId = java.util.stream.IntStream.range(0, 256)
            .mapToObj(i -> String.valueOf((char) (new java.util.Random().nextInt(26) + 'a')))
            .reduce("", String::concat);
        DeleteBucketLifecycleRequest request = new DeleteBucketLifecycleRequest(bucketNameForTest, inValidRuleId);
        obsClient.deleteBucketLifecycle(request);
    }

    @Test
    public void should_throw_exception_when_deleteBucketLifecycle_ruleId_length_is_0() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        DeleteBucketLifecycleRequest request = new DeleteBucketLifecycleRequest(bucketNameForTest, "");
        obsClient.deleteBucketLifecycle(request);
    }

    /*
    测试场景分析：
    setBucketLifecycle(final SetBucketLifecycleRequest request)：
        1.正常场景：
            1）.调用setBucketLifecycle,ruleId不为空且符合1-255的长度范围，返回码200
        2.边界值：
            1）.ruleId不为空长度为256，抛出异常
            2）.ruleId为空字符串，抛出异常
        3.异常场景：
            1）.request.getBucketName()为空，抛出bucketName is null异常
            2）.request.getLifecycleConfig()为空，抛出LifecycleConfiguration is null异常
    */
    @Test
    public void should_return_200_when_setBucketLifecycle_with_valid_ruleId_success() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath("").withQueryStringParameter("lifecycle"))
            .respond(response().withStatusCode(responseCodeForTest));

        LifecycleConfiguration config = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = config.new Rule();
        rule.setEnabled(true);
        rule.setId("rule1");
        rule.setPrefix("prefix1");
        LifecycleConfiguration.Transition transition = config.new Transition();
        transition.setDays(30);
        transition.setObjectStorageClass(StorageClassEnum.WARM);
        rule.getTransitions().add(transition);
        config.addRule(rule);

        String validRuleId = java.util.stream.IntStream.range(0, 255)
            .mapToObj(i -> String.valueOf((char) (new java.util.Random().nextInt(26) + 'a')))
            .reduce("", String::concat);
        SetBucketLifecycleRequest request = new SetBucketLifecycleRequest(bucketNameForTest, validRuleId, config);
        HeaderResponse response = obsClient.setBucketLifecycle(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    @Test
    public void should_throw_exception_when_setBucketLifecycle_bucketName_is_null() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        LifecycleConfiguration config = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = config.new Rule();
        rule.setEnabled(true);
        rule.setId("rule1");
        rule.setPrefix("prefix1");
        LifecycleConfiguration.Transition transition = config.new Transition();
        transition.setDays(30);
        transition.setObjectStorageClass(StorageClassEnum.WARM);
        rule.getTransitions().add(transition);
        config.addRule(rule);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setBucketLifecycle(new SetBucketLifecycleRequest(null, "aaa", config));
    }

    @Test
    public void should_throw_exception_when_setBucketLifecycle_LifecycleConfiguration_is_null() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("LifecycleConfiguration is null");
        obsClient.setBucketLifecycle(new SetBucketLifecycleRequest(bucketNameForTest, "aaa", null));
    }

    @Test
    public void should_throw_exception_when_setBucketLifecycle_ruleId_length_is_256() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        LifecycleConfiguration config = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = config.new Rule();
        rule.setEnabled(true);
        rule.setId("rule1");
        rule.setPrefix("prefix1");
        LifecycleConfiguration.Transition transition = config.new Transition();
        transition.setDays(30);
        transition.setObjectStorageClass(StorageClassEnum.WARM);
        rule.getTransitions().add(transition);
        config.addRule(rule);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        String inValidRuleId = java.util.stream.IntStream.range(0, 256)
            .mapToObj(i -> String.valueOf((char) (new java.util.Random().nextInt(26) + 'a')))
            .reduce("", String::concat);
        SetBucketLifecycleRequest request = new SetBucketLifecycleRequest(bucketNameForTest, inValidRuleId, config);
        obsClient.setBucketLifecycle(request);
    }

    @Test
    public void should_throw_exception_when_setBucketLifecycle_ruleId_length_is_0() throws IOException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        LifecycleConfiguration config = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = config.new Rule();
        rule.setEnabled(true);
        rule.setId("rule1");
        rule.setPrefix("prefix1");
        LifecycleConfiguration.Transition transition = config.new Transition();
        transition.setDays(30);
        transition.setObjectStorageClass(StorageClassEnum.WARM);
        rule.getTransitions().add(transition);
        config.addRule(rule);
        SetBucketLifecycleRequest request = new SetBucketLifecycleRequest(bucketNameForTest, "", config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rule-id length should be between 1 and 255 characters.");
        obsClient.setBucketLifecycle(request);
    }

}