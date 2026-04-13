/**
 * Copyright 2019 Huawei Technologies Co.,Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package com.obs.integrated_test;
import com.obs.test.TestTools;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.BucketMetadataInfoRequest;
import com.obs.services.model.BucketMetadataInfoResult;
import com.obs.services.model.BucketStoragePolicyConfiguration;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.CopyObjectRequest;
import com.obs.services.model.CopyObjectResult;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.RestoreObjectRequest;
import com.obs.services.model.SetObjectMetadataRequest;
import com.obs.services.model.SseKmsHeader;
import com.obs.services.model.StorageClassEnum;
import com.obs.services.model.UploadFileRequest;
import com.obs.test.tools.PropertiesTools;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Locale;

import static com.obs.test.TestTools.genTestFile;
import static org.junit.Assert.assertEquals;

public class DeepArchiveIT {
    private static final ArrayList<String> createdBuckets = new ArrayList<>();
    private static String beforeBucket;
    @Rule
    public TestName testName = new TestName();

    private static ObsClient obsClient = TestTools.getCustomPipelineEnvironment();

    @BeforeClass
    public static void create_demo_bucket() {
        beforeBucket = "create-demo-bucket";
        CreateBucketRequest request = new CreateBucketRequest();
        request.setBucketType(BucketTypeEnum.OBJECT);
        request.setBucketName(beforeBucket + "-bucket-deep-archive");
        request.setBucketStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        HeaderResponse response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket-deep-archive");

        request.setBucketName(beforeBucket + "-bucket-standard");
        request.setBucketStorageClass(StorageClassEnum.STANDARD);
        response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket-standard");

        request.setBucketName(beforeBucket + "-bucket-warm");
        request.setBucketStorageClass(StorageClassEnum.WARM);
        response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket-warm");

        request.setBucketName(beforeBucket + "-bucket-cold");
        request.setBucketStorageClass(StorageClassEnum.COLD);
        response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket-cold");

        request.setBucketName(beforeBucket + "-bucket-standard-test");
        request.setBucketStorageClass(StorageClassEnum.STANDARD);
        response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket-standard-test");
    }

    @AfterClass
    public static void delete_created_buckets() {
        for (String bucket : createdBuckets) {
            obsClient.deleteBucket(bucket);
        }
    }

    @Test
    public void test_set_bucket_store_type() {
        BucketStoragePolicyConfiguration storgePolicy = new BucketStoragePolicyConfiguration();
        storgePolicy.setBucketStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        HeaderResponse response = obsClient.setBucketStoragePolicy(beforeBucket + "-bucket-standard",
                storgePolicy);
        assertEquals(200, response.getStatusCode());

        BucketMetadataInfoRequest bucketMetadataInfoRequest = new BucketMetadataInfoRequest();
        bucketMetadataInfoRequest.setBucketName(beforeBucket + "-bucket-standard");
        BucketMetadataInfoResult bucketMetadataInfoResult = obsClient.getBucketMetadata(bucketMetadataInfoRequest);
        assertEquals(StorageClassEnum.DEEP_ARCHIVE, bucketMetadataInfoResult.getBucketStorageClass());

        response = obsClient.setBucketStoragePolicy(beforeBucket + "-bucket-warm", storgePolicy);
        assertEquals(200, response.getStatusCode());
        response = obsClient.setBucketStoragePolicy(beforeBucket + "-bucket-cold", storgePolicy);
        assertEquals(200, response.getStatusCode());

        storgePolicy.setBucketStorageClass(StorageClassEnum.STANDARD);
        response = obsClient.setBucketStoragePolicy(beforeBucket + "-bucket-cold", storgePolicy);
        assertEquals(200, response.getStatusCode());
        storgePolicy.setBucketStorageClass(StorageClassEnum.WARM);
        response = obsClient.setBucketStoragePolicy(beforeBucket + "-bucket-warm", storgePolicy);
        assertEquals(200, response.getStatusCode());
        storgePolicy.setBucketStorageClass(StorageClassEnum.COLD);
        response = obsClient.setBucketStoragePolicy(beforeBucket + "-bucket-standard", storgePolicy);
        assertEquals(200, response.getStatusCode());
    }

    @Test
    public void test_upload_download_file_for_deep_archive() throws IOException {
        PutObjectRequest request = new PutObjectRequest();
        request.setBucketName(beforeBucket + "-bucket-standard-test");
        request.setObjectKey("test_object");
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        request.setMetadata(metadata);
        request.setInput(new ByteArrayInputStream("Hello OBS".getBytes(StandardCharsets.UTF_8)));
        PutObjectResult result = obsClient.putObject(request);
        assertEquals(200, result.getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test", "test_object");

        SseKmsHeader kms = new SseKmsHeader();
        request.setSseKmsHeader(kms);
        request.setObjectKey("test_object_kms");
        result = obsClient.putObject(request);
        assertEquals(200, result.getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test", "test_object_kms");

        genTestFile("test_upload_file_for_deep_archive", 1024 * 100);
        UploadFileRequest uploadFileRequest = new UploadFileRequest(beforeBucket + "-bucket-standard-test", "test_object_upload_file");
        uploadFileRequest.setUploadFile("test_upload_file_for_deep_archive");
        uploadFileRequest.setObjectMetadata(metadata);
        CompleteMultipartUploadResult completeMultipartUploadResult = obsClient.uploadFile(uploadFileRequest);
        assertEquals(200, completeMultipartUploadResult.getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test", "test_object_upload_file");

        UploadFileRequest uploadFileDeepArchiveBucketRequest = new UploadFileRequest(beforeBucket +
                "-bucket-deep-archive", "test_upload_file_for_bucket_deep_archive");
        uploadFileDeepArchiveBucketRequest.setUploadFile("test_upload_file_for_deep_archive");
        CompleteMultipartUploadResult completeMultipartUploadDeepArchiveBucketResult
                = obsClient.uploadFile(uploadFileDeepArchiveBucketRequest);
        assertEquals(200, completeMultipartUploadDeepArchiveBucketResult.getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-deep-archive",
                "test_upload_file_for_bucket_deep_archive");

        PutObjectRequest putObjectRequest = new PutObjectRequest();
        putObjectRequest.setBucketName(beforeBucket + "-bucket-deep-archive");
        putObjectRequest.setObjectKey("test_object_deep_archive_bucket");
        putObjectRequest.setInput(new ByteArrayInputStream("Hello OBS".getBytes("UTF-8")));
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        assertEquals(200, putObjectResult.getStatusCode());
        try {
            ObsObject object = obsClient.getObject(beforeBucket + "-bucket-deep-archive",
                    "test_object_deep_archive_bucket");
        } catch (ObsException ex) {
            assertEquals("InvalidObjectState", ex.getErrorCode());
        }
        obsClient.deleteObject(beforeBucket + "-bucket-deep-archive", "test_object_deep_archive_bucket");
    }

    @Test
    public void test_set_store_type() throws UnsupportedEncodingException {
        PutObjectRequest request = new PutObjectRequest();
        request.setBucketName(beforeBucket + "-bucket-standard-test");
        request.setObjectKey("test_object");
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setObjectStorageClass(StorageClassEnum.STANDARD);
        request.setMetadata(metadata);
        request.setObjectKey("test_object_standard_to_deep_archive");
        request.setInput(new ByteArrayInputStream("Hello OBS".getBytes("UTF-8")));
        PutObjectResult putObjectResult = obsClient.putObject(request);
        assertEquals(200, putObjectResult.getStatusCode());
        SetObjectMetadataRequest setObjectMetadataRequest = new SetObjectMetadataRequest(beforeBucket +
                "-bucket-standard-test", "test_object_standard_to_deep_archive");
        setObjectMetadataRequest.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        assertEquals(200, obsClient.setObjectMetadata(setObjectMetadataRequest).getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test",
                "test_object_standard_to_deep_archive");

        metadata.setObjectStorageClass(StorageClassEnum.COLD);
        request.setMetadata(metadata);
        request.setObjectKey("test_object_cold_to_deep_archive");
        PutObjectResult putObjectColdObjectResult = obsClient.putObject(request);
        assertEquals(200, putObjectColdObjectResult.getStatusCode());
        SetObjectMetadataRequest setObjectMetadataStandardRequest = new SetObjectMetadataRequest(beforeBucket +
                "-bucket-standard-test", "test_object_cold_to_deep_archive");
        setObjectMetadataStandardRequest.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        try {
            obsClient.setObjectMetadata(setObjectMetadataStandardRequest);
        } catch (ObsException ex) {
            assertEquals("InvalidObjectState", ex.getErrorCode());
        }
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test", "test_object_cold_to_deep_archive");

        metadata.setObjectStorageClass(StorageClassEnum.WARM);
        request.setMetadata(metadata);
        request.setObjectKey("test_object_warm_to_deep_archive");
        PutObjectResult putObjectWarnResult = obsClient.putObject(request);
        assertEquals(200, putObjectWarnResult.getStatusCode());
        SetObjectMetadataRequest setObjectMetadataDeepArchiveRequest = new SetObjectMetadataRequest(beforeBucket + "-bucket-standard-test", "test_object_warm_to_deep_archive");
        setObjectMetadataDeepArchiveRequest.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        assertEquals(200, obsClient.setObjectMetadata(setObjectMetadataDeepArchiveRequest).getStatusCode());

        setObjectMetadataDeepArchiveRequest.setObjectStorageClass(StorageClassEnum.COLD);
        try {
            obsClient.setObjectMetadata(setObjectMetadataDeepArchiveRequest);
        } catch (ObsException ex) {
            assertEquals("InvalidObjectState", ex.getErrorCode());
        }

        setObjectMetadataDeepArchiveRequest.setObjectStorageClass(StorageClassEnum.WARM);
        try {
            obsClient.setObjectMetadata(setObjectMetadataDeepArchiveRequest);
        } catch (ObsException ex) {
            assertEquals("InvalidObjectState", ex.getErrorCode());
        }

        setObjectMetadataDeepArchiveRequest.setObjectStorageClass(StorageClassEnum.STANDARD);
        try {
            obsClient.setObjectMetadata(setObjectMetadataDeepArchiveRequest);
        } catch (ObsException ex) {
            assertEquals("InvalidObjectState", ex.getErrorCode());
        }
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test", "test_object_warm_to_deep_archive");
    }

    @Test
    public void test_copy_object_for_deep_archive() throws UnsupportedEncodingException {
        PutObjectRequest request = new PutObjectRequest();
        request.setBucketName(beforeBucket + "-bucket-standard-test");
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setObjectStorageClass(StorageClassEnum.STANDARD);
        request.setMetadata(metadata);
        request.setObjectKey("test_object_standard_copy_to_deep_archive");
        request.setInput(new ByteArrayInputStream("Hello OBS".getBytes("UTF-8")));
        obsClient.putObject(request);

        metadata.setObjectStorageClass(StorageClassEnum.WARM);
        request.setMetadata(metadata);
        request.setObjectKey("test_object_warm_copy_to_deep_archive");
        obsClient.putObject(request);

        metadata.setObjectStorageClass(StorageClassEnum.COLD);
        request.setMetadata(metadata);
        request.setObjectKey("test_object_cold_copy_to_deep_archive");
        obsClient.putObject(request);

        CopyObjectRequest copyObjectRequest = new CopyObjectRequest();
        copyObjectRequest.setSourceBucketName(beforeBucket + "-bucket-standard-test");
        copyObjectRequest.setDestinationBucketName(beforeBucket + "-bucket-standard-test");
        copyObjectRequest.setSourceObjectKey("test_object_standard_copy_to_deep_archive");
        copyObjectRequest.setDestinationObjectKey("test_object_standard_copy_to_deep_archive_result");
        ObjectMetadata newObjectMetadata = new ObjectMetadata();
        newObjectMetadata.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        copyObjectRequest.setNewObjectMetadata(newObjectMetadata);
        CopyObjectResult result = obsClient.copyObject(copyObjectRequest);
        assertEquals(200, result.getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test",
                "test_object_standard_copy_to_deep_archive");
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test",
                "test_object_standard_copy_to_deep_archive_result");

        copyObjectRequest.setSourceObjectKey("test_object_warm_copy_to_deep_archive");
        copyObjectRequest.setDestinationObjectKey("test_object_warm_copy_to_deep_archive_result");
        result = obsClient.copyObject(copyObjectRequest);
        assertEquals(200, result.getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test",
                "test_object_warm_copy_to_deep_archive");
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test",
                "test_object_warm_copy_to_deep_archive_result");

        copyObjectRequest.setSourceObjectKey("test_object_cold_copy_to_deep_archive");
        copyObjectRequest.setDestinationObjectKey("test_object_cold_copy_to_deep_archive_result");
        try {
            result = obsClient.copyObject(copyObjectRequest);
        } catch (ObsException ex) {
            assertEquals("InvalidObjectState", ex.getErrorCode());
        }
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test",
                "test_object_cold_copy_to_deep_archive");
    }

    @Test
    public void test_restore_object_from_deep_archive() throws UnsupportedEncodingException {
        PutObjectRequest request = new PutObjectRequest();
        request.setBucketName(beforeBucket + "-bucket-standard-test");
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        request.setMetadata(metadata);
        request.setObjectKey("test_restore_object_from_deep_archive");
        request.setInput(new ByteArrayInputStream("Hello OBS".getBytes("UTF-8")));
        obsClient.putObject(request);
        RestoreObjectRequest restoreObjectRequest = new RestoreObjectRequest();
        restoreObjectRequest.setObjectKey("test_restore_object_from_deep_archive");
        restoreObjectRequest.setBucketName(beforeBucket + "-bucket-standard-test");
        restoreObjectRequest.setDays(1);
        RestoreObjectRequest.RestoreObjectStatus status = obsClient.restoreObject(restoreObjectRequest);
        assertEquals(202, status.getStatusCode());
        obsClient.deleteObject(beforeBucket + "-bucket-standard-test",
                "test_restore_object_from_deep_archive");
    }
}
