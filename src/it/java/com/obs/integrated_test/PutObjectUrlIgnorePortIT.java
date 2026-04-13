package com.obs.integrated_test;
import com.obs.test.TestTools;

import com.obs.services.ObsClient;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.UploadFileRequest;
import com.obs.test.tools.PrepareTestBucket;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Locale;

import static com.obs.test.TestTools.genTestFile;
import static org.junit.Assert.assertEquals;

public class PutObjectUrlIgnorePortIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    @Test
    public void test_put_object_url_ignore_port() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String objectKey = "test_putObject_001";

        PutObjectRequest request = new PutObjectRequest();
        request.setObjectKey(objectKey);
        request.setBucketName(bucketName);
        request.setInput(new ByteArrayInputStream("testObject".getBytes(StandardCharsets.UTF_8)));
        PutObjectResult putResult = obsClient.putObject(request);

        assertEquals(true, putResult.getObjectUrl().contains("443") || putResult.getObjectUrl().contains("80"));

        request.setIsIgnorePort(true);
        putResult = obsClient.putObject(request);
        assertEquals(false, putResult.getObjectUrl().contains("443") || putResult.getObjectUrl().contains("80"));
    }

    @Test
    public void test_upload_file_url_ignore_port() throws IOException {
        genTestFile("test_uploadFile_with_urlIgnorePort_001", 1024 * 100);
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String objectKey = "test_putObject_002";

        UploadFileRequest request = new UploadFileRequest(bucketName, objectKey);
        request.setUploadFile("test_uploadFile_with_urlIgnorePort_001");
        CompleteMultipartUploadResult upload_result = obsClient.uploadFile(request);

        assertEquals(true, upload_result.getObjectUrl().contains("443") || upload_result.getObjectUrl().contains("80"));

        request.setIsIgnorePort(true);
        upload_result = obsClient.uploadFile(request);
        assertEquals(false, upload_result.getObjectUrl().contains("443") || upload_result.getObjectUrl().contains("80"));
    }
}
