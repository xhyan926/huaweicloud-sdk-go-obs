/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test;

import static com.obs.services.internal.Constants.SYMLINK_HEADER;
import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.services.ObsClient;
import com.obs.services.model.GetObjectRequest;
import com.obs.services.model.ListBucketsResult;
import com.obs.services.model.ListObjectsRequest;
import com.obs.services.model.symlink.GetSymlinkRequest;
import com.obs.services.model.symlink.GetSymlinkResult;
import com.obs.test.TestTools;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

public class JavaSdkEncodedAuthorizationMessageIT {

    @org.junit.Rule
    public ExpectedException expectedException = ExpectedException.none();
    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";
    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "bucket";
    public static String objectKeyForTest = "object";
    private final ListObjectsRequest listObjectsRequestForTest = new ListObjectsRequest(bucketNameForTest);

    // 定义 EncodedAuthorizationMessage 常量
    private static final String ENCODED_AUTHORIZATION_MESSAGE = "HY0L8G110e3rcfgxfKVP+AK+S33eYp/rHQ4I0kJed9WC+osCXOXLakWL01p81EEDCztcri5QsEcRUKbDRzRtPyszmkvixO8yo2mIZv/SmGuWyAkGZ3/21vJiSwi7ZII1ILkkGsxiIHcf30HbywQfiyf1J5IyJfMqaNSZE1QbpF/9TC91XELxiaQz8W/+WXo9xFq4B081d50mkraguyrVLGRU+bX4QQx0462IbXo35FpMFNHT7LTGhiDDeyksAscrwaYmLp+pt/";

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

    @Test
    public void tc_alpha_java_js_sdk_EncodedAuthorizationMessage_001() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        Integer responseCodeForTest = 403;
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        String xmlBody = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n" +
                "<Error>\n" +
                "    <Code>AccessDenied</Code>\n" +
                "    <Message>Access Denied</Message>\n" +
                "    <EncodedAuthorizationMessage>" + ENCODED_AUTHORIZATION_MESSAGE + "</EncodedAuthorizationMessage>\n" +
                "    <RequestId>000001980BE1F5BC53066C556F3A0B25</RequestId>\n" +
                "    <HostId>Z9v+cC1sRnaNlw6x0vi8pxxYA0VnKxbYHUPAFpnxkX8sLV44u5b02Z+ai1n2wCnR</HostId>\n" +
                "</Error>";

        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withBody(xmlBody));
        try {
            obsClient.listObjects(listObjectsRequestForTest);
        } catch (com.obs.services.exception.ObsException e) {
            Assert.assertEquals(responseCodeForTest.intValue(), e.getResponseCode());
            Assert.assertEquals("AccessDenied", e.getErrorCode());
            Assert.assertEquals(ENCODED_AUTHORIZATION_MESSAGE, e.getEncodedAuthorizationMessage());
        }

    }
}