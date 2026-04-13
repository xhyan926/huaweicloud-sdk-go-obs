/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.internal;

import static com.obs.test.SSLTestUtils.trustAllManager;
import static org.mockito.ArgumentMatchers.any;

import com.obs.services.ObsClient;
import com.obs.services.model.ListBucketsRequest;
import com.obs.services.model.ListBucketsResult;
import com.obs.test.TestTools;

import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.mockito.MockedStatic;
import org.mockito.Mockito;

import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;


import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;

public class SetSSLContextTest {
    static final String TEST_TLS_VERSION = "TLSv1.2";
    static final String TEST_TLS_NOT_SUPPORTED_ERROR_MESSAGE = "Test error message: "
        + TEST_TLS_VERSION + " not allowed in test!";
    protected boolean isMockedGetInstanceCalled = false;
    protected boolean isMockedGetSocketFactoryCalled = false;
    protected ListBucketsRequest listBucketsRequest = new ListBucketsRequest();
    protected SSLContext sslContext;
    @Before
    public void prepare() throws NoSuchAlgorithmException, KeyManagementException {
        sslContext = SSLContext.getInstance(TEST_TLS_VERSION);
        sslContext.init(null, new TrustManager[] { trustAllManager }, new SecureRandom());
    }
    @Test
    public void test_SetSSLContextTest_01() throws Exception {
        isMockedGetInstanceCalled = false;
        String sslProvider = "SunJSSE";
        SSLContext sslContextWithProvider = SSLContext.getInstance(TEST_TLS_VERSION, sslProvider);
        try (MockedStatic<SSLContext> mocked = Mockito.mockStatic(SSLContext.class)) {
            SSLContext mockedSSLContext = Mockito.mock(SSLContext.class);
            // 通过mock模拟SSLContext.getInstance函数，
            // 根据用户设置的sslProvider调用SSLContext.getInstance(String protocol, String provider)成功，
            // 且参数中的sslProvider 与 用户设置的sslProvider 相同
            mocked.when(() -> SSLContext.getInstance(any(String.class), any(String.class))).thenAnswer(invocation -> {
                isMockedGetInstanceCalled = true;
                Assert.assertEquals(sslProvider, invocation.getArgument(1).toString());
                return sslContextWithProvider;
            });
            ObsClient obsClient =
                TestTools.getPipelineEnvironmentWithCustomisedSSLContext(mockedSSLContext, sslProvider);
            assert obsClient != null;
            ListBucketsResult listBucketsResult = obsClient.listBucketsV2(listBucketsRequest);
            Assert.assertEquals(200, listBucketsResult.getStatusCode());
            obsClient.close();
            Assert.assertTrue(isMockedGetInstanceCalled);
        }
    }
    @Test
    public void test_SetSSLContextTest_02() throws Exception {
        isMockedGetInstanceCalled = false;
        String wrongProvider = "wrongProvider";
        SSLContext sslContextWithProvider = SSLContext.getInstance(TEST_TLS_VERSION);
        try (MockedStatic<SSLContext> mocked = Mockito.mockStatic(SSLContext.class)) {
            SSLContext mockedSSLContext = Mockito.mock(SSLContext.class);
            // 通过mock模拟SSLContext.getInstance函数，模拟用户设置了无效的的sslProvider时，
            // 自动调用SSLContext.getInstance(String protocol)成功
            mocked.when(() -> SSLContext.getInstance(any(String.class))).thenAnswer(invocation -> {
                isMockedGetInstanceCalled = true;
                return sslContextWithProvider;
            });
            ObsClient obsClient =
                TestTools.getPipelineEnvironmentWithCustomisedSSLContext(mockedSSLContext, wrongProvider);
            assert obsClient != null;
            ListBucketsResult listBucketsResult = obsClient.listBucketsV2(listBucketsRequest);
            Assert.assertEquals(200, listBucketsResult.getStatusCode());
            obsClient.close();
            Assert.assertTrue(isMockedGetInstanceCalled);
        }
    }
    @Test
    public void test_SetSSLContextTest_03() throws Exception {
        isMockedGetInstanceCalled = false;
        SSLContext sslContextWithProvider = SSLContext.getInstance(TEST_TLS_VERSION);
        try (MockedStatic<SSLContext> mocked = Mockito.mockStatic(SSLContext.class)) {
            SSLContext mockedSSLContext = Mockito.mock(SSLContext.class);
            // 通过mock模拟SSLContext.getInstance函数，模拟用户未设置sslProvider时，
            // 自动调用SSLContext.getInstance(String protocol)成功
            mocked.when(() -> SSLContext.getInstance(any(String.class))).thenAnswer(invocation -> {
                isMockedGetInstanceCalled = true;
                return sslContextWithProvider;
            });
            ObsClient obsClient =
                TestTools.getPipelineEnvironmentWithCustomisedSSLContext(mockedSSLContext);
            assert obsClient != null;
            ListBucketsResult listBucketsResult = obsClient.listBucketsV2(listBucketsRequest);
            Assert.assertEquals(200, listBucketsResult.getStatusCode());
            obsClient.close();
            Assert.assertTrue(isMockedGetInstanceCalled);
        }
    }
    @Test
    public void test_SetSSLContextTest_04() throws Exception {
        isMockedGetSocketFactoryCalled = false;
        isMockedGetInstanceCalled = false;
        try (MockedStatic<SSLContext> mocked = Mockito.mockStatic(SSLContext.class)) {
            SSLContext mockedSSLContext = Mockito.mock(SSLContext.class);
            // 通过mock模拟SSLContext.getInstance函数，模拟用户未设置sslProvider时，
            // 自动调用SSLContext.getInstance(String protocol)，但是报错，然后使用用户设置的SSLContext的场景
            mocked.when(() -> SSLContext.getInstance(any(String.class))).thenAnswer(invocation -> {
                isMockedGetInstanceCalled = true;
                throw new NoSuchAlgorithmException(TEST_TLS_NOT_SUPPORTED_ERROR_MESSAGE);
            });
            mocked.when(mockedSSLContext::getSocketFactory).thenAnswer(invocation -> {
                isMockedGetSocketFactoryCalled = true;
                return sslContext.getSocketFactory();
            });
            ObsClient obsClient = TestTools.getPipelineEnvironmentWithCustomisedSSLContext(mockedSSLContext);
            assert obsClient != null;
            ListBucketsResult listBucketsResult = obsClient.listBucketsV2(listBucketsRequest);
            Assert.assertEquals(200, listBucketsResult.getStatusCode());
            obsClient.close();
            Assert.assertTrue(isMockedGetInstanceCalled);
            Assert.assertTrue(isMockedGetSocketFactoryCalled);
        }
    }
    @Test
    public void test_SetSSLContextTest_05() throws Exception {
        isMockedGetInstanceCalled = true;
        try (MockedStatic<SSLContext> mocked = Mockito.mockStatic(SSLContext.class)) {
            // 通过mock模拟SSLContext.getInstance函数，模拟用户未设置sslProvider时，
            // 自动调用SSLContext.getInstance(String protocol)，但是报错，且未设置用户SSLContext，
            // okhttp内部默认创建一个SSLContext成功的场景
            mocked.when(() -> SSLContext.getInstance(any(String.class))).thenAnswer(invocation -> {
                if (isMockedGetInstanceCalled) {
                    return sslContext;
                } else {
                    isMockedGetInstanceCalled = true;
                    throw new NoSuchAlgorithmException(TEST_TLS_NOT_SUPPORTED_ERROR_MESSAGE);
                }
            });
            ObsClient obsClient = TestTools.getPipelineEnvironmentWithCustomisedSSLContext(null);
            assert obsClient != null;
            ListBucketsResult listBucketsResult = obsClient.listBucketsV2(listBucketsRequest);
            Assert.assertEquals(200, listBucketsResult.getStatusCode());
            obsClient.close();
            Assert.assertTrue(isMockedGetInstanceCalled);
        }

    }
    @Test
    public void test_SetSSLContextTest_06() throws Exception {
        try (MockedStatic<SSLContext> mocked = Mockito.mockStatic(SSLContext.class)) {
            // 通过mock模拟SSLContext.getInstance函数，模拟用户未设置sslProvider时，
            // 自动调用SSLContext.getInstance(String protocol)，但是报错，且未设置用户SSLContext，
            // okhttp内部默认创建一个SSLContext也失败的场景
            mocked.when(() -> SSLContext.getInstance(any(String.class))).thenAnswer(invocation -> {
                throw new NoSuchAlgorithmException(TEST_TLS_NOT_SUPPORTED_ERROR_MESSAGE);
            });
            TestTools.getPipelineEnvironmentWithCustomisedSSLContext(null);
            Assert.fail();
        } catch (java.lang.AssertionError testException) {
            Assert.assertTrue(testException.getMessage().contains(TEST_TLS_NOT_SUPPORTED_ERROR_MESSAGE));
        }

    }
}
