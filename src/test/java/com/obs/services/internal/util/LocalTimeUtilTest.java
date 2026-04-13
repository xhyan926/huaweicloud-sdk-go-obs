package com.obs.services.internal.util;

import static com.obs.services.internal.Constants.REQUEST_TIME_TOO_SKEWED_CODE;
import static com.obs.test.SSLTestUtils.trustAllManager;
import static com.obs.test.TestTools.genTestFile;
import static com.obs.test.TestTools.getEndpointWithNoPrefix;
import static com.obs.test.TestTools.printException;
import static com.obs.test.TestTools.printObsException;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.utils.LocalTimeUtil;
import com.obs.services.model.HttpMethodEnum;
import com.obs.services.model.ListBucketsRequest;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PostSignatureRequest;
import com.obs.services.model.PostSignatureResponse;
import com.obs.services.model.TemporarySignatureRequest;
import com.obs.services.model.TemporarySignatureResponse;
import com.obs.services.model.UploadFileRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import com.obs.test.tools.PropertiesTools;
import okhttp3.Call;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.IOException;
import java.io.StringWriter;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Locale;
import java.util.Map;

import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;

public class LocalTimeUtilTest {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    // 加上这个时间差(24 hour)后，请求应该会产生异常，错误码为RequestTimeTooSkewed
    public long errorTimeDiff = -24 * 60 * 60 * 1000L;

    protected String retryInfo = "Retrying connection that failed with RequestTimeTooSkewed error";

    // 关闭EnableAutoRetryForSkewedTime，并且通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后PutObjectWithInputStream，
    // 应该会产生异常，错误码为RequestTimeTooSkewed
    @Test
    public void testPutObjectWithInputStreamForRequestTimeTooSkewed() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        String testObjectToPut = objectKey + "-testObjectToPut";

        try (ObsClient obsClient = TestTools.getPipelineEnvironment();
                ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(testObjectToPut.getBytes())) {
            assert obsClient != null;
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(false);
            obsClient.putObject(bucketName, objectKey, byteArrayInputStream);
            fail();
        } catch (ObsException e) {
            assertIsRequestTimeTooSkewed(e);
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
        // 测试v4协议
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V4();
                ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(testObjectToPut.getBytes())) {
            assert obsClient != null;
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(false);
            obsClient.putObject(bucketName, objectKey, byteArrayInputStream);
            fail();
        } catch (ObsException e) {
            assertIsRequestTimeTooSkewed(e);
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
    }

    // 关闭EnableAutoRetryForSkewedTime，并且通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后HeadBucket，
    // 应该会产生异常，错误码为RequestTimeTooSkewed
    @Test
    public void testHeadBucketForRequestTimeTooSkewed() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(false);
            obsClient.headBucket(bucketName);
            fail();
        } catch (ObsException e) {
            assertIsRequestTimeTooSkewed(e);
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }

        // 测试v4协议
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V4()) {
            assert obsClient != null;
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(false);
            obsClient.headBucket(bucketName);
            fail();
        } catch (ObsException e) {
            assertIsRequestTimeTooSkewed(e);
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
    }

    public void assertIsRequestTimeTooSkewed(ObsException e) {
        assertEquals(e.getResponseCode(), 403);
        String errorCode = e.getErrorCode();
        if (errorCode != null) {
            assertEquals(errorCode, REQUEST_TIME_TOO_SKEWED_CODE);
        } else {
            Map<String, String> headers = e.getResponseHeaders();
            assertNotNull(headers);
            for (Map.Entry<String, String> header : headers.entrySet()) {
                if (header.getKey().equalsIgnoreCase("error-code")) {
                    assertEquals(header.getValue(), REQUEST_TIME_TOO_SKEWED_CODE);
                    return;
                }
            }
            fail();
        }
    }

    // 开启EnableAutoRetryForSkewedTime，并且通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后PutObjectWithInputStream，
    // 产生RequestTimeTooSkewed，但经过SDK自动重试后请求成功
    @Test
    public void testPutObjectWithInputStreamForRequestTimeTooSkewed_Retry() {
        // 初始化 log, 用于检测是否有重试
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        String testObjectToPut = objectKey + "-testObjectToPut";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment();
                ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(testObjectToPut.getBytes())) {
            writer.getBuffer().setLength(0);
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.putObject(bucketName, objectKey, byteArrayInputStream);
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }

        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V4();
                ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(testObjectToPut.getBytes())) {
            writer.getBuffer().setLength(0);
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.putObject(bucketName, objectKey, byteArrayInputStream);
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
    }

    // 开启EnableAutoRetryForSkewedTime，并且通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后HeadBucket，
    // 产生RequestTimeTooSkewed异常，但经过SDK自动重试后请求成功
    @Test
    public void testHeadBucketForRequestTimeTooSkewed_Retry() {
        // 初始化 log, 用于检测是否有重试
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.headBucket(bucketName);
            assertTrue(writer.toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (Throwable e) {
            // System.out.println("writer in testHeadBucketForRequestTimeTooSkewed_Retry:" + writer);
            printException(e);
            throw new RuntimeException(e);
        }

        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V4()) {
            writer.getBuffer().setLength(0);
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.headBucket(bucketName);
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
    }

    // 开启EnableAutoRetryForSkewedTime，并且通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，然后ListBucket，
    // 产生RequestTimeTooSkewed异常，但经过SDK自动重试后请求成功
    @Test
    public void testListBucketForRequestTimeTooSkewed_Retry() {
        // 初始化 log, 用于检测是否有重试
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            writer.getBuffer().setLength(0);
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            ListBucketsRequest listBucketsRequest = new ListBucketsRequest();
            obsClient.listBuckets(listBucketsRequest);
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }

        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V4()) {
            writer.getBuffer().setLength(0);
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            ListBucketsRequest listBucketsRequest = new ListBucketsRequest();
            obsClient.listBuckets(listBucketsRequest);
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
    }

    // 开启EnableAutoRetryForSkewedTime，并且通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后PutObjectWithInputStream，产生RequestTimeTooSkewed后SDK自动重试
    // 然后GetObject，产生RequestTimeTooSkewed后SDK自动重试
    @Test
    public void testPutObjectAndGetObjectForRequestTimeTooSkewed_Retry() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        String testFileName = bucketName + "testFile";
        // 10 mb test file
        long fileSizeInBytes = 10 * 1024 * 1024L;
        // 初始化 log, 用于检测是否有重试
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        File testFile = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            writer.getBuffer().setLength(0);
            testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.putObject(bucketName, objectKey, testFile);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            ObsObject obsObject = obsClient.getObject(bucketName, objectKey);
            obsObject.getObjectContent().close();
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }

        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V4()) {
            writer.getBuffer().setLength(0);
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            obsClient.putObject(bucketName, objectKey, testFile);
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            ObsObject obsObject = obsClient.getObject(bucketName, objectKey);
            obsObject.getObjectContent().close();
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
    }

    // 通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后通过临时鉴权进行PutObject，产生RequestTimeTooSkewed
    // 然后通过DateHeaderUtil.setTimeDiffInMs调整时间正常后重新尝试，不会产生RequestTimeTooSkewed
    @Test
    public void testPutObjectWithTemporarySignatureWithWrongDate() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        String testFileName = bucketName + "testFile";
        // 10 mb test file
        long fileSizeInBytes = 10 * 1024 * 1024L;
        // URL有效期，20分钟
        long expireSeconds = 1200;
        TemporarySignatureRequest request = new TemporarySignatureRequest(HttpMethodEnum.PUT, expireSeconds);
        request.setBucketName(bucketName);
        request.setObjectKey(objectKey);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment();
                ObsClient obsClientV4 = TestTools.getPipelineEnvironment_V4()) {
            // obsClientV4用于测试v4签名场景
            assert obsClient != null;
            assert obsClientV4 != null;
            File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);

            // 通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            // 然后通过临时鉴权进行PutObject，产生RequestTimeTooSkewed
            TemporarySignatureResponse response = obsClient.createTemporarySignature(request);
            Call call = testOkhttpPutObject(response, testFile);
            assertCallFailWithRequestTimeTooSkewed(call);

            TemporarySignatureResponse responseV4 = obsClientV4.createTemporarySignature(request);
            call = testOkhttpPutObject(responseV4, testFile);
            assertCallFailWithRequestTimeTooSkewed(call);

            LocalTimeUtil.setTimeDiffInMs(0);
            // 然后通过DateHeaderUtil.setTimeDiffInMs调整时间正常后重新尝试，不会产生RequestTimeTooSkewed
            response = obsClient.createTemporarySignature(request);
            call = testOkhttpPutObject(response, testFile);
            assertCallSuccess(call);

            responseV4 = obsClient.createTemporarySignature(request);
            call = testOkhttpPutObject(responseV4, testFile);
            assertCallSuccess(call);

            System.out.println("testPutObjectWithTemporarySignatureWithWrongDate successfully");
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException(e);
        } catch (KeyManagementException e) {
            throw new RuntimeException(e);
        }
    }

    public Call testOkhttpPutObject(TemporarySignatureResponse response, File testFile)
        throws NoSuchAlgorithmException, KeyManagementException {
        Request.Builder builder = new Request.Builder();
        for (Map.Entry<String, String> entry : response.getActualSignedRequestHeaders().entrySet()) {
            builder.header(entry.getKey(), entry.getValue());
        }
        SSLContext sslContext = SSLContext.getInstance("TLSv1.2");
        sslContext.init(null, new TrustManager[]{trustAllManager}, new SecureRandom());
        // 使用PUT请求上传对象
        Request httpRequest = builder.url(response.getSignedUrl()).put(RequestBody.create(testFile, null)).build();
        OkHttpClient httpClient =
                new OkHttpClient.Builder().followRedirects(false).retryOnConnectionFailure(false).cache(null)
                    .sslSocketFactory(sslContext.getSocketFactory(), trustAllManager).build();

        return httpClient.newCall(httpRequest);
    }

    public void assertCallFailWithRequestTimeTooSkewed(Call call) throws IOException {
        try (Response res = call.execute()) {
            assertEquals(403, res.code());
            assert res.body() != null;
            assertTrue(res.body().string().contains(REQUEST_TIME_TOO_SKEWED_CODE));
        }
    }

    public void assertCallSuccess(Call call) throws IOException {
        try (Response res = call.execute()) {
            assertEquals(200, res.code());
            assert res.body() != null;
            assertFalse(res.body().string().contains(REQUEST_TIME_TOO_SKEWED_CODE));
        }
    }

    // 通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后进行PostObject，产生RequestTimeTooSkewed
    // 然后通过DateHeaderUtil.setTimeDiffInMs调整时间正常后重新尝试，不会产生RequestTimeTooSkewed
    @Test
    public void testPostObjectWithWrongDate() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        // URL有效期，20分钟
        long expireSeconds = 1200;
        PostSignatureRequest postSignatureRequest = new PostSignatureRequest();
        postSignatureRequest.setExpires(expireSeconds);
        Map<String, Object> formParams = new HashMap<>();
        String contentType = "text/plain";
        formParams.put("content-type", contentType);
        postSignatureRequest.setFormParams(formParams);
        String policyExpired = "Invalid according to Policy: Policy expired.";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS();
                ObsClient obsClientV2 = TestTools.getPipelineEnvironment_V2()) {
            // obsClientV2用于测试v2签名场景
            assert obsClient != null;
            assert obsClientV2 != null;
            String accessKey = PropertiesTools.getInstance(TestTools.getPropertiesFile())
                    .getProperties("environment.ak");
            String endpointWithNoPrefix = getEndpointWithNoPrefix();

            // 通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景
            LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
            PostSignatureResponse postSignatureResponseExpireOBS = obsClient.createPostSignature(postSignatureRequest);
            PostSignatureResponse postSignatureResponseExpireV2 = obsClientV2.createPostSignature(postSignatureRequest);
            // 通过DateHeaderUtil.setTimeDiffInMs模拟时间正常的场景
            LocalTimeUtil.setTimeDiffInMs(0);
            PostSignatureResponse postSignatureResponseNormalOBS = obsClient.createPostSignature(postSignatureRequest);
            PostSignatureResponse postSignatureResponseNormalV2 = obsClientV2.createPostSignature(postSignatureRequest);
            formParams.put("key", objectKey);
            formParams.put("accesskeyid", accessKey);
            String postUrl = "http://" + bucketName + "." + endpointWithNoPrefix;
            formParams.put("policy", postSignatureResponseExpireOBS.getPolicy());
            formParams.put("signature", postSignatureResponseExpireOBS.getSignature());

            try (Response response = getOkHttpPostCall(postUrl, formParams, contentType).execute()) {
                assertEquals(403, response.code());
                assert response.body() != null;
                assertTrue(response.body().string().contains(policyExpired));
            }

            formParams.replace("policy", postSignatureResponseNormalOBS.getPolicy());
            formParams.replace("signature", postSignatureResponseNormalOBS.getSignature());
            try (Response response = getOkHttpPostCall(postUrl, formParams, contentType).execute()) {
                assertEquals(204, response.code());
            }

            formParams.remove("accesskeyid");
            formParams.put("AwsAccesskeyid", accessKey);
            formParams.replace("policy", postSignatureResponseExpireV2.getPolicy());
            formParams.replace("signature", postSignatureResponseExpireV2.getSignature());
            try (Response response = getOkHttpPostCall(postUrl, formParams, contentType).execute()) {
                long code = response.code();
                assert response.body() != null;
                String responseBody = response.body().string();

                assertEquals("code should be 403, body is:" + responseBody,403, code);
                assertTrue(responseBody.contains(policyExpired));
            }

            formParams.replace("policy", postSignatureResponseNormalV2.getPolicy());
            formParams.replace("signature", postSignatureResponseNormalV2.getSignature());
            try (Response response = getOkHttpPostCall(postUrl, formParams, contentType).execute()) {
                assertEquals(204, response.code());
            }

            System.out.println("testPostObjectWithWrongDate successfully");
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        }
    }

    private static Call getOkHttpPostCall(String postUrl, Map<String, Object> formFields, String contentType)
            throws IOException {
        String boundary = "9431149156168";
        OkHttpClient client = new OkHttpClient();

        StringBuffer strBuf = new StringBuffer();
        Iterator<Map.Entry<String, Object>> iter = formFields.entrySet().iterator();
        int i = 0;

        while (iter.hasNext()) {
            Map.Entry<String, Object> entry = iter.next();
            String inputName = entry.getKey();
            Object inputValue = entry.getValue();

            if (inputValue == null) {
                continue;
            }

            if (i == 0) {
                strBuf.append("--").append(boundary).append("\r\n");
                strBuf.append("Content-Disposition: form-data; name=\"" + inputName + "\"\r\n\r\n");
                strBuf.append(inputValue);
            } else {
                strBuf.append("\r\n").append("--").append(boundary).append("\r\n");
                strBuf.append("Content-Disposition: form-data; name=\"" + inputName + "\"\r\n\r\n");
                strBuf.append(inputValue);
            }

            i++;
        }
        strBuf.append("\r\n").append("--").append(boundary).append("\r\n");
        strBuf.append("Content-Disposition: form-data; name=\"file\"; " + "filename=\"testFile.txt\"\r\n");
        strBuf.append("Content-Type: " + contentType + "\r\n\r\n");
        strBuf.append("testFile values.");
        strBuf.append("\r\n--" + boundary + "--\r\n");

        Request request =
                new Request.Builder()
                        .url(postUrl)
                        .post(RequestBody.create(strBuf.toString(), null))
                        .addHeader("Content-Type", "multipart/form-data; boundary=" + boundary)
                        .build();

        return client.newCall(request);
    }

    // 异步通过DateHeaderUtil.setTimeDiffInMs模拟时间相差16min的场景，
    // 然后调用UploadFile，sdk会自动重试RequestTimeTooSkewed
    @Test
    public void testUploadFileWithWrongDate() {
        // 初始化 log, 用于检测是否有重试
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        String testFileName = bucketName + "testFile";
        // 500 mb test file
        long fileSizeInBytes = 500 * 1024 * 1024L;
        Thread testSetWrongDate =
                new Thread(
                        () -> {
                            int testTime = 100;
                            int sleepTime = 3000;
                            while (--testTime > 0) {
                                LocalTimeUtil.setTimeDiffInMs(errorTimeDiff);
                                System.out.println("set wrong time every " + sleepTime + " ms");
                                try {
                                    Thread.sleep(sleepTime);
                                } catch (InterruptedException e) {
                                    printException(e);
                                    throw new RuntimeException(e);
                                }
                            }
                        });
        testSetWrongDate.start();
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            obsClient.getLocalTimeUtil().setEnableAutoRetryForSkewedTime(true);
            File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
            UploadFileRequest uploadFileRequest =
                    new UploadFileRequest(bucketName, objectKey, testFile.getPath(), 1024 * 1024, 32, true);
            obsClient.uploadFile(uploadFileRequest);
            assertTrue(writer.getBuffer().toString().contains(retryInfo));
            System.out.println("testUploadFileWithWrongDate successfully");
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            printException(e);
            throw new RuntimeException(e);
        } finally {
            testSetWrongDate.interrupt();
        }
    }
}
