/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.services.model.dis;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.aitool.AIGenerated;
import com.obs.services.ObsClient;
import com.obs.services.model.HeaderResponse;
import com.obs.test.TestTools;

import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

import java.util.Arrays;

public class DisPolicyTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-dispolicy";

    @BeforeClass
    public static void setMockServer() {
        mockServer = ClientAndServer.startClientAndServer();
        System.setProperty(PROXY_HOST_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_PROPERTY_NAME, "" + mockServer.getLocalPort());
        System.setProperty(PROXY_HOST_S_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_S_PROPERTY_NAME, "" + mockServer.getLocalPort());
    }

    @AfterClass
    public static void clearEnv() {
        mockServer.close();
        System.clearProperty(PROXY_HOST_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_PROPERTY_NAME);
        System.clearProperty(PROXY_HOST_S_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_S_PROPERTY_NAME);
    }

    // ------------------------------ setBucketDisPolicy 测试 ------------------------------
    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试设置DIS通知策略请求为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_setBucketDisPolicyRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("SetBucketDisPolicyRequest is null");
        obsClient.setBucketDisPolicy(null);
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试设置DIS通知策略桶名为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_setBucketDisPolicyRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        DisPolicyRule rule = new DisPolicyRule("rule1", "stream_name", "projectId",
            Arrays.asList("ObjectCreated:*"), "agency");
        DisPolicyConfiguration config = new DisPolicyConfiguration(Arrays.asList(rule));
        SetBucketDisPolicyRequest request = new SetBucketDisPolicyRequest(null, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setBucketDisPolicy(request);
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试设置DIS通知策略正常请求返回201")
    public void should_succeed_when_setBucketDisPolicy_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 201;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        DisPolicyRule rule = new DisPolicyRule("rule1", "stream_name", "projectId",
            Arrays.asList("ObjectCreated:*", "ObjectRemoved:*"), "dis_agency");
        rule.setPrefix("input/");
        rule.setSuffix(".txt");
        DisPolicyConfiguration config = new DisPolicyConfiguration(Arrays.asList(rule));
        SetBucketDisPolicyRequest request =
            new SetBucketDisPolicyRequest(bucketNameForTest, config);
        HeaderResponse response = obsClient.setBucketDisPolicy(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试设置DIS通知策略更新已有策略返回200")
    public void should_succeed_when_setBucketDisPolicy_returns_200() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        DisPolicyRule rule = new DisPolicyRule("rule1", "stream_name", "projectId",
            Arrays.asList("ObjectCreated:*"), "dis_agency");
        DisPolicyConfiguration config = new DisPolicyConfiguration(Arrays.asList(rule));
        SetBucketDisPolicyRequest request =
            new SetBucketDisPolicyRequest(bucketNameForTest, config);
        HeaderResponse response = obsClient.setBucketDisPolicy(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    // ------------------------------ getBucketDisPolicy 测试 ------------------------------
    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试获取DIS通知策略请求为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_getBucketDisPolicyRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("GetBucketDisPolicyRequest is null");
        obsClient.getBucketDisPolicy(null);
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试获取DIS通知策略桶名为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_getBucketDisPolicyRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        GetBucketDisPolicyRequest request = new GetBucketDisPolicyRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.getBucketDisPolicy(request);
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试获取DIS通知策略正常请求验证字段解析")
    public void should_succeed_when_getBucketDisPolicy_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseJson = "{"
                + "\"rules\":[{"
                + "\"id\":\"rule1\","
                + "\"stream\":\"stream_name\","
                + "\"project\":\"projectId\","
                + "\"events\":[\"ObjectCreated:*\",\"ObjectRemoved:*\"],"
                + "\"prefix\":\"input/\","
                + "\"suffix\":\".txt\","
                + "\"agency\":\"dis_agency\""
                + "}]"
                + "}";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("disPolicy"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/json")
                        .withBody(validResponseJson));

        GetBucketDisPolicyRequest request = new GetBucketDisPolicyRequest(bucketNameForTest);
        GetBucketDisPolicyResult result = obsClient.getBucketDisPolicy(request);

        Assert.assertNotNull(result);
        Assert.assertNotNull(result.getDisPolicyConfiguration());
        Assert.assertNotNull(result.getDisPolicyConfiguration().getRules());
        Assert.assertEquals(1, result.getDisPolicyConfiguration().getRules().size());

        DisPolicyRule rule = result.getDisPolicyConfiguration().getRules().get(0);
        Assert.assertEquals("rule1", rule.getId());
        Assert.assertEquals("stream_name", rule.getStream());
        Assert.assertEquals("projectId", rule.getProject());
        Assert.assertEquals(2, rule.getEvents().size());
        Assert.assertEquals("ObjectCreated:*", rule.getEvents().get(0));
        Assert.assertEquals("ObjectRemoved:*", rule.getEvents().get(1));
        Assert.assertEquals("input/", rule.getPrefix());
        Assert.assertEquals(".txt", rule.getSuffix());
        Assert.assertEquals("dis_agency", rule.getAgency());
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试获取DIS通知策略多规则场景")
    public void should_succeed_when_getBucketDisPolicy_with_multiple_rules() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseJson = "{"
                + "\"rules\":[{"
                + "\"id\":\"rule1\","
                + "\"stream\":\"stream_a\","
                + "\"project\":\"projectId1\","
                + "\"events\":[\"ObjectCreated:*\"],"
                + "\"agency\":\"agency1\""
                + "},{"
                + "\"id\":\"rule2\","
                + "\"stream\":\"stream_b\","
                + "\"project\":\"projectId2\","
                + "\"events\":[\"ObjectCreated:Put\",\"ObjectRemoved:*\"],"
                + "\"prefix\":\"folder/\","
                + "\"suffix\":\".log\","
                + "\"agency\":\"agency2\""
                + "}]"
                + "}";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("disPolicy"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/json")
                        .withBody(validResponseJson));

        GetBucketDisPolicyRequest request = new GetBucketDisPolicyRequest(bucketNameForTest);
        GetBucketDisPolicyResult result = obsClient.getBucketDisPolicy(request);

        Assert.assertNotNull(result);
        Assert.assertEquals(2, result.getDisPolicyConfiguration().getRules().size());

        DisPolicyRule rule1 = result.getDisPolicyConfiguration().getRules().get(0);
        Assert.assertEquals("rule1", rule1.getId());
        Assert.assertEquals("stream_a", rule1.getStream());
        Assert.assertEquals("agency1", rule1.getAgency());
        Assert.assertNull(rule1.getPrefix());

        DisPolicyRule rule2 = result.getDisPolicyConfiguration().getRules().get(1);
        Assert.assertEquals("rule2", rule2.getId());
        Assert.assertEquals("stream_b", rule2.getStream());
        Assert.assertEquals(2, rule2.getEvents().size());
        Assert.assertEquals("folder/", rule2.getPrefix());
        Assert.assertEquals(".log", rule2.getSuffix());
    }

    // ------------------------------ deleteBucketDisPolicy 测试 ------------------------------
    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试删除DIS通知策略请求为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_deleteBucketDisPolicyRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("DeleteBucketDisPolicyRequest is null");
        obsClient.deleteBucketDisPolicy(null);
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试删除DIS通知策略桶名为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_deleteBucketDisPolicyRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        DeleteBucketDisPolicyRequest request = new DeleteBucketDisPolicyRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.deleteBucketDisPolicy(request);
    }

    @Test
    @AIGenerated(author = "zhanghaoliang", date = "2026-04-16",
        description = "测试删除DIS通知策略正常请求返回204")
    public void should_succeed_when_deleteBucketDisPolicy_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 204;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("DELETE")
                        .withPath("")
                        .withQueryStringParameter("disPolicy"))
                .respond(response().withStatusCode(responseCodeForTest));

        DeleteBucketDisPolicyRequest request =
            new DeleteBucketDisPolicyRequest(bucketNameForTest);
        HeaderResponse response = obsClient.deleteBucketDisPolicy(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }
}
