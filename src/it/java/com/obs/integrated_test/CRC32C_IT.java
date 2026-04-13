package com.obs.integrated_test;

import static com.obs.services.internal.Constants.CommonHeaders.HASH_CRC32C;
import static com.obs.test.TestTools.downloadFileWithRetry;
import static com.obs.test.TestTools.genTestFile;
import static com.obs.test.TestTools.printObsException;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNotSame;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.AbortMultipartUploadRequest;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CompleteMultipartUploadRequest;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.CopyObjectResult;
import com.obs.services.model.CopyPartRequest;
import com.obs.services.model.CopyPartResult;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.InitiateMultipartUploadRequest;
import com.obs.services.model.InitiateMultipartUploadResult;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PartEtag;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.UploadPartRequest;
import com.obs.services.model.UploadPartResult;
import com.obs.test.TestTools;
import com.obs.test.tools.CRC32C;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Locale;

public class CRC32C_IT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    /***
     * 1、用obs api创建桶
     * 2、带crc32c头域为正确的crc32值，上传1MB对象(使用obs协议)
     * 3、下载校验对象(使用obs协议)
     * 4、带crc32c头域为错误的crc32值，上传1MB对象(使用obs协议)
     *
     * @throws IOException
     */
    @Test
    public void tc_putObject_with_crc32() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT) + "pfs";
        String testFileName = bucketName + "testFile";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
        String objectKey = bucketName + "objectKey_001";
        long testFileCrc32 = CRC32C_Calculator.CRC32C_FromFile(testFile);
        String testFileCrc32UnsignedString = Long.toUnsignedString(testFileCrc32);
        ObsClient[] obsClients = new ObsClient[1];
        obsClients[0] = TestTools.getPipelineEnvironment_OBS();
        String[] protocolPrefixes = new String[1];
        protocolPrefixes[0] = "x-obs-";
        for (int i = 0; i < obsClients.length; ++i) {
            ObsClient obsClient = obsClients[i];
            String protocolPrefix = protocolPrefixes[i];
            try {
                assert obsClient != null;
                CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketName);
                createBucketRequest.setBucketType(BucketTypeEnum.PFS);
                ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
                Assert.assertEquals(200, obsBucket.getStatusCode());
                {
                    PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
                    putObjectRequest.addUserHeaders(protocolPrefix + HASH_CRC32C, testFileCrc32UnsignedString);
                    // 带crc32头域为正确的crc32值，上传1MB对象
                    PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
                    String responseCRC32 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, testFileCrc32UnsignedString);
                }
                {
                    ObsObject obsObject = obsClient.getObject(bucketName, objectKey);
                    long crc32_get_calculate = CRC32C_Calculator.CRC32C_FromInputStream(obsObject.getObjectContent());
                    obsObject.getObjectContent().close();

                    String crc32_get_calculateInUnsignedString = Long.toUnsignedString(crc32_get_calculate);
                    Assert.assertEquals(crc32_get_calculateInUnsignedString, testFileCrc32UnsignedString);
                    String responseCRC32 = (String) obsObject.getMetadata().getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, testFileCrc32UnsignedString);
                    assertEquals(responseCRC32, crc32_get_calculateInUnsignedString);
                    assertEquals(testFileCrc32, crc32_get_calculate);
                }
                try {
                    // generate wrong crc32UnsignedString
                    String wrongtestFilecrc32UnsignedString = Long.toUnsignedString(testFileCrc32 + 1);
                    PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
                    putObjectRequest.addUserHeaders(protocolPrefix + HASH_CRC32C, wrongtestFilecrc32UnsignedString);
                    // 带crc32头域为错误的crc32值，上传1MB对象
                    obsClient.putObject(putObjectRequest);
                    // 会报错的话，就不会触发这个断言
                    fail();
                } catch (ObsException e) {
                    if (e.getResponseCode() != 400) {
                        printObsException(e);
                    }
                    assertEquals(400, e.getResponseCode());
                    assertEquals("InvalidCRC32C", e.getErrorCode());
                }
            } finally {
                try {
                    TestTools.delete_bucket(obsClient, bucketName);
                } catch (Throwable ignore) {
                }
            }
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc32c头域为正确的crc32c值，上传1MB对象(使用obs协议)
     * 3、下载校验对象(使用obs协议)，同时验证获取元数据可以获取到crc32c
     * 4、带crc32c头域为错误的crc32c值，上传1MB对象(使用obs协议)
     *
     * @throws IOException
     */
    @Test
    public void tc_getObject_with_crc32() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT) + "pfs";
        String testFileName = bucketName + "testFile";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
        String objectKey = bucketName + "objectKey_001";
        long localFileCrc32 = CRC32C_Calculator.CRC32C_FromFile(testFile);
        String localFileCrc32UnsignedString = Long.toUnsignedString(localFileCrc32);
        ObsClient[] obsClients = new ObsClient[1];
        obsClients[0] = TestTools.getPipelineEnvironment_OBS();
        String[] protocolPrefixes = new String[1];
        protocolPrefixes[0] = "x-obs-";
        for (int i = 0; i < obsClients.length; ++i) {
            ObsClient obsClient = obsClients[i];
            String protocolPrefix = protocolPrefixes[i];
            try {
                assert obsClient != null;
                CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketName);
                createBucketRequest.setBucketType(BucketTypeEnum.PFS);
                ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
                Assert.assertEquals(200, obsBucket.getStatusCode());
                {
                    PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
                    putObjectRequest.addUserHeaders(protocolPrefix + HASH_CRC32C, localFileCrc32UnsignedString);
                    // 带crc32头域为正确的crc32值，上传1MB对象
                    PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
                    String responseCRC32 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, localFileCrc32UnsignedString);
                }
                {
                    ObsObject obsObject = obsClient.getObject(bucketName, objectKey);
                    long crc32_get_calculate = CRC32C_Calculator.CRC32C_FromInputStream(obsObject.getObjectContent());
                    obsObject.getObjectContent().close();

                    String crc32_get_calculateInUnsignedString = Long.toUnsignedString(crc32_get_calculate);
                    Assert.assertEquals(crc32_get_calculateInUnsignedString, localFileCrc32UnsignedString);
                    String responseCRC32 = (String) obsObject.getMetadata().getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, localFileCrc32UnsignedString);
                    assertEquals(responseCRC32, crc32_get_calculateInUnsignedString);
                    assertEquals(localFileCrc32, crc32_get_calculate);
                }
                {
                    ObjectMetadata objectMetadata = obsClient.getObjectMetadata(bucketName, objectKey);
                    String responseCRC32 = (String) objectMetadata.getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, localFileCrc32UnsignedString);
                }
            } finally {
                try {
                    TestTools.delete_bucket(obsClient, bucketName);
                } catch (Throwable ignore) {
                }
            }
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc32c头域为正确的crc32c值，上传1MB对象(使用obs协议)
     * 3、针对拷贝对象接口，服务端返回body体中会新增crc32c的值
     *
     * @throws IOException
     */
    @Test
    public void tc_copyObject_with_crc32() throws IOException {
        String bucketName = System.getenv("CRC32C_BucketName");
        String testFileName = bucketName + "testFile";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
        String objectKey = bucketName + "objectKey_001";
        String copy_objectKey = bucketName + "objectKey_001_copy";
        long localFileCrc32 = CRC32C_Calculator.CRC32C_FromFile(testFile);
        String localFileCrc32UnsignedString = Long.toUnsignedString(localFileCrc32);
        ObsClient[] obsClients = new ObsClient[1];
        obsClients[0] = TestTools.getPipelineEnvironment_OBS();
        String[] protocolPrefixes = new String[1];
        protocolPrefixes[0] = "x-obs-";
        for (int i = 0; i < obsClients.length; ++i) {
            ObsClient obsClient = obsClients[i];
            String protocolPrefix = protocolPrefixes[i];
            assert obsClient != null;
            try {
                {
                    PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
                    putObjectRequest.addUserHeaders(protocolPrefix + HASH_CRC32C, localFileCrc32UnsignedString);
                    // 带crc32头域为正确的crc32值，上传1MB对象
                    PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
                    String responseCRC32 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, localFileCrc32UnsignedString);
                }
                {
                    CopyObjectResult copyObjectResult =
                            obsClient.copyObject(bucketName, objectKey, bucketName, copy_objectKey);
                    String responseCRC32 = copyObjectResult.getCRC32C();
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, localFileCrc32UnsignedString);
                }
            } finally {
                try {
                    obsClient.deleteObject(bucketName, objectKey);
                    obsClient.deleteObject(bucketName, copy_objectKey);
                } catch (Throwable ignore) {
                }
            }
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc32c头域为正确的crc32c值，分段上传1MB对象，并且合并，响应中也有crc32c(使用obs协议)
     * 3、针对拷贝段接口，服务端返回body体中会新增crc32c的值
     *
     * @throws IOException
     */
    @Test
    public void tc_multipart_with_crc32() throws IOException {
        String bucketName = System.getenv("CRC32C_BucketName");
        String testFileName = bucketName + "testFile";
        // 1 mb test file
        final long partSize = 1024 * 1024;
        int partCount = 3;
        File testFile = genTestFile(temporaryFolder, testFileName, partCount * partSize);
        String objectKey = bucketName + "objectKey_001";
        String copy_objectKey = bucketName + "objectKey_001_copy";
        long localFileCrc32 = CRC32C_Calculator.CRC32C_FromFile(testFile);
        String localFileCrc32UnsignedString = Long.toUnsignedString(localFileCrc32);
        ObsClient[] obsClients = new ObsClient[1];
        obsClients[0] = TestTools.getPipelineEnvironment_OBS();
        String[] protocolPrefixes = new String[1];
        protocolPrefixes[0] = "x-obs-";
        for (int i = 0; i < obsClients.length; ++i) {
            ArrayList<PartEtag> eTags = new ArrayList<>();
            ObsClient obsClient = obsClients[i];
            String protocolPrefix = protocolPrefixes[i];
            assert obsClient != null;
            String uploadID;
            {
                InitiateMultipartUploadRequest initiateMultipartUploadRequest =
                        new InitiateMultipartUploadRequest(bucketName, objectKey);
                InitiateMultipartUploadResult initiateMultipartUploadResult =
                        obsClient.initiateMultipartUpload(initiateMultipartUploadRequest);
                Assert.assertEquals(200, initiateMultipartUploadResult.getStatusCode());
                uploadID = initiateMultipartUploadResult.getUploadId();
            }
            try {
                for (int j = 1; j < partCount; j++) {
                    long offset = (j - 1) * partSize;
                    long partCRC32C = CRC32C_Calculator.CRC32C_FromFile(testFile, offset, partSize);
                    String partCRC32CUnsignedString = Long.toUnsignedString(partCRC32C);
                    UploadPartRequest uploadPartRequest = new UploadPartRequest();
                    uploadPartRequest.setBucketName(bucketName);
                    uploadPartRequest.setObjectKey(objectKey);
                    uploadPartRequest.setPartNumber(j);
                    uploadPartRequest.setUploadId(uploadID);
                    uploadPartRequest.setPartSize(partSize);
                    uploadPartRequest.setOffset(offset);
                    uploadPartRequest.setFile(testFile);
                    uploadPartRequest.addUserHeaders(protocolPrefix + HASH_CRC32C, partCRC32CUnsignedString);
                    UploadPartResult uploadPartResult = obsClient.uploadPart(uploadPartRequest);
                    Assert.assertEquals(200, uploadPartResult.getStatusCode());
                    String responseCRC32 = (String) uploadPartResult.getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, partCRC32CUnsignedString);
                    eTags.add(new PartEtag(uploadPartResult.getEtag(), j));
                }
                try (FileInputStream fileInputStream = new FileInputStream(testFile)) {
                    long offset = (partCount - 1) * partSize;
                    Assert.assertEquals(offset, fileInputStream.skip(offset));
                    long partCRC32C = CRC32C_Calculator.CRC32C_FromFile(testFile, offset, partSize);
                    String partCRC32CUnsignedString = Long.toUnsignedString(partCRC32C);
                    PutObjectRequest putObjectRequest =
                            new PutObjectRequest(bucketName, copy_objectKey, fileInputStream);
                    putObjectRequest.addUserHeaders(protocolPrefix + HASH_CRC32C, partCRC32CUnsignedString);
                    // 带crc32头域为正确的crc32值，上传1MB对象
                    PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
                    String responseCRC32 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC32C);
                    assertNotNull(responseCRC32);
                    assertEquals(responseCRC32, partCRC32CUnsignedString);

                    CopyPartRequest copyPartRequest =
                            new CopyPartRequest(uploadID, bucketName, copy_objectKey, bucketName, objectKey, partCount);
                    CopyPartResult copyPartResult = obsClient.copyPart(copyPartRequest);
                    String responseCRC32_copy = copyPartResult.getCRC32C();
                    assertNotNull(responseCRC32_copy);
                    assertEquals(responseCRC32, responseCRC32_copy);
                    eTags.add(new PartEtag(copyPartResult.getEtag(), partCount));
                }

                CompleteMultipartUploadRequest completeMultipartUploadRequest =
                        new CompleteMultipartUploadRequest(bucketName, objectKey, uploadID, eTags);
                completeMultipartUploadRequest.addUserHeaders(
                        protocolPrefix + HASH_CRC32C, localFileCrc32UnsignedString);
                CompleteMultipartUploadResult completeMultipartUploadResult =
                        obsClient.completeMultipartUpload(completeMultipartUploadRequest);
                Assert.assertEquals(200, completeMultipartUploadResult.getStatusCode());
            } finally {
                try {
                    AbortMultipartUploadRequest abortMultipartUploadRequest =
                            new AbortMultipartUploadRequest(bucketName, objectKey, uploadID);
                    obsClient.abortMultipartUpload(abortMultipartUploadRequest);
                    obsClient.deleteObject(bucketName, objectKey);
                    obsClient.deleteObject(bucketName, copy_objectKey);
                } catch (Throwable ignore) {
                }
            }
        }
    }

    static class CRC32C_Calculator {
        public static long CRC32C_FromFile(File file) throws IOException {
            long crc32 = 0;
            try (FileInputStream fileInputStream = new FileInputStream(file)) {
                crc32 = CRC32C_FromInputStream(fileInputStream);
            }
            return crc32;
        }

        public static long CRC32C_FromFile(File file, long offset, long length) throws IOException {
            long crc32 = 0;
            try (FileInputStream fileInputStream = new FileInputStream(file)) {
                crc32 = CRC32C_FromInputStream(fileInputStream, offset, length);
            }
            return crc32;
        }

        public static long CRC3C_FromFilePath(String filePath) throws IOException {
            return CRC32C_FromFile(new File(filePath));
        }

        public static long CRC32C_FromInputStream(InputStream inputStream) throws IOException {
            byte[] buffer = new byte[1024];
            CRC32C crc32Instance = new CRC32C();

            int bytesRead;
            while ((bytesRead = inputStream.read(buffer)) > 0) {
                crc32Instance.update(buffer, 0, bytesRead);
            }
            return crc32Instance.getValue();
        }

        public static long CRC32C_FromInputStream(InputStream in, long offset, long sizeToReadTotal)
                throws IOException {
            in.skip(offset);
            CRC32C crc32Instance = new CRC32C();
            int bufferSize = 4096;
            byte[] b = new byte[bufferSize];
            int l;
            long sizeToRead = Long.min(sizeToReadTotal, bufferSize);
            long bytesReadTotal = 0;
            while (sizeToRead > 0 && (l = in.read(b, 0, (int) sizeToRead)) > 0) {
                crc32Instance.update(b, 0, l);
                sizeToReadTotal -= l;
                bytesReadTotal += l;
                sizeToRead = Long.min(sizeToReadTotal, bufferSize);
            }
            return crc32Instance.getValue();
        }
    }
}
