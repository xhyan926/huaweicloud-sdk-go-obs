/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.services.model.cors;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.aitool.AIGenerated;
import com.obs.services.ObsClient;
import com.obs.services.model.OptionsInfoRequest;
import com.obs.services.model.OptionsInfoResult;
import com.obs.test.TestTools;

import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

import java.util.Arrays;

public class OptionsBucketTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-options";

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

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试OPTIONS桶预检请求为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_optionsBucketRequest_is_null() {
        ObsClient obsClient = getObsClient();

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("OptionsInfoRequest is null");
        obsClient.optionsBucket((OptionsInfoRequest) null);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试OPTIONS桶预检请求桶名为null时抛出IllegalArgumentException")
    public void should_throw_exception_when_optionsBucketRequest_bucketName_is_null() {
        ObsClient obsClient = getObsClient();

        OptionsInfoRequest request = new OptionsInfoRequest();

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.optionsBucket(request);
    }

    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "测试OPTIONS桶预检正常请求验证CORS响应头解析")
    public void should_succeed_when_optionsBucket_with_valid_parameters() {
        ObsClient obsClient = getObsClient();

        setupMockOptionsResponse();

        OptionsInfoRequest request = new OptionsInfoRequest(bucketNameForTest);
        request.setOrigin("http://www.example.com");
        request.setRequestMethod(Arrays.asList("GET", "PUT"));
        request.setRequestHeaders(Arrays.asList("Authorization", "Content-Type"));

        OptionsInfoResult result = obsClient.optionsBucket(request);

        Assert.assertNotNull(result);
        Assert.assertEquals("http://www.example.com", result.getAllowOrigin());
        Assert.assertEquals(200, result.getMaxAge());
        Assert.assertTrue(result.getAllowMethods().contains("GET"));
        Assert.assertTrue(result.getAllowMethods().contains("PUT"));
        Assert.assertTrue(result.getAllowHeaders().contains("Authorization"));
        Assert.assertTrue(result.getAllowHeaders().contains("Content-Type"));
        Assert.assertTrue(result.getExposeHeaders().contains("x-obs-request-id"));
        Assert.assertTrue(result.getExposeHeaders().contains("x-obs-id-2"));
    }

    // ------------------------------ 辅助方法 ------------------------------
    private static ObsClient getObsClient() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        return obsClient;
    }

    private static void setupMockOptionsResponse() {
        mockServer.reset();
        mockServer.when(request().withMethod("OPTIONS").withPath(""))
                .respond(response()
                        .withStatusCode(200)
                        .withHeader("Access-Control-Allow-Origin", "http://www.example.com")
                        .withHeader("Access-Control-Allow-Methods", "GET")
                        .withHeader("Access-Control-Allow-Methods", "PUT")
                        .withHeader("Access-Control-Allow-Headers", "Authorization")
                        .withHeader("Access-Control-Allow-Headers", "Content-Type")
                        .withHeader("Access-Control-Max-Age", "200")
                        .withHeader("Access-Control-Expose-Headers", "x-obs-request-id")
                        .withHeader("Access-Control-Expose-Headers", "x-obs-id-2"));
    }
}
