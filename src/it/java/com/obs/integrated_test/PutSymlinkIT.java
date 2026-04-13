/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.services.ObsClient;
import com.obs.services.model.AccessControlList;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.symlink.PutSymlinkRequest;
import com.obs.test.TestTools;

import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

public class PutSymlinkIT {

    @org.junit.Rule
    public ExpectedException expectedException = ExpectedException.none();
    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";
    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "bucket";
    public static String objectKeyForTest = "object";
    public static String symlinkTargetForTest = "symlinkTarget";
    public static AccessControlList accessControlListForTest = AccessControlList.REST_CANNED_PUBLIC_READ;
    public static ObjectMetadata objectMetadataForTest = new ObjectMetadata();

    private final PutSymlinkRequest putSymlinkRequestForTest =
        new PutSymlinkRequest(bucketNameForTest,
            objectKeyForTest,
            symlinkTargetForTest,
            accessControlListForTest,
            objectMetadataForTest);
    private final PutSymlinkRequest putSymlinkRequestWithoutMetadataAndAcl =
            new PutSymlinkRequest(bucketNameForTest,objectKeyForTest,symlinkTargetForTest);

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
    public void should_throw_exception_when_putSymlinkRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(java.lang.IllegalArgumentException.class);
        expectedException.expectMessage("PutSymlinkRequest is null");
        obsClient.putSymlink(null);
    }
    @Test
    public void should_throw_exception_when_putSymlinkRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(java.lang.IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        putSymlinkRequestForTest.setBucketName(null);
        obsClient.putSymlink(putSymlinkRequestForTest);
    }
    @Test
    public void should_throw_exception_when_putSymlinkRequest_objectKey_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        expectedException.expect(java.lang.IllegalArgumentException.class);
        expectedException.expectMessage("objectKey is null");
        putSymlinkRequestForTest.setBucketName(bucketNameForTest);
        putSymlinkRequestForTest.setObjectKey(null);
        obsClient.putSymlink(putSymlinkRequestForTest);
    }
    @Test
    public void should_succeed_when_putSymlink() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        Integer responseCodeForTest = 200;
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        mockServer
            .when(request().withMethod("PUT").withPath(""))
            .respond(response().withStatusCode(responseCodeForTest));
        putSymlinkRequestForTest.setBucketName(bucketNameForTest);
        putSymlinkRequestForTest.setObjectKey(objectKeyForTest);
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequestForTest);
        Assert.assertEquals(responseCodeForTest.intValue(), headerResponse.getStatusCode());
    }
    @Test
    public void should_succeed_when_metadata_and_acl_is_null_putSymlink(){
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        Integer responseCodeForTest = 200;
        // 清除旧的模拟请求、响应规则
        mockServer.reset();
        // 设置新的模拟请求、响应规则
        mockServer
                .when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequestWithoutMetadataAndAcl);
        Assert.assertEquals(responseCodeForTest.intValue(), headerResponse.getStatusCode());
    }
}
