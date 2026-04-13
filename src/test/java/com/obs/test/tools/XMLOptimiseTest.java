package com.obs.test.tools;

import static com.obs.test.TestTools.genTestFile;
import static com.obs.test.TestTools.printObsException;
import static org.junit.Assert.fail;

import com.obs.services.ObsClientAsync;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.task.UploadFileTask;
import com.obs.services.internal.utils.CRC64;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.ListBucketsRequest;
import com.obs.services.model.ListBucketsResult;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.TaskCallback;
import com.obs.services.model.UploadFileRequest;
import com.obs.test.TestTools;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Optional;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

public class XMLOptimiseTest
{
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    // uploadFile 接口 包含了CompleteMultipartUploadXMLBuilder
    @Test
    public void testUploadFileWithXmlBuildOptimised001() {
        // 1、加密上传成功，加密结果和对象元数据显示符合预期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "-objectKey";
        String testFileName = bucketName + "testFile";
        // test file size in bytes
        long testFileSizeInBytes = 500 * 100 * 1024;
        int testTimes = 4;
        int testThreads = 4;
        List<UploadFileTask> uploadFileTasks = new ArrayList<>();
        AtomicBoolean isSuccess = new AtomicBoolean(true);
        TaskCallback<CompleteMultipartUploadResult, UploadFileRequest> completeCallback =
                new TaskCallback<CompleteMultipartUploadResult, UploadFileRequest>() {
                    @Override
                    public void onSuccess(CompleteMultipartUploadResult result) {
                        // 失败一次视为全部失败
                        if (isSuccess.get()) {
                            isSuccess.set(true);
                        }
                    }

                    @Override
                    public void onException(ObsException exception, UploadFileRequest singleRequest) {
                        // 失败一次, 打印报错
                        isSuccess.set(false);
                        printObsException(exception);
                    }
                };
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            File testFile = genTestFile(temporaryFolder, testFileName, testFileSizeInBytes);
            assert obsClientAsync != null;
            obsClientAsync.setExecutorService(Executors.newFixedThreadPool(4,
                    r -> new Thread(r, "async-obs-thread")));
            CRC64 crc64OfLocalFile = CRC64.fromFile(testFile);
            System.out.println("testUploadFileWithXmlBuildOptimised001 crc64OfLocalFile:" + crc64OfLocalFile);
            for (int j = 0; j < testTimes; j++)
            {
                uploadFileTasks.clear();
                for (int i = 0; i < testThreads; i++)
                {
                    UploadFileRequest uploadFileRequest =
                            new UploadFileRequest(bucketName, objectKey + i, testFile.getPath(),
                                    100 * 1024, 4, true);
                    uploadFileRequest.setNeedCalculateCRC64(true);
                    uploadFileTasks.add(obsClientAsync.uploadFileAsync(uploadFileRequest, completeCallback));
                }
                for (int i = 0; i < testThreads; i++)
                {
                    UploadFileTask uploadFileTask = uploadFileTasks.get(i);
                    Optional<CompleteMultipartUploadResult> resultOptional = uploadFileTask.getResult();
                    System.out.println("end test upload thread " + i);
                    Assert.assertTrue(isSuccess.get());
                    Assert.assertTrue(resultOptional.isPresent());
                    Assert.assertEquals(resultOptional.get().getStatusCode(), 200);
                    ObjectMetadata objectMetadata = obsClientAsync.getObjectMetadata(bucketName, objectKey + i);
                    Assert.assertEquals(crc64OfLocalFile.toString(), objectMetadata.getCrc64());
                }
            }
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

    }

    @Test
    public void testOptimisedXmlParse_ListBuckets() {
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            assert obsClientAsync != null;
            ListBucketsRequest listBucketsRequest = new ListBucketsRequest();
            ExecutorService executorService = Executors.newFixedThreadPool(4,
                    r -> new Thread(r, "async-obs-thread"));
            final int testTime = 8;
            final int testThreads = 8;
            AtomicBoolean isAllSuccess = new AtomicBoolean(true);
            for (int i = 0; i < testThreads; i++) {
                executorService.submit(() -> {
                    for (int j = 0; j < testTime; j++){
                        try {
                            List<ObsBucket> buckets = obsClientAsync.listBuckets(listBucketsRequest);
                            ListBucketsResult listBucketsResult = obsClientAsync.listBucketsV2(listBucketsRequest);
                            Assert.assertEquals(buckets.size(), listBucketsResult.getBuckets().size());
                            for (int k = 0; k < buckets.size(); k++)
                            {
                                Assert.assertEquals(buckets.get(k).getBucketName(),
                                        listBucketsResult.getBuckets().get(k).getBucketName());
                                Assert.assertEquals(buckets.get(k).getLocation(),
                                        listBucketsResult.getBuckets().get(k).getLocation());
                            }
                        } catch (ObsException e) {
                            printObsException(e);
                            if (isAllSuccess.get()) {
                                isAllSuccess.set(false);
                            }
                        }
                    }
                });
            }
            executorService.shutdown();
            executorService.awaitTermination(30 * 1000, TimeUnit.MILLISECONDS);
            Assert.assertTrue(isAllSuccess.get());
        } catch (ObsException e) {
            printObsException(e);
            fail();
        } catch (IOException e) {
            throw new RuntimeException(e);
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }
    }
}
