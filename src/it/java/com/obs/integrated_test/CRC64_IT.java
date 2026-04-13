package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.ObsClientAsync;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.handler.XmlResponsesSaxParser;
import com.obs.services.internal.task.UploadFileTask;
import com.obs.services.internal.utils.CRC64;
import com.obs.services.internal.utils.CRC64InputStream;
import com.obs.services.internal.utils.CallCancelHandler;
import com.obs.services.internal.utils.ServiceUtils;
import com.obs.services.model.AppendObjectRequest;
import com.obs.services.model.AppendObjectResult;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CompleteMultipartUploadRequest;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.CopyObjectRequest;
import com.obs.services.model.CopyObjectResult;
import com.obs.services.model.CopyPartRequest;
import com.obs.services.model.CopyPartResult;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.DownloadFileResult;
import com.obs.services.model.GetObjectRequest;
import com.obs.services.model.InitiateMultipartUploadRequest;
import com.obs.services.model.InitiateMultipartUploadResult;
import com.obs.services.model.ModifyObjectRequest;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PartEtag;
import com.obs.services.model.ProgressListener;
import com.obs.services.model.ProgressStatus;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.TaskCallback;
import com.obs.services.model.UploadFileRequest;
import com.obs.services.model.UploadPartRequest;
import com.obs.services.model.UploadPartResult;
import com.obs.services.model.fs.NewFileRequest;
import com.obs.services.model.fs.WriteFileRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

import static com.obs.services.internal.Constants.CommonHeaders.HASH_CRC64ECMA;
import static com.obs.test.TestTools.areByteArraysEqual;
import static com.obs.test.TestTools.downloadFileWithRetry;
import static com.obs.test.TestTools.genTestFile;
import static com.obs.test.TestTools.getPipeLineTestSecureRandom;
import static com.obs.test.TestTools.getPipelineEnvironment;
import static com.obs.test.TestTools.getTestRandomIntInRange;
import static com.obs.test.TestTools.printObsException;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNotSame;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.fail;

public class CRC64_IT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，上传1MB对象(使用obs协议)
     * 3、下载校验对象(使用obs协议)
     * 4、带crc64头域为错误的crc64值，上传1MB对象(使用obs协议)
     *
     * @throws IOException
     */
    @Test
    public void tc_putObject_with_crc64_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
        // 1 mb test file
        File testFileGet = genTestFile(temporaryFolder, testFileName + "Get", 0);
        String objectKey = bucketName + "objectKey_001";
        long testFileCrc64;
        try (FileInputStream fileInputStream = new FileInputStream(testFile);
                CRC64InputStream crc64InputStream = new CRC64InputStream(fileInputStream)) {
            byte[] buffer = new byte[65536];
            while (crc64InputStream.read(buffer) != -1) {}
            testFileCrc64 = crc64InputStream.getCrc64().getValue();
        }
        String testFileCrc64UnsignedString = CRC64.toString(testFileCrc64);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.addUserHeaders("x-obs-" + HASH_CRC64ECMA, testFileCrc64UnsignedString);
            // 带crc64头域为正确的crc64值，上传1MB对象
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            String responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
            // 上传对象,sdk自动计算crc64
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(true);
            putObjectRequest.setFile(testFile);
            putObjectResult = obsClient.putObject(putObjectRequest);
            responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
        CRC64InputStream crc64InputStream = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS();
                FileOutputStream fileOutputStream = new FileOutputStream(testFileGet)) {
            ObsObject obsObject = obsClient.getObject(bucketName, objectKey);
            crc64InputStream = new CRC64InputStream(obsObject.getObjectContent());
            int len = 0;
            byte[] buffer = new byte[65536];
            while ((len = crc64InputStream.read(buffer)) != -1) {
                fileOutputStream.write(buffer, 0, len);
            }
            long crc64_get_calculate = crc64InputStream.getCrc64().getValue();
            String crc64_get_calculateInUnsignedString = CRC64.toString(crc64_get_calculate);
            Assert.assertEquals(crc64_get_calculateInUnsignedString, testFileCrc64UnsignedString);
            String responseCrc64 = (String) obsObject.getMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
            assertEquals(responseCrc64, crc64_get_calculateInUnsignedString);
            assertEquals(testFileCrc64, crc64_get_calculate);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } finally {
            if (crc64InputStream != null) {
                crc64InputStream.close();
            }
        }
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            // generate wrong crc64UnsignedString
            String wrongtestFileCrc64UnsignedString = CRC64.toString(testFileCrc64 + 1);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.addUserHeaders("x-obs-" + HASH_CRC64ECMA, wrongtestFileCrc64UnsignedString);
            // 带crc64头域为错误的crc64值，上传1MB对象
            obsClient.putObject(putObjectRequest);
            // 会报错的话，就不会触发这个断言
            fail();
        } catch (ObsException e) {
            assertNotNull(e);
            assertEquals(400, e.getResponseCode());
            assertEquals("InvalidCRC64", e.getErrorCode());
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，上传1MB对象(使用s3协议)
     * 3、下载校验对象(使用s3协议)
     * 4、带crc64头域为错误的crc64值，上传1MB对象(使用s3协议)
     *
     * @throws IOException
     */
    @Test
    public void tc_putObject_with_crc64_002() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
        // 1 mb test file
        File testFileGet = genTestFile(temporaryFolder, testFileName + "Get", 0);
        String objectKey = bucketName + "objectKey_001";
        long testFileCrc64;
        try (FileInputStream fileInputStream = new FileInputStream(testFile);
                CRC64InputStream crc64InputStream = new CRC64InputStream(fileInputStream)) {
            byte[] buffer = new byte[65536];
            while (crc64InputStream.read(buffer) != -1) {}
            testFileCrc64 = crc64InputStream.getCrc64().getValue();
        }
        String testFileCrc64UnsignedString = CRC64.toString(testFileCrc64);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V2()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.addUserHeaders("x-amz-" + HASH_CRC64ECMA, testFileCrc64UnsignedString);
            // 带crc64头域为正确的crc64值，上传1MB对象
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            String responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
            // 上传对象,sdk自动计算crc64
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(true);
            putObjectRequest.setFile(testFile);
            putObjectResult = obsClient.putObject(putObjectRequest);
            responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
        CRC64InputStream crc64InputStream = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V2();
             FileOutputStream fileOutputStream = new FileOutputStream(testFileGet)) {
            ObsObject obsObject = obsClient.getObject(bucketName, objectKey);
            crc64InputStream = obsObject.getObjectContentWithCRC64();
            int len;
            byte[] buffer = new byte[65536];
            while ((len = crc64InputStream.read(buffer)) != -1) {
                fileOutputStream.write(buffer, 0, len);
            }
            long crc64_get_calculate = crc64InputStream.getCrc64().getValue();
            String crc64_get_calculateInUnsignedString = CRC64.toString(crc64_get_calculate);
            Assert.assertEquals(crc64_get_calculateInUnsignedString, testFileCrc64UnsignedString);
            String responseCrc64 = (String) obsObject.getMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
            assertEquals(responseCrc64, crc64_get_calculateInUnsignedString);
            assertEquals(testFileCrc64, crc64_get_calculate);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } finally {
            if (crc64InputStream != null) {
                crc64InputStream.close();
            }
        }
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_V2()) {
            assert obsClient != null;
            // generate wrong crc64UnsignedString
            String wrongtestFileCrc64UnsignedString = CRC64.toString(testFileCrc64 + 1);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.addUserHeaders("x-amz-" + HASH_CRC64ECMA, wrongtestFileCrc64UnsignedString);
            // 带crc64头域为错误的crc64值，上传1MB对象
            obsClient.putObject(putObjectRequest);
            // 会报错的话，就不会触发这个断言
            fail();
        } catch (ObsException e) {
            assertNotNull(e);
            assertEquals(400, e.getResponseCode());
            assertEquals("InvalidCRC64", e.getErrorCode());
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，上传空字符串
     * 3、下载校验对象
     *
     * @throws IOException
     */
    @Test
    public void tc_putObject_with_crc64_003() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // empty test file
        File testFile = genTestFile(temporaryFolder, testFileName, 0);
        // test file to Get
        File testFileGet = genTestFile(temporaryFolder, testFileName + "Get", 0);
        String objectKey = bucketName + "objectKey_001";
        long testFileCrc64;
        try (CRC64InputStream crc64InputStream = new CRC64InputStream(new FileInputStream(testFile))) {
            byte[] buffer = new byte[65536];
            while (crc64InputStream.read(buffer) != -1) {}
            testFileCrc64 = crc64InputStream.getCrc64().getValue();
        }
        String testFileCrc64UnsignedString = CRC64.toString(testFileCrc64);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.addUserHeaders("x-obs-" + HASH_CRC64ECMA, testFileCrc64UnsignedString);
            // 带crc64头域为正确的crc64值，上传对象
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            String responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
            // 上传对象,sdk自动计算crc64
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(true);
            putObjectRequest.setFile(testFile);
            putObjectResult = obsClient.putObject(putObjectRequest);
            responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
        InputStream inputStream = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS();
                FileOutputStream fileOutputStream = new FileOutputStream(testFileGet)) {
            ObsObject obsObject = obsClient.getObject(bucketName, objectKey);
            inputStream = obsObject.getObjectContent();
            int len = 0;
            byte[] buffer = new byte[65536];
            while ((len = inputStream.read(buffer)) != -1) {
                fileOutputStream.write(buffer, 0, len);
            }
            long crc64_get_calculate = CRC64.fromFile(testFileGet).getValue();
            String crc64_get_calculateInUnsignedString = CRC64.toString(crc64_get_calculate);
            Assert.assertEquals(crc64_get_calculateInUnsignedString, testFileCrc64UnsignedString);
            String responseCrc64 = (String) obsObject.getMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
            assertEquals(responseCrc64, crc64_get_calculateInUnsignedString);
            assertEquals(testFileCrc64, crc64_get_calculate);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } finally {
            if (inputStream != null) {
                inputStream.close();
            }
        }
    }

    /***
     * 1、用obs api创建POSIX桶
     * 2、带crc64头域为正确的crc64值，上传1MB对象
     * 3、上传对象失败，返回405
     *
     * @throws IOException
     */
    @Test
    public void tc_putObject_with_crc64_004() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT) + "-pfs";
        String testFileName = bucketName + "_testFile";
        // empty test file
        File testFile = genTestFile(temporaryFolder, testFileName, 0);
        String objectKey = bucketName + "_objectKey_001";
        ObsClient obsClient = TestTools.getPipelineEnvironment_OBS();
        try {
            assert obsClient != null;
            CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketName);
            createBucketRequest.setBucketType(BucketTypeEnum.PFS);
            obsClient.createBucket(createBucketRequest);
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(true);
            // 上传对象,sdk自动计算crc64
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            assertNull(putObjectResult);
        } catch (ObsException e) {
            assertEquals(405, e.getResponseCode());
            assertEquals("MethodNotAllowed", e.getErrorCode());
        } finally {
            assert obsClient != null;
            obsClient.deleteObject(bucketName, objectKey);
            obsClient.deleteBucket(bucketName);
            obsClient.close();
        }
    }

    /***
     * 1、用obs api创建桶，初始化上传段任务
     * 2、带crc64头域为正确的crc64值，上传1MB多段
     * 3、带crc64头域为错误的crc64值，再上传1MB多段，上传多段成功，返回头域带crc64值，与客户端算的一致
     * 4、合并段，合并成功，返回头域带crc64值，与客户端算的一致
     * 5、下载对象成功，返回头域带crc64值，与客户端算的一致
     *
     * @throws IOException
     */
    @Test
    public void tc_uploadPart_with_crc64_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        // 2 mb test file for multipart test
        long fileSizeInBytes = 2 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        File testFileGet = genTestFile(temporaryFolder, testFileName + "_Get", 0L);
        String objectKey = bucketName + "_objectKey_001";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            InitiateMultipartUploadRequest initiateMultipartUploadRequest =
                    new InitiateMultipartUploadRequest(bucketName, objectKey);
            InitiateMultipartUploadResult initiateMultipartUploadResult =
                    obsClient.initiateMultipartUpload(initiateMultipartUploadRequest);
            String uploadId = initiateMultipartUploadResult.getUploadId();
            // 分段大小
            long currPartSize = 1024 * 1024;
            // 计算需要上传的段数
            long partCount = 2;
            final List<PartEtag> partETags = Collections.synchronizedList(new ArrayList<>());
            final List<CRC64> partsCRC64 = Collections.synchronizedList(new ArrayList<>());
            final List<Long> partsSize = Collections.synchronizedList(new ArrayList<>());
            UploadPartRequest uploadPartRequest = new UploadPartRequest();
            UploadPartResult uploadPartResult;
            // 上传段1
            // 分段在文件中的起始位置
            long offset = 0;
            // 分段号
            int partNumber = 1;
            uploadPartRequest.setBucketName(bucketName);
            uploadPartRequest.setObjectKey(objectKey);
            uploadPartRequest.setUploadId(uploadId);
            uploadPartRequest.setFile(testFile);
            uploadPartRequest.setPartSize(currPartSize);
            uploadPartRequest.setOffset(offset);
            uploadPartRequest.setPartNumber(partNumber);
            uploadPartRequest.setNeedCalculateCRC64(true);
            // 分段1上传成功
            uploadPartResult = obsClient.uploadPart(uploadPartRequest);
            String uploadPartResult1ServerCrc64 = (String) uploadPartResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertEquals(uploadPartResult1ServerCrc64, uploadPartResult.getClientCalculatedCRC64().toString());
            partETags.add(new PartEtag(uploadPartResult.getEtag(), uploadPartResult.getPartNumber()));
            partsCRC64.add(uploadPartResult.getClientCalculatedCRC64());
            partsSize.add(currPartSize);
            CRC64 crc64Part2;
            try (FileInputStream fileInputStream = new FileInputStream(testFile)) {
                uploadPartRequest.setOffset(offset + currPartSize);
                uploadPartRequest.setPartNumber(2);
                uploadPartRequest.setNeedCalculateCRC64(false);
                uploadPartRequest.setFile(testFile);
                crc64Part2 = CRC64.fromInputStream(fileInputStream, 0, currPartSize);
                // 设置一个错的crc64值
                uploadPartRequest.addUserHeaders("x-amz-" + HASH_CRC64ECMA, CRC64.toString(crc64Part2.getValue() + 1));
                uploadPartRequest.addUserHeaders("x-obs-" + HASH_CRC64ECMA, CRC64.toString(crc64Part2.getValue() + 1));
                obsClient.uploadPart(uploadPartRequest);
                fail("part 2 with wrong crc64 should fail.");
            } catch (ObsException e) {
                assertNotNull(e);
                assertEquals(400, e.getResponseCode());
                assertEquals("InvalidCRC64", e.getErrorCode());
            }
            CRC64 crc64Total = partsCRC64.get(0);
            // 合并段
            CompleteMultipartUploadRequest completeMultipartUploadRequest =
                    new CompleteMultipartUploadRequest(bucketName, objectKey, uploadId, partETags);
            String completeMultipartUploadRequestCrc64 = CRC64.toString(crc64Total.getValue());
            completeMultipartUploadRequest.addUserHeaders(
                    "x-obs-" + HASH_CRC64ECMA, completeMultipartUploadRequestCrc64);
            completeMultipartUploadRequest.addUserHeaders(
                    "x-amz-" + HASH_CRC64ECMA, completeMultipartUploadRequestCrc64);
            CompleteMultipartUploadResult completeMultipartUploadResult =
                    obsClient.completeMultipartUpload(completeMultipartUploadRequest);
            assertNotNull(completeMultipartUploadResult);
            Map<String, Object> completeMultipartUploadResultResponseHeaders =
                    completeMultipartUploadResult.getResponseHeaders();
            assertNotNull(completeMultipartUploadResultResponseHeaders);
            String completeMultipartUploadResultCRC64 =
                    (String) completeMultipartUploadResultResponseHeaders.get(HASH_CRC64ECMA);
            assertNotNull(completeMultipartUploadResultCRC64);
            assertEquals(completeMultipartUploadResultCRC64, completeMultipartUploadRequestCrc64);
            assertEquals(uploadPartResult1ServerCrc64, completeMultipartUploadResultCRC64);
            System.out.println("completeMultipartUploadResultCRC64:" + completeMultipartUploadResultCRC64);
            DownloadFileResult downloadFileResult =
                    downloadFileWithRetry(obsClient, bucketName, objectKey, testFileGet, currPartSize);
            assertNotNull(downloadFileResult);
            String downloadFileResultCRC64 =
                    (String) downloadFileResult.getObjectMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertEquals(downloadFileResultCRC64, completeMultipartUploadRequestCrc64);
            assertEquals(completeMultipartUploadResultCRC64, downloadFileResultCRC64);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，append上传1MB对象成功
     * 3、下载校验对象成功，返回头域带crc64值，与客户端算的一致
     * 4、带crc64头域为错误的crc64值，append上传1MB对象失败。报400
     *
     * @throws IOException
     */
    @Test
    public void tc_appendObject_with_crc64_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        // 2 mb test file for multipart test
        long fileSizeInBytes = 2 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        File testFileGet = genTestFile(temporaryFolder, testFileName + "_Get", 0L);
        String objectKey = bucketName + "_objectKey_001";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            AppendObjectRequest appendObjectRequest = new AppendObjectRequest(bucketName);
            appendObjectRequest.setObjectKey(objectKey);
            appendObjectRequest.setNeedCalculateCRC64(true);
            appendObjectRequest.setFile(testFile);
            AppendObjectResult appendObjectResult = obsClient.appendObject(appendObjectRequest);
            assertEquals(
                    appendObjectResult.getClientCalculatedCRC64().toString(),
                    appendObjectResult.getResponseHeaders().get(HASH_CRC64ECMA));

            DownloadFileResult downloadFileResult =
                    downloadFileWithRetry(obsClient, bucketName, objectKey, testFileGet, fileSizeInBytes);
            assertNotNull(downloadFileResult);
            String downloadFileResultCRC64 =
                    (String) downloadFileResult.getObjectMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertEquals(downloadFileResultCRC64, appendObjectResult.getResponseHeaders().get(HASH_CRC64ECMA));

            // 设置一个错的crc64值
            String wrongCrc64 = CRC64.toString(appendObjectResult.getClientCalculatedCRC64().getValue() + 1);
            appendObjectRequest.setNeedCalculateCRC64(false);
            appendObjectRequest.setFile(testFile);
            appendObjectRequest.addUserHeaders("x-amz-" + HASH_CRC64ECMA, wrongCrc64);
            appendObjectRequest.addUserHeaders("x-obs-" + HASH_CRC64ECMA, wrongCrc64);
            appendObjectRequest.setPosition(appendObjectResult.getNextPosition());
            try {
                obsClient.appendObject(appendObjectRequest);
                fail("appendObject with wrong crc64 should fail.");
            } catch (ObsException e) {
                assertEquals(400, e.getResponseCode());
                assertEquals("InvalidCRC64", e.getErrorCode());
            }

        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，上传对象成功
     * 3、getObject成功，返回头域带crc64值，与客户端算的一致，是整个对象的crc64值
     *
     * @throws IOException
     */
    @Test
    public void tc_get_crc64_by_GetObject_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // random test file
        long testFileSizeInBytes = getTestRandomIntInRange(5 * 1024 * 1024, 10 * 1024 * 1024);
        File testFile = genTestFile(temporaryFolder, testFileName, testFileSizeInBytes);
        String objectKey = bucketName + "objectKey_001";
        String responseCrc64 = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            // 上传对象,sdk自动计算crc64
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(true);
            putObjectRequest.setFile(testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
        InputStream inputStream = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            GetObjectRequest getObjectRequest = new GetObjectRequest(bucketName, objectKey);
            getObjectRequest.setRangeStart(0L);
            getObjectRequest.setRangeEnd((long) getTestRandomIntInRange(0, (int) testFileSizeInBytes));
            ObsObject obsObject = obsClient.getObject(getObjectRequest);
            inputStream = obsObject.getObjectContent();
            String responseCrc64OfGetObject = (String) obsObject.getMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64OfGetObject);
            assertEquals(responseCrc64, responseCrc64OfGetObject);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } finally {
            if (inputStream != null) {
                inputStream.close();
            }
        }
    }

    /***
     * 1、创桶成功
     * 2、上传对象成功
     * 3、获取对象元数据成功，返回头域带crc64值，与客户端算的一致
     *
     * @throws IOException
     */
    @Test
    public void tc_get_crc64_by_GetObjectMetadata_001() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // random test file
        long testFileSizeInBytes = getTestRandomIntInRange(5 * 1024 * 1024, 10 * 1024 * 1024);
        File testFile = genTestFile(temporaryFolder, testFileName, testFileSizeInBytes);
        String objectKey = bucketName + "objectKey_001";
        String responseCrc64 = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            // 上传对象,sdk自动计算crc64
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(true);
            putObjectRequest.setFile(testFile);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            ObjectMetadata objectMetadata = obsClient.getObjectMetadata(bucketName, objectKey);
            String responseCrc64OfGetObjectMetadata = (String) objectMetadata.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64OfGetObjectMetadata);
            assertEquals(responseCrc64, responseCrc64OfGetObjectMetadata);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，putObject上传5GB对象
     * 3、下载校验对象，返回头域带crc64值，与客户端算的一致
     * 4、不携带crc64头域，putObject上传5GB对象
     * 5、对比两次上传耗时
     *
     * @throws IOException
     */
    @Test
    public void SceneCase_1719192209901() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // random test file
        long testFileSizeInBytes = 5 * 1024 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, testFileSizeInBytes);
        String objectKey = bucketName + "objectKey_001";
        String responseCrc64;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            // 上传对象,sdk自动计算crc64
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(true);
            putObjectRequest.setFile(testFile);
            long start = System.currentTimeMillis();
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            long end = System.currentTimeMillis();
            System.out.println("crc64 put 5 gb cost " + (end - start) / 1000.0 + "seconds.");
            responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);

            ObjectMetadata objectMetadata = obsClient.getObjectMetadata(bucketName, objectKey);
            String responseCrc64OfGetObjectMetadata = (String) objectMetadata.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64OfGetObjectMetadata);
            assertEquals(responseCrc64, responseCrc64OfGetObjectMetadata);

            // 上传对象,流式计算crc64
            putObjectRequest.setUserHeaders(null);
            putObjectRequest.setNeedCalculateCRC64(false);
            CRC64InputStream crc64InputStream = new CRC64InputStream(new FileInputStream(testFile));
            putObjectRequest.setInput(crc64InputStream);
            start = System.currentTimeMillis();
            obsClient.putObject(putObjectRequest);
            end = System.currentTimeMillis();
            System.out.println("stream crc64 put 5 gb cost " + (end - start) / 1000.0 + "seconds.");
            assertEquals(responseCrc64, crc64InputStream.getCrc64().toString());

            // 上传对象,sdk不自动计算crc64
            putObjectRequest.setFile(testFile);
            start = System.currentTimeMillis();
            obsClient.putObject(putObjectRequest);
            end = System.currentTimeMillis();
            System.out.println("no crc64 put 5 gb cost " + (end - start) / 1000.0 + "seconds.");
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    // putObject with crc64 in user header
    @Test
    public void tc_putObject_with_crc64_005() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // 10 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 10 * 1024 * 1024);
        String objectKey = bucketName + "objectKey_001";
        long testFileCrc64;
        try (FileInputStream fileInputStream = new FileInputStream(testFile);
                CRC64InputStream crc64InputStream = new CRC64InputStream(fileInputStream)) {
            byte[] buffer = new byte[65536];
            while (crc64InputStream.read(buffer) != -1) {}
            testFileCrc64 = crc64InputStream.getCrc64().getValue();
        }
        String testFileCrc64UnsignedString = CRC64.toString(testFileCrc64);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.addUserHeaders("x-amz-" + HASH_CRC64ECMA, testFileCrc64UnsignedString);
            putObjectRequest.addUserHeaders("x-obs-" + HASH_CRC64ECMA, testFileCrc64UnsignedString);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            String responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    // putObject with automatic calculated crc64
    @Test
    public void tc_putObject_with_crc64_006() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // 10 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 10 * 1024 * 1024);
        String objectKey = bucketName + "objectKey_001";
        long testFileCrc64;
        try (FileInputStream fileInputStream = new FileInputStream(testFile);
                CRC64InputStream crc64InputStream = new CRC64InputStream(fileInputStream)) {
            byte[] buffer = new byte[65536];
            while (crc64InputStream.read(buffer) != -1) {}
            testFileCrc64 = crc64InputStream.getCrc64().getValue();
        }
        String testFileCrc64UnsignedString = CRC64.toString(testFileCrc64);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.setNeedCalculateCRC64(true);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            String responseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(responseCrc64);
            assertEquals(responseCrc64, testFileCrc64UnsignedString);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    // test crc64 inputStream
    @Test
    public void tc_crc64() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // 10 kb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 10 * 1024);
        long testFileCrc64;
        try (FileInputStream fileInputStream = new FileInputStream(testFile);
                CRC64InputStream crc64InputStream = new CRC64InputStream(fileInputStream)) {
            byte[] buffer = new byte[65536];
            while (crc64InputStream.read(buffer) != -1) {}
            testFileCrc64 = crc64InputStream.getCrc64().getValue();
        }
        String testFileCrc64UnsignedString = CRC64.toString(testFileCrc64);
        System.out.println("testFileCrc64 is         " + testFileCrc64);
        System.out.println("testFileCrc64Unsigned is " + testFileCrc64UnsignedString);
        try (CRC64InputStream crc64InputStream = new CRC64InputStream(new FileInputStream(testFile))) {
            while (crc64InputStream.read() != -1) {}
            Assert.assertEquals(testFileCrc64, crc64InputStream.getCrc64().getValue());
            Assert.assertEquals(testFileCrc64UnsignedString, CRC64.toString(crc64InputStream.getCrc64().getValue()));
        }
    }

    /***
     * TestInputStream is used for simulating the situation when skip returns not equal to parameter
     */
    class TestInputStream extends FileInputStream {

        public TestInputStream(File file) throws FileNotFoundException {
            super(file);
        }

        @Override
        public long skip(long n) {
            return 0;
        }
    }
    // test crc64 inputStream
    @Test
    public void tc_crc64_InputStreamSkipFailed() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // 10 kb test file
        long fileSizeInBytes = 10 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        try (TestInputStream crc64InputStream = new TestInputStream(testFile);
             CRC64InputStream crc64InputStream1 = new CRC64InputStream(new FileInputStream(testFile))) {
            assertEquals(crc64InputStream1.skip(fileSizeInBytes), fileSizeInBytes);
            CRC64.fromInputStream(crc64InputStream, fileSizeInBytes * 1024, 1);
            fail("InputStream skip size not equal to parameter, should throw exception");
        } catch (IOException e) {
            Assert.assertTrue(e.toString().contains("Failed to skip the input stream to the specified"));
        }
    }

    /***
     * 1、用obs api创建桶，初始化上传段任务
     * 2、带crc64头域为正确的crc64值，上传全部多段
     * 3、本地通过所有的分段crc64计算整个对象的crc64
     * 4、带整个对象的crc64进行合并段，合并成功，返回头域带crc64值，与客户端算的一致
     * 5、下载对象成功，返回头域带crc64值，与客户端算的一致
     *
     * @throws IOException
     */
    @Test
    public void tc_uploadPart_with_crc64_002() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        // 1000 Mb test file for multipart test
        long fileSizeInBytes = 1000 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        File testFileGet = genTestFile(temporaryFolder, testFileName + "_Get", 0L);
        String objectKey = bucketName + "_objectKey_001";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            InitiateMultipartUploadRequest initiateMultipartUploadRequest =
                    new InitiateMultipartUploadRequest(bucketName, objectKey);
            InitiateMultipartUploadResult initiateMultipartUploadResult =
                    obsClient.initiateMultipartUpload(initiateMultipartUploadRequest);
            String uploadId = initiateMultipartUploadResult.getUploadId();
            long partSize = 1024 * 1024L; // 每段上传1MB
            // 初始化线程池
            ExecutorService executorService = Executors.newFixedThreadPool(32);
            // 计算需要上传的段数
            long partCount =
                    fileSizeInBytes % partSize == 0 ? fileSizeInBytes / partSize : fileSizeInBytes / partSize + 1;
            final List<PartEtag> partETags = Collections.synchronizedList(new ArrayList<>());
            final List<CRC64> uploadPartResultClientCalculatedCRC64s =
                    Collections.synchronizedList(new ArrayList<>((int) partCount));

            for (int i = 0; i < partCount; i++) {
                uploadPartResultClientCalculatedCRC64s.add(null);
            }
            // 执行并发上传段
            for (int i = 0; i < partCount; i++) {
                // 分段在文件中的起始位置
                final long offset = i * partSize;
                // 分段大小
                final long currPartSize = (i + 1 == partCount) ? fileSizeInBytes - offset : partSize;
                // 分段号
                final int partNumber = i + 1;
                executorService.execute(
                        () -> {
                            UploadPartRequest uploadPartRequest = new UploadPartRequest();
                            uploadPartRequest.setBucketName(bucketName);
                            uploadPartRequest.setObjectKey(objectKey);
                            uploadPartRequest.setUploadId(uploadId);
                            uploadPartRequest.setFile(testFile);
                            uploadPartRequest.setPartSize(currPartSize);
                            uploadPartRequest.setOffset(offset);
                            uploadPartRequest.setPartNumber(partNumber);
                            uploadPartRequest.setNeedCalculateCRC64(true);
                            if (getPipeLineTestSecureRandom().nextBoolean()) {
                                // 随机添加进度回调，测试进度回调是否影响crc64计算
                                uploadPartRequest.setProgressListener(
                                        status -> {
                                            // 获取上传平均速率
                                            System.out.println("AverageSpeed:" + status.getAverageSpeed());
                                            // 获取上传进度百分比
                                            System.out.println("TransferPercentage:" + status.getTransferPercentage());
                                        });
                                uploadPartRequest.setProgressInterval(currPartSize / 10);
                            }
                            UploadPartResult uploadPartResult;
                            try {
                                uploadPartResult = obsClient.uploadPart(uploadPartRequest);
                                System.out.println("Part#" + partNumber + " done\n");
                                partETags.add(
                                        new PartEtag(uploadPartResult.getEtag(), uploadPartResult.getPartNumber()));
                                uploadPartResultClientCalculatedCRC64s.set(
                                        partNumber - 1, uploadPartResult.getClientCalculatedCRC64());
                                assertEquals(uploadPartResult.getClientCalculatedCRC64().toString(),
                                        uploadPartResult.getResponseHeaders().get(HASH_CRC64ECMA));
                            } catch (ObsException e) {
                                printObsException(e);
                            }
                        });
            }
            // 等待上传完成
            executorService.shutdown();
            while (!executorService.isTerminated()) {
                try {
                    executorService.awaitTermination(5, TimeUnit.SECONDS);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
            CRC64 crc64Total = new CRC64(uploadPartResultClientCalculatedCRC64s.get(0));
            // 合并所有crc64
            for (int i = 1; i < partCount; i++) {
                CRC64 uploadPartResultClientCalculatedCRC64I = uploadPartResultClientCalculatedCRC64s.get(i);
                assertNotNull(
                        i + " uploadPartResultClientCalculatedCRC64 should not be null.",
                        uploadPartResultClientCalculatedCRC64I);
                long currentPartSize =
                        (fileSizeInBytes % partSize == 0)
                                ? partSize
                                : ((i + 1 == partCount) ? fileSizeInBytes % partSize : partSize);
                crc64Total.combineWithAnotherCRC64(uploadPartResultClientCalculatedCRC64I, currentPartSize);
            }
            // 合并段
            CompleteMultipartUploadRequest completeMultipartUploadRequest =
                    new CompleteMultipartUploadRequest(bucketName, objectKey, uploadId, partETags);
            String completeMultipartUploadRequestCrc64 = CRC64.toString(crc64Total.getValue());
            completeMultipartUploadRequest.addUserHeaders(
                    "x-obs-" + HASH_CRC64ECMA, completeMultipartUploadRequestCrc64);
            completeMultipartUploadRequest.addUserHeaders(
                    "x-amz-" + HASH_CRC64ECMA, completeMultipartUploadRequestCrc64);
            CompleteMultipartUploadResult completeMultipartUploadResult =
                    obsClient.completeMultipartUpload(completeMultipartUploadRequest);
            assertNotNull(completeMultipartUploadResult);
            Map<String, Object> completeMultipartUploadResultResponseHeaders =
                    completeMultipartUploadResult.getResponseHeaders();
            assertNotNull(completeMultipartUploadResultResponseHeaders);
            String completeMultipartUploadResultCRC64 =
                    (String) completeMultipartUploadResultResponseHeaders.get(HASH_CRC64ECMA);
            assertNotNull(completeMultipartUploadResultCRC64);
            assertEquals(completeMultipartUploadResultCRC64, completeMultipartUploadRequestCrc64);
            System.out.println("completeMultipartUploadResultCRC64:" + completeMultipartUploadResultCRC64);
            DownloadFileResult downloadFileResult =
                    downloadFileWithRetry(obsClient, bucketName, objectKey, testFileGet, partSize * 10);
            assertNotNull(downloadFileResult);
            String downloadFileResultCRC64 =
                    (String) downloadFileResult.getObjectMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertEquals(downloadFileResultCRC64, completeMultipartUploadRequestCrc64);
            assertEquals(completeMultipartUploadResultCRC64, downloadFileResultCRC64);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    // test crc64 not supported situation
    @Test
    public void tc_crc64_not_supported() {
        try {
            WriteFileRequest writeFileRequest = new WriteFileRequest("", "");
            writeFileRequest.isNeedCalculateCRC64();
            writeFileRequest.setNeedCalculateCRC64(true);
            fail();
        } catch (IllegalArgumentException e) {
        }
        try {
            NewFileRequest writeFileRequest = new NewFileRequest("", "");
            writeFileRequest.isNeedCalculateCRC64();
            writeFileRequest.setNeedCalculateCRC64(true);
            fail();
        } catch (IllegalArgumentException e) {
        }
        try {
            ModifyObjectRequest writeFileRequest = new ModifyObjectRequest("", "");
            writeFileRequest.isNeedCalculateCRC64();
            writeFileRequest.setNeedCalculateCRC64(true);
            fail();
        } catch (IllegalArgumentException e) {
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，append上传对象成功testAppendTime次
     * 3、下载校验对象成功，返回头域带crc64值，与客户端算的一致
     *
     * @throws IOException
     */
    @Test
    public void tc_appendObject_with_crc64_002() throws IOException {
        tc_appendObjects_with_crc64_with_progress(null);
    }

    // 整个对象复制，返回复制目标对象的crc64，不是源对象的crc64
    @Test
    public void tc_copy_object_with_crc64_01() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "objectKey_001";
        String copyObjectKey = bucketName + "copyObjectKey_001";
        String testFileName = bucketName + "testFile";
        long fileSizeInBytes = 10 * 1024 * 1024L;
        // 10 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        try (ObsClient obsClient = getPipelineEnvironment()) {
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.setNeedCalculateCRC64(true);
            assertNotNull(obsClient);
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            String putObjectResultCRC64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(putObjectResultCRC64);
            System.out.println("putObjectResultCRC64:" + putObjectResultCRC64);
            ObjectMetadata objectMetadata = obsClient.getObjectMetadata(bucketName, objectKey);
            String objectMetadataCRC64 = (String) objectMetadata.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(objectMetadataCRC64);
            System.out.println("objectMetadataCRC64:" + objectMetadataCRC64);
            assertEquals(objectMetadataCRC64, putObjectResultCRC64);

            // 整个对象复制，返回复制目标对象的crc64，不是源对象的crc64
            CopyObjectRequest copyObjectRequest =
                    new CopyObjectRequest(bucketName, objectKey, bucketName, copyObjectKey);
            CopyObjectResult copyObjectResult = obsClient.copyObject(copyObjectRequest);
            System.out.println("copyObjectResultCRC64:" + copyObjectResult.getCRC64());
            assertEquals(objectMetadataCRC64, copyObjectResult.getCRC64());
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    @Test
    public void tc_copy_part_with_crc64_01() throws IOException {
        long start = System.currentTimeMillis();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        // 5 gb test file for multipart test
        long fileSizeInBytes = 50 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        long end = System.currentTimeMillis();
        System.out.println("genTestFile cost seconds:" + ((end - start) / 1000.0));
        String objectKey = bucketName + "_objectKey_001";
        String copyObjectKey = bucketName + "_copyObjectKey_001";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            long partSize = 10 * 1024 * 1024L; // 每段上传10MB
            int uploadFileTaskNum = 32;
            UploadFileRequest uploadFileRequest =
                    new UploadFileRequest(bucketName, objectKey, testFile.getPath(), partSize, uploadFileTaskNum, true);
            uploadFileRequest.setNeedCalculateCRC64(true);
            start = System.currentTimeMillis();
            CompleteMultipartUploadResult completeMultipartUploadResult = obsClient.uploadFile(uploadFileRequest);
            end = System.currentTimeMillis();
            System.out.println("uploadFile cost seconds:" + ((end - start) / 1000.0));
            assertNotNull(completeMultipartUploadResult);
            Map<String, Object> completeMultipartUploadResultResponseHeaders =
                    completeMultipartUploadResult.getResponseHeaders();
            assertNotNull(completeMultipartUploadResultResponseHeaders);
            String completeMultipartUploadResultCRC64 =
                    (String) completeMultipartUploadResultResponseHeaders.get(HASH_CRC64ECMA);
            assertNotNull(completeMultipartUploadResultCRC64);
            System.out.println("completeMultipartUploadResultCRC64:" + completeMultipartUploadResultCRC64);

            InitiateMultipartUploadRequest initiateMultipartUploadRequest =
                    new InitiateMultipartUploadRequest(bucketName, copyObjectKey);
            InitiateMultipartUploadResult initiateMultipartUploadResult =
                    obsClient.initiateMultipartUpload(initiateMultipartUploadRequest);
            String uploadId = initiateMultipartUploadResult.getUploadId();
            CopyPartRequest copyPartRequest = new CopyPartRequest();
            copyPartRequest.setUploadId(uploadId);
            copyPartRequest.setSourceBucketName(bucketName);
            copyPartRequest.setSourceObjectKey(objectKey);
            copyPartRequest.setDestinationBucketName(bucketName);
            copyPartRequest.setDestinationObjectKey(copyObjectKey);
            copyPartRequest.setByteRangeStart(0L);
            copyPartRequest.setByteRangeEnd(100 * 1024L);
            copyPartRequest.setPartNumber(1);
            CopyPartResult copyPartResult = obsClient.copyPart(copyPartRequest);
            System.out.println(copyPartResult.getCrc64());
            // 分段复制，只有开启桶crc64开关才会计算分段的crc64，默认即使源对象有crc64也不计算
            assertNull(copyPartResult.getCrc64());

            // code coverage
            System.out.println(copyPartResult);
            XmlResponsesSaxParser.CopyPartResultHandler copyPartResultHandler = new XmlResponsesSaxParser.CopyPartResultHandler(null);
            copyPartResultHandler.endCRC64("");
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }
    @Test
    public void tc_uploadFile_with_crc64_01() throws IOException {
        long start = System.currentTimeMillis();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        // 1000 gb test file for multipart test
        long fileSizeInBytes = 1000 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        File testFileGet = genTestFile(temporaryFolder, testFileName + "_Get", 0L);
        long end = System.currentTimeMillis();
        System.out.println("genTestFile cost seconds:" + ((end - start) / 1000.0));
        String objectKey = bucketName + "_objectKey_001";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            long partSize = 1024 * 1024L; // 每段上传1MB
            int uploadFileTaskNum = 32;
            UploadFileRequest uploadFileRequest =
                    new UploadFileRequest(bucketName, objectKey, testFile.getPath(), partSize, uploadFileTaskNum, true);
            uploadFileRequest.setNeedCalculateCRC64(true);
            start = System.currentTimeMillis();
            CompleteMultipartUploadResult completeMultipartUploadResult = obsClient.uploadFile(uploadFileRequest);
            end = System.currentTimeMillis();
            System.out.println("uploadFile cost seconds:" + ((end - start) / 1000.0));
            assertNotNull(completeMultipartUploadResult);
            Map<String, Object> completeMultipartUploadResultResponseHeaders =
                    completeMultipartUploadResult.getResponseHeaders();
            assertNotNull(completeMultipartUploadResultResponseHeaders);
            String completeMultipartUploadResultCRC64 =
                    (String) completeMultipartUploadResultResponseHeaders.get(HASH_CRC64ECMA);
            assertNotNull(completeMultipartUploadResultCRC64);
            System.out.println("completeMultipartUploadResultCRC64:" + completeMultipartUploadResultCRC64);
            start = System.currentTimeMillis();
            DownloadFileResult downloadFileResult =
                    downloadFileWithRetry(obsClient, bucketName, objectKey, testFileGet, partSize * 10);
            end = System.currentTimeMillis();
            System.out.println("downloadFile cost seconds:" + ((end - start) / 1000.0));
            assertNotNull(downloadFileResult);
            String downloadFileResultCRC64 =
                    (String) downloadFileResult.getObjectMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertEquals(completeMultipartUploadResultCRC64, downloadFileResultCRC64);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    /***
     * 1、不开启crc64，开启checkpoint，断点续传上传对象，异步计算源文件md5-1
     * 2、开启checkpoint，断点续传下载对象，异步计算md5-2
     * 3、不开启crc64，不开启checkpoint，断点续传上传对象
     * 4、开启checkpoint，断点续传下载对象，异步计算md5-3
     * 5、对比 md5-1 md5-2 是否一致， md5-1 md5-3是否一致
     *
     * @throws IOException
     */
    @Test
    public void tc_uploadFile_with_no_crc64_01() throws IOException {
        long start = System.currentTimeMillis();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        // 1000 mb test file for multipart test
        long fileSizeInBytes = 1000 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        File testFileGetWithCheckpoint =
                genTestFile(temporaryFolder, testFileName + "_GetWithCheckpoint", 0L);
        File testFileGetWithoutCheckpoint =
                genTestFile(temporaryFolder, testFileName + "_GetWithoutCheckpoint", 0L);
        long end = System.currentTimeMillis();
        System.out.println("genTestFile cost seconds:" + ((end - start) / 1000.0));
        String objectKey = bucketName + "_objectKey_001";
        ExecutorService executor = Executors.newFixedThreadPool(4);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            long partSize = 1024 * 1024L; // 每段上传1MB
            int uploadFileTaskNum = 32;
            UploadFileRequest uploadFileRequest =
                    new UploadFileRequest(bucketName, objectKey, testFile.getPath(),
                            partSize, uploadFileTaskNum, true);
            uploadFileRequest.setNeedCalculateCRC64(false);
            start = System.currentTimeMillis();
            obsClient.uploadFile(uploadFileRequest);
            end = System.currentTimeMillis();
            System.out.println("uploadFileWithCheckpoint cost seconds:" + ((end - start) / 1000.0));

            start = System.currentTimeMillis();
            downloadFileWithRetry(obsClient, bucketName, objectKey, testFileGetWithCheckpoint, partSize * 10);
            end = System.currentTimeMillis();
            System.out.println("downloadFileWithCheckpoint cost seconds:" + ((end - start) / 1000.0));

            // 定义两个Future对象，用于存储异步计算的MD5值,用于后面校验上传、下载是否正常
            Future<byte[]> futureTestFileMD5 =
                    executor.submit(
                            () -> {
                                // 使用try-with-resources语句，确保文件流在使用完毕后能够被正确关闭
                                try (FileInputStream fileInputStream = new FileInputStream(testFile)) {
                                    return ServiceUtils.computeMD5Hash(fileInputStream);
                                } catch (NoSuchAlgorithmException | IOException e) {
                                    throw new RuntimeException(e);
                                }
                            });
            Future<byte[]> futureTestFileDownloadWithCheckpointMD5 =
                    executor.submit(
                            () -> {
                                // 使用try-with-resources语句，确保文件流在使用完毕后能够被正确关闭
                                try (FileInputStream fileInputStream = new FileInputStream(testFileGetWithCheckpoint)) {
                                    return ServiceUtils.computeMD5Hash(fileInputStream);
                                } catch (NoSuchAlgorithmException | IOException e) {
                                    throw new RuntimeException(e);
                                }
                            });

            uploadFileRequest.setEnableCheckpoint(false);
            start = System.currentTimeMillis();
            obsClient.uploadFile(uploadFileRequest);
            end = System.currentTimeMillis();
            System.out.println("uploadFileWithOutCheckpoint cost seconds:" + ((end - start) / 1000.0));

            start = System.currentTimeMillis();
            downloadFileWithRetry(obsClient, bucketName, objectKey, testFileGetWithoutCheckpoint, partSize * 10);
            end = System.currentTimeMillis();
            System.out.println("downloadFileWithCheckpoint cost seconds:" + ((end - start) / 1000.0));
            Future<byte[]> futureTestFileDownloadWithoutCheckpointMD5 =
                    executor.submit(
                            () -> {
                                // 使用try-with-resources语句，确保文件流在使用完毕后能够被正确关闭
                                try (FileInputStream fileInputStream = new FileInputStream(testFileGetWithoutCheckpoint)) {
                                    return ServiceUtils.computeMD5Hash(fileInputStream);
                                } catch (NoSuchAlgorithmException | IOException e) {
                                    throw new RuntimeException(e);
                                }
                            });
            byte[] testFileMD5 = futureTestFileMD5.get();
            assertTrue(areByteArraysEqual(testFileMD5, futureTestFileDownloadWithCheckpointMD5.get()));
            assertTrue(areByteArraysEqual(testFileMD5, futureTestFileDownloadWithoutCheckpointMD5.get()));
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (ExecutionException | InterruptedException e) {
            throw new RuntimeException(e);
        } finally {
            executor.shutdown();
        }
    }

    @Test
    public void tc_test_CRC64DeepCopy() {
        CRC64 crc64 = new CRC64();
        crc64.update(1);
        long oldValue = crc64.getValue();
        CRC64 crc64Clone = new CRC64(crc64);
        assertNotSame(crc64, crc64Clone);
        Assert.assertEquals(crc64.getValue(), crc64Clone.getValue());
        crc64.update(2);
        long newValue = crc64.getValue();
        assertEquals(oldValue, crc64Clone.getValue());
        assertEquals(newValue, crc64.getValue());
    }


    protected static class TestTaskCallback
            implements TaskCallback<CompleteMultipartUploadResult, UploadFileRequest>
    {
        public AtomicBoolean isAllowFailure = new AtomicBoolean(true);
        /**
         * Callback when the task is executed successfully.
         *
         * @param result Callback parameter. Generally, the return type of a specific
         *               operation is used.
         */
        @Override
        public void onSuccess(CompleteMultipartUploadResult result) {
            System.out.println(result);
        }
        /**
         * Callback when an exception is thrown during task execution.
         *
         * @param exception     Exception information
         * @param singleRequest The request that causes an exception
         */
        @Override
        public void onException(ObsException exception, UploadFileRequest singleRequest)
        {
            if (!isAllowFailure.get()) {
                printObsException(exception);
                fail();
            } else {
                System.out.println("UploadFileRequest failed, isCancelled ? " +
                        singleRequest.getCancelHandler().isCancelled());
            }
        }
    }
    protected static class TestProgressListener implements ProgressListener {
        public int progressToAbort;
        public String requestType;
        @Override
        public void progressChanged(ProgressStatus status)
        {
            int transferPercentage = status.getTransferPercentage();
            if (transferPercentage > progressToAbort) {
                progressToAbort = 101;
                throw new RuntimeException("abort due to test");
            }
            // 获取上传平均速率
            System.out.println(requestType +"AverageSpeed:" + status.getAverageSpeed());
            // 获取上传进度百分比
            System.out.println("TransferPercentage:" + transferPercentage);
        }
    }
    // 异步断点续传带crc64，设置进度回调超过一定数值时暂停
    // 继续上传，进度不能低于暂停时的数值，证明续传机制有效
    // 下载校验crc64, 也验证下载的进度能不能续传
    @Test
    public void tc_uploadFileAsync_with_crc64_01() throws IOException {
        long start = System.currentTimeMillis();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        // 1000 mb test file for multipart test
        long fileSizeInBytes = 1000 * 1024 * 1024L + getTestRandomIntInRange(-1000, 1000);
        File testFile = genTestFile(temporaryFolder, testFileName, fileSizeInBytes);
        File testFileGet = genTestFile(temporaryFolder, testFileName + "_Get", 0L);
        long end = System.currentTimeMillis();
        System.out.println("genTestFile cost seconds:" + ((end - start) / 1000.0));
        String objectKey = bucketName + "_objectKey_001";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            assert obsClientAsync != null;
            long partSize = getTestRandomIntInRange(1024 * 1024, 2 * 1024 * 1024); // 每段上传1 - 2MB
            int progressToCancel = 20; // 大于该进度时暂停上传，开始续传时进度如果不低于该进度，说明
            int uploadFileTaskNum = 32;
            UploadFileRequest uploadFileRequest =
                    new UploadFileRequest(bucketName, objectKey, testFile.getPath(), partSize,
                            uploadFileTaskNum, true);
            uploadFileRequest.setNeedCalculateCRC64(true);
            uploadFileRequest.setNeedAbortUploadFileAfterCancel(false);
            CallCancelHandler cancelHandler = new CallCancelHandler();
            cancelHandler.setMaxCallCapacity(0);
            uploadFileRequest.setCancelHandler(cancelHandler);
            uploadFileRequest.setEnableCheckpoint(true);
            uploadFileRequest.setProgressListener(
                    status -> {
                        // 获取上传平均速率
                        System.out.println("AverageSpeed:" + status.getAverageSpeed());
                        // 获取上传进度百分比
                        System.out.println("TransferPercentage:" + status.getTransferPercentage());
                        if (status.getTransferPercentage() >= progressToCancel) {
                            // 上传进度大于等于progressToCancel时暂停
                            cancelHandler.cancel();
                        }
                    });
            // 每上传100KB数据反馈上传进度
            uploadFileRequest.setProgressInterval(fileSizeInBytes / 100);
            TestTaskCallback completeCallback = new TestTaskCallback();
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest, completeCallback);
            Optional<CompleteMultipartUploadResult> resultOptional = uploadFileTask.getResult();
            if (resultOptional.isPresent()) {
                fail("resultOptional of uploadFileTask should not be present, cause request is canceled.");
            } else {
                System.out.println("resultOptional of uploadFileTask should not be present,"
                        + " cause request is canceled.");
            }
            start = System.currentTimeMillis();
            // retry upload, check if crc64 is recorded to checkpoint and reloaded from checkpoint
            completeCallback.isAllowFailure.set(false);
            uploadFileRequest.setProgressListener(
                    status -> {
                        // 获取上传平均速率
                        System.out.println("AverageSpeed:" + status.getAverageSpeed());
                        // 获取上传进度百分比
                        int transferPercentage = status.getTransferPercentage();
                        System.out.println("TransferPercentage:" + transferPercentage);
                        // 确保断点续传机制还生效，进度不能丢失
                        assertTrue(transferPercentage >= progressToCancel);
                    });
            uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest, completeCallback);
            resultOptional = uploadFileTask.getResult();
            if (!resultOptional.isPresent()) {
                fail("resultOptional of uploadFileTask is not present.");
            }
            CompleteMultipartUploadResult completeMultipartUploadResult = resultOptional.get();
            end = System.currentTimeMillis();
            System.out.println("uploadFile cost seconds:" + ((end - start) / 1000.0));
            assertNotNull(completeMultipartUploadResult);
            Map<String, Object> completeMultipartUploadResultResponseHeaders =
                    completeMultipartUploadResult.getResponseHeaders();
            assertNotNull(completeMultipartUploadResultResponseHeaders);
            String completeMultipartUploadResultCRC64 =
                    (String) completeMultipartUploadResultResponseHeaders.get(HASH_CRC64ECMA);
            assertNotNull(completeMultipartUploadResultCRC64);
            System.out.println("completeMultipartUploadResultCRC64:" + completeMultipartUploadResultCRC64);

            start = System.currentTimeMillis();
            TestProgressListener testProgressListener = new TestProgressListener();
            testProgressListener.progressToAbort = 20;
            testProgressListener.requestType = "downloadFile";
            DownloadFileResult downloadFileResult =
                    downloadFileWithRetry(obsClientAsync, bucketName, objectKey, testFileGet, partSize * 10,
                            true, true, testProgressListener, fileSizeInBytes / 100);
            end = System.currentTimeMillis();
            System.out.println("downloadFile cost seconds:" + ((end - start) / 1000.0));
            assertNotNull(downloadFileResult);
            CRC64 downloadFileResultCombinedCRC64 = downloadFileResult.getCombinedCRC64();
            assertNotNull(downloadFileResultCombinedCRC64);
            String downloadFileResultCRC64 =
                    (String) downloadFileResult.getObjectMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertEquals(completeMultipartUploadResultCRC64, downloadFileResultCRC64);
            assertEquals(downloadFileResultCombinedCRC64.toString(), downloadFileResultCRC64);
            CRC64 crc64 = CRC64.fromFile(testFileGet);
            assertEquals(crc64.toString(), downloadFileResultCRC64);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }

    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，上传1MB对象, 带进度回调
     * 3、下载校验对象
     *
     * @throws IOException
     */
    @Test
    public void tc_putObject_with_crc64_001_with_progress() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        // 100 mb test file
        long testFileSizeInBytes = 100 * 1024 * 1024L;
        File testFile = genTestFile(temporaryFolder, testFileName, testFileSizeInBytes);
        // 100 mb test file
        String objectKey = bucketName + "objectKey_001";
        String putObjectResponseCrc64;
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, objectKey, testFile);
            putObjectRequest.setNeedCalculateCRC64(true);
            putObjectRequest.setProgressListener(
                    status -> {
                        // 获取上传平均速率
                        System.out.println("AverageSpeed:" + status.getAverageSpeed());
                        // 获取上传进度百分比
                        System.out.println("TransferPercentage:" + status.getTransferPercentage());
                    });
            // 每上传10%数据反馈上传进度
            putObjectRequest.setProgressInterval(testFileSizeInBytes / 10);
            // 带crc64头域为正确的crc64值，上传1MB对象
            PutObjectResult putObjectResult = obsClient.putObject(putObjectRequest);
            putObjectResponseCrc64 = (String) putObjectResult.getResponseHeaders().get(HASH_CRC64ECMA);
            assertNotNull(putObjectResponseCrc64);
        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }


    /***
     * 1、用obs api创建桶
     * 2、带crc64头域为正确的crc64值，append带进度回调上传对象, 成功testAppendTime次
     * 3、下载校验对象成功，返回头域带crc64值，与客户端算的一致
     *
     * @throws IOException
     */
    @Test
    public void tc_appendObject_with_crc64_003() throws IOException {
        tc_appendObjects_with_crc64_with_progress(
                status -> {
                    // 获取上传平均速率
                    System.out.println("AverageSpeed:" + status.getAverageSpeed());
                    // 获取上传进度百分比
                    System.out.println("TransferPercentage:" + status.getTransferPercentage());
                });
    }
    public void tc_appendObjects_with_crc64_with_progress(ProgressListener progressListener)
            throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "_testFile";
        File testFileGet = genTestFile(temporaryFolder, testFileName + "_Get", 0L);
        String objectKey = bucketName + "_objectKey_001";
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;
            int testAppendTime = getTestRandomIntInRange(10, 20);
            System.out.println("tc_appendObject_with_crc64_002 testAppendTime:" + testAppendTime);
            // random test file size for appendObject test
            long fileSizeInBytes = getTestRandomIntInRange(2 * 1024 * 1024, 10 * 1024 * 1024);
            AppendObjectRequest appendObjectRequest = new AppendObjectRequest(bucketName);
            appendObjectRequest.setObjectKey(objectKey);
            appendObjectRequest.setNeedCalculateCRC64(true);
            appendObjectRequest.setProgressListener(progressListener);
            appendObjectRequest.setProgressInterval(fileSizeInBytes / 10);
            CRC64 crc64Total = null;
            for (int i = 0; i < testAppendTime; ++i) {
                File testFile = genTestFile(temporaryFolder, testFileName + i, fileSizeInBytes);
                appendObjectRequest.setFile(testFile);
                AppendObjectResult appendObjectResult = obsClient.appendObject(appendObjectRequest);
                assertEquals(
                        appendObjectResult.getClientCalculatedCRC64().toString(),
                        appendObjectResult.getResponseHeaders().get(HASH_CRC64ECMA));
                appendObjectRequest.setPosition(appendObjectResult.getNextPosition());
                appendObjectRequest.setCrc64BeforeAppend(
                        (String) appendObjectResult.getResponseHeaders().get(HASH_CRC64ECMA));
                if (crc64Total == null) {
                    crc64Total = new CRC64(appendObjectResult.getClientCalculatedCRC64());
                } else {
                    crc64Total.combineWithAnotherCRC64(CRC64.fromFile(testFile), testFile.length());
                }
            }

            DownloadFileResult downloadFileResult =
                    downloadFileWithRetry(obsClient, bucketName, objectKey, testFileGet, fileSizeInBytes);
            assertNotNull(downloadFileResult);
            String downloadFileResultCRC64 =
                    (String) downloadFileResult.getObjectMetadata().getResponseHeaders().get(HASH_CRC64ECMA);
            assertEquals(downloadFileResultCRC64, crc64Total.toString());

        } catch (ObsException e) {
            printObsException(e);
            fail();
        }
    }


}
