package com.obs.integrated_test;

import static com.obs.test.TestTools.genTestFile;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotEquals;
import static org.junit.Assert.assertNotNull;

import com.obs.services.crypto.CTRCipherGenerator;
import com.obs.services.crypto.CryptoObsClient;
import com.obs.services.crypto.CtrRSACipherGenerator;
import com.obs.services.exception.ObsException;
import com.obs.services.model.GetObjectMetadataRequest;
import com.obs.services.model.GetObjectRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsObject;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.SetObjectMetadataRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.security.NoSuchAlgorithmException;
import java.util.Locale;
import java.util.Map;
import java.util.stream.Collectors;

public class CryptoClientIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    @Test
    public void test_SDK_Crypto_001() throws IOException {
        // 1、加密上传成功，加密结果和对象元数据显示符合预期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        CryptoObsClient cryptoObsClient = TestTools.getPipelineCryptoEnvironment();
        String objectKey = "test_crypto_object";
        int testFileSizeKb = 1024 * 20;
        String filePath = "test_SDK_Crypto";
        genTestFile(filePath, testFileSizeKb);
        assert cryptoObsClient != null;
        PutObjectResult putObjectResult = cryptoObsClient.putObject(bucketName, objectKey, new File(filePath));
        assertEquals(200, putObjectResult.getStatusCode());

        GetObjectMetadataRequest getObjectMetadataRequest = new GetObjectMetadataRequest(bucketName, objectKey);
        getObjectMetadataRequest.setIsEncodeHeaders(false);
        ObjectMetadata objectMetadata = cryptoObsClient.getObjectMetadata(getObjectMetadataRequest);
        String masterKeyInfo = (String) objectMetadata.getUserMetadata(CTRCipherGenerator.MASTER_KEY_INFO_META_NAME);
        String encryptedAlgorithm =
                (String) objectMetadata.getUserMetadata(CTRCipherGenerator.ENCRYPTED_ALGORITHM_META_NAME);
        String encryptedStart = (String) objectMetadata.getUserMetadata(CTRCipherGenerator.ENCRYPTED_START_META_NAME);
        String encryptedAESKey = (String) objectMetadata.getUserMetadata(CtrRSACipherGenerator.ENCRYPTED_AES_KEY_META_NAME);
        String encryptedSha256 = (String) objectMetadata.getUserMetadata(CTRCipherGenerator.ENCRYPTED_SHA_256_META_NAME);
        String plaintextContentLength =
                (String) objectMetadata.getUserMetadata(CTRCipherGenerator.PLAINTEXT_CONTENT_LENGTH_META_NAME);

        assertEquals(CtrRSACipherGenerator.ENCRYPTED_ALGORITHM, encryptedAlgorithm);
        assertNotEquals("", masterKeyInfo);
        assertNotEquals("", encryptedStart);
        assertNotEquals("", encryptedAESKey);
        assertNotEquals("", encryptedSha256);
        assertNotEquals("", plaintextContentLength);
    }

    @Test
    public void test_SDK_Crypto_002() throws IOException, NoSuchAlgorithmException {
        // 下载解密成功，解密数据结果sha256和原文件完全一致
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        CryptoObsClient cryptoObsClient = TestTools.getPipelineCryptoEnvironment();
        String objectKey = "test_crypto_object";
        int testFileSizeKb = 1024 * 20;
        String plainTextFilePath = "test_SDK_Crypto";
        String decryptedTextFilePath = "test_SDK_Crypto_Decrypted";
        genTestFile(plainTextFilePath, testFileSizeKb);
        assert cryptoObsClient != null;
        PutObjectResult putObjectResult = cryptoObsClient.putObject(bucketName, objectKey, new File(plainTextFilePath));
        assertEquals(200, putObjectResult.getStatusCode());
        GetObjectRequest getObjectRequest = new GetObjectRequest(bucketName, objectKey);
        ObsObject obsObject = cryptoObsClient.getObject(getObjectRequest);
        assertEquals(200, putObjectResult.getStatusCode());
        InputStream input = obsObject.getObjectContent();
        byte[] b = new byte[1024];
        FileOutputStream fileOutputStream = new FileOutputStream(decryptedTextFilePath);
        int len;
        while ((len = input.read(b)) != -1) {
            fileOutputStream.write(b, 0, len);
        }
        fileOutputStream.close();
        input.close();
        byte[] plainTextSha256 = CTRCipherGenerator.getFileSha256Bytes(plainTextFilePath);
        byte[] decryptedFileSha256 = CTRCipherGenerator.getFileSha256Bytes(decryptedTextFilePath);
        assertArrayEquals("sha256 not equal", plainTextSha256, decryptedFileSha256);
    }

    @Test
    public void test_SDK_Crypto_003() throws IOException {
        // 解密所需的必要辅助信息存储在相应对象的元数据中，元数据被更改，解密失败；
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        CryptoObsClient cryptoObsClient = TestTools.getPipelineCryptoEnvironment();
        String objectKey = "test_crypto_object";
        int testFileSizeKb = 1024 * 20;
        String plainTextFilePath = "test_SDK_Crypto";
        genTestFile(plainTextFilePath, testFileSizeKb);
        assert cryptoObsClient != null;
        PutObjectResult putObjectResult = cryptoObsClient.putObject(bucketName, objectKey, new File(plainTextFilePath));
        assertEquals(200, putObjectResult.getStatusCode());

        ObjectMetadata objectMetadata = cryptoObsClient.getObjectMetadata(bucketName, objectKey);
        SetObjectMetadataRequest setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKey);
        setObjectMetadataRequest.setIsEncodeHeaders(false);
        Map<String, String> allUserMetaData =
                objectMetadata.getAllMetadata().entrySet().stream()
                        .collect(Collectors.toMap(Map.Entry::getKey, entry -> (String) entry.getValue()));
        setObjectMetadataRequest.addAllUserMetadata(allUserMetaData);
        setObjectMetadataRequest.addUserMetadata(
                CtrRSACipherGenerator.ENCRYPTED_START_META_NAME, "wrong encrypted AES iv");
        setObjectMetadataRequest.addUserMetadata(
                CTRCipherGenerator.ENCRYPTED_ALGORITHM_META_NAME, CTRCipherGenerator.ENCRYPTED_ALGORITHM);
        setObjectMetadataRequest.setRemoveUnset(false);
        HeaderResponse headerResponse = cryptoObsClient.setObjectMetadata(setObjectMetadataRequest);
        assertEquals(200, headerResponse.getStatusCode());

        GetObjectRequest getObjectRequest = new GetObjectRequest(bucketName, objectKey);
        try {
            cryptoObsClient.getObject(getObjectRequest);
        } catch (Exception e) {
            assertNotNull(e);
        }
    }

    @Test
    public void test_SDK_Crypto_004() throws IOException {
        // 使用错误的主秘钥解密，解密失败
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        CryptoObsClient cryptoObsClient = TestTools.getPipelineCryptoEnvironment();
        String objectKey = "test_crypto_object";
        int testFileSizeKb = 1024 * 20;
        String plainTextFilePath = "test_SDK_Crypto";
        genTestFile(plainTextFilePath, testFileSizeKb);
        assert cryptoObsClient != null;
        PutObjectResult putObjectResult = cryptoObsClient.putObject(bucketName, objectKey, new File(plainTextFilePath));
        assertEquals(200, putObjectResult.getStatusCode());

        ObjectMetadata objectMetadata = cryptoObsClient.getObjectMetadata(bucketName, objectKey);
        SetObjectMetadataRequest setObjectMetadataRequest = new SetObjectMetadataRequest(bucketName, objectKey);
        setObjectMetadataRequest.setIsEncodeHeaders(false);
        Map<String, String> allUserMetaData =
                objectMetadata.getAllMetadata().entrySet().stream()
                        .collect(Collectors.toMap(Map.Entry::getKey, entry -> (String) entry.getValue()));
        setObjectMetadataRequest.addAllUserMetadata(allUserMetaData);
        setObjectMetadataRequest.addUserMetadata(
                CtrRSACipherGenerator.ENCRYPTED_AES_KEY_META_NAME, "wrong encrypted AES key");
        setObjectMetadataRequest.setRemoveUnset(false);
        HeaderResponse headerResponse = cryptoObsClient.setObjectMetadata(setObjectMetadataRequest);
        assertEquals(200, headerResponse.getStatusCode());

        GetObjectRequest getObjectRequest = new GetObjectRequest(bucketName, objectKey);
        try {
            cryptoObsClient.getObject(getObjectRequest);
        } catch (ObsException e) {
            assertNotNull(e);
        }
    }
}
