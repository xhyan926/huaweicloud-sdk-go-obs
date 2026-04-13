package com.obs.integrated_test;
import com.obs.test.TestTools;

import com.obs.services.ObsClient;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.fs.ContentSummaryFsRequest;
import com.obs.services.model.fs.ContentSummaryFsResult;
import com.obs.services.model.fs.ListContentSummaryFsRequest;
import com.obs.services.model.fs.ListContentSummaryFsResult;
import com.obs.test.tools.PrepareTestBucket;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.ByteArrayInputStream;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

import static org.junit.Assert.assertEquals;

public class ContentSummarySupIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    @Test
    public void test_multi_list_content_summary() throws InterruptedException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)+"pfs";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketName);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        HeaderResponse response1 = obsClient.createBucket(createBucketRequest);
        Assert.assertEquals(200, response1.getStatusCode());

        try
        {
            ListContentSummaryFsRequest listContentSummaryFsRequest = new ListContentSummaryFsRequest();
            listContentSummaryFsRequest.setBucketName(bucketName);
            listContentSummaryFsRequest.setMaxKeys(20);
            List<ListContentSummaryFsRequest.DirLayer> dirLayers = new ArrayList<>();

            ListContentSummaryFsRequest.DirLayer dir_test = new ListContentSummaryFsRequest.DirLayer();
            dir_test.setKey("test_directory001/");
            dir_test.setInode(123L);
            dirLayers.add(dir_test);
            dir_test = new ListContentSummaryFsRequest.DirLayer();
            dir_test.setKey("test_directory002/");
            dirLayers.add(dir_test);
            listContentSummaryFsRequest.setDirLayers(dirLayers);
            ListContentSummaryFsResult listContentSummaryFsResult = obsClient.listContentSummaryFs(listContentSummaryFsRequest);
            assertEquals(123L, listContentSummaryFsResult.getDirContentSummaries().get(0).getInode());
            assertEquals("test_directory002/", listContentSummaryFsResult.getErrorResults().get(0).getKey());
            assertEquals("404", listContentSummaryFsResult.getErrorResults().get(0).getStatusCode());
            assertEquals("NoSuchKey", listContentSummaryFsResult.getErrorResults().get(0).getErrorCode());

            // 创建文件夹预埋数据
            obsClient.putObject(bucketName, "test_directory001/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory002/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/test_subsub001/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/testkey1",
                    new ByteArrayInputStream("testObject1".getBytes(StandardCharsets.UTF_8)));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/testkey2",
                    new ByteArrayInputStream("testObject2".getBytes(StandardCharsets.UTF_8)));
            obsClient.putObject(bucketName, "test_directory001/testkey4",
                    new ByteArrayInputStream("testObject5".getBytes(StandardCharsets.UTF_8)));
            Thread.sleep(10000L);
            dirLayers.clear();
            dir_test = new ListContentSummaryFsRequest.DirLayer();
            dir_test.setKey("test_directory001/");
            dirLayers.add(dir_test);

            listContentSummaryFsRequest.setDirLayers(dirLayers);
            listContentSummaryFsResult = obsClient.listContentSummaryFs(listContentSummaryFsRequest);
            assertEquals("test_directory001/", listContentSummaryFsResult.getDirContentSummaries().get(0).getKey());
            assertEquals(true, listContentSummaryFsResult.getDirContentSummaries().get(0).getInode() != 0);
            assertEquals("test_directory001/test_sub001/", listContentSummaryFsResult.getDirContentSummaries().get(0).getSubDir().get(0).getName());
            assertEquals(1L, listContentSummaryFsResult.getDirContentSummaries().get(0).getSubDir().get(0).getDirCount());
            assertEquals(2L, listContentSummaryFsResult.getDirContentSummaries().get(0).getSubDir().get(0).getFileCount());
            assertEquals(22L, listContentSummaryFsResult.getDirContentSummaries().get(0).getSubDir().get(0).getFileSize());
            assertEquals(true, listContentSummaryFsResult.getDirContentSummaries().get(0).getSubDir().get(0).getInode() != 0);
        }
        finally
        {

            obsClient.deleteObject(bucketName,"test_directory002/");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/test_subsub001/");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/testkey1");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/testkey2");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/");
            obsClient.deleteObject(bucketName,"test_directory001/testkey4");
            obsClient.deleteObject(bucketName,"test_directory001/");
            obsClient.deleteBucket(bucketName);
        }
    }

    @Test
    public void test_get_content_summary() throws InterruptedException {
        String bucketName = "test-get-content-summary-pfs";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketName);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);
        HeaderResponse response1 = obsClient.createBucket(createBucketRequest);
        Assert.assertEquals(200, response1.getStatusCode());
        try
        {
            obsClient.putObject(bucketName, "test_directory001/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory002/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/test_subsub001/", new ByteArrayInputStream(new byte[0]));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/testkey1",
                    new ByteArrayInputStream("testObject1".getBytes(StandardCharsets.UTF_8)));
            obsClient.putObject(bucketName, "test_directory001/test_sub001/testkey2",
                    new ByteArrayInputStream("testObject2".getBytes(StandardCharsets.UTF_8)));
            obsClient.putObject(bucketName, "test_directory001/testkey4",
                    new ByteArrayInputStream("testObject5".getBytes(StandardCharsets.UTF_8)));
            ContentSummaryFsRequest contentSummaryFsRequest = new ContentSummaryFsRequest();
            contentSummaryFsRequest.setBucketName(bucketName);
            contentSummaryFsRequest.setDirName("test_directory001/");
            Thread.sleep(15000L);
            ContentSummaryFsResult contentSummaryFsResult = obsClient.getContentSummaryFs(contentSummaryFsRequest);
            assertEquals("test_directory001/", contentSummaryFsResult.getContentSummary().getName());
            assertEquals(1L, contentSummaryFsResult.getContentSummary().getDirCount());
            assertEquals(1L, contentSummaryFsResult.getContentSummary().getFileCount());
            assertEquals(11L, contentSummaryFsResult.getContentSummary().getFileSize());
            assertEquals(true, contentSummaryFsResult.getContentSummary().getInode() != 0);
        }
        finally
        {
            obsClient.deleteObject(bucketName,"test_directory002/");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/test_subsub001/");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/testkey1");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/testkey2");
            obsClient.deleteObject(bucketName,"test_directory001/test_sub001/");
            obsClient.deleteObject(bucketName,"test_directory001/testkey4");
            obsClient.deleteObject(bucketName,"test_directory001/");
            obsClient.deleteBucket(bucketName);
        }
    }
}
