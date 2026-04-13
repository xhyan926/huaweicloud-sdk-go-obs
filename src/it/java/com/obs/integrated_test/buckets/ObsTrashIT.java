/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.integrated_test.buckets;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.trash.BucketTrashConfiguration;
import com.obs.services.model.trash.DeleteBucketTrashRequest;
import com.obs.services.model.trash.GetBucketTrashRequest;
import com.obs.services.model.trash.GetBucketTrashResult;
import com.obs.services.model.trash.SetBucketTrashRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.Locale;

public class ObsTrashIT {
    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    // 设置桶回收站API
    @Test
    public void test_SDK_fs_trash_001() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNamePosix = bucketName + "-posix";

        int reservedDays = 4;
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        BucketTrashConfiguration bucketTrashConfiguration = new BucketTrashConfiguration(reservedDays);
        SetBucketTrashRequest setBucketTrashRequest = new SetBucketTrashRequest(bucketNamePosix, bucketTrashConfiguration);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNamePosix);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());

            HeaderResponse headerResponse = obsClient.setBucketTrash(setBucketTrashRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());
            GetBucketTrashRequest getBucketTrashRequest = new GetBucketTrashRequest();
            getBucketTrashRequest.setBucketName(bucketNamePosix);
            GetBucketTrashResult getBucketTrashResult = obsClient.getBucketTrash(getBucketTrashRequest);
            assertEquals(getBucketTrashResult.getTrashConfiguration().getReservedDays(), reservedDays);

            try {
                bucketTrashConfiguration.setReservedDays(0);
                obsClient.setBucketTrash(setBucketTrashRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertTrue(400 <= e.getResponseCode());
            }

            bucketTrashConfiguration.setReservedDays(1);
            headerResponse = obsClient.setBucketTrash(setBucketTrashRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            bucketTrashConfiguration.setReservedDays(30);
            headerResponse = obsClient.setBucketTrash(setBucketTrashRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            try {
                bucketTrashConfiguration.setReservedDays(-1);
                obsClient.setBucketTrash(setBucketTrashRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertTrue(400 <= e.getResponseCode());
            }
        } finally {
            try {
                TestTools.delete_bucket(obsClient, bucketNamePosix);
            } catch (Throwable ignore) {
            }
        }
    }

    // 删除桶回收站API
    @Test
    public void test_SDK_fs_trash_002() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNamePosix = bucketName + "-posix";

        int reservedDays = 4;
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        BucketTrashConfiguration bucketTrashConfiguration = new BucketTrashConfiguration(reservedDays);
        SetBucketTrashRequest setBucketTrashRequest = new SetBucketTrashRequest(bucketNamePosix, bucketTrashConfiguration);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNamePosix);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());

            HeaderResponse headerResponse = obsClient.setBucketTrash(setBucketTrashRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());
            GetBucketTrashRequest getBucketTrashRequest = new GetBucketTrashRequest();
            getBucketTrashRequest.setBucketName(bucketNamePosix);
            GetBucketTrashResult getBucketTrashResult = obsClient.getBucketTrash(getBucketTrashRequest);
            assertEquals(getBucketTrashResult.getTrashConfiguration().getReservedDays(), reservedDays);

            DeleteBucketTrashRequest deleteBucketTrashRequest = new DeleteBucketTrashRequest();
            deleteBucketTrashRequest.setBucketName(bucketNamePosix);
            headerResponse = obsClient.deleteBucketTrash(deleteBucketTrashRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            try {
                obsClient.getBucketTrash(getBucketTrashRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertEquals(404, e.getResponseCode());
            }
        } finally {
            try {
                TestTools.delete_bucket(obsClient, bucketNamePosix);
            } catch (Throwable ignore) {
            }
        }
    }

    // obs-对象桶不支持回收站配置
    @Test
    public void test_SDK_fs_trash_003() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT) + "-pfs";
        int reservedDays = 4;
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        BucketTrashConfiguration bucketTrashConfiguration = new BucketTrashConfiguration(reservedDays);
        SetBucketTrashRequest setBucketTrashRequest = new SetBucketTrashRequest(bucketName,
            bucketTrashConfiguration);
        try {
            obsClient.setBucketTrash(setBucketTrashRequest);
            fail();
        } catch (ObsException e) {
            Assert.assertTrue(400 <= e.getResponseCode());
        }

        GetBucketTrashRequest getBucketTrashRequest = new GetBucketTrashRequest();
        getBucketTrashRequest.setBucketName(bucketName);
        try {
            obsClient.getBucketTrash(getBucketTrashRequest);
            fail();
        } catch (ObsException e) {
            Assert.assertTrue(400 <= e.getResponseCode());
        }
    }

    // 支持S3接口方式调用trash的3个api
    @Test
    public void test_SDK_fs_trash_004() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNamePosix = bucketName + "-posix";

        int reservedDays = 4;
        ObsClient obsClient = TestTools.getPipelineEnvironment_V2();
        assert obsClient != null;
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNamePosix);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());

            BucketTrashConfiguration bucketTrashConfiguration = new BucketTrashConfiguration(reservedDays);
            SetBucketTrashRequest setBucketTrashRequest = new SetBucketTrashRequest(bucketNamePosix,
                bucketTrashConfiguration);
            HeaderResponse headerResponse = obsClient.setBucketTrash(setBucketTrashRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            GetBucketTrashRequest getBucketTrashRequest = new GetBucketTrashRequest();
            getBucketTrashRequest.setBucketName(bucketNamePosix);
            GetBucketTrashResult getBucketTrashResult = obsClient.getBucketTrash(getBucketTrashRequest);
            assertEquals(getBucketTrashResult.getTrashConfiguration().getReservedDays(), reservedDays);

            DeleteBucketTrashRequest deleteBucketTrashRequest = new DeleteBucketTrashRequest();
            deleteBucketTrashRequest.setBucketName(bucketNamePosix);
            headerResponse = obsClient.deleteBucketTrash(deleteBucketTrashRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());
        } finally {
            try {
                TestTools.delete_bucket(obsClient, bucketNamePosix);
            } catch (Throwable ignore) {
            }
        }
    }
}
