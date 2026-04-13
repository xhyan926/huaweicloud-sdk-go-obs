package com.obs.integrated_test;

import com.obs.services.BasicObsCredentialsProvider;
import com.obs.services.IObsCredentialsProvider;
import com.obs.services.ObsClientAsync;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.task.UploadFileTask;
import com.obs.services.internal.utils.CallCancelHandler;
import com.obs.services.model.BucketMetadataInfoRequest;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.ListMultipartUploadsRequest;
import com.obs.services.model.MultipartUploadListing;
import com.obs.services.model.TaskCallback;
import com.obs.services.model.UploadFileRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import okhttp3.Call;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import static com.obs.test.TestTools.genTestFile;
import static com.obs.test.TestTools.removeTestFile;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertTrue;

public class CancelUploadFileIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    private static final String TEST_FILE_NORMAL = "testFileNormalInCancelUploadFileTest";
    private static final int KB_SIZEOF_TEST_FILE_NORMAL = 1024 * 100;

    private TaskCallback<CompleteMultipartUploadResult, UploadFileRequest> completeCallback =
            new TaskCallback<CompleteMultipartUploadResult, UploadFileRequest>() {
                @Override
                public void onSuccess(CompleteMultipartUploadResult result) {
                    System.out.println("uploadFileAsync Successfully！");
                    System.out.println("HTTP StatusCode:" + result.getStatusCode());
                    System.out.println("ObjectUrl:" + result.getObjectUrl());
                    System.out.println("Etag:" + result.getEtag());
                }

                @Override
                public void onException(ObsException e, UploadFileRequest singleRequest) {
                    System.out.println("uploadFileAsync failed");
                    if (singleRequest.getCancelHandler() != null) {
                        System.out.println(
                                "UploadFileRequest isCanceled ? " + singleRequest.getCancelHandler().isCancelled());
                    }
                    // 请求失败,打印http状态码
                    System.out.println("HTTP Code:" + e.getResponseCode());
                    // 请求失败,打印服务端错误码
                    System.out.println("Error Code:" + e.getErrorCode());
                    // 请求失败,打印详细错误信息
                    System.out.println("Error Message:" + e.getErrorMessage());
                    // 请求失败,打印请求id
                    System.out.println("Request ID:" + e.getErrorRequestId());
                    System.out.println("Host ID:" + e.getErrorHostId());
                    // 遍历Map的entry,打印所有报错相关头域
                    Map<String, String> headers = e.getResponseHeaders();
                    if (headers != null) {
                        for (Map.Entry<String, String> header : headers.entrySet()) {
                            if (header.getKey().contains("error")) {
                                System.out.println(
                                        "errorHeaderKey:"
                                                + header.getKey()
                                                + ", errorHeaderValue:"
                                                + header.getValue());
                            }
                        }
                    }
                    try (StringWriter sw = new StringWriter();
                            PrintWriter pw = new PrintWriter(sw)) {
                        e.printStackTrace(pw);
                        System.out.println("StackTrace is: " + sw);
                    } catch (IOException ex) {
                        System.out.println("Failed to print StackTrace.");
                    }
                }
            };

    // 暂停上传成功，但不取消上传任务
    @Test
    public void test_cancel_upload_01() throws IOException {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_01";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            CallCancelHandler callCancelHandler1 = new CallCancelHandler();
            uploadFileRequest1.setCancelHandler(callCancelHandler1);
            // 每上传10MB数据反馈上传进度
            uploadFileRequest1.setProgressInterval(10 * 1024 * 1024L);
            uploadFileRequest1.setEnableCheckpoint(true);
            // 暂停上传时，自动取消分段上传任务
            uploadFileRequest1.setNeedAbortUploadFileAfterCancel(false);
            uploadFileRequest1.setProgressListener(
                    status -> {
                        if (status.getTransferPercentage() > 20) {
                            // 当上传进度超过20时，取消上传
                            callCancelHandler1.cancel();
                        }
                    });
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();
            // upload canceled
            assertTrue(uploadFileRequest1.getCancelHandler().isCancelled());
            // upload canceled，result is null
            assertFalse(result.isPresent());
            // checkpoint 还存在
            Assert.assertTrue(new File(checkpointFilePath).exists());

            ListMultipartUploadsRequest listMultipartUploadsRequest = new ListMultipartUploadsRequest(bucketName);
            MultipartUploadListing multipartUploadListing =
                    obsClientAsync.listMultipartUploads(listMultipartUploadsRequest);
            // 确保暂停的分片上传任务还在分片上传列表中
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().size() > 0);
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().get(0).getObjectKey().equals(objectKey));

        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 暂停上传成功，同时取消上传任务
    @Test
    public void test_cancel_upload_02() throws IOException {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_02";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            CallCancelHandler callCancelHandler1 = new CallCancelHandler();
            uploadFileRequest1.setCancelHandler(callCancelHandler1);
            // 每上传10MB数据反馈上传进度
            uploadFileRequest1.setProgressInterval(10 * 1024 * 1024L);
            uploadFileRequest1.setEnableCheckpoint(true);
            // 暂停上传时，自动取消分段上传任务
            uploadFileRequest1.setNeedAbortUploadFileAfterCancel(true);
            uploadFileRequest1.setProgressListener(
                    status -> {
                        if (status.getTransferPercentage() > 20) {
                            // 当上传进度超过20时，取消上传
                            callCancelHandler1.cancel();
                        }
                    });
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();
            // upload canceled
            assertTrue(uploadFileRequest1.getCancelHandler().isCancelled());
            // upload canceled，result is null
            assertFalse(result.isPresent());
            // checkpoint 不存在
            Assert.assertFalse(new File(checkpointFilePath).exists());

            ListMultipartUploadsRequest listMultipartUploadsRequest = new ListMultipartUploadsRequest(bucketName);
            MultipartUploadListing multipartUploadListing =
                    obsClientAsync.listMultipartUploads(listMultipartUploadsRequest);
            // 确保暂停的分片上传任务已不在分片上传列表中
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().isEmpty());

        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 存在的call数量超过默认最大值，暂停上传成功
    @Test
    public void test_cancel_upload_03() throws IOException {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_03";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            CallCancelHandler callCancelHandler1 = new CallCancelHandler();
            // 设置最大call容量为1
            callCancelHandler1.setMaxCallCapacity(0);
            uploadFileRequest1.setTaskNum(5);
            uploadFileRequest1.setCancelHandler(callCancelHandler1);
            uploadFileRequest1.setProgressInterval(1024 * 1024L);
            uploadFileRequest1.setPartSize(1024 * 1024);
            uploadFileRequest1.setEnableCheckpoint(true);
            // 暂停上传时，不取消分段上传任务
            uploadFileRequest1.setNeedAbortUploadFileAfterCancel(false);
            final int progressToCancel = 20;
            uploadFileRequest1.setProgressListener(
                    status -> {
                        // 获取上传平均速率
                        System.out.println("AverageSpeed:" + status.getAverageSpeed());
                        // 获取上传进度百分比
                        System.out.println("TransferPercentage:" + status.getTransferPercentage());
                        if (status.getTransferPercentage() > progressToCancel) {
                            // 当上传进度超过progressToCancel时，取消上传
                            callCancelHandler1.cancel();
                        }
                    });
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();
            assertTrue(uploadFileRequest1.getCancelHandler().isCancelled());
            // upload canceled，result is null
            assertFalse(result.isPresent());
            uploadFileRequest1.setProgressListener(
                    status -> {
                        // 获取上传平均速率
                        System.out.println("AverageSpeed:" + status.getAverageSpeed());
                        // 获取上传进度百分比
                        System.out.println("TransferPercentage:" + status.getTransferPercentage());
                        // 确保断点续传机制还生效，进度不能丢失
                        assertTrue(status.getTransferPercentage() >= progressToCancel);
                    });
            uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            result = uploadFileTask.getResult();
            // upload succeeded，result is not null
            assertTrue(result.isPresent());
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 断点续传，上传成功（不暂停）
    @Test
    public void test_cancel_upload_04() throws IOException {
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_04";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            CallCancelHandler callCancelHandler1 = new CallCancelHandler();
            uploadFileRequest1.setCancelHandler(callCancelHandler1);
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();
        } finally {
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 断点续传，上传成功（不暂停， 10000分段）
    @Test
    public void test_cancel_upload_05() throws IOException {
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, 100 * 10000);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_04";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            TestCancelHandler defaultOBSHttpCancelHandler1 = new TestCancelHandler();
            uploadFileRequest1.setCancelHandler(defaultOBSHttpCancelHandler1);
            uploadFileRequest1.setTaskNum(32);
            // slice to 10000 parts
            uploadFileRequest1.setPartSize(100 * 1024);
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();
        } finally {
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 暂停上传成功，但不取消上传任务（10000分段）
    @Test
    public void test_cancel_upload_06() throws IOException {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, 100 * 10000);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_06";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            TestCancelHandler defaultOBSHttpCancelHandler1 = new TestCancelHandler();
            uploadFileRequest1.setCancelHandler(defaultOBSHttpCancelHandler1);
            uploadFileRequest1.setTaskNum(32);
            // slice to 10000 parts
            uploadFileRequest1.setPartSize(100 * 1024);
            uploadFileRequest1.setEnableCheckpoint(true);
            uploadFileRequest1.setProgressListener(
                    status -> {
                        if (status.getTransferPercentage() > 5) {
                            // 当上传进度超过20时，取消上传
                            defaultOBSHttpCancelHandler1.cancel();
                        }
                    });
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();

            // upload canceled
            assertTrue(uploadFileRequest1.getCancelHandler().isCancelled());
            // upload canceled，result is null
            assertFalse(result.isPresent());
            // checkpoint 还存在
            Assert.assertTrue(new File(checkpointFilePath).exists());

            ListMultipartUploadsRequest listMultipartUploadsRequest = new ListMultipartUploadsRequest(bucketName);
            MultipartUploadListing multipartUploadListing =
                    obsClientAsync.listMultipartUploads(listMultipartUploadsRequest);
            // 确保暂停的分片上传任务还在分片上传列表中
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().size() > 0);
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().get(0).getObjectKey().equals(objectKey));
        } finally {
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 暂停上传成功，同时取消上传任务（10000分段）
    @Test
    public void test_cancel_upload_07() throws IOException {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, 100 * 10000);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_02";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            CallCancelHandler callCancelHandler1 = new CallCancelHandler();
            uploadFileRequest1.setCancelHandler(callCancelHandler1);
            // 每上传10MB数据反馈上传进度
            uploadFileRequest1.setProgressInterval(100 * 1024L);
            uploadFileRequest1.setTaskNum(32);
            // slice to 10000 parts
            uploadFileRequest1.setPartSize(10 * 1024);
            uploadFileRequest1.setEnableCheckpoint(true);
            // 暂停上传时，自动取消分段上传任务
            uploadFileRequest1.setNeedAbortUploadFileAfterCancel(true);
            uploadFileRequest1.setProgressListener(
                    status -> {
                        if (status.getTransferPercentage() > 5) {
                            // 当上传进度超过20时，取消上传
                            callCancelHandler1.cancel();
                        }
                    });
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();
            // upload canceled
            assertTrue(uploadFileRequest1.getCancelHandler().isCancelled());
            // upload canceled，result is null
            assertFalse(result.isPresent());
            // checkpoint 不存在
            Assert.assertFalse(new File(checkpointFilePath).exists());

            ListMultipartUploadsRequest listMultipartUploadsRequest = new ListMultipartUploadsRequest(bucketName);
            MultipartUploadListing multipartUploadListing =
                    obsClientAsync.listMultipartUploads(listMultipartUploadsRequest);
            // 确保暂停的分片上传任务已不在分片上传列表中
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().isEmpty());

        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 暂停上传成功，但不取消上传任务，继续上传，可以成功上传
    @Test
    public void test_cancel_upload_08() throws IOException {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test_cancel_upload_03";
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            CallCancelHandler callCancelHandler1 = new CallCancelHandler();
            uploadFileRequest1.setTaskNum(5);
            uploadFileRequest1.setCancelHandler(callCancelHandler1);
            uploadFileRequest1.setProgressInterval(1024 * 1024L);
            uploadFileRequest1.setPartSize(1024 * 1024);
            uploadFileRequest1.setEnableCheckpoint(true);
            // 暂停上传时，不取消分段上传任务
            uploadFileRequest1.setNeedAbortUploadFileAfterCancel(false);
            uploadFileRequest1.setProgressListener(
                    status -> {
                        if (status.getTransferPercentage() > 20) {
                            // 当上传进度超过20时，取消上传
                            callCancelHandler1.cancel();
                        }
                    });
            assert obsClientAsync != null;
            UploadFileTask uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();

            // upload canceled
            assertTrue(uploadFileRequest1.getCancelHandler().isCancelled());
            // upload canceled，result is null
            assertFalse(result.isPresent());
            // checkpoint 还存在
            Assert.assertTrue(new File(checkpointFilePath).exists());

            ListMultipartUploadsRequest listMultipartUploadsRequest = new ListMultipartUploadsRequest(bucketName);
            MultipartUploadListing multipartUploadListing =
                    obsClientAsync.listMultipartUploads(listMultipartUploadsRequest);
            // 确保暂停的分片上传任务还在分片上传列表中
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().size() > 0);
            Assert.assertTrue(multipartUploadListing.getMultipartTaskList().get(0).getObjectKey().equals(objectKey));

            uploadFileRequest1.setProgressListener(null);
            // 继续上传
            uploadFileTask = obsClientAsync.uploadFileAsync(uploadFileRequest1, completeCallback);
            result = uploadFileTask.getResult();
            // upload succeeded，result is not null
            assertTrue(result.isPresent());
            // 对象已上传成功
            assertTrue(obsClientAsync.doesObjectExist(bucketName, objectKey));
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    // 验证设置queryInterval、executor功能可用,验证client的初始化函数可用
    @Test
    public void test_cancel_upload_09() throws IOException {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync1 = TestTools.getPipelineEnvironmentForAsyncClient()) {

            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = bucketName;
            assert obsClientAsync1 != null;
            int queryInterval = 3000;
            obsClientAsync1.setQueryInterval(queryInterval);
            assertEquals(queryInterval,obsClientAsync1.getQueryInterval());
            ExecutorService newFixedThreadPool = Executors.newFixedThreadPool(200,
                    r -> new Thread(r, "newTestThread"));
            UploadFileRequest uploadFileRequest1 = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            // 异步上传
            UploadFileTask uploadFileTask = obsClientAsync1.uploadFileAsync(uploadFileRequest1, completeCallback);
            obsClientAsync1.setExecutorService(newFixedThreadPool);
            uploadFileTask.waitUntilFinished();
            Assert.assertTrue(uploadFileTask.isTaskFinished());
            Optional<CompleteMultipartUploadResult> result = uploadFileTask.getResult();
            // upload succeeded，result is not null
            assertTrue(result.isPresent());
            boolean cancelExecuted = uploadFileTask.cancel();
            // no cancelHandler, cancel failed
            assertFalse(cancelExecuted);
            String testEndpoint = "test-endpoint-not-existed";
            ObsConfiguration configuration = new ObsConfiguration();
            configuration.setEndPoint(testEndpoint);
            IObsCredentialsProvider obsCredentialsProvider = new BasicObsCredentialsProvider(testEndpoint, testEndpoint);
            ObsClientAsync obsClientAsync2 = new ObsClientAsync(testEndpoint);
            obsClientAsync2.close();
            ObsClientAsync obsClientAsync3 = new ObsClientAsync(configuration);
            obsClientAsync3.close();
            ObsClientAsync obsClientAsync4 = new ObsClientAsync(testEndpoint, testEndpoint, testEndpoint);
            obsClientAsync4.close();
            ObsClientAsync obsClientAsync5 = new ObsClientAsync(testEndpoint, testEndpoint, testEndpoint, testEndpoint);
            obsClientAsync5.close();
            ObsClientAsync obsClientAsync6 = new ObsClientAsync(testEndpoint, testEndpoint, testEndpoint, configuration);
            obsClientAsync6.close();
            ObsClientAsync obsClientAsync7 = new ObsClientAsync(obsCredentialsProvider, testEndpoint);
            obsClientAsync7.close();
            ObsClientAsync obsClientAsync8 = new ObsClientAsync(obsCredentialsProvider, configuration);
            obsClientAsync8.close();
            assertEquals(uploadFileTask.getObsClient(), obsClientAsync1);
            uploadFileTask.setBucketName(testEndpoint);
            assertEquals(uploadFileTask.getBucketName(), testEndpoint);
            obsClientAsync1.getBucketMetadata(new BucketMetadataInfoRequest(testEndpoint));
        } catch (ObsException e) {
            System.out.println(e);
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }
    public class TestCancelHandler extends CallCancelHandler {
        @Override
        public void setCall(Call call) {
            super.setCall(call);
            // call队列容量不能超过设置的容量限制
            if(super.calls.size() > getMaxCallCapacity()){
                System.out.println("CallCancelHandler calls size is:" + super.calls.size() + ". MaxCallCapacity is " + getMaxCallCapacity());
            }
            Assert.assertTrue(super.calls.size() <= getMaxCallCapacity());
        }
    }
}
