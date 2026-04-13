/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.integrated_test;

import static com.obs.test.SSLTestUtils.trustAllManager;

import com.obs.services.ObsClient;
import com.obs.services.internal.utils.AccessLoggerUtils;
import com.obs.services.model.HttpMethodEnum;
import com.obs.services.model.PostSignatureRequest;
import com.obs.services.model.PostSignatureResponse;
import com.obs.services.model.TemporarySignatureRequest;
import com.obs.services.model.TemporarySignatureResponse;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import okhttp3.Call;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;

public class MemoryLeakIT {

    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();
    @Rule
    public final ExpectedException exception = ExpectedException.none();

    /***
     * 1、createTemporarySignature创建临时url
     * 2、检测access日志已打印，length为0，无堆积现象
     *
     * @throws IOException
     */
    @Test
    public void tc_createTemporarySignatureForListObject()
        throws IOException, InvocationTargetException, NoSuchMethodException, IllegalAccessException,
        NoSuchAlgorithmException, KeyManagementException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        // URL有效期，3600秒
        long expireSeconds = 3600L;
        TemporarySignatureRequest request = new TemporarySignatureRequest(HttpMethodEnum.GET, expireSeconds);
        request.setBucketName(bucketName);

        TemporarySignatureResponse response = obsClient.createTemporarySignature(request);
        // 检测access日志已打印，length为0，无堆积现象
        StringBuilder logBuilder = getSdkLog();
        Assert.assertNotNull(logBuilder);
        Assert.assertEquals(0, logBuilder.length());
        Request.Builder builder = new Request.Builder();
        for (Map.Entry<String, String> entry : response.getActualSignedRequestHeaders().entrySet()) {
            builder.header(entry.getKey(), entry.getValue());
        }
        SSLContext sslContext = SSLContext.getInstance("TLSv1.2");
        sslContext.init(null, new TrustManager[]{trustAllManager}, new SecureRandom());
        // 使用GET请求获取对象列表
        Request httpRequest = builder.url(response.getSignedUrl()).get().build();
        OkHttpClient httpClient =
            new OkHttpClient.Builder().sslSocketFactory(sslContext.getSocketFactory(), trustAllManager)
                .followRedirects(false)
                .retryOnConnectionFailure(false)
                .cache(null)
                .build();
        Call c = httpClient.newCall(httpRequest);
        Response res = c.execute();
        Assert.assertEquals(200, res.code());
        res.close();
    }

    /***
     * 1、createPostSignature创建临时url
     * 2、检测access日志已打印，length为0，无堆积现象
     *
     */
    @Test
    public void tc_createPostSignatureForListObject()
        throws InvocationTargetException, NoSuchMethodException, IllegalAccessException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        // 生成基于表单上传的请求
        PostSignatureRequest request = new PostSignatureRequest();
        // 设置表单参数
        Map<String, Object> formParams = new HashMap<>();
        // 设置对象访问权限为公共读
        formParams.put("x-obs-acl", "public-read");
        // 设置对象MIME类型
        formParams.put("content-type", "text/plain");

        request.setFormParams(formParams);
        // 设置表单上传请求有效期，单位：秒
        request.setExpires(3600);
        obsClient.createPostSignature(request);
        // 检测access日志已打印，length为0，无堆积现象
        StringBuilder logBuilder = getSdkLog();
        Assert.assertNotNull(logBuilder);
        Assert.assertEquals(0, logBuilder.length());
    }

    private StringBuilder getSdkLog() throws NoSuchMethodException, InvocationTargetException, IllegalAccessException {
        Method getLog = AccessLoggerUtils.class.getDeclaredMethod("getLog");
        getLog.setAccessible(true);
        Object actual = getLog.invoke(AccessLoggerUtils.class);
        return (StringBuilder) actual;
    }
}
