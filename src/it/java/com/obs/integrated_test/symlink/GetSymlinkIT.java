/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test.symlink;

import static com.obs.test.TestTools.computeMd5Base64;
import static com.obs.test.TestTools.computeMd5Etag;
import static com.obs.test.TestTools.genTestFile;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.utils.ServiceUtils;
import com.obs.services.model.*;
import com.obs.services.model.symlink.GetSymlinkRequest;
import com.obs.services.model.symlink.GetSymlinkResult;
import com.obs.services.model.symlink.PutSymlinkRequest;
import com.obs.test.TestTools;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;
import java.io.File;
import java.io.IOException;
import java.security.NoSuchAlgorithmException;
import java.util.Date;
import java.util.Locale;

public class GetSymlinkIT {

    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public TestName testName = new TestName();

    public static String replaceSpacesWithPlus(String input){
        if (input == null) {
            return null;
        }
        return input.replace(" ", "+");
    }


    public void cleanUp(ObsClient obsClient,String bucketName){
        TestTools.deleteObjects(obsClient, bucketName);
    }

    @Test
    public void tc_obs_symlink_get_001() throws IOException, NoSuchAlgorithmException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String standardObjectKey = bucketName + "-STANDARD";
        File testFile = genTestFile(temporaryFolder, standardObjectKey, 1024 * 1024);
        String testFileMD5Base64 = computeMd5Base64(testFile);
        String testFileMD5Etag = computeMd5Etag(testFile);
        ObjectMetadata objectMetadata1 = new ObjectMetadata();
        ObjectMetadata objectSymlinkMetadata1 = new ObjectMetadata();
        cleanUp(obsClient,bucketName);
        String symlinkEtag = "\"" + ServiceUtils.toHex(ServiceUtils.computeMD5Hash(standardObjectKey.getBytes())) + "\"";

        {
            // 预设对象和软链接对象的元数据
            objectMetadata1.setContentEncoding("ContentEncoding1");
            objectMetadata1.setContentDisposition("ContentDisposition1");
            objectMetadata1.setCacheControl("CacheControl1");
            objectMetadata1.setContentType("ContentType1");
            objectMetadata1.setLastModified(new Date());
            objectMetadata1.setContentMd5(testFileMD5Base64);
            objectMetadata1.setEtag(testFileMD5Etag);
            objectSymlinkMetadata1.setContentEncoding("SymlinkContentEncoding1");
            objectSymlinkMetadata1.setContentDisposition("SymlinkContentDisposition1");
            objectSymlinkMetadata1.setCacheControl("SymlinkCache-Control1");
            objectSymlinkMetadata1.setContentType("binary/octet-stream");
            objectSymlinkMetadata1.setLastModified(new Date());
            objectSymlinkMetadata1.setObjectStorageClass(StorageClassEnum.WARM);
            objectSymlinkMetadata1.setContentMd5("");
            objectSymlinkMetadata1.setEtag("");
        }

        {
//            1、分别上传普通对象、软链接对象，特别注意生成的两个对象的
//            Content-Encoding、Content-Disposition、Cache-Control、Content-Type、Last-Modified、x-obs-object-type、x-obs-storage-class、Content-length、Content-Md5、Etag都设定为不同的值；
//            普通对象创建成功；创建软连接对象成功，接口响应200，ETag头域=md5.update(多段对象名称.tobyte)
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, standardObjectKey, testFile);
            putObjectRequest.setAcl(AccessControlList.REST_CANNED_PUBLIC_READ_WRITE);
            putObjectRequest.setMetadata(objectMetadata1);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                    bucketName, standardObjectKey + ".symlink", standardObjectKey,
                    AccessControlList.REST_CANNED_PUBLIC_READ_WRITE, objectSymlinkMetadata1);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
            Assert.assertEquals(symlinkEtag, headerResponse.getResponseHeaders().get("ETag"));
        }

        {
//            2、使用getsymlink接口获取步骤1软连接对象的元数据
//            getsymlink接口均响应成功，对象元数据（Content-Encoding、Content-Disposition、Cache-Control、Content-Type、Last-Modified、x-obs-storage-class、x-obs-symlink-target、x-obs-object-type、Content-length）为软连接对象后自身的；
            GetSymlinkRequest getSymlinkRequest =
                    new GetSymlinkRequest(bucketName, standardObjectKey + ".symlink");
            GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());
            Assert.assertEquals(objectSymlinkMetadata1.getContentEncoding(),
                    getSymlinkResult.getResponseHeaders().get("Content-Encoding"));
            Assert.assertEquals(objectSymlinkMetadata1.getContentDisposition(),
                    getSymlinkResult.getResponseHeaders().get("Content-Disposition"));
            Assert.assertEquals(objectSymlinkMetadata1.getCacheControl(),
                    getSymlinkResult.getResponseHeaders().get("Cache-Control"));
            Assert.assertEquals(standardObjectKey, getSymlinkResult.getSymlinkTarget());
            Assert.assertEquals("Symlink",
                    getSymlinkResult.getResponseHeaders().get("object-type"));
            Assert.assertEquals("0", getSymlinkResult.getResponseHeaders().get("Content-length"));
        }

        {
//            3、使用headObject接口获取步骤1软连接对象的元数据
//            headObject接口均响应成功，对象元数据（Content-Encoding、Content-Disposition、Cache-Control、Content-Type、Last-Modified、x-obs-object-type）
//            为软连接对象后自身的，对象数据相关（x-obs-storage-class、Content-length、Content-Md5、Etag）是目标对象最新版本的数据；
            ObjectMetadata metadataStandardObjectSymlink =
                    obsClient.getObjectMetadata(bucketName, standardObjectKey + ".symlink");
            Assert.assertEquals(200, metadataStandardObjectSymlink.getStatusCode());
            Assert.assertEquals(objectSymlinkMetadata1.getContentEncoding(),
                    metadataStandardObjectSymlink.getResponseHeaders().get("Content-Encoding"));
            Assert.assertEquals(objectSymlinkMetadata1.getContentDisposition(),
                    metadataStandardObjectSymlink.getResponseHeaders().get("Content-Disposition"));
            Assert.assertEquals(objectSymlinkMetadata1.getCacheControl(),
                    metadataStandardObjectSymlink.getResponseHeaders().get("Cache-Control"));
            Assert.assertEquals(objectSymlinkMetadata1.getContentType(),
                    metadataStandardObjectSymlink.getResponseHeaders().get("Content-Type"));
            Assert.assertEquals("Symlink",
                    metadataStandardObjectSymlink.getResponseHeaders().get("object-type"));
            // 对象数据相关（x-obs-storage-class、Content-length、Content-Md5、Etag）是目标对象的；
            Assert.assertEquals(String.valueOf(testFile.length()),
                    metadataStandardObjectSymlink.getResponseHeaders().get("Content-length"));
            Assert.assertEquals(testFileMD5Base64,
                    replaceSpacesWithPlus(String.valueOf(metadataStandardObjectSymlink.getResponseHeaders().get("Content-Md5"))));
            Assert.assertEquals(testFileMD5Etag,
                    metadataStandardObjectSymlink.getResponseHeaders().get("ETag"));
        }
        {
            // 4、使用getObject接口获取步骤1软连接对象的元数据和数据
            // 预期：getObject接口均响应成功，对象元数据（Content-Encoding、Content-Disposition、Cache-Control、Content-Type、
            // Last-Modified、x-obs-object-type）为软连接对象后自身的，对象数据相关（x-obs-storage-class、Content-length、Content-Md5、Etag）是目标对象最新版本的数据。
            ObsObject obsObject = obsClient.getObject(bucketName, standardObjectKey + ".symlink");
            obsObject.getObjectContent().close();
            ObjectMetadata objectMetadataStandardObject = obsObject.getMetadata();
            Assert.assertEquals(200, objectMetadataStandardObject.getStatusCode());
            Assert.assertEquals(objectSymlinkMetadata1.getContentEncoding(),
                    objectMetadataStandardObject.getResponseHeaders().get("Content-Encoding"));
            Assert.assertEquals(objectSymlinkMetadata1.getCacheControl(),
                    objectMetadataStandardObject.getResponseHeaders().get("Cache-Control"));
            Assert.assertEquals(objectSymlinkMetadata1.getContentType(),
                    objectMetadataStandardObject.getResponseHeaders().get("Content-Type"));
            Assert.assertEquals("Symlink",
                    objectMetadataStandardObject.getResponseHeaders().get("object-type"));
            // 对象数据相关（、Content-length、Content-Md5、Etag）是目标对象的；
            Assert.assertEquals(String.valueOf(testFile.length()),
                    objectMetadataStandardObject.getResponseHeaders().get("Content-length"));
            Assert.assertEquals(testFileMD5Base64,
                    replaceSpacesWithPlus(String.valueOf(objectMetadataStandardObject.getResponseHeaders().get("Content-Md5"))));
            Assert.assertEquals(testFileMD5Etag,
                    objectMetadataStandardObject.getResponseHeaders().get("ETag"));
        }

        {
//            5、删除软连接对象
//            删除软链接对象均成功，响应204。
            HeaderResponse response1 =
                    obsClient.deleteObject(bucketName, standardObjectKey + ".symlink");
            Assert.assertEquals(204, response1.getStatusCode());
        }

    }

    @Test
    public void tc_obs_symlink_get_002() throws IOException, NoSuchAlgorithmException {
        // 软连接指向一个多版本对象他的名字和一个正常对象相同
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String standardObjectKey = bucketName + "-STANDARD";
        String versionObjectKey = bucketName + "-MultipleVersion";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, standardObjectKey, 1024 * 1024);
        String testFileMD5Base64 = computeMd5Base64(testFile);
        String testFileMD5Etag = computeMd5Etag(testFile);
        ObjectMetadata objectMetadata1 = new ObjectMetadata();
        objectMetadata1.setContentMd5(testFileMD5Base64);
        ObjectMetadata objectSymlinkMetadata1 = new ObjectMetadata();
        cleanUp(obsClient,bucketName);
        {
            // 预设对象和软链接对象的元数据
            objectMetadata1.setContentEncoding("ContentEncoding1");
            objectMetadata1.setContentDisposition("ContentDisposition1");
            objectMetadata1.setCacheControl("CacheControl1");
            objectMetadata1.setContentType("ContentType1");
            objectMetadata1.setLastModified(new Date());
            objectMetadata1.setContentMd5(testFileMD5Base64);
            objectMetadata1.setEtag(testFileMD5Etag);
            objectSymlinkMetadata1.setContentEncoding("SymlinkContentEncoding1");
            objectSymlinkMetadata1.setContentDisposition("SymlinkContentDisposition1");
            objectSymlinkMetadata1.setCacheControl("SymlinkCache-Control1");
            objectSymlinkMetadata1.setContentType("binary/octet-stream");
            objectSymlinkMetadata1.setLastModified(new Date());
            objectSymlinkMetadata1.setObjectStorageClass(StorageClassEnum.WARM);
            objectSymlinkMetadata1.setContentMd5("");
        }

        {
//            1、在测试桶上传一个多版本对象
//            上传多版本对象成功，响应200
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, versionObjectKey, testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());
            PutObjectRequest putObjectRequestAgain = new PutObjectRequest(bucketName, versionObjectKey, testFile);
            PutObjectResult putObjectResultAgain = obsClient.putObject(putObjectRequestAgain);
            Assert.assertEquals(200, putObjectResultAgain.getStatusCode());
        }
        {
            //2、上传一个正常对象，携带x-obs-symlink-target指向该对象，软连接对象名称与步骤一多版本对象名称相同，创建软连接对象
//            上传正常对象成功，创建软连接对象成功，接口都响应200，ETag头域=md5.update(目标对象名称.tobyte)
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, standardObjectKey, testFile);
            putObjectRequest.setMetadata(objectMetadata1);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                    bucketName, versionObjectKey, standardObjectKey,
                    AccessControlList.REST_CANNED_PUBLIC_READ_WRITE,objectSymlinkMetadata1);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
//            3、使用getSymlink接口(不指定versionId)，获取软连接对象的元数据
            GetSymlinkRequest getSymlinkRequest =
                    new GetSymlinkRequest(bucketName, versionObjectKey);
            GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());
            Assert.assertEquals(objectSymlinkMetadata1.getContentEncoding(),
                    getSymlinkResult.getResponseHeaders().get("Content-Encoding"));
            Assert.assertEquals(objectSymlinkMetadata1.getContentDisposition(),
                    getSymlinkResult.getResponseHeaders().get("Content-Disposition"));
            Assert.assertEquals(objectSymlinkMetadata1.getCacheControl(),
                    getSymlinkResult.getResponseHeaders().get("Cache-Control"));
            Assert.assertEquals(standardObjectKey, getSymlinkResult.getSymlinkTarget());
            Assert.assertEquals("Symlink",
                    getSymlinkResult.getResponseHeaders().get("object-type"));
            Assert.assertEquals("0", getSymlinkResult.getResponseHeaders().get("Content-length"));
        }

        {
//            4、使用headObject接口，获取软连接对象的元数据
//            headObject接口均响应成功，接口响应200，对象元数据（Content-Encoding、Content-Disposition、Cache-Control、Content-Type、Last-Modified、x-obs-object-type）为软连接对象后自身的，对象数据相关（x-obs-storage-class、Content-length、Content-Md5、Etag）是目标对象的；
            ObjectMetadata  getSymlinkResult =
                    obsClient.getObjectMetadata(bucketName, versionObjectKey);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());
            Assert.assertEquals(objectSymlinkMetadata1.getContentEncoding(),
                    getSymlinkResult.getResponseHeaders().get("Content-Encoding"));
            Assert.assertEquals(objectSymlinkMetadata1.getContentDisposition(),
                    getSymlinkResult.getResponseHeaders().get("Content-Disposition"));
            Assert.assertEquals(objectSymlinkMetadata1.getCacheControl(),
                    getSymlinkResult.getResponseHeaders().get("Cache-Control"));
            Assert.assertEquals("Symlink",
                    getSymlinkResult.getResponseHeaders().get("object-type"));
            // 对象数据相关（x-obs-storage-class、Content-length、Content-Md5、Etag）是目标对象的；
            Assert.assertEquals(String.valueOf(testFile.length()),
                    getSymlinkResult.getResponseHeaders().get("Content-length"));
            Assert.assertEquals(testFileMD5Base64,
                    replaceSpacesWithPlus(String.valueOf(getSymlinkResult.getResponseHeaders().get("Content-Md5"))));
            Assert.assertEquals(testFileMD5Etag,
                    getSymlinkResult.getResponseHeaders().get("ETag"));
        }
        {
//            5、使用getObject接口，获取软连接对象的元数据和数据
//            getObject接口均响应成功，接口响应200，对象元数据（Content-Encoding、Content-Disposition、Cache-Control、Content-Type、Last-Modified、x-obs-object-type）为软连接对象后自身的，对象数据相关（x-obs-storage-class、Content-length、Content-Md5、Etag）是目标对象的。
            ObsObject obsObject = obsClient.getObject(bucketName, versionObjectKey);
            obsObject.getObjectContent().close();
            Assert.assertEquals(200, obsObject.getMetadata().getStatusCode());
        }

        {
//            6、调用deleteObject接口，删除最新版本的软连接对象；
//            删除对象成功，响应204；
            HeaderResponse response =
                    obsClient.deleteObject(bucketName, versionObjectKey);
            Assert.assertEquals(204, response.getStatusCode());
        }

        {
//            7、使用getSymlink接口(不指定versionId)，获取软连接对象的元数据
//            7、getsSymlink接口均响应404，错误信息“The specified key does not exist”，在响应header中返回x-obs-delete-marker = true以及版本ID : x-obs-version-id。删除标记没有关联数据，因此也没有软链接指向的TargetObject，即没有x-obs-symlink-target头域；
            try {
                obsClient.getSymlink(
                        new GetSymlinkRequest(bucketName, versionObjectKey));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
        }
        {
//            8、使用headObject接口(不指定versionId)，获取软连接对象的元数据
//            8、headObject接口均响应404，错误信息“The specified key does not exist”；
            try {
                obsClient.getObjectMetadata(bucketName, versionObjectKey);
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
        }
        {
//            9、使用getObject接口(不指定versionId)，获取软连接对象的元数据和数据
//            getObject接口均响应404，错误信息“The specified key does not exist”；
            try {
                obsClient.getObject(bucketName, versionObjectKey);
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
        }
    }
    @Test
    public void tc_obs_symlink_get_003(){
        // 通过不存在软连接对象名称，获取软连接对象信息，响应404
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String standardObjectKey = bucketName + "-STANDARD";
        cleanUp(obsClient,bucketName);
        {
//            1、使用getsymlink接口获取不存在的软连接对象的元数据
//            getsymlink接口均响应404，错误信息“The specified key does not exist”，
            try {
                obsClient.getSymlink(
                        new GetSymlinkRequest(bucketName, standardObjectKey));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
        }
    }
    @Test
    public void tc_obs_symlink_get_004() throws IOException{
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String standardObjectKey = bucketName + "-STANDARD";
        File testFile = genTestFile(temporaryFolder, standardObjectKey, 1024 * 1024);
        ObjectMetadata objectMetadata1 = new ObjectMetadata();
        cleanUp(obsClient,bucketName);
        {
//            1、上传普通对象
//            普通对象创建成功；接口响应200
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, standardObjectKey, testFile);
            putObjectRequest.setAcl(AccessControlList.REST_CANNED_PUBLIC_READ_WRITE);
            putObjectRequest.setMetadata(objectMetadata1);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());
        }

        {
//            2、通过getsymlink接口查询该对象信息；
//            getsymlink接口响应400,错误信息“The specified key does not exist”.
            try {
                obsClient.getSymlink(
                        new GetSymlinkRequest(bucketName, standardObjectKey));
                fail();
            } catch (ObsException e) {
                assertEquals(400, e.getResponseCode());
            }
        }
    }

    @Test
    public void tc_obs_symlink_get_005() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String standardObjectKey = bucketName + "-STANDARD";
        File testFile = genTestFile(temporaryFolder, standardObjectKey, 1024 * 1024);
        ObjectMetadata objectMetadata1 = new ObjectMetadata();
        ObjectMetadata objectSymlinkMetadata1=new ObjectMetadata();
        cleanUp(obsClient,bucketName);
        {
//            1、上传一个正常对象a，创建软连接对象b链接到a
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, standardObjectKey, testFile);
            putObjectRequest.setAcl(AccessControlList.REST_CANNED_PUBLIC_READ_WRITE);
            putObjectRequest.setMetadata(objectMetadata1);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                    bucketName, standardObjectKey + ".symlink", standardObjectKey,
                    AccessControlList.REST_CANNED_PUBLIC_READ_WRITE, objectSymlinkMetadata1);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }

        {
//            2、判断b对象是否为软连接对象
//            确实是symlink类型对象
            ObjectMetadata metadataStandardObjectSymlink =
                    obsClient.getObjectMetadata(bucketName, standardObjectKey + ".symlink");
            Assert.assertEquals("Symlink",
                    metadataStandardObjectSymlink.getResponseHeaders().get("object-type"));

        }

        {
//            3、使用headobject 获取对象a和软连接对象b的元数据，对比其中的etag是否相同
//            调用headobject查询a、b都成功，响应状态码都是200。a、b对象的etag对比相同
            ObjectMetadata metadataStandardObjectSymlink =
                    obsClient.getObjectMetadata(bucketName, standardObjectKey + ".symlink");
            ObjectMetadata metadataStandardObject =
                    obsClient.getObjectMetadata(bucketName, standardObjectKey );
            Assert.assertEquals(metadataStandardObject.getEtag(),metadataStandardObjectSymlink.getEtag());
        }
    }
}
