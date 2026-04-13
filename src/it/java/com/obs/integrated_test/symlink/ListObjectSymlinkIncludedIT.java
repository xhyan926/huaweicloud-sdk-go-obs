/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test.symlink;


import static com.obs.test.TestTools.genTestFile;
import static org.junit.Assert.assertEquals;

import com.obs.services.ObsClient;
import com.obs.services.model.*;
import com.obs.services.model.symlink.PutSymlinkRequest;
import com.obs.test.TestTools;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.util.Locale;

public class ListObjectSymlinkIncludedIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public TestName testName = new TestName();

    public void cleanUp(ObsClient obsClient,String bucketName){
        TestTools.deleteObjects(obsClient, bucketName);
    }

    @Test
    public void tc_obs_symlink_list_symlink_001() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String normalObjectKey = bucketName + "-normal";
        String appendObjectKey = bucketName + "-append";
        String uploadFileObjectKey = bucketName + "-uploadFile";
        File testFile = genTestFile(temporaryFolder, normalObjectKey, 1024 * 1024);
        ObjectMetadata objectSymlinkMetadata1 = new ObjectMetadata();
        cleanUp(obsClient,bucketName);
        {
//            1、上传普通对象a，再上传append对象b，再上传软链接对象c链接到对象a，最后断点续传上传对象d
//            上传对象均成功，接口响应200；
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, normalObjectKey, testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            AppendObjectRequest appendObjectRequest = new AppendObjectRequest(bucketName);
            appendObjectRequest.setObjectKey(appendObjectKey);
            appendObjectRequest.setFile(testFile);
            AppendObjectResult appendObjectResult = obsClient.appendObject(appendObjectRequest);
            Assert.assertEquals(200, appendObjectResult.getStatusCode());
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                    bucketName, normalObjectKey + ".symlink", normalObjectKey,
                    AccessControlList.REST_CANNED_PUBLIC_READ_WRITE, objectSymlinkMetadata1);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
            UploadFileRequest uploadFileRequest = new UploadFileRequest(bucketName,uploadFileObjectKey,testFile.getPath());
            CompleteMultipartUploadResult uploadFileResult = obsClient.uploadFile(uploadFileRequest);
            Assert.assertEquals(200,uploadFileResult.getStatusCode());
        }

        {
//            2、设置maxKey为3，列举桶内对象，再通过nextmarker分页列举，查看所有列举出的对象类型
//            列举测试桶内对象成功，接口响应200，响应参数满足以下断言:
//            1）普通a、多段对象d的ObjectType字段为NORMAL；
//            2）追加写对象b的ObjectType字段为APPENDABLE；
//            3）软连接对象c的ObjectType字段为SYMLINK；
            ListObjectsRequest listObjectsRequest = new ListObjectsRequest(bucketName);
            listObjectsRequest.setMaxKeys(3);
            ObjectListing objectListing ;
            objectListing = obsClient.listObjects(listObjectsRequest);
            Assert.assertEquals(200,objectListing.getStatusCode());
            for(ObsObject object:objectListing.getObjects()){
                String objectKey = object.getObjectKey();
                if(objectKey.equals(normalObjectKey)){
                    assertEquals(ObjectTypeEnum.NORMAL,object.getObjectType());
                }
                if(objectKey.equals(uploadFileObjectKey)){
                    assertEquals(ObjectTypeEnum.NORMAL,object.getObjectType());
                }
                if(objectKey.equals(appendObjectKey)){
                    assertEquals(ObjectTypeEnum.APPENDABLE,object.getObjectType());
                }
                if(objectKey.equals(normalObjectKey + ".symlink")){
                    assertEquals(ObjectTypeEnum.SYMLINK,object.getObjectType());
                }
            }
            listObjectsRequest.setMarker(objectListing.getNextMarker());
            objectListing = obsClient.listObjects(listObjectsRequest);
            Assert.assertEquals(200,objectListing.getStatusCode());
            for(ObsObject object:objectListing.getObjects()){
                String objectKey = object.getObjectKey();
                if(objectKey.equals(normalObjectKey)){
                    assertEquals(ObjectTypeEnum.NORMAL,object.getObjectType());
                }
                if(objectKey.equals(uploadFileObjectKey)){
                    assertEquals(ObjectTypeEnum.NORMAL,object.getObjectType());
                }
                if(objectKey.equals(appendObjectKey)){
                    assertEquals(ObjectTypeEnum.APPENDABLE,object.getObjectType());
                }
                if(objectKey.equals(normalObjectKey + ".symlink")){
                    assertEquals(ObjectTypeEnum.SYMLINK,object.getObjectType());
                }
            }
        }
        {
//            3、删除软链接对象c
//            删除成功，返回204
            DeleteObjectRequest deleteObjectRequest = new DeleteObjectRequest(bucketName,normalObjectKey + ".symlink");
            DeleteObjectResult deleteObjectResult = obsClient.deleteObject(deleteObjectRequest);
            assertEquals(204,deleteObjectResult.getStatusCode());
        }
        {
//            4、列举桶内对象，查看对象类型
//            4、列举测试桶内对象成功，接口响应200，响应参数中不存在Type为SYMLINK类型的对象
            ListObjectsRequest listObjectsRequest = new ListObjectsRequest(bucketName);
            ObjectListing objectListing = obsClient.listObjects(listObjectsRequest);
            Assert.assertEquals(200,objectListing.getStatusCode());
            for(ObsObject object:objectListing.getObjects()) {
                Assert.assertNotEquals(ObjectTypeEnum.SYMLINK,object.getObjectType());
            }
        }
    }
}
