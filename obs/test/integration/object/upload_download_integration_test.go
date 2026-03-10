//go:build integration

package object

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// TestUploadDownload_ShouldUploadAndDownloadSuccessfully_GivenSmallFile 测试小文件上传下载
func TestUploadDownload_ShouldUploadAndDownloadSuccessfully_GivenSmallFile(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("small-file-test.txt")
	content := "这是用于测试小文件上传和下载的内容。"

	t.Run("ShouldUploadSuccessfully_GivenValidSmallFile", func(t *testing.T) {
		// 上传小文件
		input := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   strings.NewReader(content),
		}

		output, err := client.TestClient().PutObject(input)
		if err != nil {
			t.Fatalf("上传小文件失败: %v", err)
		}

		// 验证上传结果
		if output == nil {
			t.Error("PutObject返回nil")
		}

		if output.VersionId == "" {
			t.Log("对象版本ID为空（可能未启用版本控制）")
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除小文件失败: %v", err)
			}
		})

		client.AddTestCase("小文件上传成功")
		t.Logf("小文件上传成功: %s", objectKey)
	})

	t.Run("ShouldDownloadSuccessfully_GivenExistingSmallFile", func(t *testing.T) {
		// 下载小文件
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(input)
		if err != nil {
			t.Fatalf("下载小文件失败: %v", err)
		}
		defer output.Body.Close()

		// 读取内容
		downloadedContent, err := io.ReadAll(output.Body)
		if err != nil {
			t.Fatalf("读取下载内容失败: %v", err)
		}

		// 验证内容完整性
		if string(downloadedContent) != content {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(content), len(downloadedContent))
		}

		// 验证元数据
		if output.ContentLength != int64(len(content)) {
			t.Errorf("内容长度不匹配，期望: %d, 实际: %d",
				len(content), output.ContentLength)
		}

		client.AddTestCase("小文件下载成功")
		t.Logf("小文件下载成功: %s, 大小: %d bytes", objectKey, len(downloadedContent))
	})
}

// TestUploadDownload_ShouldUploadAndDownloadSuccessfully_GivenLargeFile 测试大文件上传下载
func TestUploadDownload_ShouldUploadAndDownloadSuccessfully_GivenLargeFile(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("large-file-test.bin")

	// 创建大文件内容（10MB）
	fileSize := 10 * 1024 * 1024
	content := bytes.Repeat([]byte("A"), fileSize)

	t.Run("ShouldUploadSuccessfully_GivenValidLargeFile", func(t *testing.T) {
		// 上传大文件
		input := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(content),
		}

		output, err := client.TestClient().PutObject(input)
		if err != nil {
			t.Fatalf("上传大文件失败: %v", err)
		}

		// 验证上传结果
		if output == nil {
			t.Error("PutObject返回nil")
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除大文件失败: %v", err)
			}
		})

		client.AddTestCase("大文件上传成功")
		t.Logf("大文件上传成功: %s, 大小: %d bytes", objectKey, fileSize)
	})

	t.Run("ShouldDownloadSuccessfully_GivenExistingLargeFile", func(t *testing.T) {
		// 下载大文件
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(input)
		if err != nil {
			t.Fatalf("下载大文件失败: %v", err)
		}
		defer output.Body.Close()

		// 读取内容
		downloadedContent, err := io.ReadAll(output.Body)
		if err != nil {
			t.Fatalf("读取下载内容失败: %v", err)
		}

		// 验证内容完整性
		if !bytes.Equal(downloadedContent, content) {
			t.Errorf("下载内容不匹配，期望长度: %d, 实际长度: %d",
				len(content), len(downloadedContent))
		}

		// 验证元数据
		if output.ContentLength != int64(fileSize) {
			t.Errorf("内容长度不匹配，期望: %d, 实际: %d",
				fileSize, output.ContentLength)
		}

		client.AddTestCase("大文件下载成功")
		t.Logf("大文件下载成功: %s, 大小: %d bytes", objectKey, len(downloadedContent))
	})
}

// TestUploadDownload_ConcurrentUploads 测试并发上传
func TestUploadDownload_ConcurrentUploads(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldHandleConcurrentUploads_GivenMultipleFiles", func(t *testing.T) {
		numConcurrent := 10
		var wg sync.WaitGroup
		errChan := make(chan error, numConcurrent)
		objectKeys := make([]string, numConcurrent)

		// 并发上传多个文件
		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				objectKey := client.GetTestObjectKey(
					client.GenerateUniqueKey(fmt.Sprintf("concurrent-upload-%d.txt", index)))
				content := fmt.Sprintf("并发上传测试内容 %d", index)

				input := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    objectKey,
					Body:   strings.NewReader(content),
				}

				_, err := client.TestClient().PutObject(input)
				if err != nil {
					errChan <- fmt.Errorf("并发上传 %d 失败: %v", index, err)
					return
				}

				objectKeys[index] = objectKey

				// 注册清理函数
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
			t.Errorf("有 %d/%d 个并发上传失败", errorCount, numConcurrent)
		}

		client.AddTestCase("并发上传测试通过")
		t.Logf("并发上传完成，成功率: %d/%d", numConcurrent-errorCount, numConcurrent)
	})
}

// TestUploadDownload_ConcurrentDownloads 测试并发下载
func TestUploadDownload_ConcurrentDownloads(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	numObjects := 5
	objectKeys := make([]string, numObjects)

	// 准备测试对象
	t.Run("ShouldPrepareTestObjects_GivenMultipleFiles", func(t *testing.T) {
		for i := 0; i < numObjects; i++ {
			objectKey := client.GetTestObjectKey(fmt.Sprintf("download-test-%d.txt", i))
			content := fmt.Sprintf("下载测试内容 %d", i)

			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    objectKey,
				Body:   strings.NewReader(content),
			}

			if _, err := client.TestClient().PutObject(input); err != nil {
				t.Fatalf("准备测试对象 %d 失败: %v", i, err)
			}

			objectKeys[i] = objectKey

			// 注册清理函数
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
		}

		client.AddTestCase("测试对象准备完成")
	})

	t.Run("ShouldHandleConcurrentDownloads_GivenMultipleFiles", func(t *testing.T) {
		var wg sync.WaitGroup
		errChan := make(chan error, numObjects)

		// 并发下载多个文件
		for i := 0; i < numObjects; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				input := &obs.GetObjectInput{
					Bucket: bucket,
					Key:    objectKeys[index],
				}

				output, err := client.TestClient().GetObject(input)
				if err != nil {
					errChan <- fmt.Errorf("并发下载 %d 失败: %v", index, err)
					return
				}
				defer output.Body.Close()

				// 读取内容验证
				_, err = io.ReadAll(output.Body)
				if err != nil {
					errChan <- fmt.Errorf("读取下载内容 %d 失败: %v", index, err)
					return
				}

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
			t.Errorf("有 %d/%d 个并发下载失败", errorCount, numObjects)
		}

		client.AddTestCase("并发下载测试通过")
		t.Logf("并发下载完成，成功率: %d/%d", numObjects-errorCount, numObjects)
	})
}

// TestUploadDownload_FileIntegrity 测试文件完整性验证
func TestUploadDownload_FileIntegrity(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("integrity-test.bin")

	// 创建测试内容（包含各种字节数据）
	testContent := make([]byte, 10000)
	for i := range testContent {
		testContent[i] = byte(i % 256)
	}

	t.Run("ShouldVerifyContentIntegrity_GivenUploadedFile", func(t *testing.T) {
		// 上传文件
		input := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(testContent),
		}

		if _, err := client.TestClient().PutObject(input); err != nil {
			t.Fatalf("上传文件失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除文件失败: %v", err)
			}
		})

		client.AddTestCase("文件上传完成")
	})

	t.Run("ShouldVerifyContentIntegrity_GivenDownloadedFile", func(t *testing.T) {
		// 下载文件
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(input)
		if err != nil {
			t.Fatalf("下载文件失败: %v", err)
		}
		defer output.Body.Close()

		// 读取内容
		downloadedContent, err := io.ReadAll(output.Body)
		if err != nil {
			t.Fatalf("读取下载内容失败: %v", err)
		}

		// 逐字节验证内容
		for i := 0; i < len(testContent); i++ {
			if i < len(downloadedContent) {
				if downloadedContent[i] != testContent[i] {
					t.Errorf("字节 %d 不匹配，期望: %d, 实际: %d",
						i, testContent[i], downloadedContent[i])
				}
			}
		}

		// 验证总长度
		if len(downloadedContent) != len(testContent) {
			t.Errorf("文件长度不匹配，期望: %d, 实际: %d",
				len(testContent), len(downloadedContent))
		}

		client.AddTestCase("文件完整性验证通过")
		t.Logf("文件完整性验证通过，文件大小: %d bytes", len(downloadedContent))
	})
}

// TestUploadDownload_MetadataOperations 测试元数据操作
func TestUploadDownload_MetadataOperations(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("metadata-test.txt")
	content := "测试元数据操作的内容。"

	t.Run("ShouldUploadWithMetadata_GivenValidFile", func(t *testing.T) {
		// 上传文件并设置元数据
		input := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   strings.NewReader(content),
			Metadata: map[string]string{
				"author":        "test-user",
				"description":    "测试文件",
				"upload-time":    "2024-01-01",
				"content-type":   "text/plain",
				"custom-header1": "value1",
				"custom-header2": "value2",
			},
		}

		if _, err := client.TestClient().PutObject(input); err != nil {
			t.Fatalf("上传文件失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除文件失败: %v", err)
			}
		})

		client.AddTestCase("带元数据的文件上传成功")
	})

	t.Run("ShouldRetrieveMetadata_GivenExistingFile", func(t *testing.T) {
		// 获取对象元数据
		input := &obs.GetObjectMetadataInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObjectMetadata(input)
		if err != nil {
			t.Fatalf("获取对象元数据失败: %v", err)
		}

		// 验证元数据
		if output.Metadata == nil {
			t.Error("元数据为空")
		}

		// 验证特定元数据字段
		expectedMetadata := map[string]string{
			"author":        "test-user",
			"description":    "测试文件",
			"upload-time":    "2024-01-01",
			"content-type":   "text/plain",
			"custom-header1": "value1",
			"custom-header2": "value2",
		}

		for key, expectedValue := range expectedMetadata {
			actualValue, exists := output.Metadata[key]
			if !exists {
				t.Errorf("元数据字段 %s 不存在", key)
			} else if actualValue != expectedValue {
				t.Errorf("元数据字段 %s 不匹配，期望: %s, 实际: %s",
					key, expectedValue, actualValue)
			}
		}

		client.AddTestCase("元数据检索成功")
		t.Logf("元数据检索成功，字段数量: %d", len(output.Metadata))
	})
}

// TestUploadDownload_RangeDownload 测试范围下载
func TestUploadDownload_RangeDownload(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("range-download-test.txt")
	content := "这是用于测试范围下载的文件内容。" +
		"它包含足够的文本内容以支持不同范围的下载。"

	t.Run("ShouldUploadRangeTestFile_GivenValidContent", func(t *testing.T) {
		// 上传测试文件
		input := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   strings.NewReader(content),
		}

		if _, err := client.TestClient().PutObject(input); err != nil {
			t.Fatalf("上传文件失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除文件失败: %v", err)
			}
		})

		client.AddTestCase("范围测试文件上传完成")
	})

	t.Run("ShouldDownloadWithRange_GivenValidRange", func(t *testing.T) {
		// 测试不同范围的下载
		testRanges := []struct {
			name  string
			range string
			start int
			end   int
		}{
			{"前半部分", "0-19", 0, 20},
			{"后半部分", "20-", 20, len(content)},
			{"中间部分", "10-29", 10, 30},
		}

		for _, test := range testRanges {
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
				Range:  test.range,
			}

			output, err := client.TestClient().GetObject(input)
			if err != nil {
				t.Fatalf("范围下载失败 (%s): %v", test.name, err)
			}
			defer output.Body.Close()

			// 读取内容
			downloadedContent, err := io.ReadAll(output.Body)
			if err != nil {
				t.Fatalf("读取下载内容失败: %v", err)
			}

			// 验证内容
			expectedContent := content[test.start:test.end]
			if string(downloadedContent) != expectedContent {
				t.Errorf("范围下载内容不匹配 (%s)，期望: %s, 实际: %s",
					test.name, expectedContent, string(downloadedContent))
			}

			t.Logf("范围下载成功 (%s): %s, 大小: %d bytes",
				test.name, test.range, len(downloadedContent))
		}

		client.AddTestCase("范围下载测试通过")
	})
}

// TestUploadDownload_ErrorHandling 测试错误处理
func TestUploadDownload_ErrorHandling(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldFail_GivenNonexistentObject", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("nonexistent-object.txt")

		// 尝试下载不存在的对象
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		_, err := client.TestClient().GetObject(input)
		if err == nil {
			t.Error("期望下载失败，但操作成功")
		}

		// 验证错误类型
		obsErr, ok := err.(obs.ObsError)
		if !ok {
			t.Fatalf("错误不是ObsError类型: %T", err)
		}

		// 验证错误信息
		if obsErr.StatusCode != 404 {
			t.Errorf("期望404错误，实际: %d", obsErr.StatusCode)
		}

		if obsErr.Code != "NoSuchKey" {
			t.Errorf("期望NoSuchKey错误代码，实际: %s", obsErr.Code)
		}

		client.AddTestCase("不存在对象错误处理测试通过")
		t.Logf("不存在对象错误: 状态码=%d, 代码=%s", obsErr.StatusCode, obsErr.Code)
	})

	t.Run("ShouldFail_GivenInvalidBucket", func(t *testing.T) {
		objectKey := "test-object.txt"
		invalidBucket := "nonexistent-bucket-12345"

		// 尝试向不存在的桶上传
		input := &obs.PutObjectInput{
			Bucket: invalidBucket,
			Key:    objectKey,
			Body:   strings.NewReader("test content"),
		}

		_, err := client.TestClient().PutObject(input)
		if err == nil {
			t.Error("期望上传失败，但操作成功")
		}

		// 验证错误类型
		obsErr, ok := err.(obs.ObsError)
		if !ok {
			t.Fatalf("错误不是ObsError类型: %T", err)
		}

		// 验证错误信息
		if obsErr.StatusCode != 404 {
			t.Errorf("期望404错误，实际: %d", obsErr.StatusCode)
		}

		client.AddTestCase("无效桶错误处理测试通过")
		t.Logf("无效桶错误: 状态码=%d, 代码=%s", obsErr.StatusCode, obsErr.Code)
	})
}

// TestUploadDownload_ContentType 测试内容类型
func TestUploadDownload_ContentType(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldUploadWithContentType_GivenTextFile", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("content-type-test.txt")
		content := "这是文本文件内容。"

		// 上传文本文件并指定内容类型
		input := &obs.PutObjectInput{
			Bucket:      bucket,
			Key:         objectKey,
			Body:        strings.NewReader(content),
			ContentType: "text/plain; charset=utf-8",
		}

		if _, err := client.TestClient().PutObject(input); err != nil {
			t.Fatalf("上传文件失败: %v", err)
		}

		// 添加清理函数
		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			_, err := client.TestClient().DeleteObject(deleteInput)
			if err != nil {
				t.Logf("删除文件失败: %v", err)
			}
		})

		client.AddTestCase("带内容类型的文件上传成功")
	})

	t.Run("ShouldVerifyContentType_GivenDownloadedFile", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("content-type-test.txt")

		// 获取对象元数据验证内容类型
		input := &obs.GetObjectMetadataInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObjectMetadata(input)
		if err != nil {
			t.Fatalf("获取对象元数据失败: %v", err)
		}

		// 验证内容类型
		expectedContentType := "text/plain"
		if !strings.HasPrefix(output.ContentType, expectedContentType) {
			t.Errorf("内容类型不匹配，期望以 %s 开头，实际: %s",
				expectedContentType, output.ContentType)
		}

		client.AddTestCase("内容类型验证通过")
		t.Logf("内容类型验证通过: %s", output.ContentType)
	})
}
