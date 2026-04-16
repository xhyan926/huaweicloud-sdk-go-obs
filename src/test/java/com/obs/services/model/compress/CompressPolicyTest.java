/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.services.model.compress;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

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

public class CompressPolicyTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-compresspolicy";

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

    // ------------------------------ setBucketCompressPolicy 测试 ------------------------------
    @Test
    public void should_throw_exception_when_setBucketCompressPolicyRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("SetBucketCompressPolicyRequest is null");
        obsClient.setBucketCompressPolicy(null);
    }

    @Test
    public void should_throw_exception_when_setBucketCompressPolicyRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        CompressPolicyRule rule = new CompressPolicyRule("rule1", "projectId", "agency",
            Arrays.asList("ObjectCreated:*"), ".zip", 0);
        CompressPolicyConfiguration config = new CompressPolicyConfiguration(Arrays.asList(rule));
        SetBucketCompressPolicyRequest request = new SetBucketCompressPolicyRequest(null, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setBucketCompressPolicy(request);
    }

    @Test
    public void should_succeed_when_setBucketCompressPolicy_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 201;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        CompressPolicyRule rule = new CompressPolicyRule("rule1", "projectId", "agency",
            Arrays.asList("ObjectCreated:*"), ".zip", 0);
        rule.setPrefix("decompress/");
        rule.setDecompresspath("after-decompress/");
        rule.setPolicytype("decompress");
        CompressPolicyConfiguration config = new CompressPolicyConfiguration(Arrays.asList(rule));
        SetBucketCompressPolicyRequest request =
            new SetBucketCompressPolicyRequest(bucketNameForTest, config);
        HeaderResponse response = obsClient.setBucketCompressPolicy(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    @Test
    public void should_succeed_when_setBucketCompressPolicy_returns_200() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        CompressPolicyRule rule = new CompressPolicyRule("rule1", "projectId", "agency",
            Arrays.asList("ObjectCreated:*"), ".zip", 1);
        CompressPolicyConfiguration config = new CompressPolicyConfiguration(Arrays.asList(rule));
        SetBucketCompressPolicyRequest request =
            new SetBucketCompressPolicyRequest(bucketNameForTest, config);
        HeaderResponse response = obsClient.setBucketCompressPolicy(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    // ------------------------------ getBucketCompressPolicy 测试 ------------------------------
    @Test
    public void should_throw_exception_when_getBucketCompressPolicyRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("GetBucketCompressPolicyRequest is null");
        obsClient.getBucketCompressPolicy(null);
    }

    @Test
    public void should_throw_exception_when_getBucketCompressPolicyRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        GetBucketCompressPolicyRequest request = new GetBucketCompressPolicyRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.getBucketCompressPolicy(request);
    }

    @Test
    public void should_succeed_when_getBucketCompressPolicy_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseJson = "{"
                + "\"rules\":[{"
                + "\"id\":\"rule1\","
                + "\"project\":\"projectId\","
                + "\"agency\":\"testagency\","
                + "\"events\":[\"ObjectCreated:*\"],"
                + "\"prefix\":\"decompress\","
                + "\"suffix\":\".zip\","
                + "\"overwrite\":0,"
                + "\"decompresspath\":\"after-decompress/\""
                + "}]"
                + "}";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("obscompresspolicy"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/json")
                        .withBody(validResponseJson));

        GetBucketCompressPolicyRequest request = new GetBucketCompressPolicyRequest(bucketNameForTest);
        GetBucketCompressPolicyResult result = obsClient.getBucketCompressPolicy(request);

        Assert.assertNotNull(result);
        Assert.assertNotNull(result.getCompressPolicyConfiguration());
        Assert.assertNotNull(result.getCompressPolicyConfiguration().getRules());
        Assert.assertEquals(1, result.getCompressPolicyConfiguration().getRules().size());

        CompressPolicyRule rule = result.getCompressPolicyConfiguration().getRules().get(0);
        Assert.assertEquals("rule1", rule.getId());
        Assert.assertEquals("projectId", rule.getProject());
        Assert.assertEquals("testagency", rule.getAgency());
        Assert.assertEquals(1, rule.getEvents().size());
        Assert.assertEquals("ObjectCreated:*", rule.getEvents().get(0));
        Assert.assertEquals("decompress", rule.getPrefix());
        Assert.assertEquals(".zip", rule.getSuffix());
        Assert.assertEquals(Integer.valueOf(0), rule.getOverwrite());
        Assert.assertEquals("after-decompress/", rule.getDecompresspath());
    }

    @Test
    public void should_succeed_when_getBucketCompressPolicy_with_multiple_rules() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseJson = "{"
                + "\"rules\":[{"
                + "\"id\":\"rule1\","
                + "\"project\":\"projectId1\","
                + "\"agency\":\"agency1\","
                + "\"events\":[\"ObjectCreated:*\"],"
                + "\"suffix\":\".zip\","
                + "\"overwrite\":0"
                + "},{"
                + "\"id\":\"rule2\","
                + "\"project\":\"projectId2\","
                + "\"agency\":\"agency2\","
                + "\"events\":[\"ObjectCreated:Put\",\"ObjectCreated:Post\"],"
                + "\"prefix\":\"folder\","
                + "\"suffix\":\".zip\","
                + "\"overwrite\":2,"
                + "\"decompresspath\":\"output/\""
                + "}]"
                + "}";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("obscompresspolicy"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/json")
                        .withBody(validResponseJson));

        GetBucketCompressPolicyRequest request = new GetBucketCompressPolicyRequest(bucketNameForTest);
        GetBucketCompressPolicyResult result = obsClient.getBucketCompressPolicy(request);

        Assert.assertNotNull(result);
        Assert.assertEquals(2, result.getCompressPolicyConfiguration().getRules().size());

        CompressPolicyRule rule1 = result.getCompressPolicyConfiguration().getRules().get(0);
        Assert.assertEquals("rule1", rule1.getId());
        Assert.assertEquals(Integer.valueOf(0), rule1.getOverwrite());

        CompressPolicyRule rule2 = result.getCompressPolicyConfiguration().getRules().get(1);
        Assert.assertEquals("rule2", rule2.getId());
        Assert.assertEquals(Integer.valueOf(2), rule2.getOverwrite());
        Assert.assertEquals(2, rule2.getEvents().size());
        Assert.assertEquals("output/", rule2.getDecompresspath());
    }

    // ------------------------------ deleteBucketCompressPolicy 测试 ------------------------------
    @Test
    public void should_throw_exception_when_deleteBucketCompressPolicyRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("DeleteBucketCompressPolicyRequest is null");
        obsClient.deleteBucketCompressPolicy(null);
    }

    @Test
    public void should_throw_exception_when_deleteBucketCompressPolicyRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        DeleteBucketCompressPolicyRequest request = new DeleteBucketCompressPolicyRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.deleteBucketCompressPolicy(request);
    }

    @Test
    public void should_succeed_when_deleteBucketCompressPolicy_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 204;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("DELETE")
                        .withPath("")
                        .withQueryStringParameter("obscompresspolicy"))
                .respond(response().withStatusCode(responseCodeForTest));

        DeleteBucketCompressPolicyRequest request =
            new DeleteBucketCompressPolicyRequest(bucketNameForTest);
        HeaderResponse response = obsClient.deleteBucketCompressPolicy(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }
}
