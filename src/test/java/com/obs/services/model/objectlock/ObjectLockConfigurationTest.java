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

public class ObjectLockConfigurationTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-objectlock";

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

    // ------------------------------ setObjectLockConfiguration 测试 ------------------------------
    @Test
    public void should_throw_exception_when_setObjectLockConfigurationRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("SetObjectLockConfigurationRequest is null");
        obsClient.setObjectLockConfiguration(null);
    }

    @Test
    public void should_throw_exception_when_setObjectLockConfigurationRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        ObjectLockConfiguration config = new ObjectLockConfiguration("Enabled", null);
        SetObjectLockConfigurationRequest request = new SetObjectLockConfigurationRequest(null, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setObjectLockConfiguration(request);
    }

    @Test
    public void should_succeed_when_setObjectLockConfiguration_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        DefaultRetention retention = new DefaultRetention("COMPLIANCE", 30, null);
        ObjectLockRule rule = new ObjectLockRule(retention);
        ObjectLockConfiguration config = new ObjectLockConfiguration("Enabled", rule);
        SetObjectLockConfigurationRequest request =
            new SetObjectLockConfigurationRequest(bucketNameForTest, config);
        HeaderResponse response = obsClient.setObjectLockConfiguration(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    // ------------------------------ getObjectLockConfiguration 测试 ------------------------------
    @Test
    public void should_throw_exception_when_getObjectLockConfigurationRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("GetObjectLockConfigurationRequest is null");
        obsClient.getObjectLockConfiguration(null);
    }

    @Test
    public void should_throw_exception_when_getObjectLockConfigurationRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        GetObjectLockConfigurationRequest request = new GetObjectLockConfigurationRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.getObjectLockConfiguration(request);
    }

    @Test
    public void should_succeed_when_getObjectLockConfiguration_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseXml = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
                + "<ObjectLockConfiguration>"
                + "  <ObjectLockEnabled>Enabled</ObjectLockEnabled>"
                + "  <Rule>"
                + "    <DefaultRetention>"
                + "      <Mode>COMPLIANCE</Mode>"
                + "      <Days>30</Days>"
                + "    </DefaultRetention>"
                + "  </Rule>"
                + "</ObjectLockConfiguration>";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("object-lock"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/xml;charset=utf-8")
                        .withBody(validResponseXml));

        GetObjectLockConfigurationRequest request =
            new GetObjectLockConfigurationRequest(bucketNameForTest);
        GetObjectLockConfigurationResult result = obsClient.getObjectLockConfiguration(request);

        Assert.assertNotNull(result);
        Assert.assertNotNull(result.getObjectLockConfiguration());
        Assert.assertEquals("Enabled", result.getObjectLockConfiguration().getObjectLockEnabled());
        Assert.assertNotNull(result.getObjectLockConfiguration().getRule());
        Assert.assertNotNull(result.getObjectLockConfiguration().getRule().getDefaultRetention());
        Assert.assertEquals("COMPLIANCE",
            result.getObjectLockConfiguration().getRule().getDefaultRetention().getMode());
        Assert.assertEquals(Integer.valueOf(30),
            result.getObjectLockConfiguration().getRule().getDefaultRetention().getDays());
    }

    @Test
    public void should_succeed_when_getObjectLockConfiguration_without_rule() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        String validResponseXml = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
                + "<ObjectLockConfiguration>"
                + "  <ObjectLockEnabled>Enabled</ObjectLockEnabled>"
                + "</ObjectLockConfiguration>";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("object-lock"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/xml;charset=utf-8")
                        .withBody(validResponseXml));

        GetObjectLockConfigurationRequest request =
            new GetObjectLockConfigurationRequest(bucketNameForTest);
        GetObjectLockConfigurationResult result = obsClient.getObjectLockConfiguration(request);

        Assert.assertNotNull(result);
        Assert.assertEquals("Enabled", result.getObjectLockConfiguration().getObjectLockEnabled());
        Assert.assertNull(result.getObjectLockConfiguration().getRule());
    }
}
