/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.services.model.objectlock;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.test.TestTools;

import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

public class ObjectRetentionTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-retention";
    public static String objectKeyForTest = "test-object-retention";

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

    // ------------------------------ setObjectRetention 测试 ------------------------------

    @Test
    public void should_throw_exception_when_setObjectRetentionRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("SetObjectRetentionRequest is null");
        obsClient.setObjectRetention(null);
    }

    @Test
    public void should_throw_exception_when_setObjectRetentionRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        ObjectRetention retention = new ObjectRetention("COMPLIANCE", 1767225600000L);
        SetObjectRetentionRequest request =
            new SetObjectRetentionRequest(null, objectKeyForTest, retention);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setObjectRetention(request);
    }

    @Test
    public void should_throw_exception_when_setObjectRetentionRequest_objectKey_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        ObjectRetention retention = new ObjectRetention("COMPLIANCE", 1767225600000L);
        SetObjectRetentionRequest request =
            new SetObjectRetentionRequest(bucketNameForTest, null, retention);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("objectKey is null");
        obsClient.setObjectRetention(request);
    }

    @Test
    public void should_succeed_when_setObjectRetention_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath("/" + objectKeyForTest))
                .respond(response().withStatusCode(responseCodeForTest));

        ObjectRetention retention = new ObjectRetention("COMPLIANCE", 1767225600000L);
        SetObjectRetentionRequest request =
            new SetObjectRetentionRequest(bucketNameForTest, objectKeyForTest, retention);
        HeaderResponse response = obsClient.setObjectRetention(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    @Test
    public void should_succeed_when_setObjectRetention_with_versionId() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath("/" + objectKeyForTest))
                .respond(response().withStatusCode(responseCodeForTest));

        ObjectRetention retention = new ObjectRetention("COMPLIANCE", 1767225600000L);
        SetObjectRetentionRequest request =
            new SetObjectRetentionRequest(bucketNameForTest, objectKeyForTest, retention, "v1");
        HeaderResponse response = obsClient.setObjectRetention(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }
}
