package com.obs;

import com.obs.log.ILogger;
import com.obs.log.LoggerBuilder;
import com.obs.services.RefreshAkSkStsTokenProvider;
import com.obs.services.model.ISecurityKey;
import com.obs.test.tools.PropertiesTools;
import okhttp3.*;
import org.junit.Before;
import org.junit.Test;

import java.io.File;
import java.io.IOException;
import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.mockito.Mockito.*;

public class RefreshAkSkStsTokenProviderTest {

    private static final ILogger log = LoggerBuilder.getLogger(RefreshAkSkStsTokenProviderTest.class);

    private static final File file = new File("./app/src/test/resource/test_data.properties");
    private RefreshAkSkStsTokenProvider tokenProvider;
    private OkHttpClient httpClientMock;
    private ScheduledExecutorService schedulerMock;
    private String iamDomain;
    private String iamUser;
    private String iamPassword;
    private String iamEndpoint;

    @Before
    public void setUp() throws IOException, NoSuchFieldException, IllegalAccessException {
        log.info("Setting up RefreshAkSkStsTokenProviderTest");

        iamDomain = PropertiesTools.getInstance(file).getProperties("environment.iamDomain");
        iamUser = PropertiesTools.getInstance(file).getProperties("environment.iamUser");
        iamPassword = PropertiesTools.getInstance(file).getProperties("environment.iamPassword");
        iamEndpoint = PropertiesTools.getInstance(file).getProperties("environment.iamEndpoint");

        httpClientMock = mock(OkHttpClient.class);
        schedulerMock = mock(ScheduledExecutorService.class);

        tokenProvider = new RefreshAkSkStsTokenProvider(iamDomain, iamUser, iamPassword, iamEndpoint, httpClientMock);

        setPrivateField(tokenProvider, "scheduler", schedulerMock);

        log.info("Completed setup for RefreshAkSkStsTokenProviderTest");
    }

    private void setPrivateField(Object object, String fieldName, Object value) throws NoSuchFieldException, IllegalAccessException {
        Field field = object.getClass().getDeclaredField(fieldName);
        field.setAccessible(true);
        field.set(object, value);
    }

    @Test
    public void testGetSecurityKey() throws Exception {
        log.info("Running testGetSecurityKey");

        // Mocking the first call to get the subject token
        String mockSubjectToken = "mockSubjectToken";
        Response subjectTokenResponse = createMockResponse(200, "OK", "X-Subject-Token", mockSubjectToken);

        // Mocking the second call to get the security token
        String mockAccessKey = "newAccessKey";
        String mockSecretKey = "newSecretKey";
        String mockSecurityToken = "newSecurityToken";
        String securityTokenResponseJson = createSecurityTokenResponseJson(mockAccessKey, mockSecretKey, mockSecurityToken);
        Response securityTokenResponse = createMockResponse(200, "OK", securityTokenResponseJson);

        // Mock OkHttpClient's behavior
        Call mockCallSubjectToken = mock(Call.class);
        Call mockCallSecurityToken = mock(Call.class);
        when(httpClientMock.newCall(any(Request.class)))
                .thenReturn(mockCallSubjectToken)
                .thenReturn(mockCallSecurityToken);
        when(mockCallSubjectToken.execute()).thenReturn(subjectTokenResponse);
        when(mockCallSecurityToken.execute()).thenReturn(securityTokenResponse);

        // Fetch and verify the security key
        ISecurityKey securityKey = tokenProvider.getSecurityKey();
        assertNotNull("Security key should not be null", securityKey);
        assertEquals("Access key mismatch", mockAccessKey, securityKey.getAccessKey());
        assertEquals("Secret key mismatch", mockSecretKey, securityKey.getSecretKey());
        assertEquals("Security token mismatch", mockSecurityToken, securityKey.getSecurityToken());

        log.info("Completed testGetSecurityKey");
    }

    @Test
    public void testScheduledTaskIsStarted() throws NoSuchMethodException, InvocationTargetException, IllegalAccessException {
        log.info("Running testScheduledTaskIsStarted");

        // Invoke the private startTokenRefreshTask method
        Method startTokenRefreshTask = tokenProvider.getClass().getDeclaredMethod("startTokenRefreshTask");
        startTokenRefreshTask.setAccessible(true);
        startTokenRefreshTask.invoke(tokenProvider);

        // Verify scheduler's periodic scheduling
        verify(schedulerMock).scheduleAtFixedRate(any(Runnable.class), eq(0L), eq(10L), eq(TimeUnit.MINUTES));

        log.info("Completed testScheduledTaskIsStarted");
    }

    // Helper methods to simplify test setup and mocking

    private Response createMockResponse(int code, String message, String headerName, String headerValue) {
        return new Response.Builder()
                .code(code)
                .message(message)
                .protocol(Protocol.HTTP_1_1)
                .request(new Request.Builder().url("http://example.com").build())
                .header(headerName, headerValue)
                .body(ResponseBody.create("", MediaType.get("application/json; charset=utf-8")))
                .build();
    }

    private Response createMockResponse(int code, String message, String responseBodyJson) {
        return new Response.Builder()
                .code(code)
                .message(message)
                .protocol(Protocol.HTTP_1_1)
                .request(new Request.Builder().url("http://example.com").build())
                .body(ResponseBody.create(responseBodyJson, MediaType.get("application/json; charset=utf-8")))
                .build();
    }

    private String createSecurityTokenResponseJson(String accessKey, String secretKey, String securityToken) {
        return String.format("{ \"access\":\"%s\", \"secret\":\"%s\", \"securitytoken\":\"%s\" }", accessKey, secretKey, securityToken);
    }

    //@Test
    public void testGetSecurityKeyWithRealRequest() {
        log.info("Running testGetSecurityKeyWithRealRequest");

        tokenProvider = new RefreshAkSkStsTokenProvider(iamDomain, iamUser, iamPassword, iamEndpoint,httpClientMock);
        ISecurityKey securityKey = tokenProvider.getSecurityKey();

        assertNotNull("Access key should not be null", securityKey.getAccessKey());
        assertNotNull("Secret key should not be null", securityKey.getSecretKey());
        assertNotNull("Security token should not be null", securityKey.getSecurityToken());

        log.info("Successfully fetched security key with access key: " + securityKey.getAccessKey());
    }

}
