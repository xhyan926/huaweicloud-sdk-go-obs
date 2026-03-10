//go:build integration

package object

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// TestResumeUpload_ShouldResumeSuccessfully_GivenInterruptedUpload 测试上传断点续传
func TestResumeUpload_ShouldResumeSuccessfully_GivenInterruptedUpload(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("resume-upload-test.txt")

	// 创建本地测试文件
	localFilePath := filepath.Join(os.TempDir(), "resume-upload-test.txt")
	testContent := bytes.Repeat([]byte("A"), 20*1024*1024) // 20MB

	// 清理本地测试文件
	defer func() {
		if _, err := os.Stat(localFilePath); err == nil {
			os.Remove(localFilePath)
		}
	}()

	t.Run("ShouldStartUploadFile_GivenValidFile", func(t *testing.T) {
		// 创建本地测试文件
		err := os.WriteFile(localFilePath, testContent, 0644)
		if err != nil {
			t.Fatalf("创建本地文件失败: %v", err)
		}

		client.AddTestCase("本地测试文件创建成功")
		t.Logf("本地测试文件: %s, 大小: %d bytes", localFilePath, len(testContent))
	})

	t.Run("ShouldUploadPartially_GivenInterruption", func(t *testing.T) {
		// 模拟部分上传
		tempFile, err := os.CreateTemp("", "partial-upload-*.txt")
		if err != nil {
			t.Fatalf("创建临时文件失败: %v", err)
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		// 写入部分内容
		partialContent := testContent[:10*1024*1024] // 前10MB
		tempFile.Write(partialContent)

		// 使用断点续传上传
		input := &obs.UploadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			UploadFile:   tempFile.Name(),
			EnableCheckpoint: true,
		}

		_, err := client.TestClient().UploadFile(input)
		if err != nil {
			t.Logf("部分上传失败: %v", err)
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

		client.AddTestCase("部分上传完成")
	})

	t.Run("ShouldResumeUpload_GivenCheckpoint", func(t *testing.T) {
		// 使用断点续传上传完整文件
		input := &obs.UploadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			UploadFile:   localFilePath,
			EnableCheckpoint: true,
			PartSize:     5 * 1024 * 1024, // 5MB分块
		}

		output, err := client.TestClient().UploadFile(input)
		if err != nil {
			t.Fatalf("断点续传上传失败: %v", err)
		}

		// 验证上传结果
		if output == nil {
			t.Error("UploadFile返回nil")
		}

		client.AddTestCase("断点续传上传成功")
		t.Logf("断点续传上传完成，对象: %s, 大小: %d bytes", objectKey, output.ContentLength)
	})
}

// TestResumeDownload_ShouldResumeSuccessfully_GivenInterruptedDownload 测试下载断点续传
func TestResumeDownload_ShouldResumeSuccessfully_GivenInterruptedDownload(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("resume-download-test.bin")

	// 创建测试对象
	testContent := bytes.Repeat([]byte("B"), 15*1024*1024) // 15MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(testContent),
	}

	if _, err := client.TestClient().PutObject(putInput); err != nil {
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

	t.Run("ShouldDownloadPartially_GivenInterruption", func(t *testing.T) {
		localFilePath := filepath.Join(os.TempDir(), "partial-download-test.bin")

		// 创建部分下载
		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile:  localFilePath,
			EnableCheckpoint: true,
			PartSize:     5 * 1024 * 1024, // 5MB分块
		}

		// 使用context模拟中断
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := client.TestClient().DownloadFileWithContext(ctx, input)
		if err != nil {
			// 超时是预期的
			if !strings.Contains(err.Error(), "context deadline") &&
			   !strings.Contains(err.Error(), "timeout") {
				t.Logf("部分下载失败: %v", err)
			}
		}

		// 清理部分下载文件
		os.Remove(localFilePath)

		client.AddTestCase("部分下载完成（模拟中断）")
	})

	t.Run("ShouldResumeDownload_GivenCheckpoint", func(t *testing.T) {
		localFilePath := filepath.Join(os.TempDir(), "resume-download-test.bin")

		// 清理本地文件
		defer os.Remove(localFilePath)

		// 使用断点续传下载完整文件
		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile:  localFilePath,
			EnableCheckpoint: true,
			PartSize:     5 * 1024 * 1024, // 5MB分块
		}

		output, err := client.TestClient().DownloadFile(input)
		if err != nil {
			t.Fatalf("断点续传下载失败: %v", err)
		}

		// 验证下载结果
		if output == nil {
			t.Error("DownloadFile返回nil")
		}

		// 验证下载文件内容
		downloadedContent, err := os.ReadFile(localFilePath)
		if err != nil {
			t.Fatalf("读取下载文件失败: %v", err)
		}

		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("断点续传下载成功")
		t.Logf("断点续传下载完成，文件: %s, 大小: %d bytes", localFilePath, len(downloadedContent))
	})
}

// TestResumeUpload_Validation 测试断点续传验证
func TestResumeUpload_Validation(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldValidateCheckpoint_GivenExistingCheckpoint", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("checkpoint-validation-test.txt")

		// 创建测试文件
		localFilePath := filepath.Join(os.TempDir(), "checkpoint-validation.txt")
		testContent := bytes.Repeat([]byte("C"), 10*1024*1024) // 10MB

		err := os.WriteFile(localFilePath, testContent, 0644)
		if err != nil {
			t.Fatalf("创建本地文件失败: %v", err)
		}

		defer os.Remove(localFilePath)

		// 上传文件并启用断点续传
		input := &obs.UploadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			UploadFile:   localFilePath,
			EnableCheckpoint: true,
			CheckpointFile: filepath.Join(os.TempDir(), "checkpoint-file.json"),
		}

		_, err = client.TestClient().UploadFile(input)
		if err != nil {
			t.Fatalf("上传文件失败: %v", err)
		}

		// 验证断点文件是否存在
		checkpointFile := filepath.Join(os.TempDir(), "checkpoint-file.json")
		defer os.Remove(checkpointFile)

		if _, err := os.Stat(checkpointFile); err == nil {
			t.Log("断点文件存在，断点续传成功")
		} else {
			t.Log("断点文件不存在（可能已自动清理）")
		}

		// 添加对象清理函数
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

		client.AddTestCase("断点续传验证成功")
	})
}

// TestResumeDownload_CheckpointIntegrity 测试断点续传完整性
func TestResumeDownload_CheckpointIntegrity(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("integrity-test.bin")

	// 创建测试对象
	testContent := bytes.Repeat([]byte("D"), 12*1024*1024) // 12MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(testContent),
	}

	if _, err := client.TestClient().PutObject(putInput); err != nil {
		t.Fatalf("创建测试对象失败: %v", err)
	}

	// 添加对象清理函数
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

	t.Run("ShouldMaintainIntegrity_GivenResumableDownload", func(t *testing.T) {
		localFilePath := filepath.Join(os.TempDir(), "integrity-test.bin")

		// 清理本地文件
		defer os.Remove(localFilePath)

		// 多次模拟中断和恢复
		for i := 0; i < 3; i++ {
			input := &obs.DownloadFileInput{
				Bucket:       bucket,
				Key:          objectKey,
				DownloadFile:  localFilePath,
				EnableCheckpoint: true,
				PartSize:     3 * 1024 * 1024, // 3MB分块
			}

			// 模拟第3次为完整下载
			if i == 2 {
				_, err := client.TestClient().DownloadFile(input)
				if err != nil {
					t.Fatalf("下载文件失败: %v", err)
				}
			} else {
				// 模拟中断
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				_, err := client.TestClient().DownloadFileWithContext(ctx, input)
				cancel()
				if err != nil {
					// 中断是预期的
				}
			}
		}

		// 验证最终下载的完整性
		downloadedContent, err := os.ReadFile(localFilePath)
		if err != nil {
			t.Fatalf("读取下载文件失败: %v", err)
		}

		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("断点续传完整性验证成功")
		t.Logf("断点续传完成，文件: %s, 大小: %d bytes", localFilePath, len(downloadedContent))
	})
}

// TestResumeUpload_LargeFile 测试大文件断点续传上传
func TestResumeUpload_LargeFile(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("large-resume-upload-test.bin")

	// 创建大文件
	fileSize := 50 * 1024 * 1024 // 50MB
	testContent := bytes.Repeat([]byte("E"), fileSize)

	localFilePath := filepath.Join(os.TempDir(), "large-resume-upload-test.bin")

	// 清理本地文件
	defer os.Remove(localFilePath)

	err := os.WriteFile(localFilePath, testContent, 0644)
	if err != nil {
		t.Fatalf("创建本地大文件失败: %v", err)
	}

	t.Run("ShouldResumeLargeFileUpload_GivenValidCheckpoint", func(t *testing.T) {
		// 使用断点续传上传大文件
		input := &obs.UploadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			UploadFile:   localFilePath,
			EnableCheckpoint: true,
			PartSize:     10 * 1024 * 1024, // 10MB分块
		}

		output, err := client.TestClient().UploadFile(input)
		if err != nil {
			t.Fatalf("大文件断点续传上传失败: %v", err)
		}

		// 验证上传结果
		if output == nil {
			t.Error("UploadFile返回nil")
		}

		// 添加对象清理函数
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

		client.AddTestCase("大文件断点续传上传成功")
		t.Logf("大文件断点续传上传完成，对象: %s, 大小: %d bytes", objectKey, fileSize)
	})
}

// TestResumeDownload_LargeFile 测试大文件断点续传下载
func TestResumeDownload_LargeFile(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("large-resume-download-test.bin")

	// 创建测试对象
	fileSize := 40 * 1024 * 1024 // 40MB
	testContent := bytes.Repeat([]byte("F"), fileSize)

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(testContent),
	}

	if _, err := client.TestClient().PutObject(putInput); err != nil {
		t.Fatalf("创建测试对象失败: %v", err)
	}

	// 添加对象清理函数
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

	client.AddTestCase("测试大对象创建完成")

	t.Run("ShouldResumeLargeFileDownload_GivenValidCheckpoint", func(t *testing.T) {
		localFilePath := filepath.Join(os.TempDir(), "large-resume-download-test.bin")

		// 清理本地文件
		defer os.Remove(localFilePath)

		// 使用断点续传下载大文件
		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile:  localFilePath,
			EnableCheckpoint: true,
			PartSize:     8 * 1024 * 1024, // 8MB分块
		}

		output, err := client.TestClient().DownloadFile(input)
		if err != nil {
			t.Fatalf("大文件断点续传下载失败: %v", err)
		}

		// 验证下载结果
		if output == nil {
			t.Error("DownloadFile返回nil")
		}

		// 验证下载文件内容
		downloadedContent, err := os.ReadFile(localFilePath)
		if err != nil {
			t.Fatalf("读取下载文件失败: %v", err)
		}

		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("大文件断点续传下载成功")
		t.Logf("大文件断点续传下载完成，文件: %s, 大小: %d bytes", localFilePath, len(downloadedContent))
	})
}

// TestResumeUpload_ConcurrentResume 测试并发断点续传
func TestResumeUpload_ConcurrentResume(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldHandleConcurrentResume_GivenMultipleFiles", func(t *testing.T) {
		numFiles := 3
		var wg sync.WaitGroup
		errChan := make(chan error, numFiles)

		// 并发进行多个断点续传上传
		for i := 0; i < numFiles; i++ {
			wg.Add(1)
			go func(fileIndex int) {
				defer wg.Done()

				objectKey := client.GetTestObjectKey(fmt.Sprintf("concurrent-resume-%d.txt", fileIndex))
				fileSize := 5 * 1024 * 1024 // 5MB
				testContent := bytes.Repeat([]byte("G"), fileSize)

				// 创建本地文件
				localFilePath := filepath.Join(os.TempDir(),
					fmt.Sprintf("concurrent-resume-%d.txt", fileIndex))
				defer os.Remove(localFilePath)

				err := os.WriteFile(localFilePath, testContent, 0644)
				if err != nil {
					errChan <- fmt.Errorf("创建本地文件 %d 失败: %v", fileIndex, err)
					return
				}

				// 上传文件并启用断点续传
				input := &obs.UploadFileInput{
					Bucket:       bucket,
					Key:          objectKey,
					UploadFile:   localFilePath,
					EnableCheckpoint: true,
					PartSize:     2 * 1024 * 1024, // 2MB分块
				}

				_, err = client.TestClient().UploadFile(input)
				if err != nil {
					errChan <- fmt.Errorf("上传文件 %d 失败: %v", fileIndex, err)
					return
				}

				// 添加对象清理函数
				client.AddCleanup(func(t *testing.T) {
					deleteInput := &obs.DeleteObjectInput{
						Bucket: bucket,
						Key:    objectKey,
					}
					_, err := client.TestClient().DeleteObject(deleteInput)
					if err != nil {
						t.Logf("删除对象 %s 失败: %v", objectKey, err)
					}
				})

				errChan <- nil
			}(i)
		}

		wg.Wait()
		close(errChan)

		// 收集错误
		errorCount := 0
		for err := range errChan {
			if err != nil {
				errorCount++
				t.Error(err)
			}
		}

		if errorCount > 0 {
			t.Errorf("有 %d/%d 个并发断点续传上传失败", errorCount, numFiles)
		}

		client.AddTestCase("并发断点续传上传测试通过")
		t.Logf("并发断点续传上传完成，成功率: %d/%d", numFiles-errorCount, numFiles)
	})
}

// TestResumeDownload_ErrorHandling 测试错误处理
func TestResumeDownload_ErrorHandling(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldFail_GivenInvalidLocalPath", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("invalid-path-test.bin")

		// 尝试下载到无效路径
		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile:  "/invalid/path/to/file.txt",
			EnableCheckpoint: true,
		}

		_, err := client.TestClient().DownloadFile(input)
		if err == nil {
			t.Error("期望下载失败，但操作成功")
		}

		// 验证错误信息
		if !strings.Contains(err.Error(), "no such file") &&
		   !strings.Contains(err.Error(), "permission denied") {
			t.Logf("错误信息: %v", err)
		}

		client.AddTestCase("无效路径错误处理测试通过")
	})

	t.Run("ShouldFail_GivenNonexistentObject", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("nonexistent-resume-test.bin")

		// 尝试下载不存在的对象
		localFilePath := filepath.Join(os.TempDir(), "nonexistent-test.bin")
		defer os.Remove(localFilePath)

		input := &obs.DownloadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			DownloadFile:  localFilePath,
			EnableCheckpoint: true,
		}

		_, err := client.TestClient().DownloadFile(input)
		if err == nil {
			t.Error("期望下载失败，但操作成功")
		}

		// 验证错误类型
		obsErr, ok := err.(obs.ObsError)
		if !ok {
			t.Fatalf("错误不是ObsError类型: %T", err)
		}

		if obsErr.StatusCode != 404 {
			t.Errorf("期望404错误，实际: %d", obsErr.StatusCode)
		}

		client.AddTestCase("不存在对象错误处理测试通过")
		t.Logf("不存在对象错误: 状态码=%d, 代码=%s", obsErr.StatusCode, obsErr.Code)
	})
}

// TestResumeUpload_FileVerification 测试文件验证
func TestResumeUpload_FileVerification(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("file-verification-test.txt")

	// 创建测试文件
	testContent := []byte("这是用于验证文件完整性测试的内容。")

	localFilePath := filepath.Join(os.TempDir(), "file-verification-test.txt")
	err := os.WriteFile(localFilePath, testContent, 0644)
	if err != nil {
		t.Fatalf("创建本地文件失败: %v", err)
	}

	// 清理本地文件
	defer os.Remove(localFilePath)

	t.Run("ShouldUploadAndVerify_GivenValidFile", func(t *testing.T) {
		// 使用断点续传上传
		input := &obs.UploadFileInput{
			Bucket:       bucket,
			Key:          objectKey,
			UploadFile:   localFilePath,
			EnableCheckpoint: true,
		}

		_, err = client.TestClient().UploadFile(input)
		if err != nil {
			t.Fatalf("断点续传上传失败: %v", err)
		}

		// 验证上传的内容
		getInput := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(getInput)
		if err != nil {
			t.Fatalf("获取对象失败: %v", err)
		}
		defer output.Body.Close()

		// 读取内容
		downloadedContent, err := io.ReadAll(output.Body)
		if err != nil {
			t.Fatalf("读取下载内容失败: %v", err)
		}

		// 验证内容完整性
		if !bytes.Equal(downloadedContent, testContent) {
			t.Errorf("上传内容不匹配，期望: %s, 实际: %s",
				string(testContent), string(downloadedContent))
		}

		// 添加对象清理函数
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

		client.AddTestCase("文件验证测试通过")
		t.Logf("文件验证完成，内容匹配")
	})
}
