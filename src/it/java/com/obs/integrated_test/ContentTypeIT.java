package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsObject;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.util.Locale;

import static com.obs.test.TestTools.genTestFile;

public class ContentTypeIT {

    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    protected static final String FILE_TYPE_WEBP = ".webp";
    protected static final String CONTENT_TYPE_WEBP = "image/webp";

    @Test
    public void tc_putObject_with_webpContentType() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        String objectKeyWebp = bucketName + "-objectKey" + FILE_TYPE_WEBP;
        String testFileName = bucketName + "testFile";
        // 1 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 1024 * 1024);
        File testFileWebp = genTestFile(temporaryFolder, testFileName + FILE_TYPE_WEBP, 1024 * 1024);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment_OBS()) {
            assert obsClient != null;

            // 对象名和文件名不带webp后缀，对象ContentType元数据不是image/webp
            obsClient.putObject(bucketName, objectKey, testFile);
            ObjectMetadata objectMetadata = obsClient.getObjectMetadata(bucketName, objectKey);
            Assert.assertNotEquals(CONTENT_TYPE_WEBP, objectMetadata.getContentType());

            // 对象名带webp后缀，文件名不带webp后缀，对象ContentType元数据是image/webp
            obsClient.putObject(bucketName, objectKeyWebp, testFile);
            objectMetadata = obsClient.getObjectMetadata(bucketName, objectKeyWebp);
            Assert.assertEquals(CONTENT_TYPE_WEBP, objectMetadata.getContentType());

            // 对象名不带webp后缀，文件名带webp后缀，对象ContentType元数据是image/webp
            obsClient.putObject(bucketName, objectKey, testFileWebp);
            objectMetadata = obsClient.getObjectMetadata(bucketName, objectKey);
            Assert.assertEquals(CONTENT_TYPE_WEBP, objectMetadata.getContentType());

            // 对象名和文件名均带webp后缀，对象ContentType元数据是image/webp
            obsClient.putObject(bucketName, objectKeyWebp, testFileWebp);
            objectMetadata = obsClient.getObjectMetadata(bucketName, objectKeyWebp);
            Assert.assertEquals(CONTENT_TYPE_WEBP, objectMetadata.getContentType());
        }
    }
}
