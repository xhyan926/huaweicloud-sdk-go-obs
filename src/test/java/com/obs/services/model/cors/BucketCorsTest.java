/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.services.model.cors;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.aitool.AIGenerated;
import com.obs.services.ObsClient;
import com.obs.services.model.BaseBucketRequest;
import com.obs.services.model.BucketCors;
import com.obs.services.model.BucketCorsRule;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.SetBucketCorsRequest;
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

public class BucketCorsTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-cors";

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

    // ------------------------------ setBucketCors 测试 ------------------------------
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试设置桶CORS配置请求为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_setBucketCorsRequest_is_null() {
        ObsClient obsClient = getObsClient();

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("SetBucketCorsRequest is null");
        obsClient.setBucketCors((SetBucketCorsRequest) null);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试设置桶CORS配置桶名为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_setBucketCorsRequest_bucketName_is_null() {
        ObsClient obsClient = getObsClient();

        BucketCors bucketCors = new BucketCors();
        BucketCorsRule rule = new BucketCorsRule();
        rule.setAllowedOrigin(Arrays.asList("http://www.example.com"));
        rule.setAllowedMethod(Arrays.asList("GET", "PUT"));
        bucketCors.setRules(Collections.singletonList(rule));

        SetBucketCorsRequest request = new SetBucketCorsRequest(null, bucketCors);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setBucketCors(request);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试设置桶CORS配置正常请求返回200")
    public void should_succeed_when_setBucketCors_with_valid_parameters() {
        ObsClient obsClient = getObsClient();

        setupMockPutResponse(200);

        BucketCors bucketCors = createBucketCors();
        SetBucketCorsRequest request = new SetBucketCorsRequest(bucketNameForTest, bucketCors);
        HeaderResponse response = obsClient.setBucketCors(request);

        Assert.assertEquals(200, response.getStatusCode());
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试设置桶CORS配置首次创建返回201")
    public void should_succeed_when_setBucketCors_returns_201() {
        ObsClient obsClient = getObsClient();

        setupMockPutResponse(201);

        BucketCors bucketCors = createBucketCors();
        SetBucketCorsRequest request = new SetBucketCorsRequest(bucketNameForTest, bucketCors);
        HeaderResponse response = obsClient.setBucketCors(request);

        Assert.assertEquals(201, response.getStatusCode());
    }

    // ------------------------------ getBucketCors 测试 ------------------------------
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试获取桶CORS配置请求为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_getBucketCorsRequest_is_null() {
        ObsClient obsClient = getObsClient();

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("BaseBucketRequest is null");
        obsClient.getBucketCors((BaseBucketRequest) null);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试获取桶CORS配置桶名为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_getBucketCorsRequest_bucketName_is_null() {
        ObsClient obsClient = getObsClient();

        BaseBucketRequest request = new BaseBucketRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.getBucketCors(request);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试获取桶CORS配置正常请求验证字段解析")
    public void should_succeed_when_getBucketCors_with_valid_parameters() {
        ObsClient obsClient = getObsClient();

        String validResponseXml = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
                + "<CORSConfiguration>"
                + "  <CORSRule>"
                + "    <ID>rule-001</ID>"
                + "    <AllowedMethod>GET</AllowedMethod>"
                + "    <AllowedMethod>PUT</AllowedMethod>"
                + "    <AllowedOrigin>http://www.example.com</AllowedOrigin>"
                + "    <AllowedHeader>x-obs-header</AllowedHeader>"
                + "    <MaxAgeSeconds>100</MaxAgeSeconds>"
                + "    <ExposeHeader>x-obs-expose-header</ExposeHeader>"
                + "  </CORSRule>"
                + "</CORSConfiguration>";

        setupMockCorsGetResponse(validResponseXml);

        BaseBucketRequest request = new BaseBucketRequest(bucketNameForTest);
        BucketCors result = obsClient.getBucketCors(request);

        Assert.assertNotNull(result);
        Assert.assertNotNull(result.getRules());
        Assert.assertEquals(1, result.getRules().size());

        BucketCorsRule rule = result.getRules().get(0);
        Assert.assertEquals("rule-001", rule.getId());
        Assert.assertEquals(2, rule.getAllowedMethod().size());
        Assert.assertEquals("GET", rule.getAllowedMethod().get(0));
        Assert.assertEquals("PUT", rule.getAllowedMethod().get(1));
        Assert.assertEquals(1, rule.getAllowedOrigin().size());
        Assert.assertEquals("http://www.example.com", rule.getAllowedOrigin().get(0));
        Assert.assertEquals(1, rule.getAllowedHeader().size());
        Assert.assertEquals("x-obs-header", rule.getAllowedHeader().get(0));
        Assert.assertEquals(100, rule.getMaxAgeSecond());
        Assert.assertEquals(1, rule.getExposeHeader().size());
        Assert.assertEquals("x-obs-expose-header", rule.getExposeHeader().get(0));
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试获取桶CORS配置多规则场景")
    public void should_succeed_when_getBucketCors_with_multiple_rules() {
        ObsClient obsClient = getObsClient();

        String validResponseXml = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
                + "<CORSConfiguration>"
                + "  <CORSRule>"
                + "    <ID>rule-001</ID>"
                + "    <AllowedMethod>PUT</AllowedMethod>"
                + "    <AllowedMethod>POST</AllowedMethod>"
                + "    <AllowedMethod>DELETE</AllowedMethod>"
                + "    <AllowedOrigin>http://www.example.com</AllowedOrigin>"
                + "    <AllowedHeader>*</AllowedHeader>"
                + "  </CORSRule>"
                + "  <CORSRule>"
                + "    <ID>rule-002</ID>"
                + "    <AllowedMethod>GET</AllowedMethod>"
                + "    <AllowedOrigin>*</AllowedOrigin>"
                + "    <MaxAgeSeconds>300</MaxAgeSeconds>"
                + "  </CORSRule>"
                + "</CORSConfiguration>";

        setupMockCorsGetResponse(validResponseXml);

        BaseBucketRequest request = new BaseBucketRequest(bucketNameForTest);
        BucketCors result = obsClient.getBucketCors(request);

        Assert.assertNotNull(result);
        Assert.assertEquals(2, result.getRules().size());

        BucketCorsRule rule1 = result.getRules().get(0);
        Assert.assertEquals("rule-001", rule1.getId());
        Assert.assertEquals(3, rule1.getAllowedMethod().size());
        Assert.assertEquals("http://www.example.com", rule1.getAllowedOrigin().get(0));
        Assert.assertEquals(1, rule1.getAllowedHeader().size());
        Assert.assertEquals("*", rule1.getAllowedHeader().get(0));

        BucketCorsRule rule2 = result.getRules().get(1);
        Assert.assertEquals("rule-002", rule2.getId());
        Assert.assertEquals(1, rule2.getAllowedMethod().size());
        Assert.assertEquals("GET", rule2.getAllowedMethod().get(0));
        Assert.assertEquals("*", rule2.getAllowedOrigin().get(0));
        Assert.assertEquals(300, rule2.getMaxAgeSecond());
    }

    // ------------------------------ deleteBucketCors 测试 ------------------------------
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试删除桶CORS配置请求为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_deleteBucketCorsRequest_is_null() {
        ObsClient obsClient = getObsClient();

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("BaseBucketRequest is null");
        obsClient.deleteBucketCors((BaseBucketRequest) null);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试删除桶CORS配置桶名为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_deleteBucketCorsRequest_bucketName_is_null() {
        ObsClient obsClient = getObsClient();

        BaseBucketRequest request = new BaseBucketRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.deleteBucketCors(request);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试删除桶CORS配置正常请求返回204")
    public void should_succeed_when_deleteBucketCors_with_valid_parameters() {
        ObsClient obsClient = getObsClient();

        Integer responseCodeForTest = 204;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("DELETE")
                        .withPath("")
                        .withQueryStringParameter("cors"))
                .respond(response().withStatusCode(responseCodeForTest));

        BaseBucketRequest request = new BaseBucketRequest(bucketNameForTest);
        HeaderResponse response = obsClient.deleteBucketCors(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    // ------------------------------ 辅助方法 ------------------------------
    private static ObsClient getObsClient() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        return obsClient;
    }

    private static void setupMockPutResponse(int statusCode) {
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(statusCode));
    }

    private static void setupMockCorsGetResponse(String xmlBody) {
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("cors"))
                .respond(response()
                        .withStatusCode(200)
                        .withHeader("Content-Type", "application/xml")
                        .withBody(xmlBody));
    }

    private static BucketCors createBucketCors() {
        BucketCors bucketCors = new BucketCors();
        BucketCorsRule rule = new BucketCorsRule();
        rule.setId("test-rule");
        rule.setAllowedOrigin(Arrays.asList("http://www.a.com", "http://www.b.com"));
        rule.setAllowedMethod(Arrays.asList("GET", "HEAD", "PUT"));
        rule.setAllowedHeader(Collections.singletonList("x-obs-header"));
        rule.setExposeHeader(Collections.singletonList("x-obs-expose-header"));
        rule.setMaxAgeSecond(100);
        bucketCors.setRules(Collections.singletonList(rule));
        return bucketCors;
    }
}
