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

package com.example.obs;

import com.obs.services.ObsClient;
import com.obs.services.ObsClientAsync;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.task.DownloadFileTask;
import com.obs.services.internal.task.UploadFileTask;
import com.obs.services.model.CompleteMultipartUploadResult;
import com.obs.services.model.DownloadFileRequest;
import com.obs.services.model.DownloadFileResult;
import com.obs.services.model.ResumableTransferHandle;
import com.obs.services.model.TaskCallback;
import com.obs.services.model.UploadFileRequest;

import java.io.IOException;
import java.util.Optional;

/**
 * 断点续传暂停/取消功能示例代码。
 *
 * <p>本示例展示如何使用 {@link ResumableTransferHandle} 对上传和下载任务进行暂停、恢复和取消操作。
 *
 * <p>包含场景：
 * <ul>
 *   <li>同步上传暂停与恢复</li>
 *   <li>同步下载暂停与恢复</li>
 *   <li>异步上传取消</li>
 *   <li>异步下载暂停和取消</li>
 * </ul>
 */
public class ResumableTransferSample {
    private static final String endPoint = "https://your-endpoint";
    private static final String ak = "*** Provide your Access Key ***";
    private static final String sk = "*** Provide your Secret Key ***";
    private static final String bucketName = "my-obs-bucket-demo";
    private static final String objectKey = "my-obs-object-key-demo";
    private static final String localFilePath = "local-file-path";
    private static final String downloadFilePath = "download-file-path";

    public static void main(String[] args) {
        ObsConfiguration config = new ObsConfiguration();
        config.setEndPoint(endPoint);
        config.setSocketTimeout(30000);
        config.setConnectionTimeout(10000);

        try {
            syncUploadWithPauseAndResume(config);
            syncDownloadWithPauseAndResume(config);
            asyncUploadWithCancel(config);
            asyncDownloadWithPauseAndCancel(config);
        } finally {
            // 资源清理由各方法内部处理
        }
    }

    /**
     * 场景1：同步上传使用 ResumableTransferHandle 实现暂停与恢复。
     *
     * <p>关键步骤：
     * <ol>
     *   <li>创建 UploadFileRequest 并设置 enableCheckpoint=true</li>
     *   <li>创建 ResumableTransferHandle 并设置到 request 上</li>
     *   <li>在进度回调中根据条件调用 controller.pause()</li>
     *   <li>暂停后调用 controller.reset() 并重新上传以恢复</li>
     * </ol>
     */
    private static void syncUploadWithPauseAndResume(ObsConfiguration config) {
        ObsClient obsClient = null;
        try {
            obsClient = new ObsClient(ak, sk, config);

            // 第一次上传：设置 controller 并在进度 >50% 时暂停
            UploadFileRequest request = new UploadFileRequest(bucketName, objectKey, localFilePath);
            request.setEnableCheckpoint(true);
            request.setPartSize(5 * 1024 * 1024L);
            request.setTaskNum(5);
            request.setNeedAbortUploadFileAfterCancel(false);

            ResumableTransferHandle controller = new ResumableTransferHandle();
            request.setTransferHandle(controller);
            request.setProgressInterval(1 * 1024 * 1024L);
            request.setProgressListener(status -> {
                System.out.println("Upload progress: " + status.getTransferPercentage() + "%");
                if (status.getTransferPercentage() > 50) {
                    System.out.println("Pausing upload at " + status.getTransferPercentage() + "%");
                    controller.pause();
                }
            });

            try {
                obsClient.uploadFile(request);
            } catch (ObsException e) {
                if (controller.isPaused()) {
                    System.out.println("Upload paused successfully, checkpoint saved.");
                } else {
                    throw e;
                }
            }

            // 恢复上传：重置 controller，移除暂停触发器
            controller.resume();
            request.setProgressListener(status ->
                    System.out.println("Resume upload progress: " + status.getTransferPercentage() + "%"));

            CompleteMultipartUploadResult result = obsClient.uploadFile(request);
            System.out.println("Resume upload completed! ObjectUrl: " + result.getObjectUrl());
        } catch (ObsException e) {
            System.out.println("Upload failed - HTTP Code: " + e.getResponseCode()
                    + ", Error Message: " + e.getErrorMessage());
        } finally {
            if (obsClient != null) {
                try {
                    obsClient.close();
                } catch (IOException e) {
                    // ignore
                }
            }
        }
    }

    /**
     * 场景2：同步下载使用 ResumableTransferHandle 实现暂停与恢复。
     *
     * <p>关键步骤：
     * <ol>
     *   <li>创建 DownloadFileRequest 并设置 enableCheckpoint=true</li>
     *   <li>创建 ResumableTransferHandle 并设置到 request 上</li>
     *   <li>在进度回调中根据条件调用 controller.pause()</li>
     *   <li>暂停后调用 controller.reset() 并重新下载以恢复</li>
     * </ol>
     */
    private static void syncDownloadWithPauseAndResume(ObsConfiguration config) {
        ObsClient obsClient = null;
        try {
            obsClient = new ObsClient(ak, sk, config);

            // 第一次下载：设置 controller 并在进度 >50% 时暂停
            DownloadFileRequest request = new DownloadFileRequest(bucketName, objectKey);
            request.setDownloadFile(downloadFilePath);
            request.setEnableCheckpoint(true);
            request.setPartSize(5 * 1024 * 1024L);
            request.setTaskNum(5);

            ResumableTransferHandle controller = new ResumableTransferHandle();
            request.setTransferHandle(controller);
            request.setProgressInterval(1 * 1024 * 1024L);
            request.setProgressListener(status -> {
                System.out.println("Download progress: " + status.getTransferPercentage() + "%");
                if (status.getTransferPercentage() > 50) {
                    System.out.println("Pausing download at " + status.getTransferPercentage() + "%");
                    controller.pause();
                }
            });

            try {
                obsClient.downloadFile(request);
            } catch (ObsException e) {
                if (controller.isPaused()) {
                    System.out.println("Download paused successfully, checkpoint saved.");
                } else {
                    throw e;
                }
            }

            // 恢复下载：重置 controller，移除暂停触发器
            controller.resume();
            request.setProgressListener(status ->
                    System.out.println("Resume download progress: " + status.getTransferPercentage() + "%"));

            DownloadFileResult result = obsClient.downloadFile(request);
            System.out.println("Resume download completed! Metadata: " + result.getObjectMetadata());
        } catch (ObsException e) {
            System.out.println("Download failed - HTTP Code: " + e.getResponseCode()
                    + ", Error Message: " + e.getErrorMessage());
        } finally {
            if (obsClient != null) {
                try {
                    obsClient.close();
                } catch (IOException e) {
                    // ignore
                }
            }
        }
    }

    /**
     * 场景3：异步上传使用 UploadFileTask.cancel() 取消上传。
     *
     * <p>关键步骤：
     * <ol>
     *   <li>创建 ObsClientAsync</li>
     *   <li>设置 ResumableTransferHandle 到 UploadFileRequest</li>
     *   <li>通过 uploadFileAsync 提交异步上传任务</li>
     *   <li>调用 controller.cancel() 或 task.cancel() 取消上传</li>
     * </ol>
     */
    private static void asyncUploadWithCancel(ObsConfiguration config) {
        ObsClientAsync obsClientAsync = null;
        try {
            obsClientAsync = new ObsClientAsync(ak, sk, config);

            UploadFileRequest request = new UploadFileRequest(bucketName, objectKey, localFilePath);
            request.setEnableCheckpoint(true);
            request.setPartSize(5 * 1024 * 1024L);
            request.setNeedAbortUploadFileAfterCancel(true);

            ResumableTransferHandle controller = new ResumableTransferHandle();
            request.setTransferHandle(controller);
            request.setProgressInterval(1 * 1024 * 1024L);

            TaskCallback<CompleteMultipartUploadResult, UploadFileRequest> callback =
                    new TaskCallback<CompleteMultipartUploadResult, UploadFileRequest>() {
                        @Override
                        public void onSuccess(CompleteMultipartUploadResult result) {
                            System.out.println("Async upload succeeded: " + result.getObjectUrl());
                        }

                        @Override
                        public void onException(ObsException e, UploadFileRequest singleRequest) {
                            System.out.println("Async upload exception: " + e.getErrorMessage());
                        }
                    };

            UploadFileTask task = obsClientAsync.uploadFileAsync(request, callback);

            // 等待一段时间后取消上传（实际使用中可根据业务逻辑触发）
            Thread.sleep(3000);
            controller.cancel();
            System.out.println("Upload cancel requested.");

            Optional<CompleteMultipartUploadResult> result = task.getResult();
            if (!result.isPresent()) {
                System.out.println("Upload was cancelled successfully.");
            }
        } catch (ObsException e) {
            System.out.println("Async upload failed - HTTP Code: " + e.getResponseCode()
                    + ", Error Message: " + e.getErrorMessage());
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } finally {
            if (obsClientAsync != null) {
                try {
                    obsClientAsync.close();
                } catch (IOException e) {
                    // ignore
                }
            }
        }
    }

    /**
     * 场景4：异步下载使用 DownloadFileTask 的 pause() 和 cancel() 方法。
     *
     * <p>关键步骤：
     * <ol>
     *   <li>创建 ObsClientAsync</li>
     *   <li>设置 ResumableTransferHandle 到 DownloadFileRequest</li>
     *   <li>通过 downloadFileAsync 提交异步下载任务</li>
     *   <li>通过 controller.pause() 暂停或 controller.cancel() 取消</li>
     *   <li>暂停后可 reset controller 并重新 downloadFileAsync 恢复下载</li>
     * </ol>
     */
    private static void asyncDownloadWithPauseAndCancel(ObsConfiguration config) {
        ObsClientAsync obsClientAsync = null;
        try {
            obsClientAsync = new ObsClientAsync(ak, sk, config);

            DownloadFileRequest request = new DownloadFileRequest(bucketName, objectKey);
            request.setDownloadFile(downloadFilePath);
            request.setEnableCheckpoint(true);
            request.setPartSize(5 * 1024 * 1024L);
            request.setTaskNum(5);

            ResumableTransferHandle controller = new ResumableTransferHandle();
            request.setTransferHandle(controller);
            request.setProgressInterval(1 * 1024 * 1024L);

            TaskCallback<DownloadFileResult, DownloadFileRequest> callback =
                    new TaskCallback<DownloadFileResult, DownloadFileRequest>() {
                        @Override
                        public void onSuccess(DownloadFileResult result) {
                            System.out.println("Async download succeeded: " + result.getObjectMetadata());
                        }

                        @Override
                        public void onException(ObsException e, DownloadFileRequest singleRequest) {
                            System.out.println("Async download exception: " + e.getErrorMessage());
                        }
                    };

            DownloadFileTask task = obsClientAsync.downloadFileAsync(request, callback);

            // 暂停下载（实际使用中可根据业务逻辑触发）
            Thread.sleep(2000);
            controller.pause();
            System.out.println("Download pause requested.");

            Optional<DownloadFileResult> pauseResult = task.getResult();
            if (!pauseResult.isPresent() && controller.isPaused()) {
                System.out.println("Download paused, checkpoint saved. Resuming...");

                // 恢复下载
                controller.resume();
                request.setProgressListener(null);
                DownloadFileTask resumeTask = obsClientAsync.downloadFileAsync(request, callback);
                Optional<DownloadFileResult> resumeResult = resumeTask.getResult();
                if (resumeResult.isPresent()) {
                    System.out.println("Resume download completed!");
                }
            }

            // 取消下载示例（另一个独立场景）
            // controller.cancel();
            // DownloadFileTask cancelTask = obsClientAsync.downloadFileAsync(cancelRequest, callback);
            // controller.cancel();
        } catch (ObsException e) {
            System.out.println("Async download failed - HTTP Code: " + e.getResponseCode()
                    + ", Error Message: " + e.getErrorMessage());
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } finally {
            if (obsClientAsync != null) {
                try {
                    obsClientAsync.close();
                } catch (IOException e) {
                    // ignore
                }
            }
        }
    }
}
