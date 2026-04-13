package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.ListBucketsRequest;
import com.obs.test.TestTools;
import org.junit.Assert;
import org.junit.Test;

import java.io.IOException;

public class RestStorageServiceIT
{

    @Test
    public void testInfoForSignatureDoesNotMatch() {
        try (ObsClient obsClient = TestTools.getWrongPipelineEnvironment())
        {
            obsClient.listBuckets(new ListBucketsRequest());
        } catch (ObsException e) {
            System.out.println("PutObjectWithSha256 failed");
            // 请求失败,打印http状态码
            System.out.println("HTTP Code:" + e.getResponseCode());
            // 请求失败,打印服务端错误码
            System.out.println("Error Code:" + e.getErrorCode());
            // 请求失败,打印详细错误信息
            System.out.println("Error Message:" + e.getErrorMessage());
            // 请求失败,打印请求id
            System.out.println("Request ID:" + e.getErrorRequestId());
            System.out.println("Host ID:" + e.getErrorHostId());
            e.printStackTrace();
            Assert.assertTrue("ErrorMessage should contains local StringToSign while SignatureDoesNotMatch",
                    e.getErrorMessage().contains("your local StringToSign is (between\"---\"):"));
            Assert.assertTrue("ErrorMessage should contains Server's StringToSign while SignatureDoesNotMatch",
                    e.getErrorMessage().contains("Please compare it to Server's StringToSign (between\"---\"):"));
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
