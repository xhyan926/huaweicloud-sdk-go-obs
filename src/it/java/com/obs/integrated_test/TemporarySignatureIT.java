package com.obs.integrated_test;
import com.obs.test.TestTools;

import static com.obs.test.SSLTestUtils.trustAllManager;

import com.obs.services.ObsClient;
import com.obs.services.model.HttpMethodEnum;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.TemporarySignatureRequest;
import com.obs.services.model.TemporarySignatureResponse;
import com.obs.test.tools.PrepareTestBucket;
import okhttp3.Call;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.util.Arrays;
import java.util.List;
import java.util.Locale;
import java.util.Map;

import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;

/**
 * TemporarySignatureTest
 *
 * @since 2023-12-05
 */
public class TemporarySignatureIT {
    /**
     * testCaseName
     */
    @Rule
    public final TestName testName = new TestName();

    /**
     * prepareTestBucket
     */
    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    /**
     * testObjectKeys
     */
    public final List<String> testObjectKeys = Arrays.asList("///a/b", "a/b///", "a/b///c/d", "//////", "///飞洒发/风格的个", "大事故干撒/是否买哦///", "是打开链接/撒官方///噶时光/个撒跟");

    @Test
    public void test_temporarySignature_for_object_contains_forwardSlash()
        throws NoSuchAlgorithmException, KeyManagementException {
        //forwardSlash is '/'
        ObsClient obsClient = com.obs.test.TestTools.getPipelineEnvironment();
        assert obsClient != null;
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        // URL有效期，3600秒
        long expireSeconds = 3600L;
        for (String testObjectKey : testObjectKeys) {
            String testStringObj = String.valueOf(System.currentTimeMillis());
            System.out.println("testStringObj:"+testStringObj);
            //upload testStringObj
            try (ByteArrayInputStream inputStream = new ByteArrayInputStream(testStringObj.getBytes(
                    StandardCharsets.UTF_8)))
            {
                PutObjectResult putObjectResult = obsClient.putObject(bucketName,testObjectKey,inputStream);
                Assert.assertEquals(200,putObjectResult.getStatusCode());
            }
            catch (IOException e)
            {
                throw new RuntimeException(e);
            }

            // createTemporarySignature and test if contains full object name
            TemporarySignatureRequest req = new TemporarySignatureRequest(HttpMethodEnum.GET, expireSeconds);
            req.setBucketName(bucketName);
            req.setObjectKey(testObjectKey);
            TemporarySignatureResponse response = obsClient.createTemporarySignature(req);

            // download object and check if equals original testStringObj
            System.out.println("Getting object using temporary signature url:");
            System.out.println(response.getSignedUrl());
            Request.Builder builder = new Request.Builder();
            for (Map.Entry<String, String> entry : response.getActualSignedRequestHeaders().entrySet()) {
                builder.header(entry.getKey(), entry.getValue());
            }

            SSLContext sslContext = SSLContext.getInstance("TLSv1.2");
            sslContext.init(null, new TrustManager[]{trustAllManager}, new SecureRandom());
            //使用GET请求下载对象
            Request httpRequest = builder.url(response.getSignedUrl()).get().build();
            OkHttpClient httpClient = new OkHttpClient.Builder().followRedirects(false).retryOnConnectionFailure(false)
                    .cache(null).sslSocketFactory(sslContext.getSocketFactory(), trustAllManager).build();

            Call c = httpClient.newCall(httpRequest);
            try (Response res = c.execute())
            {
                System.out.println("Status:" + res.code());
                String downloadedTestStringObj = "";
                if (res.body() != null) {
                    downloadedTestStringObj = res.body().string();
                }
                System.out.println("downloaded testStringObj:" + downloadedTestStringObj);
                System.out.println("\n");
                Assert.assertEquals(testStringObj,downloadedTestStringObj);
            }
            catch (IOException e)
            {
                throw new RuntimeException(e);
            }
        }
    }
}
