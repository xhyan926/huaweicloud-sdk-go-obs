/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test;

import static com.obs.services.internal.Constants.SYMLINK_HEADER;
import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.services.ObsClient;
import com.obs.services.model.symlink.GetSymlinkRequest;
import com.obs.services.model.symlink.GetSymlinkResult;
import com.obs.test.TestTools;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

public class GetSymlinkIT {

    @org.junit.Rule
    public ExpectedException expectedException = ExpectedException.none();
    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";
    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "bucket";
    public static String objectKeyForTest = "object";
    public static String versionIdForTest = "versionId";
    private final GetSymlinkRequest getSymlinkRequestForTest =
        new GetSymlinkRequest(bucketNameForTest, objectKeyForTest, versionIdForTest);
    private final GetSymlinkRequest getGetSymlinkRequestWithoutVersionIdForTest =
            new GetSymlinkRequest(bucketNameForTest,objectKeyForTest);

    @BeforeClass
    public static void setMockServer(){
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
    public void should_succeed_when_getSymlink_with_version_id() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        Integer responseCodeForTest = 200;
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        String symlinkTargetForTest = "symlink_target_for_test";
        mockServer
            .when(request().withMethod("GET").withPath(""))
            .respond(response()
                .withStatusCode(responseCodeForTest)
                .withHeader("x-obs-" + SYMLINK_HEADER, symlinkTargetForTest));
        getSymlinkRequestForTest.setVersionId(versionIdForTest);
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequestForTest);
        Assert.assertEquals(responseCodeForTest.intValue(), getSymlinkResult.getStatusCode());
        Assert.assertEquals(symlinkTargetForTest, getSymlinkResult.getSymlinkTarget());
    }
    @Test
    public void should_succeed_when_getSymlink_without_version_id() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        Integer responseCodeForTest = 200;
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        mockServer
            .when(request().withMethod("GET").withPath(""))
            .respond(response().withStatusCode(responseCodeForTest));
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getGetSymlinkRequestWithoutVersionIdForTest);
        Assert.assertEquals(responseCodeForTest.intValue(), getSymlinkResult.getStatusCode());
    }
}
