/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test.high_performance;

import static com.obs.test.TestTools.genTestFile;

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.GetObjectRequest;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.SetObjectMetadataRequest;
import com.obs.services.model.StorageClassEnum;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;

public class HighPerformanceObjectIT {
    private static final StorageClassEnum HIGH_PERFORMANCE_STORAGE_CLASS_ENUM = StorageClassEnum.HIGH_PERFORMANCE;
    @Rule
    public TestName testName = new TestName();
    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));
    private static final String TEST_HIGH_PERFORMANCE_BUCKET_NAME = "test-performance-object-bucket";
    private static final String TEST_HIGH_PERFORMANCE_OBJECT_NAME = "testHighPerformanceObjectName";
    @Test
    public void tc_PutHighPerformanceObject() throws IOException {
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        ObsClient obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        Assert.assertNotNull(obsClient);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        createBucketRequest.setBucketStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, TEST_HIGH_PERFORMANCE_OBJECT_NAME, 1024 * 1024L);
            PutObjectRequest putObjectRequest =
                new PutObjectRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME, TEST_HIGH_PERFORMANCE_OBJECT_NAME);
            ObjectMetadata objectMetadata = new ObjectMetadata();
            objectMetadata.setObjectStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
            putObjectRequest.setMetadata(objectMetadata);
            putObjectRequest.setFile(testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());
        } finally {
            TestTools.delete_bucket(obsClient, TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        }
    }

    @Test
    public void tc_GetHighPerformanceObject() throws IOException {
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        ObsClient obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        Assert.assertNotNull(obsClient);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        createBucketRequest.setBucketStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, TEST_HIGH_PERFORMANCE_OBJECT_NAME, 1024 * 1024L);
            PutObjectRequest putObjectRequest =
                new PutObjectRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME, TEST_HIGH_PERFORMANCE_OBJECT_NAME);
            ObjectMetadata objectMetadata = new ObjectMetadata();
            objectMetadata.setObjectStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
            putObjectRequest.setMetadata(objectMetadata);
            putObjectRequest.setFile(testFile);

            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            GetObjectRequest getObjectRequest =
                new GetObjectRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME, TEST_HIGH_PERFORMANCE_OBJECT_NAME);
            ObsObject obsObject = obsClient.getObject(getObjectRequest);
            Assert.assertEquals(200, obsObject.getMetadata().getStatusCode());
            obsObject.getObjectContent().close();
        } finally {
            TestTools.delete_bucket(obsClient, TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        }
    }


    @Test
    public void tc_SetHighPerformanceObjectStorageClass() throws IOException {
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        ObsClient obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        Assert.assertNotNull(obsClient);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        createBucketRequest.setBucketStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        try {
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, TEST_HIGH_PERFORMANCE_OBJECT_NAME, 1024 * 1024L);
            PutObjectRequest putObjectRequest =
                new PutObjectRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME, TEST_HIGH_PERFORMANCE_OBJECT_NAME);
            putObjectRequest.setFile(testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            ObjectMetadata objectMetadata = new ObjectMetadata();
            objectMetadata.setObjectStorageClass(HIGH_PERFORMANCE_STORAGE_CLASS_ENUM);
            putObjectRequest.setMetadata(objectMetadata);
            SetObjectMetadataRequest setObjectMetadataRequest =
                new SetObjectMetadataRequest(TEST_HIGH_PERFORMANCE_BUCKET_NAME, TEST_HIGH_PERFORMANCE_OBJECT_NAME);
            ObjectMetadata setObjectMetadataResult = obsClient.setObjectMetadata(setObjectMetadataRequest);
            Assert.assertEquals(200, setObjectMetadataResult.getStatusCode());
        } finally {
            TestTools.delete_bucket(obsClient, TEST_HIGH_PERFORMANCE_BUCKET_NAME);
        }
    }
}
