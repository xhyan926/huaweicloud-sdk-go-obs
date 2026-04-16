/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.services.model.mirrorback;

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
import java.util.Collections;

public class MirrorBackToSourceTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-mirrorback";

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

    // ------------------------------ setBucketMirrorBackToSource 测试 ------------------------------
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试设置镜像回源策略请求为null时抛出异常")
    public void should_throw_exception_when_setBucketMirrorBackToSourceRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("SetBucketMirrorBackToSourceRequest is null");
        obsClient.setBucketMirrorBackToSource(null);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试设置镜像回源策略桶名为null时抛出异常")
    public void should_throw_exception_when_setBucketMirrorBackToSourceRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        MirrorBackToSourceConfiguration config = createSimpleConfiguration();
        SetBucketMirrorBackToSourceRequest request = new SetBucketMirrorBackToSourceRequest(null, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setBucketMirrorBackToSource(request);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试设置镜像回源策略返回201创建成功")
    public void should_succeed_when_setBucketMirrorBackToSource_returns_201() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 201;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        MirrorBackToSourceConfiguration config = createSimpleConfiguration();
        SetBucketMirrorBackToSourceRequest request =
            new SetBucketMirrorBackToSourceRequest(bucketNameForTest, config);
        HeaderResponse response = obsClient.setBucketMirrorBackToSource(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试设置镜像回源策略返回200更新成功")
    public void should_succeed_when_setBucketMirrorBackToSource_returns_200() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        MirrorBackToSourceConfiguration config = createSimpleConfiguration();
        SetBucketMirrorBackToSourceRequest request =
            new SetBucketMirrorBackToSourceRequest(bucketNameForTest, config);
        HeaderResponse response = obsClient.setBucketMirrorBackToSource(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    // ------------------------------ getBucketMirrorBackToSource 测试 ------------------------------
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试获取镜像回源策略请求为null时抛出异常")
    public void should_throw_exception_when_getBucketMirrorBackToSourceRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("GetBucketMirrorBackToSourceRequest is null");
        obsClient.getBucketMirrorBackToSource(null);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试获取镜像回源策略桶名为null时抛出异常")
    public void should_throw_exception_when_getBucketMirrorBackToSourceRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        GetBucketMirrorBackToSourceRequest request = new GetBucketMirrorBackToSourceRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.getBucketMirrorBackToSource(request);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试获取镜像回源策略完整JSON响应解析")
    public void should_succeed_when_getBucketMirrorBackToSource_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseJson = "{"
                + "\"rules\":[{"
                + "\"id\":\"rule1\","
                + "\"condition\":{\"httpErrorCodeReturnedEquals\":\"404\",\"objectKeyPrefixEquals\":\"prefix/\"},"
                + "\"redirect\":{"
                + "\"agency\":\"test-agency\","
                + "\"publicSource\":{\"sourceEndpoint\":{\"master\":[\"https://source1.com\"],\"slave\":[\"https://source2.com\"]}},"
                + "\"passQueryString\":true,"
                + "\"mirrorFollowRedirect\":false"
                + "}"
                + "}]"
                + "}";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("mirrorBackToSource"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/json")
                        .withBody(validResponseJson));

        GetBucketMirrorBackToSourceRequest request =
            new GetBucketMirrorBackToSourceRequest(bucketNameForTest);
        GetBucketMirrorBackToSourceResult result = obsClient.getBucketMirrorBackToSource(request);

        Assert.assertNotNull(result);
        Assert.assertNotNull(result.getMirrorBackToSourceConfiguration());
        Assert.assertNotNull(result.getMirrorBackToSourceConfiguration().getRules());
        Assert.assertEquals(1, result.getMirrorBackToSourceConfiguration().getRules().size());

        MirrorBackToSourceRule rule = result.getMirrorBackToSourceConfiguration().getRules().get(0);
        Assert.assertEquals("rule1", rule.getId());

        MirrorBackCondition condition = rule.getCondition();
        Assert.assertNotNull(condition);
        Assert.assertEquals("404", condition.getHttpErrorCodeReturnedEquals());
        Assert.assertEquals("prefix/", condition.getObjectKeyPrefixEquals());

        MirrorBackRedirect redirect = rule.getRedirect();
        Assert.assertNotNull(redirect);
        Assert.assertEquals("test-agency", redirect.getAgency());
        Assert.assertTrue(redirect.getPassQueryString());
        Assert.assertFalse(redirect.getMirrorFollowRedirect());

        MirrorBackPublicSource publicSource = redirect.getPublicSource();
        Assert.assertNotNull(publicSource);
        MirrorBackSourceEndpoint endpoint = publicSource.getSourceEndpoint();
        Assert.assertNotNull(endpoint);
        Assert.assertEquals(1, endpoint.getMaster().size());
        Assert.assertEquals("https://source1.com", endpoint.getMaster().get(0));
        Assert.assertEquals(1, endpoint.getSlave().size());
        Assert.assertEquals("https://source2.com", endpoint.getSlave().get(0));
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试获取镜像回源策略包含HTTP头传递规则的解析")
    public void should_succeed_when_getBucketMirrorBackToSource_with_httpHeader_config() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseJson = "{"
                + "\"rules\":[{"
                + "\"id\":\"rule2\","
                + "\"condition\":{\"httpErrorCodeReturnedEquals\":\"404\"},"
                + "\"redirect\":{"
                + "\"agency\":\"agency2\","
                + "\"publicSource\":{\"sourceEndpoint\":{\"master\":[\"https://src.com\"]}},"
                + "\"mirrorHttpHeader\":{"
                + "\"passAll\":true,"
                + "\"pass\":[\"header1\",\"header2\"],"
                + "\"remove\":[\"header3\"],"
                + "\"set\":[{\"key\":\"x-custom-header\",\"value\":\"custom-value\"}]"
                + "},"
                + "\"replaceKeyPrefixWith\":\"prefix/\","
                + "\"mirrorAllowHttpMethod\":[\"GET\",\"HEAD\"]"
                + "}"
                + "}]"
                + "}";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("mirrorBackToSource"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/json")
                        .withBody(validResponseJson));

        GetBucketMirrorBackToSourceRequest request =
            new GetBucketMirrorBackToSourceRequest(bucketNameForTest);
        GetBucketMirrorBackToSourceResult result = obsClient.getBucketMirrorBackToSource(request);

        Assert.assertNotNull(result);
        MirrorBackToSourceRule rule = result.getMirrorBackToSourceConfiguration().getRules().get(0);
        MirrorBackRedirect redirect = rule.getRedirect();

        MirrorBackHttpHeader httpHeader = redirect.getMirrorHttpHeader();
        Assert.assertNotNull(httpHeader);
        Assert.assertTrue(httpHeader.getPassAll());
        Assert.assertEquals(Arrays.asList("header1", "header2"), httpHeader.getPass());
        Assert.assertEquals(Collections.singletonList("header3"), httpHeader.getRemove());
        Assert.assertEquals(1, httpHeader.getSet().size());
        Assert.assertEquals("x-custom-header", httpHeader.getSet().get(0).getKey());
        Assert.assertEquals("custom-value", httpHeader.getSet().get(0).getValue());

        Assert.assertEquals("prefix/", redirect.getReplaceKeyPrefixWith());
        Assert.assertEquals(Arrays.asList("GET", "HEAD"), redirect.getMirrorAllowHttpMethod());
    }

    // ------------------------------ deleteBucketMirrorBackToSource 测试 ------------------------------
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试删除镜像回源策略请求为null时抛出异常")
    public void should_throw_exception_when_deleteBucketMirrorBackToSourceRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("DeleteBucketMirrorBackToSourceRequest is null");
        obsClient.deleteBucketMirrorBackToSource(null);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试删除镜像回源策略桶名为null时抛出异常")
    public void should_throw_exception_when_deleteBucketMirrorBackToSourceRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        DeleteBucketMirrorBackToSourceRequest request = new DeleteBucketMirrorBackToSourceRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.deleteBucketMirrorBackToSource(request);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试删除镜像回源策略返回204成功")
    public void should_succeed_when_deleteBucketMirrorBackToSource_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 204;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("DELETE")
                        .withPath("")
                        .withQueryStringParameter("mirrorBackToSource"))
                .respond(response().withStatusCode(responseCodeForTest));

        DeleteBucketMirrorBackToSourceRequest request =
            new DeleteBucketMirrorBackToSourceRequest(bucketNameForTest);
        HeaderResponse response = obsClient.deleteBucketMirrorBackToSource(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    private MirrorBackToSourceConfiguration createSimpleConfiguration() {
        MirrorBackCondition condition = new MirrorBackCondition("404", "prefix/");
        MirrorBackSourceEndpoint endpoint = new MirrorBackSourceEndpoint(
            Collections.singletonList("https://source.example.com"), null);
        MirrorBackPublicSource publicSource = new MirrorBackPublicSource(endpoint);
        MirrorBackRedirect redirect = new MirrorBackRedirect();
        redirect.setAgency("test-agency");
        redirect.setPublicSource(publicSource);
        redirect.setPassQueryString(true);
        redirect.setMirrorFollowRedirect(false);
        MirrorBackToSourceRule rule = new MirrorBackToSourceRule("rule1", condition, redirect);
        return new MirrorBackToSourceConfiguration(Collections.singletonList(rule));
    }

    // ------------------------------ Model 类覆盖率补充测试 ------------------------------

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试Model类默认构造函数和setter覆盖")
    public void should_cover_model_setters_and_default_constructors() {
        // MirrorBackCondition: default constructor + setters
        MirrorBackCondition condition = new MirrorBackCondition();
        condition.setHttpErrorCodeReturnedEquals("404");
        condition.setObjectKeyPrefixEquals("prefix/");
        Assert.assertEquals("404", condition.getHttpErrorCodeReturnedEquals());
        Assert.assertEquals("prefix/", condition.getObjectKeyPrefixEquals());
        Assert.assertTrue(condition.toString().contains("404"));

        // MirrorBackSourceEndpoint: setters
        MirrorBackSourceEndpoint endpoint = new MirrorBackSourceEndpoint();
        endpoint.setMaster(Collections.singletonList("https://master.example.com"));
        endpoint.setSlave(Collections.singletonList("https://slave.example.com"));
        Assert.assertEquals(1, endpoint.getMaster().size());
        Assert.assertEquals(1, endpoint.getSlave().size());
        Assert.assertTrue(endpoint.toString().contains("master"));

        // MirrorBackPublicSource: default constructor + setter
        MirrorBackPublicSource publicSource = new MirrorBackPublicSource();
        publicSource.setSourceEndpoint(endpoint);
        Assert.assertNotNull(publicSource.getSourceEndpoint());
        Assert.assertTrue(publicSource.toString().contains("sourceEndpoint"));

        // MirrorBackRedirect: all setters
        MirrorBackRedirect redirect = new MirrorBackRedirect();
        redirect.setAgency("test-agency");
        redirect.setPublicSource(publicSource);
        redirect.setRetryConditions(Arrays.asList("4XX", "500"));
        redirect.setPassQueryString(true);
        redirect.setMirrorFollowRedirect(true);
        redirect.setReplaceKeyWith("prefix${key}suffix");
        redirect.setReplaceKeyPrefixWith("backup/");
        redirect.setVpcEndpointURN("urn:xxx");
        redirect.setRedirectWithoutReferer(true);
        redirect.setMirrorAllowHttpMethod(Collections.singletonList("HEAD"));

        Assert.assertEquals("test-agency", redirect.getAgency());
        Assert.assertEquals(2, redirect.getRetryConditions().size());
        Assert.assertEquals("prefix${key}suffix", redirect.getReplaceKeyWith());
        Assert.assertEquals("backup/", redirect.getReplaceKeyPrefixWith());
        Assert.assertEquals("urn:xxx", redirect.getVpcEndpointURN());
        Assert.assertTrue(redirect.getRedirectWithoutReferer());
        Assert.assertEquals(1, redirect.getMirrorAllowHttpMethod().size());
        Assert.assertTrue(redirect.toString().contains("test-agency"));

        // MirrorBackHttpHeader: all setters
        MirrorBackHttpHeader httpHeader = new MirrorBackHttpHeader();
        httpHeader.setPassAll(true);
        httpHeader.setPass(Arrays.asList("content-type"));
        httpHeader.setRemove(Arrays.asList("authorization"));
        MirrorBackHttpHeaderSet headerSet = new MirrorBackHttpHeaderSet();
        headerSet.setKey("x-custom");
        headerSet.setValue("test-value");
        httpHeader.setSet(Collections.singletonList(headerSet));

        Assert.assertTrue(httpHeader.getPassAll());
        Assert.assertEquals(1, httpHeader.getPass().size());
        Assert.assertEquals(1, httpHeader.getRemove().size());
        Assert.assertEquals("x-custom", httpHeader.getSet().get(0).getKey());
        Assert.assertEquals("test-value", httpHeader.getSet().get(0).getValue());
        Assert.assertTrue(httpHeader.toString().contains("passAll"));
        Assert.assertTrue(headerSet.toString().contains("x-custom"));

        // MirrorBackHttpHeaderSet: parameterized constructor
        MirrorBackHttpHeaderSet headerSet2 = new MirrorBackHttpHeaderSet("key2", "value2");
        Assert.assertEquals("key2", headerSet2.getKey());
        Assert.assertEquals("value2", headerSet2.getValue());

        // MirrorBackToSourceRule: default constructor + setters
        MirrorBackToSourceRule rule = new MirrorBackToSourceRule();
        rule.setId("rule-test");
        rule.setCondition(condition);
        rule.setRedirect(redirect);
        Assert.assertEquals("rule-test", rule.getId());
        Assert.assertNotNull(rule.getCondition());
        Assert.assertNotNull(rule.getRedirect());
        Assert.assertTrue(rule.toString().contains("rule-test"));

        // MirrorBackToSourceConfiguration: default constructor + setter
        MirrorBackToSourceConfiguration config = new MirrorBackToSourceConfiguration();
        config.setRules(Collections.singletonList(rule));
        Assert.assertEquals(1, config.getRules().size());
        Assert.assertTrue(config.toString().contains("rules"));
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "测试Request/Result类setter覆盖")
    public void should_cover_request_and_result_setters() {
        // SetBucketMirrorBackToSourceRequest: setter
        SetBucketMirrorBackToSourceRequest setRequest =
            new SetBucketMirrorBackToSourceRequest("bucket", new MirrorBackToSourceConfiguration());
        MirrorBackToSourceConfiguration newConfig = new MirrorBackToSourceConfiguration();
        setRequest.setMirrorBackToSourceConfiguration(newConfig);
        Assert.assertSame(newConfig, setRequest.getMirrorBackToSourceConfiguration());

        // GetBucketMirrorBackToSourceRequest: default constructor + setter
        GetBucketMirrorBackToSourceRequest getRequest = new GetBucketMirrorBackToSourceRequest();
        getRequest.setBucketName("test-bucket");
        Assert.assertEquals("test-bucket", getRequest.getBucketName());

        // GetBucketMirrorBackToSourceResult: setter
        GetBucketMirrorBackToSourceResult result = new GetBucketMirrorBackToSourceResult();
        MirrorBackToSourceConfiguration resultConfig = new MirrorBackToSourceConfiguration();
        result.setMirrorBackToSourceConfiguration(resultConfig);
        Assert.assertSame(resultConfig, result.getMirrorBackToSourceConfiguration());

        // DeleteBucketMirrorBackToSourceRequest: default constructor + setter
        DeleteBucketMirrorBackToSourceRequest deleteRequest = new DeleteBucketMirrorBackToSourceRequest();
        deleteRequest.setBucketName("test-bucket");
        Assert.assertEquals("test-bucket", deleteRequest.getBucketName());
    }
}
