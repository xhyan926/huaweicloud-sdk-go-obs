/**
 * Copyright 2019 Huawei Technologies Co.,Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package com.obs.integrated_test;

import com.obs.services.ObsClientAsync;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.task.DownloadFileTask;
import com.obs.services.internal.task.UploadFileTask;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.DownloadFileRequest;
import com.obs.services.model.DownloadFileResult;
import com.obs.services.model.ResumableTransferHandle;
import com.obs.services.model.TaskCallback;
import com.obs.services.model.UploadFileRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
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

import static com.obs.test.TestTools.genTestFile;
import static com.obs.test.TestTools.removeTestFile;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertTrue;

/**
 * 集成测试：验证 ResumableTransferHandle 对断点续传上传/下载的暂停和取消控制能力。
 *
 * <p>测试场景覆盖：
 * <ul>
 *   <li>上传暂停（保留 checkpoint）与恢复</li>
 *   <li>上传取消</li>
 *   <li>下载暂停（保留 checkpoint）与恢复</li>
 *   <li>下载取消</li>
 *   <li>不设置 handle 时上传正常完成（兼容性）</li>
 * </ul>
 */
public class ResumableTransferIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    private static final String TEST_FILE_NORMAL = "testFileInTransferHandleTest";
    private static final int KB_SIZEOF_TEST_FILE_NORMAL = 1024 * 100;
    private static final String DOWNLOAD_FILE_SUFFIX = ".download";

    private final TaskCallback<CompleteMultipartUploadResult, UploadFileRequest> uploadCallback =
            new TaskCallback<CompleteMultipartUploadResult, UploadFileRequest>() {
                @Override
                public void onSuccess(CompleteMultipartUploadResult result) {
                    System.out.println("uploadFileAsync succeeded!");
                    System.out.println("HTTP StatusCode:" + result.getStatusCode());
                    System.out.println("ObjectUrl:" + result.getObjectUrl());
                }

                @Override
                public void onException(ObsException e, UploadFileRequest singleRequest) {
                    System.out.println("uploadFileAsync failed");
                    printObsException(e);
                }
            };

    private TaskCallback<DownloadFileResult, DownloadFileRequest> downloadCallback =
            new TaskCallback<DownloadFileResult, DownloadFileRequest>() {
                @Override
                public void onSuccess(DownloadFileResult result) {
                    System.out.println("downloadFileAsync succeeded!");
                    System.out.println("ObjectMetadata:" + result.getObjectMetadata());
                }

                @Override
                public void onException(ObsException e, DownloadFileRequest singleRequest) {
                    System.out.println("downloadFileAsync failed");
                    printObsException(e);
                }
            };

    private static void printObsException(ObsException e) {
        System.out.println("HTTP Code:" + e.getResponseCode());
        System.out.println("Error Code:" + e.getErrorCode());
        System.out.println("Error Message:" + e.getErrorMessage());
        System.out.println("Request ID:" + e.getErrorRequestId());
        System.out.println("Host ID:" + e.getErrorHostId());
        Map<String, String> headers = e.getResponseHeaders();
        if (headers != null) {
            for (Map.Entry<String, String> header : headers.entrySet()) {
                if (header.getKey().contains("error")) {
                    System.out.println("errorHeaderKey:" + header.getKey()
                            + ", errorHeaderValue:" + header.getValue());
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

    /**
     * 场景1：异步上传过程中直接调用 ResumableTransferHandle.pause() 暂停上传，
     * 验证 checkpoint 文件保留，且 isPaused 为 true。
     */
    @Test
    public void test_pause_upload_with_checkpoint() throws Exception {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test-pause-upload-with-checkpoint";

            UploadFileRequest request = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            ResumableTransferHandle handle = new ResumableTransferHandle();
            request.setTransferHandle(handle);
            request.setProgressInterval(1024 * 1024L);
            request.setEnableCheckpoint(true);
            request.setNeedAbortUploadFileAfterCancel(false);

            assert obsClientAsync != null;
            UploadFileTask task = obsClientAsync.uploadFileAsync(request, uploadCallback);

            // 直接从外部调用暂停，等待上传开始后暂停
            Thread.sleep(1000);
            handle.pause();

            Optional<CompleteMultipartUploadResult> result = task.getResult();

            assertTrue("Handle should be paused", handle.isPaused());
            assertFalse("Result should be empty after pause", result.isPresent());
            assertTrue("Checkpoint file should exist", new File(checkpointFilePath).exists());
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    /**
     * 场景2：暂停上传后重置 handle，再次上传完成，验证恢复上传成功。
     */
    @Test
    public void test_pause_upload_then_resume() throws Exception {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test-pause-upload-then-resume";

            UploadFileRequest request = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            ResumableTransferHandle handle = new ResumableTransferHandle();
            request.setTransferHandle(handle);
            request.setTaskNum(5);
            request.setProgressInterval(1024 * 1024L);
            request.setPartSize(1024 * 1024);
            request.setEnableCheckpoint(true);
            request.setNeedAbortUploadFileAfterCancel(false);
            request.setProgressListener(
                    status -> {
                        if (status.getTransferPercentage() > 20) {
                            handle.pause();
                        }
                    });

            assert obsClientAsync != null;
            UploadFileTask task = obsClientAsync.uploadFileAsync(request, uploadCallback);
            Optional<CompleteMultipartUploadResult> result = task.getResult();

            assertTrue("Handle should be paused", handle.isPaused());
            assertFalse("First upload result should be empty", result.isPresent());
            assertTrue("Checkpoint should exist", new File(checkpointFilePath).exists());

            // 重置 handle，移除暂停触发器，恢复上传
            handle.resume();
            request.setProgressListener(null);
            UploadFileTask resumeTask = obsClientAsync.uploadFileAsync(request, uploadCallback);
            Optional<CompleteMultipartUploadResult> resumeResult = resumeTask.getResult();

            assertTrue("Resume upload result should be present", resumeResult.isPresent());
            assertTrue("Object should exist after resume",
                    obsClientAsync.doesObjectExist(bucketName, objectKey));
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    /**
     * 场景3：异步上传过程中直接调用 ResumableTransferHandle.cancel() 取消上传。
     */
    @Test
    public void test_cancel_upload_via_handle() throws Exception {
        String checkpointFilePath = TEST_FILE_NORMAL + ".uploadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test-cancel-upload-via-handle";

            UploadFileRequest request = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            ResumableTransferHandle handle = new ResumableTransferHandle();
            request.setTransferHandle(handle);
            request.setProgressInterval(1024 * 1024L);
            request.setEnableCheckpoint(true);
            request.setNeedAbortUploadFileAfterCancel(true);

            assert obsClientAsync != null;
            UploadFileTask task = obsClientAsync.uploadFileAsync(request, uploadCallback);

            // 直接从外部调用取消，等待上传开始后取消
            Thread.sleep(1000);
            handle.cancel();

            Optional<CompleteMultipartUploadResult> result = task.getResult();

            assertTrue("Handle should be cancelled", handle.isCancelled());
            assertFalse("Result should be empty after cancel", result.isPresent());
            assertFalse("Checkpoint should be deleted after cancel with abort",
                    new File(checkpointFilePath).exists());
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    /**
     * 场景4：异步下载过程中直接调用 ResumableTransferHandle.pause() 暂停下载，
     * 验证 checkpoint 文件保留，且 isPaused 为 true。
     */
    @Test
    public void test_pause_download_with_checkpoint() throws Exception {
        String downloadFilePath = TEST_FILE_NORMAL + DOWNLOAD_FILE_SUFFIX;
        String checkpointFilePath = downloadFilePath + ".downloadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test-pause-download-with-checkpoint";

            // 先上传一个对象
            UploadFileRequest uploadRequest = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            assert obsClientAsync != null;
            UploadFileTask uploadTask = obsClientAsync.uploadFileAsync(uploadRequest, uploadCallback);
            Optional<CompleteMultipartUploadResult> uploadResult = uploadTask.getResult();
            assertTrue("Upload should succeed first", uploadResult.isPresent());

            // 异步下载
            DownloadFileRequest downloadRequest = new DownloadFileRequest(bucketName, objectKey);
            downloadRequest.setDownloadFile(downloadFilePath);
            downloadRequest.setPartSize(1024 * 1024);
            downloadRequest.setTaskNum(5);
            downloadRequest.setEnableCheckpoint(true);

            ResumableTransferHandle handle = new ResumableTransferHandle();
            downloadRequest.setTransferHandle(handle);
            downloadRequest.setProgressInterval(1024 * 1024L);

            DownloadFileTask downloadTask = obsClientAsync.downloadFileAsync(downloadRequest, downloadCallback);

            // 直接从外部调用暂停，等待下载开始后暂停
            Thread.sleep(1000);
            handle.pause();

            Optional<DownloadFileResult> downloadResult = downloadTask.getResult();

            assertTrue("Handle should be paused", handle.isPaused());
            assertFalse("Download result should be empty after pause", downloadResult.isPresent());
            assertTrue("Download checkpoint file should exist",
                    new File(checkpointFilePath).exists());
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(downloadFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    /**
     * 场景5：暂停下载后重置 handle，再次下载完成，验证恢复下载成功。
     */
    @Test
    public void test_pause_download_then_resume() throws Exception {
        String downloadFilePath = TEST_FILE_NORMAL + DOWNLOAD_FILE_SUFFIX;
        String checkpointFilePath = downloadFilePath + ".downloadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test-pause-download-then-resume";

            // 先上传
            UploadFileRequest uploadRequest = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            assert obsClientAsync != null;
            UploadFileTask uploadTask = obsClientAsync.uploadFileAsync(uploadRequest, uploadCallback);
            assertTrue("Upload should succeed first", uploadTask.getResult().isPresent());

            // 异步下载
            DownloadFileRequest downloadRequest = new DownloadFileRequest(bucketName, objectKey);
            downloadRequest.setDownloadFile(downloadFilePath);
            downloadRequest.setPartSize(1024 * 1024);
            downloadRequest.setTaskNum(5);
            downloadRequest.setEnableCheckpoint(true);

            ResumableTransferHandle handle = new ResumableTransferHandle();
            downloadRequest.setTransferHandle(handle);
            downloadRequest.setProgressInterval(1024 * 1024L);

            DownloadFileTask downloadTask = obsClientAsync.downloadFileAsync(downloadRequest, downloadCallback);

            // 直接从外部调用暂停
            Thread.sleep(1000);
            handle.pause();

            Optional<DownloadFileResult> downloadResult = downloadTask.getResult();

            assertTrue("Handle should be paused", handle.isPaused());
            assertFalse("First download result should be empty", downloadResult.isPresent());

            // 重置 handle，恢复下载
            handle.resume();
            DownloadFileTask resumeTask = obsClientAsync.downloadFileAsync(downloadRequest, downloadCallback);
            Optional<DownloadFileResult> resumeResult = resumeTask.getResult();

            assertTrue("Resume download result should be present", resumeResult.isPresent());
            assertTrue("Downloaded file should exist", new File(downloadFilePath).exists());
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(downloadFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    /**
     * 场景6：异步下载过程中直接调用 ResumableTransferHandle.cancel() 取消下载。
     */
    @Test
    public void test_cancel_download_via_handle() throws Exception {
        String downloadFilePath = TEST_FILE_NORMAL + DOWNLOAD_FILE_SUFFIX;
        String checkpointFilePath = downloadFilePath + ".downloadFile_record";
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test-cancel-download-via-handle";

            // 先上传
            UploadFileRequest uploadRequest = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            assert obsClientAsync != null;
            UploadFileTask uploadTask = obsClientAsync.uploadFileAsync(uploadRequest, uploadCallback);
            assertTrue("Upload should succeed first", uploadTask.getResult().isPresent());

            // 异步下载
            DownloadFileRequest downloadRequest = new DownloadFileRequest(bucketName, objectKey);
            downloadRequest.setDownloadFile(downloadFilePath);
            downloadRequest.setPartSize(1024 * 1024);
            downloadRequest.setTaskNum(5);
            downloadRequest.setEnableCheckpoint(true);

            ResumableTransferHandle handle = new ResumableTransferHandle();
            downloadRequest.setTransferHandle(handle);
            downloadRequest.setProgressInterval(1024 * 1024L);

            DownloadFileTask downloadTask = obsClientAsync.downloadFileAsync(downloadRequest, downloadCallback);

            // 直接从外部调用取消，等待下载开始后取消
            Thread.sleep(1000);
            handle.cancel();

            Optional<DownloadFileResult> downloadResult = downloadTask.getResult();

            assertTrue("Handle should be cancelled", handle.isCancelled());
            assertFalse("Download result should be empty after cancel", downloadResult.isPresent());
        } finally {
            removeTestFile(checkpointFilePath);
            removeTestFile(downloadFilePath);
            removeTestFile(TEST_FILE_NORMAL);
        }
    }

    /**
     * 场景7：不设置 handle 时上传正常完成，验证向后兼容性。
     */
    @Test
    public void test_upload_without_handle_still_works() throws Exception {
        try (ObsClientAsync obsClientAsync = TestTools.getPipelineEnvironmentForAsyncClient()) {
            removeTestFile(TEST_FILE_NORMAL);
            genTestFile(TEST_FILE_NORMAL, KB_SIZEOF_TEST_FILE_NORMAL);
            String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
            String objectKey = "test-upload-without-handle";

            UploadFileRequest request = new UploadFileRequest(bucketName, objectKey, TEST_FILE_NORMAL);
            // 不设置 TransferHandle，验证兼容性
            assert obsClientAsync != null;
            UploadFileTask task = obsClientAsync.uploadFileAsync(request, uploadCallback);
            Optional<CompleteMultipartUploadResult> result = task.getResult();

            assertTrue("Upload result should be present without handle", result.isPresent());
            assertTrue("Object should exist", obsClientAsync.doesObjectExist(bucketName, objectKey));
        } finally {
            removeTestFile(TEST_FILE_NORMAL);
        }
    }
}
