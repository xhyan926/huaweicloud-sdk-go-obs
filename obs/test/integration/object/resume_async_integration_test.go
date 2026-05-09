//go:build integration

package object

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"obs-sdk-go/obs"
	"obs-sdk-go/obs/test/integration"
	// "obs-sdk-go/obs/test/integration"
)

// TestUploadFileAsync_ShouldCompleteSuccessfully_GivenValidFile 测试完整异步上传并等待完成
func TestUploadFileAsync_ShouldCompleteSuccessfully_GivenValidFile(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("async-upload-test.txt")

	// 创建本地测试文件
	localFilePath := filepath.Join(os.TempDir(), "async-upload-test.txt")
	testContent := bytes.Repeat([]byte("A"), 10*1024*1024) // 10MB

	// 清理本地测试文件
	defer func() {
		if _, err := os.Stat(localFilePath); err == nil {
			os.Remove(localFilePath)
		}
	}()

	t.Run("ShouldCreateTestFile_GivenValidContent", func(t *testing.T) {
		err := os.WriteFile(localFilePath, testContent, 0644)
		if err != nil {
			t.Fatalf("创建本地文件失败: %v", err)
		}

		client.AddTestCase("本地测试文件创建成功")
		t.Logf("本地测试文件: %s, 大小: %d bytes", localFilePath, len(testContent))
	})

	t.Run("ShouldUploadAsyncAndComplete_GivenValidInput", func(t *testing.T) {
		input := &obs.UploadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			UploadFile:       localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
		}

		// 启动异步上传
		task, err := client.TestClient().UploadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步上传失败: %v", err)
		}

		// 验证初始状态
		if task.Status() != obs.TransferStatusRunning {
			t.Errorf("期望状态RUNNING，实际: %s", task.Status())
		}

		// 等待任务完成
		<-task.Done()

		// 验证最终状态
		if task.Status() != obs.TransferStatusCompleted {
			t.Errorf("期望状态COMPLETED，实际: %s", task.Status())
		}

		// 获取结果
		result, err := task.GetResult()
		if err != nil {
			t.Fatalf("获取上传结果失败: %v", err)
		}

		output, ok := result.(*obs.CompleteMultipartUploadOutput)
		if !ok {
			t.Fatal("结果类型错误，期望CompleteMultipartUploadOutput")
		}

		if output == nil {
			t.Fatal("UploadFileAsync返回nil")
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		client.AddTestCase("异步上传完成成功")
		t.Logf("异步上传完成，对象: %s, 大小: %d bytes", objectKey, output.ContentLength)
	})

	t.Run("ShouldVerifyUploadedObject_GivenCompletedUpload", func(t *testing.T) {
		// 验证上传的对象
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(input)
		if err != nil {
			t.Fatalf("获取上传对象失败: %v", err)
		}
		defer output.Body.Close()

		// 读取内容
		downloadedContent, err := io.ReadAll(output.Body)
		if err != nil {
			t.Fatalf("读取下载内容失败: %v", err)
		}

		// 验证内容完整性
		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("上传对象验证成功")
	})
}

// TestUploadFileAsync_ShouldPauseAndResume_GivenCheckpoint 测试暂停后恢复上传
func TestUploadFileAsync_ShouldPauseAndResume_GivenCheckpoint(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("pause-resume-upload-test.bin")

	// 创建本地测试文件
	localFilePath := filepath.Join(os.TempDir(), "pause-resume-upload-test.bin")
	checkpointFile := localFilePath + ".uploadfile_record"
	testContent := bytes.Repeat([]byte("B"), 15*1024*1024) // 15MB

	// 清理本地测试文件
	defer func() {
		if _, err := os.Stat(localFilePath); err == nil {
			os.Remove(localFilePath)
		}
		if _, err := os.Stat(checkpointFile); err == nil {
			os.Remove(checkpointFile)
		}
	}()

	t.Run("ShouldCreateTestFile_GivenValidContent", func(t *testing.T) {
		err := os.WriteFile(localFilePath, testContent, 0644)
		if err != nil {
			t.Fatalf("创建本地文件失败: %v", err)
		}

		client.AddTestCase("本地测试文件创建成功")
	})

	t.Run("ShouldPauseAndResumeUpload_GivenRunningTask", func(t *testing.T) {
		input := &obs.UploadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			UploadFile:       localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
			CheckpointFile:   checkpointFile,
		}

		// 启动异步上传
		task, err := client.TestClient().UploadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步上传失败: %v", err)
		}

		// 等待一小段时间让任务开始
		time.Sleep(500 * time.Millisecond)

		// 暂停任务
		err = task.Pause()
		if err != nil {
			t.Fatalf("暂停任务失败: %v", err)
		}

		// 验证暂停状态
		if task.Status() != obs.TransferStatusPaused {
			t.Errorf("期望状态PAUSED，实际: %s", task.Status())
		}

		// 验证checkpoint文件存在
		if _, err := os.Stat(checkpointFile); os.IsNotExist(err) {
			t.Error("checkpoint文件不存在")
		}

		// 恢复任务
		err = task.Resume()
		if err != nil {
			t.Fatalf("恢复任务失败: %v", err)
		}

		// 验证恢复后状态
		if task.Status() != obs.TransferStatusRunning {
			t.Errorf("期望状态RUNNING，实际: %s", task.Status())
		}

		// 等待任务完成
		<-task.Done()

		// 验证最终状态
		if task.Status() != obs.TransferStatusCompleted {
			t.Errorf("期望状态COMPLETED，实际: %s", task.Status())
		}

		// 获取结果
		result, err := task.GetResult()
		if err != nil {
			t.Fatalf("获取上传结果失败: %v", err)
		}

		output, ok := result.(*obs.CompleteMultipartUploadOutput)
		if !ok || output == nil {
			t.Fatal("结果类型错误或为nil")
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		client.AddTestCase("暂停和恢复上传成功")
		t.Logf("暂停和恢复上传完成，对象: %s, 大小: %d bytes", objectKey, output.ContentLength)
	})
}

// TestUploadFileAsync_ShouldCancelAndClean_GivenRunningTask 测试取消上传
func TestUploadFileAsync_ShouldCancelAndClean_GivenRunningTask(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("cancel-upload-test.bin")

	// 创建本地测试文件
	localFilePath := filepath.Join(os.TempDir(), "cancel-upload-test.bin")
	checkpointFile := localFilePath + ".uploadfile_record"
	testContent := bytes.Repeat([]byte("C"), 20*1024*1024) // 20MB

	// 清理本地测试文件
	defer func() {
		if _, err := os.Stat(localFilePath); err == nil {
			os.Remove(localFilePath)
		}
		if _, err := os.Stat(checkpointFile); err == nil {
			os.Remove(checkpointFile)
		}
	}()

	t.Run("ShouldCreateTestFile_GivenValidContent", func(t *testing.T) {
		err := os.WriteFile(localFilePath, testContent, 0644)
		if err != nil {
			t.Fatalf("创建本地文件失败: %v", err)
		}

		client.AddTestCase("本地测试文件创建成功")
	})

	t.Run("ShouldCancelUpload_GivenRunningTask", func(t *testing.T) {
		input := &obs.UploadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			UploadFile:       localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
			CheckpointFile:   checkpointFile,
		}

		// 启动异步上传
		task, err := client.TestClient().UploadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步上传失败: %v", err)
		}

		// 等待一小段时间让任务开始
		time.Sleep(500 * time.Millisecond)

		// 取消任务
		err = task.Cancel()
		if err != nil {
			t.Fatalf("取消任务失败: %v", err)
		}

		// 等待任务完成（取消）
		<-task.Done()

		// 验证取消状态
		if task.Status() != obs.TransferStatusCancelled {
			t.Errorf("期望状态CANCELLED，实际: %s", task.Status())
		}

		// 获取结果应该包含错误
		_, err = task.GetResult()
		if err == nil {
			t.Error("期望返回错误，但实际为nil")
		}

		// 验证checkpoint文件被清理
		if _, err := os.Stat(checkpointFile); err == nil {
			t.Error("checkpoint文件应该被清理但仍然存在")
		}

		client.AddTestCase("取消上传成功")
	})
}

// TestDownloadFileAsync_ShouldCompleteSuccessfully_GivenValidObject 测试完整异步下载并等待完成
func TestDownloadFileAsync_ShouldCompleteSuccessfully_GivenValidObject(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("async-download-test.bin")

	// 创建测试对象
	testContent := bytes.Repeat([]byte("D"), 12*1024*1024) // 12MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(testContent),
	}

	_, err := client.TestClient().PutObject(putInput)
	if err != nil {
		t.Fatalf("创建测试对象失败: %v", err)
	}

	// 添加清理函数
	client.AddCleanup(func(t *testing.T) {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		_, err := client.TestClient().DeleteObject(deleteInput)
		if err != nil {
			t.Logf("删除对象失败: %v", err)
		}
	})

	client.AddTestCase("测试对象创建完成")

	t.Run("ShouldDownloadAsyncAndComplete_GivenValidObject", func(t *testing.T) {
		localFilePath := filepath.Join(os.TempDir(), "async-download-test.bin")
		checkpointFile := localFilePath + ".downloadfile_record"

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
			if _, err := os.Stat(checkpointFile); err == nil {
				os.Remove(checkpointFile)
			}
			if _, err := os.Stat(localFilePath + ".tmp"); err == nil {
				os.Remove(localFilePath + ".tmp")
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			DownloadFile:     localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
			CheckpointFile:   checkpointFile,
		}

		// 启动异步下载
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步下载失败: %v", err)
		}

		// 验证初始状态
		if task.Status() != obs.TransferStatusRunning {
			t.Errorf("期望状态RUNNING，实际: %s", task.Status())
		}

		// 等待任务完成
		<-task.Done()

		// 验证最终状态
		if task.Status() != obs.TransferStatusCompleted {
			t.Errorf("期望状态COMPLETED，实际: %s", task.Status())
		}

		// 获取结果
		result, err := task.GetResult()
		if err != nil {
			t.Fatalf("获取下载结果失败: %v", err)
		}

		output, ok := result.(*obs.GetObjectMetadataOutput)
		if !ok {
			t.Fatal("结果类型错误，期望GetObjectMetadataOutput")
		}

		if output == nil {
			t.Fatal("DownloadFileAsync返回nil")
		}

		// 验证下载的文件内容
		downloadedContent, err := os.ReadFile(localFilePath)
		if err != nil {
			t.Fatalf("读取下载文件失败: %v", err)
		}

		// 验证内容完整性
		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("异步下载完成成功")
		t.Logf("异步下载完成，文件: %s, 大小: %d bytes", localFilePath, output.ContentLength)
	})
}

// TestDownloadFileAsync_ShouldPauseAndResume_GivenCheckpoint 测试暂停后恢复下载
func TestDownloadFileAsync_ShouldPauseAndResume_GivenCheckpoint(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("pause-resume-download-test.bin")

	// 创建测试对象
	testContent := bytes.Repeat([]byte("E"), 18*1024*1024) // 18MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(testContent),
	}

	_, err := client.TestClient().PutObject(putInput)
	if err != nil {
		t.Fatalf("创建测试对象失败: %v", err)
	}

	// 添加清理函数
	client.AddCleanup(func(t *testing.T) {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		_, err := client.TestClient().DeleteObject(deleteInput)
		if err != nil {
			t.Logf("删除对象失败: %v", err)
		}
	})

	client.AddTestCase("测试对象创建完成")

	t.Run("ShouldPauseAndResumeDownload_GivenRunningTask", func(t *testing.T) {
		localFilePath := filepath.Join(os.TempDir(), "pause-resume-download-test.bin")
		checkpointFile := localFilePath + ".downloadfile_record"

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
			if _, err := os.Stat(checkpointFile); err == nil {
				os.Remove(checkpointFile)
			}
			if _, err := os.Stat(localFilePath + ".tmp"); err == nil {
				os.Remove(localFilePath + ".tmp")
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			DownloadFile:     localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
			CheckpointFile:   checkpointFile,
		}

		// 启动异步下载
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步下载失败: %v", err)
		}

		// 等待一小段时间让任务开始
		time.Sleep(500 * time.Millisecond)

		// 暂停任务
		err = task.Pause()
		if err != nil {
			t.Fatalf("暂停任务失败: %v", err)
		}

		// 验证暂停状态
		if task.Status() != obs.TransferStatusPaused {
			t.Errorf("期望状态PAUSED，实际: %s", task.Status())
		}

		// 验证checkpoint文件存在
		if _, err := os.Stat(checkpointFile); os.IsNotExist(err) {
			t.Error("checkpoint文件不存在")
		}

		// 恢复任务
		err = task.Resume()
		if err != nil {
			t.Fatalf("恢复任务失败: %v", err)
		}

		// 验证恢复后状态
		if task.Status() != obs.TransferStatusRunning {
			t.Errorf("期望状态RUNNING，实际: %s", task.Status())
		}

		// 等待任务完成
		<-task.Done()

		// 验证最终状态
		if task.Status() != obs.TransferStatusCompleted {
			t.Errorf("期望状态COMPLETED，实际: %s", task.Status())
		}

		// 验证下载的文件内容
		downloadedContent, err := os.ReadFile(localFilePath)
		if err != nil {
			t.Fatalf("读取下载文件失败: %v", err)
		}

		// 验证内容完整性
		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("暂停和恢复下载成功")
		t.Logf("暂停和恢复下载完成，文件大小: %d bytes", len(downloadedContent))
	})
}

// TestDownloadFileAsync_ShouldCancelAndClean_GivenRunningTask 测试取消下载
func TestDownloadFileAsync_ShouldCancelAndClean_GivenRunningTask(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("cancel-download-test.bin")

	// 创建测试对象
	testContent := bytes.Repeat([]byte("F"), 25*1024*1024) // 25MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(testContent),
	}

	_, err := client.TestClient().PutObject(putInput)
	if err != nil {
		t.Fatalf("创建测试对象失败: %v", err)
	}

	// 添加清理函数
	client.AddCleanup(func(t *testing.T) {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		_, err := client.TestClient().DeleteObject(deleteInput)
		if err != nil {
			t.Logf("删除对象失败: %v", err)
		}
	})

	client.AddTestCase("测试对象创建完成")

	t.Run("ShouldCancelDownload_GivenRunningTask", func(t *testing.T) {
		localFilePath := filepath.Join(os.TempDir(), "cancel-download-test.bin")
		checkpointFile := localFilePath + ".downloadfile_record"
		tempFile := localFilePath + ".tmp"

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
			if _, err := os.Stat(checkpointFile); err == nil {
				os.Remove(checkpointFile)
			}
			if _, err := os.Stat(tempFile); err == nil {
				os.Remove(tempFile)
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			DownloadFile:     localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
			CheckpointFile:   checkpointFile,
		}

		// 启动异步下载
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步下载失败: %v", err)
		}

		// 等待一小段时间让任务开始
		time.Sleep(500 * time.Millisecond)

		// 取消任务
		err = task.Cancel()
		if err != nil {
			t.Fatalf("取消任务失败: %v", err)
		}

		// 等待任务完成（取消）
		<-task.Done()

		// 验证取消状态
		if task.Status() != obs.TransferStatusCancelled {
			t.Errorf("期望状态CANCELLED，实际: %s", task.Status())
		}

		// 获取结果应该包含错误
		_, err = task.GetResult()
		if err == nil {
			t.Error("期望返回错误，但实际为nil")
		}

		// 验证临时文件被清理
		if _, err := os.Stat(tempFile); err == nil {
			t.Error("临时文件应该被清理但仍然存在")
		}

		// 验证checkpoint文件被清理
		if _, err := os.Stat(checkpointFile); err == nil {
			t.Error("checkpoint文件应该被清理但仍然存在")
		}

		// 验证目标文件不存在
		if _, err := os.Stat(localFilePath); err == nil {
			t.Error("目标文件不应该存在")
		}

		client.AddTestCase("取消下载成功")
	})
}

// TestTransferTask_ControlFunctions 测试TransferTask控制功能
func TestTransferTask_ControlFunctions(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("task-control-test.bin")

	t.Run("ShouldReturnCorrectStatus_GivenTaskStateChanges", func(t *testing.T) {
		// 创建测试对象
		testContent := bytes.Repeat([]byte("G"), 8*1024*1024) // 8MB

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(testContent),
		}

		_, err := client.TestClient().PutObject(putInput)
		if err != nil {
			t.Fatalf("创建测试对象失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		localFilePath := filepath.Join(os.TempDir(), "task-control-test.bin")

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
			if _, err := os.Stat(localFilePath + ".tmp"); err == nil {
				os.Remove(localFilePath + ".tmp")
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile: localFilePath,
			PartSize:     3 * 1024 * 1024, // 3MB分块
			TaskNum:      2,
		}

		// 启动异步任务
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步任务失败: %v", err)
		}

		// 验证初始状态为RUNNING
		status := task.Status()
		if status != obs.TransferStatusRunning {
			t.Errorf("期望初始状态RUNNING，实际: %s", status)
		}

		// 等待一小段时间让任务开始
		time.Sleep(300 * time.Millisecond)

		// 暂停任务
		err = task.Pause()
		if err != nil {
			t.Fatalf("暂停任务失败: %v", err)
		}

		// 验证状态变为PAUSED
		status = task.Status()
		if status != obs.TransferStatusPaused {
			t.Errorf("期望状态PAUSED，实际: %s", status)
		}

		// 恢复任务
		err = task.Resume()
		if err != nil {
			t.Fatalf("恢复任务失败: %v", err)
		}

		// 验证状态变为RUNNING
		status = task.Status()
		if status != obs.TransferStatusRunning {
			t.Errorf("期望状态RUNNING，实际: %s", status)
		}

		// 等待任务完成
		<-task.Done()

		// 验证最终状态为COMPLETED
		status = task.Status()
		if status != obs.TransferStatusCompleted {
			t.Errorf("期望最终状态COMPLETED，实际: %s", status)
		}

		// 验证Done channel已关闭
		select {
		case <-task.Done():
			// 正确，channel已关闭
		default:
			t.Error("Done channel未关闭")
		}

		client.AddTestCase("任务状态转换验证成功")
	})

	t.Run("ShouldCloseDoneChannel_GivenTaskCompletion", func(t *testing.T) {
		objectKey2 := client.GetTestObjectKey("done-channel-test.bin")

		// 创建测试对象
		testContent := bytes.Repeat([]byte("H"), 5*1024*1024) // 5MB

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey2,
			Body:   bytes.NewReader(testContent),
		}

		_, err := client.TestClient().PutObject(putInput)
		if err != nil {
			t.Fatalf("创建测试对象失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey2,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		localFilePath := filepath.Join(os.TempDir(), "done-channel-test.bin")

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey2,
			DownloadFile: localFilePath,
			PartSize:     2 * 1024 * 1024, // 2MB分块
			TaskNum:      2,
		}

		// 启动异步任务
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步任务失败: %v", err)
		}

		// 等待任务完成
		<-task.Done()

		// 验证Done channel已关闭（多次读取不应该阻塞）
		select {
		case <-task.Done():
			// 正确
		case <-time.After(1 * time.Second):
			t.Error("读取Done channel超时")
		}

		client.AddTestCase("Done channel关闭验证成功")
	})

	t.Run("ShouldReturnCorrectResult_GivenTaskCompletion", func(t *testing.T) {
		objectKey3 := client.GetTestObjectKey("get-result-test.bin")

		// 创建测试对象
		testContent := bytes.Repeat([]byte("I"), 6*1024*1024) // 6MB

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey3,
			Body:   bytes.NewReader(testContent),
		}

		_, err := client.TestClient().PutObject(putInput)
		if err != nil {
			t.Fatalf("创建测试对象失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey3,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		localFilePath := filepath.Join(os.TempDir(), "get-result-test.bin")

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey3,
			DownloadFile: localFilePath,
			PartSize:     2 * 1024 * 1024, // 2MB分块
			TaskNum:      2,
		}

		// 启动异步任务
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步任务失败: %v", err)
		}

		// 等待任务完成
		<-task.Done()

		// 获取结果
		result, err := task.GetResult()
		if err != nil {
			t.Fatalf("获取任务结果失败: %v", err)
		}

		// 验证结果类型
		output, ok := result.(*obs.GetObjectMetadataOutput)
		if !ok {
			t.Fatal("结果类型错误，期望GetObjectMetadataOutput")
		}

		if output == nil {
			t.Fatal("GetResult返回nil")
		}

		// 验证结果内容
		if output.ContentLength != int64(len(testContent)) {
			t.Errorf("ContentLength不匹配，期望: %d, 实际: %d",
				len(testContent), output.ContentLength)
		}

		client.AddTestCase("GetResult返回正确结果")
	})
}

// TestTransferTask_ErrorConditions 测试边界条件和错误处理
func TestTransferTask_ErrorConditions(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()

	t.Run("ShouldFail_GivenPauseCompletedTask", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("pause-completed-test.bin")

		// 创建测试对象
		testContent := bytes.Repeat([]byte("J"), 4*1024*1024) // 4MB

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(testContent),
		}

		_, err := client.TestClient().PutObject(putInput)
		if err != nil {
			t.Fatalf("创建测试对象失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		localFilePath := filepath.Join(os.TempDir(), "pause-completed-test.bin")

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile: localFilePath,
			PartSize:     2 * 1024 * 1024, // 2MB分块
			TaskNum:      2,
		}

		// 启动异步任务
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步任务失败: %v", err)
		}

		// 等待任务完成
		<-task.Done()

		// 尝试暂停已完成的任务
		err = task.Pause()
		if err == nil {
			t.Error("期望暂停已完成任务失败，但实际成功")
		}

		if !strings.Contains(err.Error(), "cannot pause task in state") {
			t.Errorf("错误消息不正确，期望包含'cannot pause task in state'，实际: %v", err)
		}

		client.AddTestCase("暂停已完成任务错误处理正确")
	})

	t.Run("ShouldFail_GivenResumeRunningTask", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("resume-running-test.bin")

		// 创建测试对象
		testContent := bytes.Repeat([]byte("K"), 5*1024*1024) // 5MB

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(testContent),
		}

		_, err := client.TestClient().PutObject(putInput)
		if err != nil {
			t.Fatalf("创建测试对象失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		localFilePath := filepath.Join(os.TempDir(), "resume-running-test.bin")

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile: localFilePath,
			PartSize:     2 * 1024 * 1024, // 2MB分块
			TaskNum:      2,
		}

		// 启动异步任务
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步任务失败: %v", err)
		}

		// 尝试恢复正在运行的任务
		err = task.Resume()
		if err == nil {
			t.Error("期望恢复运行中任务失败，但实际成功")
		}

		if !strings.Contains(err.Error(), "cannot resume task in state") {
			t.Errorf("错误消息不正确，期望包含'cannot resume task in state'，实际: %v", err)
		}

		// 清理：取消任务
		task.Cancel()
		<-task.Done()

		client.AddTestCase("恢复运行中任务错误处理正确")
	})

	t.Run("ShouldFail_GivenCancelCompletedTask", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("cancel-completed-test.bin")

		// 创建测试对象
		testContent := bytes.Repeat([]byte("L"), 4*1024*1024) // 4MB

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(testContent),
		}

		_, err := client.TestClient().PutObject(putInput)
		if err != nil {
			t.Fatalf("创建测试对象失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		localFilePath := filepath.Join(os.TempDir(), "cancel-completed-test.bin")

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile: localFilePath,
			PartSize:     2 * 1024 * 1024, // 2MB分块
			TaskNum:      2,
		}

		// 启动异步任务
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步任务失败: %v", err)
		}

		// 等待任务完成
		<-task.Done()

		// 尝试取消已完成的任务
		err = task.Cancel()
		if err == nil {
			t.Error("期望取消已完成任务失败，但实际成功")
		}

		if !strings.Contains(err.Error(), "cannot cancel task in state") {
			t.Errorf("错误消息不正确，期望包含'cannot cancel task in state'，实际: %v", err)
		}

		client.AddTestCase("取消已完成任务错误处理正确")
	})

	t.Run("ShouldHandleConcurrentControl_GivenSameTask", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("concurrent-control-test.bin")

		// 创建测试对象
		testContent := bytes.Repeat([]byte("M"), 7*1024*1024) // 7MB

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(testContent),
		}

		_, err := client.TestClient().PutObject(putInput)
		if err != nil {
			t.Fatalf("创建测试对象失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		localFilePath := filepath.Join(os.TempDir(), "concurrent-control-test.bin")

		// 清理本地文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
		}()

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile: localFilePath,
			PartSize:     2 * 1024 * 1024, // 2MB分块
			TaskNum:      2,
		}

		// 启动异步任务
		task, err := client.TestClient().DownloadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步任务失败: %v", err)
		}

		// 等待一小段时间
		time.Sleep(300 * time.Millisecond)

		// 并发控制任务
		var wg sync.WaitGroup
		errChan := make(chan error, 4)

		// 并发暂停
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- task.Pause()
		}()

		// 并发恢复
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			errChan <- task.Resume()
		}()

		// 并发查询状态
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = task.Status()
		}()

		// 并发查询状态
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = task.Status()
		}()

		wg.Wait()
		close(errChan)

		// 等待任务完成或取消
		<-task.Done()

		// 验证至少有一些操作成功
		successCount := 0
		for err := range errChan {
			if err == nil {
				successCount++
			}
		}

		if successCount == 0 {
			t.Error("所有并发控制操作都失败了")
		}

		client.AddTestCase("并发控制任务处理正确")
	})
}

// TestResumeAsync_CheckpointRecovery 测试checkpoint恢复后的异步任务控制
func TestResumeAsync_CheckpointRecovery(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("checkpoint-recovery-test.bin")

	// 创建本地测试文件
	localFilePath := filepath.Join(os.TempDir(), "checkpoint-recovery-test.bin")
	checkpointFile := localFilePath + ".uploadfile_record"
	testContent := bytes.Repeat([]byte("N"), 16*1024*1024) // 16MB

	// 清理本地测试文件
	defer func() {
		if _, err := os.Stat(localFilePath); err == nil {
			os.Remove(localFilePath)
		}
		if _, err := os.Stat(checkpointFile); err == nil {
			os.Remove(checkpointFile)
		}
	}()

	t.Run("ShouldCreateTestFile_GivenValidContent", func(t *testing.T) {
		err := os.WriteFile(localFilePath, testContent, 0644)
		if err != nil {
			t.Fatalf("创建本地文件失败: %v", err)
		}

		client.AddTestCase("本地测试文件创建成功")
	})

	t.Run("ShouldPauseAndCancel_GivenFirstAttempt", func(t *testing.T) {
		input := &obs.UploadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			UploadFile:       localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
			CheckpointFile:   checkpointFile,
		}

		// 启动异步上传
		task, err := client.TestClient().UploadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步上传失败: %v", err)
		}

		// 等待一小段时间让任务开始
		time.Sleep(500 * time.Millisecond)

		// 暂停任务
		err = task.Pause()
		if err != nil {
			t.Fatalf("暂停任务失败: %v", err)
		}

		// 取消任务
		err = task.Cancel()
		if err != nil {
			t.Fatalf("取消任务失败: %v", err)
		}

		<-task.Done()

		// 验证checkpoint文件存在（因为被暂停过）
		if _, err := os.Stat(checkpointFile); os.IsNotExist(err) {
			t.Skip("checkpoint文件不存在，跳过恢复测试")
		}

		client.AddTestCase("第一次任务取消完成，checkpoint已创建")
	})

	t.Run("ShouldResumeFromCheckpoint_GivenSecondAttempt", func(t *testing.T) {
		// 使用相同的input重新启动上传，应该从checkpoint恢复
		input := &obs.UploadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			UploadFile:       localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: true,
			CheckpointFile:   checkpointFile,
		}

		// 从checkpoint恢复上传
		task, err := client.TestClient().UploadFileAsync(input)
		if err != nil {
			t.Fatalf("从checkpoint恢复上传失败: %v", err)
		}

		// 等待任务完成
		<-task.Done()

		// 验证最终状态
		if task.Status() != obs.TransferStatusCompleted {
			t.Errorf("期望状态COMPLETED，实际: %s", task.Status())
		}

		// 获取结果
		result, err := task.GetResult()
		if err != nil {
			t.Fatalf("获取上传结果失败: %v", err)
		}

		output, ok := result.(*obs.CompleteMultipartUploadOutput)
		if !ok || output == nil {
			t.Fatal("结果类型错误或为nil")
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除对象失败: %v", err)
			}
		})

		// 验证checkpoint文件被清理
		if _, err := os.Stat(checkpointFile); err == nil {
			t.Error("完成后checkpoint文件应该被清理")
		}

		client.AddTestCase("从checkpoint恢复并完成上传")
		t.Logf("从checkpoint恢复上传完成，对象: %s, 大小: %d bytes", objectKey, output.ContentLength)
	})

	t.Run("ShouldVerifyUploadedObject_GivenCheckpointRecovery", func(t *testing.T) {
		// 验证上传的对象
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(input)
		if err != nil {
			t.Fatalf("获取上传对象失败: %v", err)
		}
		defer output.Body.Close()

		// 读取内容
		downloadedContent, err := io.ReadAll(output.Body)
		if err != nil {
			t.Fatalf("读取下载内容失败: %v", err)
		}

		// 验证内容完整性
		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("checkpoint恢复后上传对象验证成功")
	})
}

// TestResumeAsync_ContextCancellation 测试上下文取消与异步任务交互
func TestResumeAsync_ContextCancellation(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)
	defer client.PrintTestCases()

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("context-cancel-test.bin")

	t.Run("ShouldHandleContextCancel_GivenAsyncUpload", func(t *testing.T) {
		// 创建本地测试文件
		localFilePath := filepath.Join(os.TempDir(), "context-cancel-test.bin")
		testContent := bytes.Repeat([]byte("O"), 15*1024*1024) // 15MB

		// 清理本地测试文件
		defer func() {
			if _, err := os.Stat(localFilePath); err == nil {
				os.Remove(localFilePath)
			}
		}()

		err := os.WriteFile(localFilePath, testContent, 0644)
		if err != nil {
			t.Fatalf("创建本地文件失败: %v", err)
		}

		// 使用带超时的context
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// 注意：UploadFileAsync不支持context，这里测试常规异步上传
		input := &obs.UploadFileInput{
			Bucket:           bucket,
			Key:              objectKey,
			UploadFile:       localFilePath,
			PartSize:         5 * 1024 * 1024, // 5MB分块
			TaskNum:          3,
			EnableCheckpoint: false,
		}

		// 启动异步上传
		task, err := client.TestClient().UploadFileAsync(input)
		if err != nil {
			t.Fatalf("启动异步上传失败: %v", err)
		}

		// 等待一小段时间后取消
		select {
		case <-ctx.Done():
			// 超时，取消任务
			task.Cancel()
		case <-task.Done():
			// 任务已完成
		}

		// 等待任务完全停止
		<-task.Done()

		// 如果任务被取消，验证状态
		if task.Status() == obs.TransferStatusCancelled {
			client.AddTestCase("上下文超时导致任务取消成功")
			return
		}

		// 如果任务已完成，清理对象
		if task.Status() == obs.TransferStatusCompleted {
			client.AddCleanup(func(t *testing.T) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    objectKey,
				}
				_, err := client.TestClient().DeleteObject(deleteInput)
				if err != nil {
					t.Logf("删除对象失败: %v", err)
				}
			})
			client.AddTestCase("异步上传在超时前完成")
		}
	})
}
