/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.internal;

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.fail;

public class CallTimeoutUnitTest {
    private static final String TIMEOUT_MESSAGE = "timeout";
    private static MockWebServer server = null;
    @BeforeClass
    public static void beforeCallTimeoutUnitTest() throws IOException {
        server = new MockWebServer();
        server.start();
        // 设置延迟为5秒
        server.enqueue(new MockResponse()
            .setBody("{message:Hello, World!}")
            .setHeadersDelay(5, TimeUnit.SECONDS) // 模拟请求耗时长场景
            .addHeader("Content-Type", "application/xml"));
    }
    @AfterClass
    public static void afterCallTimeoutUnitTest() throws IOException {
        server.shutdown();
    }

    @Test
    public void should_obsConfiguration_setCallTimeout_successfully() {
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        int testCallTimeout = 10;
        obsConfiguration.setCallTimeout(testCallTimeout);
        Assert.assertEquals(testCallTimeout, obsConfiguration.getCallTimeout());
    }

    @Test
    public void should_obsConfiguration_getCallTimeout_successfully() {
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        Assert.assertEquals(obsConfiguration.getCallTimeout(), ObsConstraint.HTTP_CALL_TIMEOUT_VALUE);
    }

    @Test
    public void should_throw_exception_when_reached_callTimeout() throws IOException {
        long start = 0;
        int testCallTimeout = 3 * 1000;
        int testOtherTimeout = 60 * 1000;
        try {
            ObsConfiguration obsConfiguration = new ObsConfiguration();
            obsConfiguration.setCallTimeout(testCallTimeout);
            obsConfiguration.setEndPoint(server.url("/").toString());
            obsConfiguration.setAuthTypeNegotiation(false);
            obsConfiguration.setHttpsOnly(false);
            obsConfiguration.setEndpointHttpPort(server.getPort());
            // 调大其他超时时间，不影响测试的CallTimeout
            obsConfiguration.setConnectionTimeout(testOtherTimeout);
            obsConfiguration.setSocketTimeout(testOtherTimeout);
            ObsClient obsClient = new ObsClient(obsConfiguration);
            start = System.currentTimeMillis();
            obsClient.listBuckets();
            // 不应该成功
            fail();
        } catch (ObsException e) {
            System.out.println("testObsClientCallTimeout e.getErrorMessage():" + e.getErrorMessage());
            e.printStackTrace();
            // 获取耗时
            long costTime = System.currentTimeMillis() - start;
            System.out.println("testObsClientCallTimeout costTime:" + costTime);
            // costTime应该大于等于callTimeout
            Assert.assertTrue(costTime >= testCallTimeout);
            // costTime与callTimeout之差不应该超过1s
            Assert.assertTrue(  (costTime - testCallTimeout) < 1000);
            // 请求耗时长,到达超时时间，应该报错
            Assert.assertTrue(e.getErrorMessage().contains(TIMEOUT_MESSAGE));
        }
    }
}
