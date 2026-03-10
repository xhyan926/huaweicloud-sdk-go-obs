//go:build integration

package object

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// TestMultipartUpload_ShouldUploadSuccessfully_GivenBasicMultipart 测试基本分块上传
func TestMultipartUpload_ShouldUploadSuccessfully_GivenBasicMultipart(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("multipart-test.txt")

	t.Run("ShouldInitiateMultipartUpload_GivenValidInput", func(t *testing.T) {
		// 初始化分块上传
		input := &obs.InitiateMultipartUploadInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().InitiateMultipartUpload(input)
		if err != nil {
			t.Fatalf("初始化分块上传失败: %v", err)
		}

		// 验证返回结果
		if output == nil {
			t.Error("InitiateMultipartUpload返回nil")
		}

		if output.UploadId == "" {
			t.Error("UploadId为空")
		}

		// 保存upload ID用于后续测试
		t.Logf("分块上传初始化成功，UploadId: %s", output.UploadId)

		// 添加取消清理函数
		client.AddCleanup(func(t *testing.T) {
			abortInput := &obs.AbortMultipartUploadInput{
				Bucket:   bucket,
				Key:      objectKey,
				UploadId: output.UploadId,
			}
			_, err := client.TestClient().AbortMultipartUpload(abortInput)
			if err != nil {
				t.Logf("取消分块上传失败: %v", err)
			}
		})

		client.AddTestCase("分块上传初始化成功")
	})
}

// TestMultipartUpload_ShouldUploadPartsSuccessfully_GivenValidUploadId 测试分块上传
func TestMultipartUpload_ShouldUploadPartsSuccessfully_GivenValidUploadId(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("multipart-parts-test.txt")
	partSize := 5 * 1024 * 1024 // 5MB
	numParts := 3

	// 初始化分块上传
	initInput := &obs.InitiateMultipartUploadInput{
		Bucket: bucket,
		Key:    objectKey,
	}

	initOutput, err := client.TestClient().InitiateMultipartUpload(initInput)
	if err != nil {
		t.Fatalf("初始化分块上传失败: %v", err)
	}

	// 添加取消清理函数
	client.AddCleanup(func(t *testing.T) {
		abortInput := &obs.AbortMultipartUploadInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: initOutput.UploadId,
		}
		_, err := client.TestClient().AbortMultipartUpload(abortInput)
		if err != nil {
			t.Logf("取消分块上传失败: %v", err)
		}
	})

	t.Run("ShouldUploadParts_GivenValidUploadId", func(t *testing.T) {
		parts := make([]obs.Part, numParts)

		// 上传多个分块
		for i := 0; i < numParts; i++ {
			partContent := bytes.Repeat([]byte("A"), partSize)
			partKey := fmt.Sprintf("multipart-parts-test-part%d.txt", i)

			// 创建临时对象作为分块源
			putInput := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    partKey,
				Body:   bytes.NewReader(partContent),
			}

			if _, err := client.TestClient().PutObject(putInput); err != nil {
				t.Fatalf("创建临时分块 %d 失败: %v", i, err)
			}

			// 上传分块
			partInput := &obs.UploadPartInput{
				Bucket:       bucket,
				Key:          objectKey,
				PartNumber:   int32(i + 1),
				UploadId:     initOutput.UploadId,
				SourceFile:   partKey,
			}

			partOutput, err := client.TestClient().UploadPart(partInput)
			if err != nil {
				t.Fatalf("上传分块 %d 失败: %v", i, err)
			}

			parts[i] = obs.Part{
				PartNumber:   partOutput.PartNumber,
				ETag:         partOutput.ETag,
			}

			// 清理临时分块对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    partKey,
			}
			client.TestClient().DeleteObject(deleteInput)
		}

		client.AddTestCase("分块上传成功")
		t.Logf("成功上传 %d 个分块", numParts)
	})

	t.Run("ShouldCompleteMultipartUpload_GivenValidParts", func(t *testing.T) {
		parts := make([]obs.Part, numParts)
		for i := 0; i < numParts; i++ {
			// 创建临时分块对象
			partContent := bytes.Repeat([]byte("B"), partSize)
			partKey := fmt.Sprintf("complete-multipart-test-part%d.txt", i)

			putInput := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    partKey,
				Body:   bytes.NewReader(partContent),
			}

			client.TestClient().PutObject(putInput)

			// 上传分块
			partInput := &obs.UploadPartInput{
				Bucket:     bucket,
				Key:        objectKey,
				PartNumber: int32(i + 1),
				UploadId:   initOutput.UploadId,
				SourceFile: partKey,
			}

			partOutput, err := client.TestClient().UploadPart(partInput)
			if err != nil {
				t.Fatalf("上传分块 %d 失败: %v", i, err)
			}

			parts[i] = obs.Part{
				PartNumber: partOutput.PartNumber,
				ETag:       partOutput.ETag,
			}

			// 清理临时分块对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    partKey,
			}
			client.TestClient().DeleteObject(deleteInput)
		}

		// 按分块号排序
		sort.Slice(parts, func(i, j int) bool {
			return parts[i].PartNumber < parts[j].PartNumber
		})

		// 完成分块上传
		completeInput := &obs.CompleteMultipartUploadInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: initOutput.UploadId,
			Parts:    parts,
		}

		output, err := client.TestClient().CompleteMultipartUpload(completeInput)
		if err != nil {
			t.Fatalf("完成分块上传失败: %v", err)
		}

		// 验证返回结果
		if output == nil {
			t.Error("CompleteMultipartUpload返回nil")
		}

		if output.VersionId == "" {
			t.Log("对象版本ID为空（可能未启用版本控制）")
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

		client.AddTestCase("分块上传完成成功")
		t.Logf("分块上传完成，对象: %s, VersionId: %s", objectKey, output.VersionId)
	})
}

// TestMultipartUpload_ListAndCancelParts 测试分块列表和取消
func TestMultipartUpload_ListAndCancelParts(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("list-cancel-multipart-test.txt")

	// 初始化分块上传
	initInput := &obs.InitiateMultipartUploadInput{
		Bucket: bucket,
		Key:    objectKey,
	}

	initOutput, err := client.TestClient().InitiateMultipartUpload(initInput)
	if err != nil {
		t.Fatalf("初始化分块上传失败: %v", err)
	}

	t.Run("ShouldListParts_GivenValidUploadId", func(t *testing.T) {
		// 上传一些分块
		for i := 0; i < 2; i++ {
			partContent := bytes.Repeat([]byte("C"), 1024*1024)
			partKey := fmt.Sprintf("list-test-part%d.txt", i)

			putInput := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    partKey,
				Body:   bytes.NewReader(partContent),
			}

			client.TestClient().PutObject(putInput)

			partInput := &obs.UploadPartInput{
				Bucket:     bucket,
				Key:        objectKey,
				PartNumber: int32(i + 1),
				UploadId:   initOutput.UploadId,
				SourceFile: partKey,
			}

			client.TestClient().UploadPart(partInput)

			// 清理临时分块对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    partKey,
			}
			client.TestClient().DeleteObject(deleteInput)
		}

		// 列出分块
		listInput := &obs.ListPartsInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: initOutput.UploadId,
		}

		listOutput, err := client.TestClient().ListParts(listInput)
		if err != nil {
			t.Fatalf("列出分块失败: %v", err)
		}

		// 验证分块列表
		if listOutput == nil {
			t.Error("ListParts返回nil")
		}

		if len(listOutput.Parts) == 0 {
			t.Error("分块列表为空")
		}

		client.AddTestCase("分块列表获取成功")
		t.Logf("成功列出 %d 个分块", len(listOutput.Parts))
	})

	t.Run("ShouldCancelMultipartUpload_GivenValidUploadId", func(t *testing.T) {
		// 取消分块上传
		abortInput := &obs.AbortMultipartUploadInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: initOutput.UploadId,
		}

		_, err := client.TestClient().AbortMultipartUpload(abortInput)
		if err != nil {
			t.Fatalf("取消分块上传失败: %v", err)
		}

		client.AddTestCase("分块上传取消成功")
		t.Log("分块上传取消成功")
	})
}

// TestMultipartUpload_LargeFile 测试大文件分块上传
func TestMultipartUpload_LargeFile(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("large-multipart-test.bin")

	// 创建大文件（50MB）
	fileSize := 50 * 1024 * 1024
	partSize := 10 * 1024 * 1024 // 10MB每个分块
	numParts := fileSize / partSize

	t.Run("ShouldInitiateLargeFileUpload_GivenValidInput", func(t *testing.T) {
		// 初始化大文件分块上传
		input := &obs.InitiateMultipartUploadInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().InitiateMultipartUpload(input)
		if err != nil {
			t.Fatalf("初始化大文件分块上传失败: %v", err)
		}

		// 添加取消清理函数
		client.AddCleanup(func(t *testing.T) {
			abortInput := &obs.AbortMultipartUploadInput{
				Bucket:   bucket,
				Key:      objectKey,
				UploadId: output.UploadId,
			}
			_, err := client.TestClient().AbortMultipartUpload(abortInput)
			if err != nil {
				t.Logf("取消大文件分块上传失败: %v", err)
			}
		})

		client.AddTestCase("大文件分块上传初始化成功")
		t.Logf("大文件分块上传初始化，UploadId: %s", output.UploadId)
	})

	t.Run("ShouldUploadLargeFileParts_GivenValidUploadId", func(t *testing.T) {
		// 重新初始化以获取UploadId
		initInput := &obs.InitiateMultipartUploadInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		initOutput, err := client.TestClient().InitiateMultipartUpload(initInput)
		if err != nil {
			t.Fatalf("初始化分块上传失败: %v", err)
		}

		// 添加取消清理函数
		client.AddCleanup(func(t *testing.T) {
			abortInput := &obs.AbortMultipartUploadInput{
				Bucket:   bucket,
				Key:      objectKey,
				UploadId: initOutput.UploadId,
			}
			_, err := client.TestClient().AbortMultipartUpload(abortInput)
			if err != nil {
				t.Logf("取消分块上传失败: %v", err)
			}
		})

		parts := make([]obs.Part, numParts)

		// 上传大文件的所有分块
		for i := 0; i < int(numParts); i++ {
			partContent := bytes.Repeat([]byte("D"), partSize)
			partKey := fmt.Sprintf("large-file-part%d.bin", i)

			// 创建临时分块对象
			putInput := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    partKey,
				Body:   bytes.NewReader(partContent),
			}

			if _, err := client.TestClient().PutObject(putInput); err != nil {
				t.Fatalf("创建临时分块 %d 失败: %v", i, err)
			}

			// 上传分块
			partInput := &obs.UploadPartInput{
				Bucket:     bucket,
				Key:        objectKey,
				PartNumber: int32(i + 1),
				UploadId:   initOutput.UploadId,
				SourceFile: partKey,
			}

			partOutput, err := client.TestClient().UploadPart(partInput)
			if err != nil {
				t.Fatalf("上传大文件分块 %d 失败: %v", i, err)
			}

			parts[i] = obs.Part{
				PartNumber: partOutput.PartNumber,
				ETag:       partOutput.ETag,
			}

			// 清理临时分块对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    partKey,
			}
			client.TestClient().DeleteObject(deleteInput)
		}

		client.AddTestCase("大文件分块上传成功")
		t.Logf("成功上传大文件的 %d 个分块", numParts)
	})

	t.Run("ShouldCompleteLargeFileUpload_GivenValidParts", func(t *testing.T) {
		// 重新初始化以获取UploadId
		initInput := &obs.InitiateMultipartUploadInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		initOutput, err := client.TestClient().InitiateMultipartUpload(initInput)
		if err != nil {
			t.Fatalf("初始化分块上传失败: %v", err)
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

		parts := make([]obs.Part, numParts)

		// 上传分块
		for i := 0; i < int(numParts); i++ {
			partContent := bytes.Repeat([]byte("E"), partSize)
			partKey := fmt.Sprintf("complete-large-part%d.bin", i)

			putInput := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    partKey,
				Body:   bytes.NewReader(partContent),
			}

			client.TestClient().PutObject(putInput)

			partInput := &obs.UploadPartInput{
				Bucket:     bucket,
				Key:        objectKey,
				PartNumber: int32(i + 1),
				UploadId:   initOutput.UploadId,
				SourceFile: partKey,
			}

			partOutput, err := client.TestClient().UploadPart(partInput)
			if err != nil {
				t.Fatalf("上传分块 %d 失败: %v", i, err)
			}

			parts[i] = obs.Part{
				PartNumber: partOutput.PartNumber,
				ETag:       partOutput.ETag,
			}

			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    partKey,
			}
			client.TestClient().DeleteObject(deleteInput)
		}

		// 按分块号排序
		sort.Slice(parts, func(i, j int) bool {
			return parts[i].PartNumber < parts[j].PartNumber
		})

		// 完成分块上传
		completeInput := &obs.CompleteMultipartUploadInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: initOutput.UploadId,
			Parts:    parts,
		}

		output, err := client.TestClient().CompleteMultipartUpload(completeInput)
		if err != nil {
			t.Fatalf("完成大文件分块上传失败: %v", err)
		}

		if output.VersionId == "" {
			t.Log("对象版本ID为空（可能未启用版本控制）")
		}

		client.AddTestCase("大文件分块上传完成成功")
		t.Logf("大文件分块上传完成，对象: %s, 大小: %d bytes", objectKey, fileSize)
	})
}

// TestMultipartUpload_ConcurrentParts 测试并发分块上传
func TestMultipartUpload_ConcurrentParts(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("concurrent-multipart-test.txt")

	// 初始化分块上传
	initInput := &obs.InitiateMultipartUploadInput{
		Bucket: bucket,
		Key:    objectKey,
	}

	initOutput, err := client.TestClient().InitiateMultipartUpload(initInput)
	if err != nil {
		t.Fatalf("初始化分块上传失败: %v", err)
	}

	// 添加取消清理函数
	client.AddCleanup(func(t *testing.T) {
		abortInput := &obs.AbortMultipartUploadInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: initOutput.UploadId,
		}
		_, err := client.TestClient().AbortMultipartUpload(abortInput)
		if err != nil {
			t.Logf("取消分块上传失败: %v", err)
		}
	})

	t.Run("ShouldUploadPartsConcurrently_GivenValidUploadId", func(t *testing.T) {
		numParts := 5
		var wg sync.WaitGroup
		errChan := make(chan error, numParts)
		parts := make([]obs.Part, numParts)

		// 并发上传分块
		for i := 0; i < numParts; i++ {
			wg.Add(1)
			go func(partNum int) {
				defer wg.Done()

				partContent := bytes.Repeat([]byte("F"), 2*1024*1024)
				partKey := fmt.Sprintf("concurrent-part%d.txt", partNum)

				// 创建临时分块对象
				putInput := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    partKey,
					Body:   bytes.NewReader(partContent),
				}

				if _, err := client.TestClient().PutObject(putInput); err != nil {
					errChan <- fmt.Errorf("创建临时分块 %d 失败: %v", partNum, err)
					return
				}

				// 上传分块
				partInput := &obs.UploadPartInput{
					Bucket:     bucket,
					Key:        objectKey,
					PartNumber: int32(partNum + 1),
					UploadId:   initOutput.UploadId,
					SourceFile: partKey,
				}

				partOutput, err := client.TestClient().UploadPart(partInput)
				if err != nil {
					errChan <- fmt.Errorf("上传分块 %d 失败: %v", partNum, err)
					return
				}

				parts[partNum] = obs.Part{
					PartNumber: partOutput.PartNumber,
					ETag:       partOutput.ETag,
				}

				// 清理临时分块对象
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    partKey,
				}
				client.TestClient().DeleteObject(deleteInput)

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
			t.Errorf("有 %d/%d 个并发分块上传失败", errorCount, numParts)
		}

		client.AddTestCase("并发分块上传成功")
		t.Logf("并发分块上传完成，成功率: %d/%d", numParts-errorCount, numParts)
	})
}

// TestMultipartUpload_PartSizeValidation 测试分块大小验证
func TestMultipartUpload_PartSizeValidation(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("part-size-validation-test.txt")

	// 测试不同分块大小
	testSizes := []struct {
		name     string
		size     int64
		expected bool
	}{
		{"小分块 (100KB)", 100 * 1024, true},
		{"标准分块 (5MB)", 5 * 1024 * 1024, true},
		{"大分块 (100MB)", 100 * 1024 * 1024, true},
	}

	for _, test := range testSizes {
		t.Run(fmt.Sprintf("ShouldValidatePartSize_Given%s", test.name), func(t *testing.T) {
			// 初始化分块上传
			initInput := &obs.InitiateMultipartUploadInput{
				Bucket: bucket,
				Key:    objectKey + "-" + strings.Replace(test.name, " ", "_", -1),
			}

			initOutput, err := client.TestClient().InitiateMultipartUpload(initInput)
			if err != nil {
				t.Fatalf("初始化分块上传失败: %v", err)
			}

			// 添加取消清理函数
			client.AddCleanup(func(t *testing.T) {
				abortInput := &obs.AbortMultipartUploadInput{
					Bucket:   bucket,
					Key:      objectKey + "-" + strings.Replace(test.name, " ", "_", -1),
					UploadId: initOutput.UploadId,
				}
				_, err := client.TestClient().AbortMultipartUpload(abortInput)
				if err != nil {
					t.Logf("取消分块上传失败: %v", err)
				}
			})

			// 创建临时分块对象
			partContent := bytes.Repeat([]byte("G"), int(test.size))
			partKey := fmt.Sprintf("size-validation-part-%d.txt", test.size)

			putInput := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    partKey,
				Body:   bytes.NewReader(partContent),
			}

			client.TestClient().PutObject(putInput)

			// 上传分块
			partInput := &obs.UploadPartInput{
				Bucket:     bucket,
				Key:        objectKey + "-" + strings.Replace(test.name, " ", "_", -1),
				PartNumber: 1,
				UploadId:   initOutput.UploadId,
				SourceFile: partKey,
			}

			_, err = client.TestClient().UploadPart(partInput)
			if err != nil {
				if test.expected {
					t.Errorf("分块大小 %s 测试失败: %v", test.name, err)
				} else {
					t.Logf("分块大小 %s 失败（预期的）: %v", test.name, err)
				}
			}

			// 清理临时分块对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    partKey,
			}
			client.TestClient().DeleteObject(deleteInput)
		})
	}
}

// TestMultipartUpload_ErrorHandling 测试错误处理
func TestMultipartUpload_ErrorHandling(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldFail_GivenInvalidUploadId", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("invalid-uploadid-test.txt")

		// 尝试使用无效的UploadId列出分块
		listInput := &obs.ListPartsInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: "invalid-upload-id-12345",
		}

		_, err := client.TestClient().ListParts(listInput)
		if err == nil {
			t.Error("期望列表分块失败，但操作成功")
		}

		// 验证错误类型
		obsErr, ok := err.(obs.ObsError)
		if !ok {
			t.Fatalf("错误不是ObsError类型: %T", err)
		}

		client.AddTestCase("无效UploadId错误处理测试通过")
		t.Logf("无效UploadId错误: 状态码=%d, 代码=%s", obsErr.StatusCode, obsErr.Code)
	})

	t.Run("ShouldFail_GivenNonexistentObject", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("nonexistent-multipart-test.txt")

		// 尝试列出不存在对象的分块
		listInput := &obs.ListPartsInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		_, err := client.TestClient().ListParts(listInput)
		if err == nil {
			t.Error("期望列表分块失败，但操作成功")
		}

		// 验证错误类型
		obsErr, ok := err.(obs.ObsError)
		if !ok {
			t.Fatalf("错误不是ObsError类型: %T", err)
		}

		client.AddTestCase("不存在对象错误处理测试通过")
		t.Logf("不存在对象错误: 状态码=%d, 代码=%s", obsErr.StatusCode, obsErr.Code)
	})
}
