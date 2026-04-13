/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.internal;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.services.internal.security.EcsSecurityUtils;

import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

import java.io.IOException;

public class EcsSecurityUtilsTest {
    private static ClientAndServer mockServer;
    public static final String METADATA_API_TOKEN_RESOURCE_PATH = "/meta-data/latest/api/token";
    public static final String OPENSTACK_SECURITY_KEY_RESOURCE_PATH = "/openstack/latest/securitykey";
    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";

    @org.junit.Rule
    public ExpectedException expectedException = ExpectedException.none();

    @BeforeClass
    public static void setMockServer(){
        // 启动 MockServer
        mockServer = ClientAndServer.startClientAndServer();
        System.setProperty(PROXY_HOST_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_PROPERTY_NAME, "" + mockServer.getLocalPort());
    }

    @AfterClass
    public static void clearEnv() {
        // 关闭 MockServer
        mockServer.close();
        System.clearProperty(PROXY_HOST_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_PROPERTY_NAME);
    }

    @Test
    public void should_getSecurityKeyInfoWithDetail_successfully() throws IOException {
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        mockServer
            .when(request().withMethod("PUT").withPath(METADATA_API_TOKEN_RESOURCE_PATH))
            .respond(response().withStatusCode(200).withBody("test_metadata_latest_api_token"));
        mockServer
            .when(request().withMethod("GET").withPath(OPENSTACK_SECURITY_KEY_RESOURCE_PATH))
            .respond(response().withStatusCode(200).withBody("test_openstack_latest_securitykey"));
        EcsSecurityUtils.getSecurityKeyInfoWithDetail(60);
    }

    @Test
    public void should_getSecurityKeyInfoWithDetail_successfully_when_getMetadataApiToken_not_supported()
            throws IOException {
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        int[] notSupportedCodes = {404, 405};
        for (int notSupportedCode : notSupportedCodes) {
            mockServer
                .when(request().withMethod("PUT").withPath(METADATA_API_TOKEN_RESOURCE_PATH))
                .respond(response().withStatusCode(notSupportedCode).withBody("test_metadata_latest_api_token"));
            mockServer
                .when(request().withMethod("GET").withPath(OPENSTACK_SECURITY_KEY_RESOURCE_PATH))
                .respond(response().withStatusCode(200).withBody("test_openstack_latest_securitykey"));
            EcsSecurityUtils.getSecurityKeyInfoWithDetail();
        }
    }

    @Test
    public void should_throw_exception_when_getMetadataApiToken_failed()
        throws IOException {
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        final int failedCode = 403;
        mockServer
            .when(request().withMethod("PUT").withPath(METADATA_API_TOKEN_RESOURCE_PATH))
            .respond(response().withStatusCode(failedCode).withBody("test_metadata_latest_api_token"));
        mockServer
            .when(request().withMethod("GET").withPath(OPENSTACK_SECURITY_KEY_RESOURCE_PATH))
            .respond(response().withStatusCode(failedCode).withBody("test_openstack_latest_securitykey"));

        expectedException.expect(java.lang.IllegalArgumentException.class);
        expectedException.expectMessage("Get X-Metadata-Token with "
            + "X-Metadata-Token-Ttl-Seconds:21600 from ECS failed");

        EcsSecurityUtils.getSecurityKeyInfoWithDetail();
    }

    @Test
    public void should_throw_exception_when_getSecurityKeyInfoWithDetail_IMDSv2_failed()
        throws IOException {
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        final int notSupportedCode = 404;
        mockServer
            .when(request().withMethod("PUT").withPath(METADATA_API_TOKEN_RESOURCE_PATH))
            .respond(response().withStatusCode(200).withBody("test_metadata_latest_api_token"));
        mockServer
            .when(request().withMethod("GET").withPath(OPENSTACK_SECURITY_KEY_RESOURCE_PATH))
            .respond(response().withStatusCode(notSupportedCode).withBody("test_openstack_latest_securitykey"));
        expectedException.expect(java.lang.IllegalArgumentException.class);
        expectedException.expectMessage("Get securityKey by X-Metadata-Token from ECS failed");

        EcsSecurityUtils.getSecurityKeyInfoWithDetail();
    }
}
