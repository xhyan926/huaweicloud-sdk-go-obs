/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test.high_performance;

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.model.BucketMetadataInfoRequest;
import com.obs.services.model.BucketMetadataInfoResult;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.StorageClassEnum;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

public class HighPerformanceBucketIT {
    private static final StorageClassEnum HIGH_PERFORMANCE_STORAGE_CLASS_ENUM = StorageClassEnum.HIGH_PERFORMANCE;
    @Rule
    public TestName testName = new TestName();
    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();
    private static final String TEST_HIGH_PERFORMANCE_BUCKET_NAME = "test-performance-object-bucket";
    @Test
    public void tc_CreateHighPerformanceBucket() {
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        ObsClient obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        Assert.assertNotNull(obsClient);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        createBucketRequest.setBucketStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());
        } finally {
            TestTools.delete_bucket(obsClient, TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        }
    }

    @Test
    public void tc_GetHighPerformanceBucketMetadata() {
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        ObsClient obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        Assert.assertNotNull(obsClient);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        createBucketRequest.setBucketStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());
            BucketMetadataInfoRequest bucketMetadataInfoRequest =
                new BucketMetadataInfoRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME);
            BucketMetadataInfoResult bucketMetadataInfoResult =
                obsClient.getBucketMetadata(bucketMetadataInfoRequest);
            Assert.assertEquals(200, bucketMetadataInfoResult.getStatusCode());
            Assert.assertEquals(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM, bucketMetadataInfoResult.getBucketStorageClass());
        } finally {
            TestTools.delete_bucket(obsClient, TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        }
    }
}
