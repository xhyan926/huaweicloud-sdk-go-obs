/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.internal.storageclass;

import static com.obs.test.TestTools.deleteObjects;
import static com.obs.test.TestTools.genTestFile;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.BucketMetadataInfoRequest;
import com.obs.services.model.BucketMetadataInfoResult;
import com.obs.services.model.BucketStoragePolicyConfiguration;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.BucketVersioningConfiguration;
import com.obs.services.model.CompleteMultipartUploadRequest;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.CopyObjectRequest;
import com.obs.services.model.CopyObjectResult;
import com.obs.services.model.CopyPartRequest;
import com.obs.services.model.CopyPartResult;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.InitiateMultipartUploadRequest;
import com.obs.services.model.InitiateMultipartUploadResult;
import com.obs.services.model.ListMultipartUploadsRequest;
import com.obs.services.model.ListPartsRequest;
import com.obs.services.model.ListPartsResult;
import com.obs.services.model.ListVersionsResult;
import com.obs.services.model.Multipart;
import com.obs.services.model.MultipartUpload;
import com.obs.services.model.MultipartUploadListing;
import com.obs.services.model.ObjectListing;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PartEtag;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.RestoreObjectRequest;
import com.obs.services.model.RestoreTierEnum;
import com.obs.services.model.SetObjectMetadataRequest;
import com.obs.services.model.StorageClassEnum;
import com.obs.services.model.VersionOrDeleteMarker;
import com.obs.services.model.VersioningStatusEnum;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

public class IntelligentTieringTest {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    /***
     * 1、使用s3协议创建智能分级存储桶
     * 2、获取桶元数据（getBucketMetadata），检查storageclass是否智能分级桶
     * 3、使用obs协议创建智能分级存储桶
     * 4、获取桶元数据(headBucket)，检查storageclass是否智能分级桶
     * 预期结果:
     * 1、返回200，创桶成功
     * 2、返回200，x-obs-storage-class是智能分级
     * 3、返回200，创桶成功
     * 4、返回200，x-amz-storage-class是智能分级
     */
    @Test
    public void test_createBucket_with_intelligent_tiering_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNameOBS = bucketName + "-obs";
        CreateBucketRequest createBucketRequest = new CreateBucketRequest();
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            createBucketRequest.setBucketName(bucketNameOBS);
            createBucketRequest.setBucketStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            assert obsClient != null;
            ObsBucket bucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, bucket.getStatusCode());
            BucketMetadataInfoResult bucketMetadataInfoResult =
                    obsClient.getBucketMetadata(new BucketMetadataInfoRequest(bucketNameOBS));
            Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, bucketMetadataInfoResult.getBucketStorageClass());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        } finally {
            deleteBucketIgnoreError(bucketNameOBS);
        }
    }

    /***
     * 1、创建标准桶1
     * 2、不设置存储类别上传对象1
     * 3、查询对象1的存储类别
     * 4、修改标准桶1的默认存储类别为智能分级存储
     * 5、不设置存储类别上传对象2
     * 6、查询对象2的存储类别
     * 7、查询桶1的存储类别
     * 预期结果:
     * 1、返回200，创桶成功
     * 2、返回200，上传对象成功
     * 3、返回200，对象是标准对象
     * 4、返回200，修改桶存储类别成功
     * 5、返回200，上传对象成功
     * 6、返回200，对象是智能分级对象
     * 7、返回200，桶1的默认存储类别为智能分级存储
     */
    @Test
    public void test_setBucketStoragePolicy_with_intelligent_tiering_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNameOBS = bucketName + "-obs";
        String objectKey = bucketName + "testObjectKey";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            CreateBucketRequest createBucketRequest = new CreateBucketRequest();
            createBucketRequest.setBucketName(bucketNameOBS);
            createBucketRequest.setBucketStorageClass(StorageClassEnum.STANDARD);
            ObsBucket bucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, bucket.getStatusCode());
            BucketMetadataInfoResult bucketMetadataInfoResult =
                    obsClient.getBucketMetadata(new BucketMetadataInfoRequest(bucketNameOBS));
            Assert.assertEquals(StorageClassEnum.STANDARD, bucketMetadataInfoResult.getBucketStorageClass());

            String testFileName = bucketNameOBS + "testFile";
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
            PutObjectResult putObjectResult = obsClient.putObject(bucketNameOBS, objectKey, testFile);
            assertEquals(200, putObjectResult.getStatusCode());
            ObjectMetadata objectMetadata = obsClient.getObjectMetadata(bucketNameOBS, objectKey);
            // 对象默认没有存储类别
            assertNull(objectMetadata.getObjectStorageClass());

            HeaderResponse response =
                    obsClient.setBucketStoragePolicy(
                            bucketNameOBS, new BucketStoragePolicyConfiguration(StorageClassEnum.INTELLIGENT_TIERING));
            assertEquals(200, response.getStatusCode());

            putObjectResult = obsClient.putObject(bucketNameOBS, objectKey, testFile);
            assertEquals(200, putObjectResult.getStatusCode());
            objectMetadata = obsClient.getObjectMetadata(bucketNameOBS, objectKey);
            Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, objectMetadata.getObjectStorageClass());

            bucketMetadataInfoResult = obsClient.getBucketMetadata(new BucketMetadataInfoRequest(bucketNameOBS));
            Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, bucketMetadataInfoResult.getBucketStorageClass());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        } finally {
            deleteBucketIgnoreError(bucketNameOBS);
        }
    }

    /***
     * 1、创建标准桶
     * 2、上传智能分级对象1
     * 3、下载对象1
     * 4、创建温桶
     * 5、上传智能分级对象2
     * 6、下载对象2
     * 7、创建冷桶
     * 8、上传智能分级对象3
     * 9、下载对象3
     * 10、创建深度归档桶
     * 11、上传智能分级对象4
     * 12、下载对象4
     * 预期结果:
     * 1、返回200，创桶成功，桶是标准桶
     * 2、返回200，上传对象1成功
     * 3、返回200，下载对象1成功，且返回信息中对象的存储类别为智能分级存储
     * 4、返回200，创桶成功，桶是温桶
     * 5、返回200，上传对象2成功
     * 6、返回200，下载对象2成功，且返回信息中对象的存储类别为智能分级存储
     * 7、返回200，创桶成功，桶是冷桶
     * 8、返回200，上传对象3成功
     * 9、返回200，下载对象3成功，且返回信息中对象的存储类别为智能分级存储
     * 10、返回200，创桶成功，桶是深度归档桶
     * 11、返回200，上传对象4成功
     * 12、返回200，下载对象4成功，且返回信息中对象的存储类别为智能分级存储
     */
    @Test
    public void test_putObject_with_intelligent_tiering_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String[] bucketNameS = {
            bucketName + "-standard",
            bucketName + "-warm",
            bucketName + "-cold",
            bucketName + "-deep-archive",
            bucketName + "-intelligent-tiering"
        };
        StorageClassEnum[] storageClasses = {
            StorageClassEnum.STANDARD,
            StorageClassEnum.WARM,
            StorageClassEnum.COLD,
            StorageClassEnum.DEEP_ARCHIVE,
            StorageClassEnum.INTELLIGENT_TIERING
        };
        String objectKey = bucketName + "testObjectKey";
        for (int i = 0; i < bucketNameS.length; i++) {
            String bucketNameOBS = bucketNameS[i];
            try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
                System.out.println("test " + storageClasses[i].getCode());
                assert obsClient != null;
                CreateBucketRequest createBucketRequest = new CreateBucketRequest();
                createBucketRequest.setBucketName(bucketNameOBS);
                createBucketRequest.setBucketStorageClass(storageClasses[i]);
                ObsBucket bucket = obsClient.createBucket(createBucketRequest);
                Assert.assertEquals(200, bucket.getStatusCode());
                BucketMetadataInfoResult bucketMetadataInfoResult =
                        obsClient.getBucketMetadata(new BucketMetadataInfoRequest(bucketNameOBS));
                Assert.assertEquals(storageClasses[i], bucketMetadataInfoResult.getBucketStorageClass());

                String testFileName = bucketNameOBS + ".testFile";
                // 1 mb test file
                File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
                PutObjectRequest putObjectRequest = new PutObjectRequest(bucketNameOBS, objectKey, testFile);
                ObjectMetadata objectMetadata = new ObjectMetadata();
                objectMetadata.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
                putObjectRequest.setMetadata(objectMetadata);
                PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
                assertEquals(200, putObjectResult.getStatusCode());
                ObjectMetadata getObjectMetadata = obsClient.getObjectMetadata(bucketNameOBS, objectKey);
                Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, getObjectMetadata.getObjectStorageClass());
            } catch (ObsException e) {
                TestTools.printObsException(e);
                throw e;
            } catch (Exception e) {
                TestTools.printException(e);
                throw e;
            } finally {
                deleteBucketIgnoreError(bucketNameOBS);
            }
        }
    }

    /***
     * 1、创建标准桶
     * 2、上传标准对象1
     * 3、上传温对象2
     * 4、上传冷对象3
     * 5、上传深度归档对象4
     * 6、复制标准对象1为智能分级对象1，并查询对象存储类别
     * 7、复制温对象2为智能分级对象2，并查询对象存储类别
     * 8、复制冷对象3(未取回)为智能分级对象3，并查询对象存储类别
     * 9、复制冷对象3(已取回)为智能分级对象3，并查询对象存储类别
     * 10、复制深度归档对象4(未取回)为智能分级对象4，并查询对象存储类别
     * 11、复制深度归档对象4(已取回)为智能分级对象4，并查询对象存储类别
     * 12、复制智能分级对象1为标准对象1，并查询对象存储类别
     * 13、复制智能分级对象2为温对象2，并查询对象存储类别
     * 14、复制智能分级对象3为冷对象3，并查询对象存储类别
     * 15、复制智能分级对象4为深度归档对象4，并查询对象存储类别
     * 预期结果:
     * 1、返回200，创桶成功
     * 2、返回200，上传标准对象1成功
     * 3、返回200，上传温对象2成功
     * 4、返回200，上传冷对象3成功
     * 5、返回200，上传深度归档对象4成功
     * 6、返回200，复制为智能分级对象1成功，查询是智能分级对象
     * 7、返回200，复制为智能分级对象2成功，查询是智能分级对象
     * 8、返回405，复制为智能分级对象3失败，查询是冷对象
     * 9、返回200，复制为智能分级对象3成功，查询是智能分级对象
     * 10、返回405，复制为智能分级对象4失败，查询是深度归档对象
     * 11、返回200，复制为智能分级对象4成功，查询是智能分级对象
     * 12、返回200，复制为标准对象1成功，查询是标准对象
     * 13、返回200，复制为温对象2成功，查询是温对象
     * 14、返回200，复制为冷对象3成功，查询是冷对象
     * 15、返回200，复制为深度归档对象4成功，查询是深度归档对象
     */
    @Test
    public void test_copyObject_with_intelligent_tiering_001() throws IOException, InterruptedException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "testObjectKey";
        String objectKeyWarm = bucketName + "testObjectKeyWarm";
        String objectKeyCold = bucketName + "testObjectKeyCold";
        String objectKeyDeepArchive = bucketName + "testObjectKeyDeepArchive";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            String testFileName = bucketName + ".testFile";
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            ObjectMetadata objectMetadata = new ObjectMetadata();
            objectMetadata.setObjectStorageClass(StorageClassEnum.STANDARD);
            putObjectRequest.setMetadata(objectMetadata);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyWarm);
            objectMetadata.setObjectStorageClass(StorageClassEnum.WARM);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyCold);
            objectMetadata.setObjectStorageClass(StorageClassEnum.COLD);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyDeepArchive);
            objectMetadata.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            String intelligentObjectKey1 = "intelligentObjectKey1";
            CopyObjectRequest copyObjectRequest =
                    new CopyObjectRequest(bucketName, objectKey, bucketName, intelligentObjectKey1);
            ObjectMetadata objectMetadataNew = new ObjectMetadata();
            objectMetadataNew.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            copyObjectRequest.setNewObjectMetadata(objectMetadataNew);
            CopyObjectResult copyObjectResult = obsClient.copyObject(copyObjectRequest);
            assertEquals(200, copyObjectResult.getStatusCode());
            ObjectMetadata getObjectMetadata = obsClient.getObjectMetadata(bucketName, intelligentObjectKey1);
            Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, getObjectMetadata.getObjectStorageClass());

            String intelligentObjectKey2 = "intelligentObjectKey2";
            copyObjectRequest.setSourceObjectKey(objectKeyWarm);
            copyObjectRequest.setDestinationObjectKey(intelligentObjectKey2);
            copyObjectResult = obsClient.copyObject(copyObjectRequest);
            assertEquals(200, copyObjectResult.getStatusCode());
            getObjectMetadata = obsClient.getObjectMetadata(bucketName, intelligentObjectKey2);
            Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, getObjectMetadata.getObjectStorageClass());

            String intelligentObjectKey3 = "intelligentObjectKey3";
            copyObjectRequest.setSourceObjectKey(objectKeyCold);
            copyObjectRequest.setDestinationObjectKey(intelligentObjectKey3);
            try {
                obsClient.copyObject(copyObjectRequest);
                fail();
            } catch (ObsException obsException) {
                assertTrue(400 <= obsException.getResponseCode());
            }
            RestoreObjectRequest restoreObjectRequest = new RestoreObjectRequest(bucketName, objectKeyCold, 1);
            if (isTestingRestore()) {
                restoreObjectRequest.setRestoreTier(RestoreTierEnum.EXPEDITED);
                obsClient.restoreObject(restoreObjectRequest);
                waitForRestore(obsClient, bucketName, objectKeyCold, 30 * 1000);
                copyObjectRequest.setSourceObjectKey(objectKeyCold);
                copyObjectRequest.setDestinationObjectKey(intelligentObjectKey3);
                copyObjectResult = obsClient.copyObject(copyObjectRequest);
                assertEquals(200, copyObjectResult.getStatusCode());
                getObjectMetadata = obsClient.getObjectMetadata(bucketName, intelligentObjectKey3);
                Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, getObjectMetadata.getObjectStorageClass());
            }

            String intelligentObjectKey4 = "intelligentObjectKey4";
            copyObjectRequest.setSourceObjectKey(objectKeyDeepArchive);
            copyObjectRequest.setDestinationObjectKey(intelligentObjectKey4);
            try {
                obsClient.copyObject(copyObjectRequest);
                fail();
            } catch (ObsException obsException) {
                assertTrue(400 <= obsException.getResponseCode());
            }

            if (isTestingRestore()) {
                restoreObjectRequest.setObjectKey(objectKeyDeepArchive);
                obsClient.restoreObject(restoreObjectRequest);
                waitForRestore(obsClient, bucketName, objectKeyDeepArchive, 10 * 60 * 1000);
                copyObjectRequest.setSourceObjectKey(objectKeyDeepArchive);
                copyObjectRequest.setDestinationObjectKey(intelligentObjectKey4);
                copyObjectResult = obsClient.copyObject(copyObjectRequest);
                assertEquals(200, copyObjectResult.getStatusCode());
                getObjectMetadata = obsClient.getObjectMetadata(bucketName, intelligentObjectKey4);
                Assert.assertEquals(StorageClassEnum.INTELLIGENT_TIERING, getObjectMetadata.getObjectStorageClass());
            }

            String testCopyDestObject = "testCopyDestObject";
            copyObjectRequest.setSourceObjectKey(intelligentObjectKey1);
            copyObjectRequest.setDestinationObjectKey(testCopyDestObject);
            objectMetadataNew.setObjectStorageClass(StorageClassEnum.STANDARD);
            copyObjectRequest.setNewObjectMetadata(objectMetadataNew);
            copyObjectResult = obsClient.copyObject(copyObjectRequest);
            assertEquals(200, copyObjectResult.getStatusCode());
            getObjectMetadata = obsClient.getObjectMetadata(bucketName, testCopyDestObject);
            assertNull(getObjectMetadata.getObjectStorageClass());

            copyObjectRequest.setSourceObjectKey(intelligentObjectKey2);
            copyObjectRequest.setDestinationObjectKey(testCopyDestObject);
            objectMetadataNew.setObjectStorageClass(StorageClassEnum.WARM);
            copyObjectRequest.setNewObjectMetadata(objectMetadataNew);
            copyObjectResult = obsClient.copyObject(copyObjectRequest);
            assertEquals(200, copyObjectResult.getStatusCode());
            getObjectMetadata = obsClient.getObjectMetadata(bucketName, testCopyDestObject);
            Assert.assertEquals(StorageClassEnum.WARM, getObjectMetadata.getObjectStorageClass());

            if (isTestingRestore()) {
                copyObjectRequest.setSourceObjectKey(intelligentObjectKey3);
            } else {
                copyObjectRequest.setSourceObjectKey(intelligentObjectKey2);
            }
            copyObjectRequest.setDestinationObjectKey(testCopyDestObject);
            objectMetadataNew.setObjectStorageClass(StorageClassEnum.COLD);
            copyObjectRequest.setNewObjectMetadata(objectMetadataNew);
            copyObjectResult = obsClient.copyObject(copyObjectRequest);
            assertEquals(200, copyObjectResult.getStatusCode());
            getObjectMetadata = obsClient.getObjectMetadata(bucketName, testCopyDestObject);
            Assert.assertEquals(StorageClassEnum.COLD, getObjectMetadata.getObjectStorageClass());

            if (isTestingRestore()) {
                copyObjectRequest.setSourceObjectKey(intelligentObjectKey4);
            } else {
                copyObjectRequest.setSourceObjectKey(intelligentObjectKey2);
            }
            copyObjectRequest.setDestinationObjectKey(testCopyDestObject);
            objectMetadataNew.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
            copyObjectRequest.setNewObjectMetadata(objectMetadataNew);
            copyObjectResult = obsClient.copyObject(copyObjectRequest);
            assertEquals(200, copyObjectResult.getStatusCode());
            getObjectMetadata = obsClient.getObjectMetadata(bucketName, testCopyDestObject);
            Assert.assertEquals(StorageClassEnum.DEEP_ARCHIVE, getObjectMetadata.getObjectStorageClass());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    /***
     * 1、创建标准桶
     * 2、上传标准对象1
     * 3、上传温对象2
     * 4、上传冷对象3
     * 5、上传深度归档对象4
     * 6、修改标准对象1为智能分级对象1，并查询对象存储类别
     * 7、修改温对象2为智能分级对象2，并查询对象存储类别
     * 8、修改冷对象3(未取回)为智能分级对象3，并查询对象存储类别
     * 9、修改冷对象3(已取回)为智能分级对象3，并查询对象存储类别
     * 10、修改深度归档对象4(未取回)为智能分级对象4，并查询对象存储类别
     * 11、修改深度归档对象4(已取回)为智能分级对象4，并查询对象存储类别
     * 12、修改智能分级对象1为标准对象1，并查询对象存储类别
     * 13、修改智能分级对象2为温对象2，并查询对象存储类别
     * 14、修改智能分级对象3为冷对象3，并查询对象存储类别
     * 15、修改智能分级对象4为深度归档对象4，并查询对象存储类别
     * 预期结果:
     * 1、返回200，创桶成功
     * 2、返回200，上传标准对象1成功
     * 3、返回200，上传温对象1成功
     * 4、返回200，上传冷对象1成功
     * 5、返回200，上传深度归档对象1成功
     * 6、返回200，修改为智能分级对象1成功，查询是智能分级对象
     * 7、返回200，修改为智能分级对象2成功，查询是智能分级对象
     * 8、返回405，修改为智能分级对象3失败，查询是冷对象
     * 9、返回200，修改为智能分级对象3成功，查询是智能分级对象
     * 10、返回405，修改为智能分级对象4失败，查询是深度归档对象
     * 11、返回200，修改为智能分级对象4成功，查询是智能分级对象
     * 12、返回200，修改为标准对象1成功，查询是标准对象
     * 13、返回200，修改为温对象2成功，查询是温对象
     * 14、返回200，修改为冷对象3成功，查询是冷对象
     * 15、返回200，修改为深度归档对象4成功，查询是深度归档对象
     */
    @Test
    public void test_setObjectMetadata_with_intelligent_tiering_001() throws IOException, InterruptedException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "testObjectKey";
        String objectKeyWarm = bucketName + "testObjectKeyWarm";
        String objectKeyCold = bucketName + "testObjectKeyCold";
        String objectKeyDeepArchive = bucketName + "testObjectKeyDeepArchive";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            String testFileName = bucketName + ".testFile";
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            ObjectMetadata objectMetadata = new ObjectMetadata();
            objectMetadata.setObjectStorageClass(StorageClassEnum.STANDARD);
            putObjectRequest.setMetadata(objectMetadata);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyWarm);
            objectMetadata.setObjectStorageClass(StorageClassEnum.WARM);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyCold);
            objectMetadata.setObjectStorageClass(StorageClassEnum.COLD);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyDeepArchive);
            objectMetadata.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            SetObjectMetadataRequest setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKey);
            setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            ObjectMetadata objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
            Assert.assertEquals(200, objectMetadata1.getStatusCode());
            ObjectMetadata objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKey);
            assertEquals(objectMetadataGet.getObjectStorageClass(), StorageClassEnum.INTELLIGENT_TIERING);

            setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKeyWarm);
            setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
            Assert.assertEquals(200, objectMetadata1.getStatusCode());
            objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKeyWarm);
            assertEquals(objectMetadataGet.getObjectStorageClass(), StorageClassEnum.INTELLIGENT_TIERING);

            setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKeyCold);
            setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            try {
                obsClient.setObjectMetadata(setObjectMetadataRequest);
                fail();
            } catch (ObsException obsException) {
                assertTrue(400 <= obsException.getResponseCode());
            }
            RestoreObjectRequest restoreObjectRequest = new RestoreObjectRequest(bucketName, objectKeyCold, 1);
            restoreObjectRequest.setRestoreTier(RestoreTierEnum.EXPEDITED);
            if (isTestingRestore()) {
                obsClient.restoreObject(restoreObjectRequest);
                waitForRestore(obsClient, bucketName, objectKeyCold, 30 * 1000);
                objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
                Assert.assertEquals(200, objectMetadata1.getStatusCode());
                objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKeyCold);
                assertEquals(objectMetadataGet.getObjectStorageClass(), StorageClassEnum.INTELLIGENT_TIERING);
            }

            setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKeyDeepArchive);
            setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            try {
                obsClient.setObjectMetadata(setObjectMetadataRequest);
                fail();
            } catch (ObsException obsException) {
                assertTrue(400 <= obsException.getResponseCode());
            }
            restoreObjectRequest = new RestoreObjectRequest(bucketName, objectKeyDeepArchive, 1);
            restoreObjectRequest.setRestoreTier(RestoreTierEnum.EXPEDITED);
            if (isTestingRestore()) {
                obsClient.restoreObject(restoreObjectRequest);
                waitForRestore(obsClient, bucketName, objectKeyDeepArchive, 10 * 60 * 1000);
                objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
                Assert.assertEquals(200, objectMetadata1.getStatusCode());
                objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKeyDeepArchive);
                assertEquals(objectMetadataGet.getObjectStorageClass(), StorageClassEnum.INTELLIGENT_TIERING);
            }

            setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKey);
            setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.STANDARD);
            objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
            Assert.assertEquals(200, objectMetadata1.getStatusCode());
            objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKey);
            assertNull(objectMetadataGet.getObjectStorageClass());

            setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKeyWarm);
            setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.WARM);
            objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
            Assert.assertEquals(200, objectMetadata1.getStatusCode());
            objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKeyWarm);
            assertEquals(objectMetadataGet.getObjectStorageClass(), StorageClassEnum.WARM);

            if (isTestingRestore()) {
                setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKeyCold);
                setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.COLD);
                objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
                Assert.assertEquals(200, objectMetadata1.getStatusCode());
                objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKeyCold);
                assertEquals(objectMetadataGet.getObjectStorageClass(), StorageClassEnum.COLD);

                setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKeyDeepArchive);
                setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
                objectMetadata1 = obsClient.setObjectMetadata(setObjectMetadataRequest);
                Assert.assertEquals(200, objectMetadata1.getStatusCode());
                objectMetadataGet = obsClient.getObjectMetadata(bucketName, objectKeyDeepArchive);
                assertEquals(objectMetadataGet.getObjectStorageClass(), StorageClassEnum.DEEP_ARCHIVE);
            }
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    /***
     * 1、创建标准桶
     * 2、初始化段任务，指定对象存储类别为智能分级
     * 3、列举桶内多段任务
     * 3、上传段1
     * 4、创建标准对象1
     * 5、拷贝对象1为段2
     * 6、创建温对象1
     * 7、拷贝温对象1为段3
     * 8、创建冷对象1
     * 9、拷贝冷对象1(未取回)为段4
     * 10、拷贝冷对象1(已取回)为段4
     * 11、创建深度归档对象1
     * 12、拷贝深度归档对象1(未取回)为段5
     * 13、拷贝深度归档对象1(已取回)为段5
     * 14、创建智能分级对象1
     * 15、拷贝智能分级对象1为段6
     * 16、列举桶中已初始化的段任务
     * 17、列举已上传的段
     * 18、合并段
     * 19、初始化段任务，指定对象存储类别为标准，复拷贝智能分级对象1为段1，合并段
     * 20、初始化段任务，指定对象存储类别为温，复拷贝智能分级对象1为段1，合并段
     * 21、初始化段任务，指定对象存储类别为冷，复拷贝智能分级对象1为段1，合并段
     * 22、初始化段任务，指定对象存储类别为深度对党，复拷贝智能分级对象1为段1，合并段
     * 预期结果:
     * 1、返回200，创桶成功
     * 2、返回200，初始化段任务成功，且存储类型是智能分级
     * 3、返回200，上传段成功，且存储类型是智能分级
     * 4、返回200，创建对象成功，且对象是标准对象
     * 5、返回200，拷贝段成功，且段的存储类别是智能分级
     * 6、返回200，创建对象成功，且对象是温对象
     * 7、返回200，拷贝段成功，且段的存储类别是智能分级
     * 8、返回200，创建对象成功，且对象是冷对象
     * 9、返回403，拷贝段失败
     * 10、返回200，拷贝段成功，且段的存储类别是智能分级
     * 11、返回200，创建对象成功，且对象是深度归档对象
     * 12、返回403，拷贝段失败
     * 13、返回200，拷贝段成功，且段的存储类别是智能分级
     * 14、返回200，创建对象成功，且对象是智能分级对象
     * 15、返回200，拷贝段成功，且段的存储类别是智能分级
     * 16、返回200，列举桶中的已初始化任务，且段任务的存储类别是智能分级
     * 17、返回200，列举已上传的段，且段的存储类别是智能分级
     * 18、返回200，合并段成功，且对象的存储类别是智能分级存储
     * 19、创建对象成功，且对象的存储类别是标准
     * 20、创建对象成功，且对象的存储类别是温
     * 21、创建对象成功，且对象的存储类别是冷
     * 22、创建对象成功，且对象的存储类别是深度归档
     */
    @Test
    public void test_multipart_with_intelligent_tiering_001() throws IOException, InterruptedException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "testObjectKey";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            InitiateMultipartUploadRequest initiateMultipartUploadRequest =
                    new InitiateMultipartUploadRequest(bucketName, objectKey);
            ObjectMetadata objectMetadata = new ObjectMetadata();
            objectMetadata.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            initiateMultipartUploadRequest.setMetadata(objectMetadata);
            InitiateMultipartUploadResult initiateMultipartUploadResult =
                    obsClient.initiateMultipartUpload(initiateMultipartUploadRequest);
            Assert.assertEquals(200, initiateMultipartUploadResult.getStatusCode());
            String uploadID = initiateMultipartUploadResult.getUploadId();

            ListMultipartUploadsRequest listMultipartUploadsRequest = new ListMultipartUploadsRequest(bucketName);
            MultipartUploadListing multipartUploadListing = obsClient.listMultipartUploads(listMultipartUploadsRequest);
            for (MultipartUpload multipartUpload : multipartUploadListing.getMultipartTaskList()) {
                Assert.assertEquals(multipartUpload.getObjectStorageClass(), StorageClassEnum.INTELLIGENT_TIERING);
            }

            String testFileName = bucketName + ".testFile";
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
            obsClient.uploadPart(bucketName, objectKey, uploadID, 1, testFile);

            String objectKeyPart = bucketName + "testObjectKeyPart";
            objectMetadata.setObjectStorageClass(StorageClassEnum.STANDARD);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKeyPart, testFile);
            putObjectRequest.setMetadata(objectMetadata);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());
            // 标准存储类型不会返回
            assertNull(obsClient.getObjectMetadata(bucketName, objectKeyPart).getObjectStorageClass());

            CopyPartRequest copyPartRequest =
                    new CopyPartRequest(uploadID, bucketName, objectKeyPart, bucketName, objectKey, 2);
            assertEquals(200, obsClient.copyPart(copyPartRequest).getStatusCode());

            objectMetadata.setObjectStorageClass(StorageClassEnum.WARM);
            putObjectRequest.setFile(testFile);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());
            // WARM会返回
            assertEquals(
                    StorageClassEnum.WARM,
                    obsClient.getObjectMetadata(bucketName, objectKeyPart).getObjectStorageClass());
            copyPartRequest.setPartNumber(3);
            assertEquals(200, obsClient.copyPart(copyPartRequest).getStatusCode());

            objectMetadata.setObjectStorageClass(StorageClassEnum.COLD);
            putObjectRequest.setFile(testFile);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());
            // COLD
            assertEquals(
                    StorageClassEnum.COLD,
                    obsClient.getObjectMetadata(bucketName, objectKeyPart).getObjectStorageClass());
            copyPartRequest.setPartNumber(4);
            try {
                obsClient.copyPart(copyPartRequest);
                fail();
            } catch (ObsException obsException) {
                assertTrue(400 <= obsException.getResponseCode());
            }
            RestoreObjectRequest restoreObjectRequest = new RestoreObjectRequest(bucketName, objectKeyPart, 1);
            restoreObjectRequest.setRestoreTier(RestoreTierEnum.EXPEDITED);
            if (isTestingRestore()) {
                obsClient.restoreObject(restoreObjectRequest);
                waitForRestore(obsClient, bucketName, objectKeyPart, 30 * 1000);
                assertEquals(200, obsClient.copyPart(copyPartRequest).getStatusCode());
            }

            objectMetadata.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
            putObjectRequest.setFile(testFile);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());
            // DEEP_ARCHIVE
            assertEquals(
                    StorageClassEnum.DEEP_ARCHIVE,
                    obsClient.getObjectMetadata(bucketName, objectKeyPart).getObjectStorageClass());
            copyPartRequest.setPartNumber(5);
            try {
                obsClient.copyPart(copyPartRequest);
                fail();
            } catch (ObsException obsException) {
                assertTrue(400 <= obsException.getResponseCode());
            }

            if (isTestingRestore()) {
                obsClient.restoreObject(restoreObjectRequest);
                waitForRestore(obsClient, bucketName, objectKeyPart, 10 * 60 * 1000);
                assertEquals(200, obsClient.copyPart(copyPartRequest).getStatusCode());
            }

            objectMetadata.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            putObjectRequest.setFile(testFile);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());
            assertEquals(
                    StorageClassEnum.INTELLIGENT_TIERING,
                    obsClient.getObjectMetadata(bucketName, objectKeyPart).getObjectStorageClass());

            if (isTestingRestore()) {
                copyPartRequest.setPartNumber(6);
            } else {
                copyPartRequest.setPartNumber(4);
            }

            assertEquals(200, obsClient.copyPart(copyPartRequest).getStatusCode());

            listMultipartUploadsRequest = new ListMultipartUploadsRequest(bucketName);
            multipartUploadListing = obsClient.listMultipartUploads(listMultipartUploadsRequest);
            for (MultipartUpload multipartUpload : multipartUploadListing.getMultipartTaskList()) {
                Assert.assertEquals(multipartUpload.getObjectStorageClass(), StorageClassEnum.INTELLIGENT_TIERING);
            }

            ListPartsRequest listPartsRequest = new ListPartsRequest(bucketName, objectKey, uploadID);
            ListPartsResult listPartsResult = obsClient.listParts(listPartsRequest);
            Assert.assertEquals(200, listPartsResult.getStatusCode());

            List<PartEtag> partETags = new ArrayList<>();
            for (Multipart multipart : listPartsResult.getMultipartList()) {
                partETags.add(new PartEtag(multipart.getEtag(), multipart.getPartNumber()));
            }
            CompleteMultipartUploadRequest completeMultipartUploadRequest =
                    new CompleteMultipartUploadRequest(bucketName, objectKey, uploadID, partETags);
            CompleteMultipartUploadResult completeMultipartUploadResult =
                    obsClient.completeMultipartUpload(completeMultipartUploadRequest);
            assertEquals(200, completeMultipartUploadResult.getStatusCode());
            assertEquals(
                    StorageClassEnum.INTELLIGENT_TIERING,
                    obsClient.getObjectMetadata(bucketName, objectKey).getObjectStorageClass());

            StorageClassEnum[] storageClassEnums = {
                StorageClassEnum.STANDARD, StorageClassEnum.WARM, StorageClassEnum.COLD, StorageClassEnum.DEEP_ARCHIVE
            };
            StorageClassEnum[] returnStorageClassEnums = {
                null, StorageClassEnum.WARM, StorageClassEnum.COLD, StorageClassEnum.DEEP_ARCHIVE
            };
            for (int i = 0; i < storageClassEnums.length; ++i) {
                initiateMultipartUploadRequest = new InitiateMultipartUploadRequest(bucketName, objectKey);
                objectMetadata.setObjectStorageClass(storageClassEnums[i]);
                initiateMultipartUploadRequest.setMetadata(objectMetadata);
                initiateMultipartUploadResult = obsClient.initiateMultipartUpload(initiateMultipartUploadRequest);
                Assert.assertEquals(200, initiateMultipartUploadResult.getStatusCode());
                uploadID = initiateMultipartUploadResult.getUploadId();
                copyPartRequest.setPartNumber(1);
                copyPartRequest.setUploadId(uploadID);
                CopyPartResult copyPartResult = obsClient.copyPart(copyPartRequest);
                assertEquals(200, copyPartResult.getStatusCode());
                completeMultipartUploadRequest.setUploadId(uploadID);
                completeMultipartUploadRequest.setObjectKey(objectKey);
                partETags.clear();
                partETags.add(new PartEtag(copyPartResult.getEtag(), 1));
                completeMultipartUploadResult = obsClient.completeMultipartUpload(completeMultipartUploadRequest);
                assertEquals(200, completeMultipartUploadResult.getStatusCode());
                assertEquals(
                        returnStorageClassEnums[i],
                        obsClient.getObjectMetadata(bucketName, objectKey).getObjectStorageClass());
            }
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    /***
     * 1、创建智能分级桶
     * 2、上传对象，存储类别继承桶配置
     * 3、上传智能分级对象
     * 4、列举桶内对象
     * 5、开启桶的多版本
     * 6、覆盖步骤2上传的对象
     * 7、列举桶的多版本对象
     * 预期结果:
     * 1、返回200，创桶成功
     * 2、返回200，上传对象成功
     * 3、返回200，上传对象成功
     * 4、返回200，存储类型是智能分级存储
     * 5、返回200，开启多版本成功
     * 6、返回200，上传对象成功
     * 7、返回200，存储类型是智能分级存储
     */
    @Test
    public void test_listObjects_with_intelligent_tiering_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNameIntelligent = bucketName + "-intelligent";
        String objectKey = bucketName + "testObjectKey";
        String objectKeyIntelligent = objectKey + "-intelligent";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNameIntelligent);
            createBucketRequest.setBucketStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            HeaderResponse response = obsClient.createBucket(createBucketRequest);
            assertEquals(200, response.getStatusCode());

            String testFileName = bucketName + ".testFile";
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketNameIntelligent, objectKey, testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyIntelligent);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            ObjectListing objectListing = obsClient.listObjects(bucketNameIntelligent);
            assertEquals(200, objectListing.getStatusCode());

            for (ObsObject obsObject : objectListing.getObjects()) {
                assertEquals(StorageClassEnum.INTELLIGENT_TIERING, obsObject.getMetadata().getObjectStorageClass());
            }

            BucketVersioningConfiguration bucketVersioningConfiguration =
                    new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED);
            response = obsClient.setBucketVersioning(bucketNameIntelligent, bucketVersioningConfiguration);
            assertEquals(200, response.getStatusCode());

            putObjectRequest.setFile(testFile);
            putObjectRequest.setObjectKey(objectKeyIntelligent);
            putObjectResult = obsClient.putObject(putObjectRequest);
            assertEquals(200, putObjectResult.getStatusCode());

            ListVersionsResult listVersionsResult = obsClient.listVersions(bucketNameIntelligent);
            assertEquals(200, listVersionsResult.getStatusCode());

            for (VersionOrDeleteMarker version : listVersionsResult.getVersions()) {
                assertEquals(StorageClassEnum.INTELLIGENT_TIERING, version.getObjectStorageClass());
            }
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        } finally {
            deleteBucketIgnoreError(bucketNameIntelligent);
        }
    }

    /***
     * 1、创建posix桶A，桶存储类别为智能分级存储
     * 2、创建posix桶B,桶存储类别默认
     * 3、上传对象到桶B，对象存储类型为智能分级
     * 预期结果:
     * 1、创建失败
     * 2、创建成功
     * 3、上传失败
     */
    @Test
    public void test_posixbucket_with_intelligent_tiering_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNamePosixIntelligent = bucketName + "-posix-intelligent";
        String bucketNamePosix = bucketName + "-posix";
        String objectKey = bucketName + "testObjectKey";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNamePosixIntelligent);
            createBucketRequest.setBucketStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            createBucketRequest.setBucketType(BucketTypeEnum.PFS);
            try {
                obsClient.createBucket(createBucketRequest);
                fail();
            } catch (ObsException e) {
            }

            createBucketRequest.setBucketName(bucketNamePosix);
            createBucketRequest.setBucketStorageClass(null);
            HeaderResponse response = obsClient.createBucket(createBucketRequest);
            assertEquals(200, response.getStatusCode());

            String testFileName = bucketName + ".testFile";
            // 1 mb test file
            File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketNamePosix, objectKey, testFile);
            ObjectMetadata objectMetadata = new ObjectMetadata();
            objectMetadata.setObjectStorageClass(StorageClassEnum.INTELLIGENT_TIERING);
            putObjectRequest.setMetadata(objectMetadata);
            try {
                obsClient.putObject(putObjectRequest);
                fail();
            } catch (ObsException e) {
            }
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        } finally {
            deleteBucketIgnoreError(bucketNamePosixIntelligent);
            deleteBucketIgnoreError(bucketNamePosix);
        }
    }

    public void deleteBucketIgnoreError(String bucket) {
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            deleteObjects(obsClient, bucket);
            obsClient.deleteBucket(bucket);
        } catch (Throwable ignore) {
        }
    }

    public static void waitForRestore(ObsClient obsClient, String bucketName, String objectKey, long sleepMs)
            throws InterruptedException {
        boolean restored;
        do {
            restored = false;
            ObjectMetadata objectMetadata = obsClient.getObjectMetadata(bucketName, objectKey);
            String restoreStateHeader = (String) objectMetadata.getResponseHeaders().get("restore");
            if (restoreStateHeader != null
                    && restoreStateHeader.startsWith("ongoing-request=\"false\", expiry-date=\"")) {
                restored = true;
            } else {
                System.out.println("object:" + objectKey + " not restored, bucket:" + bucketName);
                Thread.sleep(sleepMs);
            }
        } while (!restored);
    }

    public boolean isTestingRestore() {
        String customParameter = System.getProperty("custom.test.restore");
        return "true".equals(customParameter);
    }
}
