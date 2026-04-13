/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test.symlink;

import static com.obs.test.TestTools.genTestFile;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.utils.ServiceUtils;
import com.obs.services.model.BucketVersioningConfiguration;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.ServerAlgorithm;
import com.obs.services.model.SseCHeader;
import com.obs.services.model.StorageClassEnum;
import com.obs.services.model.VersioningStatusEnum;
import com.obs.services.model.symlink.GetSymlinkRequest;
import com.obs.services.model.symlink.GetSymlinkResult;
import com.obs.services.model.symlink.PutSymlinkRequest;
import com.obs.test.TestTools;

import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.security.SecureRandom;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;
import java.util.Random;

public class PutSymlinkIT {

    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public TestName testName = new TestName();

    public void cleanUp(ObsClient obsClient,String bucketName){
        TestTools.deleteObjects(obsClient, bucketName);
    }

    private static String sse_c_base64;

    @BeforeClass
    public static void prepareForTest() {
        // 设置SSE-C方式下使用的密钥，用于加解密对象，该值是密钥进行base64encode后的值
        byte[] sse_c = new byte[32];
        new SecureRandom().nextBytes(sse_c);
        sse_c_base64 = ServiceUtils.toBase64(sse_c);
    }

    @Test
    public void tc_obs_symlink_create_001() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String normalObjectKey = bucketName + "-normal";
        String encObjectKey = bucketName + "-sse_obs";
        String  symlinkObjectKey = bucketName + "-symlink";
        String versionObjectKey = bucketName + "-version";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, normalObjectKey, 1024 * 1024);
        cleanUp(obsClient,bucketName);
        {
            // 1.上传一个普通对象，对它创建软连接对象
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, normalObjectKey, testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, normalObjectKey + ".symlink", normalObjectKey);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 2.上传一个多版本对象，对它创建软连接对象
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, versionObjectKey, testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, versionObjectKey + ".version.symlink", versionObjectKey);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 3.上传一个加密对象，对它创建软连接对象
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, encObjectKey, testFile);
            //设置自定义头域，设置使用sse-obs加密
            HashMap<String, String> userHeaders = new HashMap<>();
            userHeaders.put("x-obs-server-side-encryption","AES256");
            putObjectRequest.setUserHeaders(userHeaders);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, encObjectKey + ".symlink", encObjectKey);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 4.创建一个软连接对象，链接的目标指向自己
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, symlinkObjectKey, symlinkObjectKey);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 5.针对上面任一软连接对象，创建软连接对象
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, symlinkObjectKey + ".symlink", encObjectKey + ".symlink");
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 6.对上面一个已经创建过软连接的对象创建软连接
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, normalObjectKey + ".symlink2", normalObjectKey);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 7、使用getsymlink接口获取6个软连接对象的元数据
            GetSymlinkRequest getSymlinkRequest =
                new GetSymlinkRequest(bucketName, normalObjectKey + ".symlink");
            GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());

            getSymlinkRequest =
                new GetSymlinkRequest(bucketName, versionObjectKey + ".version.symlink");
            getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());

            getSymlinkRequest =
                new GetSymlinkRequest(bucketName, encObjectKey + ".symlink");
            getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());

            getSymlinkRequest =
                new GetSymlinkRequest(bucketName, symlinkObjectKey);
            getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());

            getSymlinkRequest =
                new GetSymlinkRequest(bucketName, symlinkObjectKey + ".symlink");
            getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());

            getSymlinkRequest =
                new GetSymlinkRequest(bucketName, normalObjectKey + ".symlink2");
            getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());

        }
        {
            // 8、使用headObject接口获取6个软连接对象的元数据
            ObjectMetadata metadataSymlink =
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".symlink");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            metadataSymlink =
                obsClient.getObjectMetadata(bucketName, versionObjectKey + ".version.symlink");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            metadataSymlink =
                obsClient.getObjectMetadata(bucketName, encObjectKey + ".symlink");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            try {
                obsClient.getObjectMetadata(bucketName, symlinkObjectKey);
                fail();
            } catch (ObsException e) {
                assertEquals(400, e.getResponseCode());
            }

            try {
                obsClient.getObjectMetadata(bucketName,symlinkObjectKey);
                fail();
            } catch (ObsException e) {
                assertEquals(400, e.getResponseCode());
            }

            metadataSymlink =
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".symlink2");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());
        }
        {
            // 9、使用getObject接口获取6个软连接对象的元数据和数据
            ObsObject obsObject = obsClient.getObject(bucketName, normalObjectKey + ".symlink");
            obsObject.getObjectContent().close();
            ObjectMetadata metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            obsObject = obsClient.getObject(bucketName, versionObjectKey + ".version.symlink");
            obsObject.getObjectContent().close();
            metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            obsObject = obsClient.getObject(bucketName, encObjectKey + ".symlink");
            obsObject.getObjectContent().close();
            metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            try {
                obsClient.getObject(bucketName, symlinkObjectKey);
                fail();
            } catch (ObsException e) {
                assertEquals(400, e.getResponseCode());
            }

            try {
                obsClient.getObject(bucketName, symlinkObjectKey + ".symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(400, e.getResponseCode());
            }

            obsObject = obsClient.getObject(bucketName, normalObjectKey + ".symlink2");
            obsObject.getObjectContent().close();
            metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());
        }
        {
            // 10、删除6个软连接对象
            // 预期：响应204
            HeaderResponse response1 =
                obsClient.deleteObject(bucketName, normalObjectKey + ".symlink");
            Assert.assertEquals(204, response1.getStatusCode());
            HeaderResponse response2 =
                obsClient.deleteObject(bucketName, versionObjectKey + ".version.symlink");
            Assert.assertEquals(204, response2.getStatusCode());
            HeaderResponse response3 =
                obsClient.deleteObject(bucketName, encObjectKey + ".symlink");
            Assert.assertEquals(204, response3.getStatusCode());
            response3 =
                obsClient.deleteObject(bucketName, symlinkObjectKey);
            Assert.assertEquals(204, response3.getStatusCode());
            response3 =
                obsClient.deleteObject(bucketName, symlinkObjectKey + ".symlink");
            Assert.assertEquals(204, response3.getStatusCode());
            response3 =
                obsClient.deleteObject(bucketName, normalObjectKey + ".symlink2");
            Assert.assertEquals(204, response3.getStatusCode());
        }
        {
            // 11、使用getsymlink接口获取6个软连接对象的元数据
            // getsymlink接口均响应404，错误信息“The specified key does not exist”
            try {
                obsClient.getSymlink(
                    new GetSymlinkRequest(bucketName, normalObjectKey + ".symlink"));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getSymlink(
                    new GetSymlinkRequest(bucketName, versionObjectKey + ".version.symlink"));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getSymlink(
                    new GetSymlinkRequest(bucketName, encObjectKey + ".symlink"));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getSymlink(
                    new GetSymlinkRequest(bucketName, symlinkObjectKey));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getSymlink(
                    new GetSymlinkRequest(bucketName, symlinkObjectKey + ".symlink"));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getSymlink(
                    new GetSymlinkRequest(bucketName, normalObjectKey + ".symlink2"));
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
        }
        {
            // 12、使用headObject接口获取6个软连接对象的元数据
            try {
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObjectMetadata(bucketName, versionObjectKey + ".version.symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObjectMetadata(bucketName, encObjectKey + ".symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObjectMetadata(bucketName, symlinkObjectKey);
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObjectMetadata(bucketName, symlinkObjectKey + ".symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".symlink2");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
        }
        {
            // 13、使用getObject接口获取6个软连接对象的元数据和数据
            try {
                obsClient.getObject(bucketName, normalObjectKey + ".symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObject(bucketName, versionObjectKey + ".version.symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObject(bucketName, encObjectKey + ".symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObject(bucketName, symlinkObjectKey);
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObject(bucketName, symlinkObjectKey + ".symlink");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
            try {
                obsClient.getObject(bucketName, normalObjectKey + ".symlink2");
                fail();
            } catch (ObsException e) {
                assertEquals(404, e.getResponseCode());
            }
        }
    }

    @Test
    public void tc_obs_symlink_create_002() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String normalObjectKey = bucketName + "-normal";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, normalObjectKey, 1024 * 1024);
        // 创建一个随机的8k长度的字符串用于测试
        int length = 8192 - "x-obs-meta-test_user_metadata".length(); // 8K 字节
        StringBuilder sb = new StringBuilder(length);
        Random random = new Random();
        String chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        for (int i = 0; i < length; i++) {
            int index = random.nextInt(chars.length());
            sb.append(chars.charAt(index));
        }
        String test_user_metadata = sb.toString();
        cleanUp(obsClient,bucketName);
        {
            // 1.上传一个普通对象，对它创建软连接对象，并携带1个随机的x-obs-meta-为前缀头域长度刚好8kb
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, normalObjectKey, testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            ObjectMetadata objectMetadata = new ObjectMetadata();
            Map<String, Object> userMetadata = new HashMap<>();
            userMetadata.put("test_user_metadata", test_user_metadata);
            objectMetadata.setUserMetadata(userMetadata);
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, normalObjectKey + ".symlink", normalObjectKey, null, objectMetadata);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 2、使用getsymlink接口获取软连接对象的元数据
            GetSymlinkRequest getSymlinkRequest =
                new GetSymlinkRequest(bucketName, normalObjectKey + ".symlink");
            GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());
            Assert.assertEquals(test_user_metadata,
                getSymlinkResult.getResponseHeaders().get("test_user_metadata"));
        }
        {
            // 3、使用headObject接口获取该软连接对象的元数据
            ObjectMetadata metadataSymlink =
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".symlink");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());
        }
        {
            // 4、使用getObject接口获取该软连接对象的元数据和数据
            ObsObject obsObject = obsClient.getObject(bucketName, normalObjectKey + ".symlink");
            obsObject.getObjectContent().close();
            ObjectMetadata metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());
        }
        {
            // 5、删除软连接对象
            // 预期：响应204
            HeaderResponse response1 =
                obsClient.deleteObject(bucketName, normalObjectKey + ".symlink");
            Assert.assertEquals(204, response1.getStatusCode());
        }
    }

    @Test
    public void tc_obs_symlink_create_003() throws IOException{
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String normalObjectKey = bucketName + "-normal";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, normalObjectKey, 1024 * 1024);
        ObjectMetadata objectMetadataStandard = new ObjectMetadata();
        ObjectMetadata objectMetadataStandard_AI = new ObjectMetadata();
        ObjectMetadata objectMetadataCold = new ObjectMetadata();
        ObjectMetadata objectMetadataDeepArchive = new ObjectMetadata();
        cleanUp(obsClient,bucketName);
        {
            // 预设对象和软链接对象的元数据
            objectMetadataStandard.setObjectStorageClass(StorageClassEnum.STANDARD);
            objectMetadataStandard_AI.setObjectStorageClass(StorageClassEnum.WARM);
            objectMetadataCold.setObjectStorageClass(StorageClassEnum.COLD);
            objectMetadataDeepArchive.setObjectStorageClass(StorageClassEnum.DEEP_ARCHIVE);
        }
        {
            // 1.上传一个正常对象，携带x-obs-symlink-target指向该对象，x-obs-storage-class为STANDARD，创建软连接对象
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, normalObjectKey, testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            Assert.assertEquals(200, putObjectResult.getStatusCode());

            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, normalObjectKey + ".standard.symlink", normalObjectKey,
                null, objectMetadataStandard);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 2.上传一个正常对象，携带x-obs-symlink-target指向该对象，x-obs-storage-class为STANDARD_AI，创建软连接对象
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, normalObjectKey + ".standard_ai.symlink", normalObjectKey,
                null, objectMetadataStandard_AI);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 3、上传一个正常对象，携带x-obs-symlink-target指向该对象，x-obs-storage-class为COLD，创建软连接对象
            PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
                bucketName, normalObjectKey + ".cold.symlink", normalObjectKey,
                null, objectMetadataCold);
            HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
        }
        {
            // 5、使用getsymlink接口获取4个软连接对象的元数据
            GetSymlinkRequest getSymlinkRequest =
                new GetSymlinkRequest(bucketName, normalObjectKey + ".standard.symlink");
            GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());
            Assert.assertNull(getSymlinkResult.getResponseHeaders().get("storage_class"));

            getSymlinkRequest =
                new GetSymlinkRequest(bucketName, normalObjectKey + ".standard_ai.symlink");
            getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());
            assert(String.valueOf(getSymlinkResult.getResponseHeaders().get("storage-class")).equals("WARM")||String.valueOf(getSymlinkResult.getResponseHeaders().get("storage-class")).equals("STANDARD_IA"));
            getSymlinkRequest =
                new GetSymlinkRequest(bucketName, normalObjectKey + ".cold.symlink");
            getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
            Assert.assertEquals(200, getSymlinkResult.getStatusCode());
            assert(String.valueOf(getSymlinkResult.getResponseHeaders().get("storage-class")).equals("COLD")||String.valueOf(getSymlinkResult.getResponseHeaders().get("storage-class")).equals("GLACIER"));
        }
        {
            // 6、使用headObject接口获取4个软连接对象的元数据
            ObjectMetadata metadataSymlink =
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".standard.symlink");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            metadataSymlink =
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".standard_ai.symlink");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            metadataSymlink =
                obsClient.getObjectMetadata(bucketName, normalObjectKey + ".cold.symlink");
            Assert.assertEquals(200, metadataSymlink.getStatusCode());
        }
        {
            // 7、使用getObject接口获取4个软连接对象的元数据和数据
            ObsObject obsObject = obsClient.getObject(bucketName, normalObjectKey + ".standard.symlink");
            obsObject.getObjectContent().close();
            ObjectMetadata metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            obsObject = obsClient.getObject(bucketName, normalObjectKey + ".standard_ai.symlink");
            obsObject.getObjectContent().close();
            metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());

            obsObject = obsClient.getObject(bucketName, normalObjectKey + ".cold.symlink");
            obsObject.getObjectContent().close();
            metadataSymlink = obsObject.getMetadata();
            Assert.assertEquals(200, metadataSymlink.getStatusCode());
        }
        {
            // 8、删除4个软连接对象
            // 预期：响应204
            HeaderResponse response1 =
                obsClient.deleteObject(bucketName, normalObjectKey + ".standard.symlink");
            Assert.assertEquals(204, response1.getStatusCode());
            HeaderResponse response2 =
                obsClient.deleteObject(bucketName, normalObjectKey + ".standard_ai.symlink");
            Assert.assertEquals(204, response2.getStatusCode());
            HeaderResponse response3 =
                obsClient.deleteObject(bucketName, normalObjectKey + ".cold.symlink");
            Assert.assertEquals(204, response3.getStatusCode());
            response3 =
                obsClient.deleteObject(bucketName, normalObjectKey + ".DeepArchive.symlink");
            Assert.assertEquals(204, response3.getStatusCode());
        }
    }

    @Test
    public void tc_obs_symlink_create_004() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        obsClient.setBucketVersioning(bucketName, new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));
        String ObjectKey = bucketName + "-normal";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, ObjectKey, 1024 * 1024);
        // 1、上传一个正常对象，携带x-obs-symlink-target指向该对象，创建软连接对象；
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, ObjectKey, testFile);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, ObjectKey + ".symlink", ObjectKey);
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());
        // 2、使用getsymlink接口获取这个软连接对象的元数据
        GetSymlinkRequest getSymlinkRequest =
            new GetSymlinkRequest(bucketName, ObjectKey + ".symlink");
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
        Assert.assertEquals(200, getSymlinkResult.getStatusCode());

        // 3、使用headObject接口获取这个软连接对象的元数据
        ObjectMetadata metadataSymlink =
            obsClient.getObjectMetadata(bucketName, ObjectKey + ".symlink");
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 4、使用getObject接口获取这个软连接对象的元数据和数据
        ObsObject obsObject = obsClient.getObject(bucketName, ObjectKey + ".symlink");
        obsObject.getObjectContent().close();
        metadataSymlink = obsObject.getMetadata();
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 5、上传一个与步骤一同名的对象；
        putObjectRequest.setFile(testFile);
        putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 6、使用getsymlink接口获取这个软连接对象的元数据
        getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
        Assert.assertEquals(200, getSymlinkResult.getStatusCode());
        // 7、使用headObject接口获取这个软连接对象的元数据
        metadataSymlink =
            obsClient.getObjectMetadata(bucketName, ObjectKey + ".symlink");
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 8、使用getObject接口获取这个软连接对象的元数据和数据
        obsObject = obsClient.getObject(bucketName, ObjectKey + ".symlink");
        obsObject.getObjectContent().close();
        metadataSymlink = obsObject.getMetadata();
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 9、删除软连接对象
        HeaderResponse response1 =
            obsClient.deleteObject(bucketName, ObjectKey + ".symlink");
        Assert.assertEquals(204, response1.getStatusCode());
    }

    @Test
    public void tc_obs_symlink_create_005() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String ObjectKey = bucketName + "-normal";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, ObjectKey, 1024 * 1024);
        // 1、上传普通对象D
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, ObjectKey, testFile);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 2、创建软连接对象A，链接到对象D
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, ObjectKey + ".symlink", ObjectKey);
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());
        // 3、上传普通对象名为A
        putObjectRequest = new PutObjectRequest(bucketName, ObjectKey + ".symlink", testFile);
        putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 4、getSymlink接口均响应400，报错信息The object is not symlink
        try {
            GetSymlinkRequest getSymlinkRequest =
                new GetSymlinkRequest(bucketName, ObjectKey + ".symlink");
            obsClient.getSymlink(getSymlinkRequest);
            fail();
        } catch (ObsException e) {
            Assert.assertEquals(400, e.getResponseCode());
        }
    }

    @Test
    public void tc_obs_symlink_create_006() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String ObjectKey = bucketName + "-normal";
        String ObjectKey2 = bucketName + "-normal2";
        cleanUp(obsClient,bucketName);

        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, ObjectKey, 1024 * 1024);
        // 1、上传一个对象A
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, ObjectKey2, testFile);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 2、上传一个对象B，再创建软链接对象A。
        putObjectRequest = new PutObjectRequest(bucketName, ObjectKey, testFile);
        putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, ObjectKey + ".symlink", ObjectKey);
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());

        // 6、使用getsymlink接口获取这个软连接对象的元数据
        GetSymlinkRequest getSymlinkRequest =
            new GetSymlinkRequest(bucketName, ObjectKey + ".symlink");
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
        Assert.assertEquals(200, getSymlinkResult.getStatusCode());
        // 7、使用headObject接口获取这个软连接对象的元数据
        ObjectMetadata  metadataSymlink =
            obsClient.getObjectMetadata(bucketName, ObjectKey + ".symlink");
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 8、使用getObject接口获取这个软连接对象的元数据和数据
        ObsObject obsObject = obsClient.getObject(bucketName, ObjectKey + ".symlink");
        obsObject.getObjectContent().close();
        metadataSymlink = obsObject.getMetadata();
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
    }

    @Test
    public void tc_obs_symlink_create_007() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String ObjectKey = bucketName + "-normal";
        String ObjectKey2 = bucketName + "-normal2";
        String symlink1 = bucketName + "-symlink1";
        String symlink2 = ObjectKey;
        cleanUp(obsClient,bucketName);

        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, ObjectKey, 1024 * 1024);
        // 1、上传一个正常对象，携带x-obs-symlink-target指向该对象，创建软连接对象
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, ObjectKey, testFile);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, symlink1, ObjectKey);
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());
        // 2、再上传一个正常对象，携带x-obs-symlink-target指向该对象，软连接对象名称与步骤一的目标对象对象名称相同，创建软连接对象
        putObjectRequest = new PutObjectRequest(bucketName, ObjectKey2, testFile);
        putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        putSymlinkRequest = new PutSymlinkRequest(
            bucketName, symlink2, ObjectKey2);
        headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());

        // 3、使用getSymlink接口获取步骤一的软连接对象的元数据
        GetSymlinkRequest getSymlinkRequest =
            new GetSymlinkRequest(bucketName, symlink1);
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
        Assert.assertEquals(200, getSymlinkResult.getStatusCode());
        // 4、使用headObject接口获取这个软连接对象的元数据
        try {
            obsClient.getObjectMetadata(bucketName, symlink1);
            fail();
        } catch (ObsException e) {
            Assert.assertEquals(400, e.getResponseCode());
        }
        // 5、使用getObject接口获取这个软连接对象的元数据和数据
        try {
            obsClient.getObject(bucketName, symlink1);
            fail();
        } catch (ObsException e) {
            Assert.assertEquals(400, e.getResponseCode());
        }
    }

    @Test
    public void tc_obs_symlink_create_008() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        cleanUp(obsClient,bucketName);

        assert obsClient != null;
        String encObjectKey = bucketName + "-normal";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, encObjectKey, 1024 * 1024);
        obsClient.setBucketVersioning(bucketName, new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));
        // 1、上传对象A
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, "A", testFile);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 2、上传对象B，创建软链接对象A指向B
        putObjectRequest = new PutObjectRequest(bucketName, "B", testFile);
        putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, "A", "B");
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());
        // 5、使用getsymlink接口获取软连接对象的元数据
        GetSymlinkRequest getSymlinkRequest =
            new GetSymlinkRequest(bucketName, "A");
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
        Assert.assertEquals(200, getSymlinkResult.getStatusCode());
        // 6、使用headObject接口获取软连接对象的元数据
        ObjectMetadata  metadataSymlink =
            obsClient.getObjectMetadata(bucketName, "A");
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 7、使用getObject接口获取软连接对象的元数据和数据
        ObsObject obsObject = obsClient.getObject(bucketName, "A");
        obsObject.getObjectContent().close();
        metadataSymlink = obsObject.getMetadata();
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
    }
    @Test
    public void tc_obs_symlink_create_009() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();
        assert obsClient != null;
        String encObjectKey = bucketName + "-normal";

        cleanUp(obsClient,bucketName);
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, encObjectKey, 1024 * 1024);
        obsClient.setBucketVersioning(bucketName, new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));
        // 1、上传对象A，创建软链接B指向A
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, "A", testFile);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, "A", "B");
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());
        // 2、上传对象B
        putObjectRequest = new PutObjectRequest(bucketName, "B", testFile);
        putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 3、使用getsymlink接口获取软连接对象的元数据
        try {
            GetSymlinkRequest getSymlinkRequest =
                new GetSymlinkRequest(bucketName, "B");
            obsClient.getSymlink(getSymlinkRequest);
        } catch (ObsException e) {
            Assert.assertEquals(400, e.getResponseCode());
        }
        // 4、使用headObject接口获取软连接对象的元数据
        ObjectMetadata  metadataSymlink =
            obsClient.getObjectMetadata(bucketName, "B");
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 5、使用getObject接口获取软连接对象的元数据和数据
        ObsObject obsObject = obsClient.getObject(bucketName, "B");
        obsObject.getObjectContent().close();
        metadataSymlink = obsObject.getMetadata();
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
    }
    @Test
    public void tc_obs_symlink_create_010() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();

        cleanUp(obsClient,bucketName);
        assert obsClient != null;
        String encObjectKey = bucketName + "-normal";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, encObjectKey, 1024 * 1024);
        // 1、在测试桶，上创建SSE-C加密对象；
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, encObjectKey, testFile);
        //设置自定义头域，设置使用sse-obs加密
        HashMap<String, String> userHeaders = new HashMap<>();
        userHeaders.put("x-obs-server-side-encryption","AES256");
        putObjectRequest.setUserHeaders(userHeaders);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 2、创建指向步骤一加密对象的软连接对象；
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, encObjectKey + ".symlink", encObjectKey);
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());
        // 3、使用getsymlink接口获取软连接对象的元数据
        GetSymlinkRequest getSymlinkRequest =
            new GetSymlinkRequest(bucketName, encObjectKey + ".symlink");
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
        Assert.assertEquals(200, getSymlinkResult.getStatusCode());

        // 4、使用headObject接口获取软连接对象的元数据
        ObjectMetadata metadataSymlink =
            obsClient.getObjectMetadata(bucketName, encObjectKey + ".symlink");
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 5、使用getObject接口获取软连接对象的元数据和数据
        ObsObject obsObject = obsClient.getObject(bucketName, encObjectKey + ".symlink");
        obsObject.getObjectContent().close();
        metadataSymlink = obsObject.getMetadata();
        Assert.assertEquals(200, metadataSymlink.getStatusCode());
        // 6、删除4个软连接对象
        // 预期：响应204
        HeaderResponse response1 =
            obsClient.deleteObject(bucketName, encObjectKey + ".symlink");
        Assert.assertEquals(204, response1.getStatusCode());
    }
    @Test
    public void tc_obs_symlink_create_011() throws IOException {
        String bucketName = "obs-sdk-symlink";
        ObsClient obsClient = TestTools.getPipelineForSymEnvironment();

        cleanUp(obsClient,bucketName);
        assert obsClient != null;
        String encObjectKey = bucketName + "-normal";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, encObjectKey, 1024 * 1024);
        // 1、在测试桶，上创建SSE-C加密对象；
        PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, encObjectKey, testFile);
        // 设置SSE-C算法加密对象
        SseCHeader sseCHeader = new SseCHeader();
        sseCHeader.setAlgorithm(ServerAlgorithm.AES256);
        // 设置SSE-C方式下使用的密钥，用于加解密对象，该值是密钥进行base64encode后的值
        sseCHeader.setSseCKeyBase64(sse_c_base64);
        putObjectRequest.setSseCHeader(sseCHeader);
        PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        // 2、创建指向步骤一加密对象的软连接对象；
        PutSymlinkRequest putSymlinkRequest = new PutSymlinkRequest(
            bucketName, encObjectKey + ".symlink", encObjectKey);
        HeaderResponse headerResponse = obsClient.putSymlink(putSymlinkRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());
        // 3、使用getsymlink接口获取软连接对象的元数据
        GetSymlinkRequest getSymlinkRequest =
            new GetSymlinkRequest(bucketName, encObjectKey + ".symlink");
        GetSymlinkResult getSymlinkResult = obsClient.getSymlink(getSymlinkRequest);
        Assert.assertEquals(200, getSymlinkResult.getStatusCode());

        // 4、不携带任何解密密钥，使用headObject接口获取软连接对象的元数据
        try {
            obsClient.getObjectMetadata(bucketName, encObjectKey + ".symlink");
            fail();
        } catch (ObsException e) {
            assertEquals(400, e.getResponseCode());
        }
        // 5、不携带任何解密密钥，使用getObject接口获取软连接对象的元数据和数据
        try {
            obsClient.getObject(bucketName, encObjectKey + ".symlink");
            fail();
        } catch (ObsException e) {
            assertEquals(400, e.getResponseCode());
        }
        // 6、删除4个软连接对象
        // 预期：响应204
        HeaderResponse response1 =
            obsClient.deleteObject(bucketName, encObjectKey + ".symlink");
        Assert.assertEquals(204, response1.getStatusCode());
    }
}
