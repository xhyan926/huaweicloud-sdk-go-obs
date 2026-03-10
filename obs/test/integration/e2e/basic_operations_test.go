//go:build integration

package e2e

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// TestBasicOperations_ShouldSucceed_GivenValidCredentials 测试基本操作能否成功执行，给定有效的凭证
func TestBasicOperations_ShouldSucceed_GivenValidCredentials(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	t.Run("认证测试", func(t *testing.T) {
		testAuthentication_ShouldConnectSuccessfully_GivenValidCredentials(t, client)
	})

	t.Run("桶操作测试", func(t *testing.T) {
		testBucketOperations_ShouldCreateAndListSuccessfully_GivenValidBucket(t, client)
	})

	t.Run("对象操作测试", func(t *testing.T) {
		testObjectOperations_ShouldUploadAndDownloadSuccessfully_GivenValidObject(t, client)
	})
}

// testAuthentication_ShouldConnectSuccessfully_GivenValidCredentials 测试认证功能
func testAuthentication_ShouldConnectSuccessfully_GivenValidCredentials(t *testing.T, client *integration.TestClient) {
	t.Run("ShouldConnectSuccessfully_GivenValidCredentials", func(t *testing.T) {
		// 验证连接是否正常
		input := &obs.GetBucketLocationInput{
			Bucket: client.GetTestBucket(),
		}

		_, err := client.TestClient().GetBucketLocation(input)
		if err != nil {
			t.Fatalf("认证失败: %v", err)
		}

		client.AddTestCase("认证成功")
		t.Log("认证测试通过")
	})
}

// testBucketOperations_ShouldCreateAndListSuccessfully_GivenValidBucket 测试桶操作
func testBucketOperations_ShouldCreateAndListSuccessfully_GivenValidBucket(t *testing.T, client *integration.TestClient) {
	bucket := client.GetTestBucket()

	t.Run("ShouldCreateBucketSuccessfully_GivenValidBucketName", func(t *testing.T) {
		input := &obs.CreateBucketInput{
			Bucket: bucket,
		}

		_, err := client.TestClient().CreateBucket(input)
		if err != nil {
			t.Fatalf("创建桶失败: %v", err)
		}

		client.AddCleanup(func(t *testing.T) {
			deleteInput := &obs.DeleteBucketInput{
				Bucket: bucket,
			}
			_, err := client.TestClient().DeleteBucket(deleteInput)
			if err != nil {
				t.Logf("删除桶失败: %v", err)
			}
		})

		client.AddTestCase("创建桶成功")
		t.Logf("成功创建桶: %s", bucket)
	})

	t.Run("ShouldListBucketsSuccessfully_GivenValidCredentials", func(t *testing.T) {
		output, err := client.TestClient().ListBuckets(&obs.ListBucketsInput{})
		if err != nil {
			t.Fatalf("列出桶失败: %v", err)
		}

		// 验证测试桶在列表中
		found := false
		for _, b := range output.Buckets {
			if b.Name == bucket {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("测试桶 %s 不在桶列表中", bucket)
		}

		client.AddTestCase("列出桶成功")
		t.Logf("成功列出桶，共 %d 个桶", len(output.Buckets))
	})

	t.Run("ShouldGetBucketInfoSuccessfully_GivenValidBucket", func(t *testing.T) {
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		output, err := client.TestClient().GetBucketLocation(input)
		if err != nil {
			t.Fatalf("获取桶信息失败: %v", err)
		}

		if output.Location == "" {
			t.Error("桶区域信息为空")
		}

		client.AddTestCase("获取桶信息成功")
		t.Logf("桶 %s 的区域: %s", bucket, output.Location)
	})
}

// testObjectOperations_ShouldUploadAndDownloadSuccessfully_GivenValidObject 测试对象操作
func testObjectOperations_ShouldUploadAndDownloadSuccessfully_GivenValidObject(t *testing.T, client *integration.TestClient) {
	bucket := client.GetTestBucket()

	// 测试小文件上传下载
	t.Run("ShouldUploadAndDownloadSuccessfully_GivenSmallFile", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("small-file.txt")
		content := "This is a small test file for basic operations."

		// 上传对象
		err := client.TestClient().PutObject(bucket, objectKey, content, nil)
		if err != nil {
			t.Fatalf("上传对象失败: %v", err)
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

		// 下载对象
		getInput := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(getInput)
		if err != nil {
			t.Fatalf("下载对象失败: %v", err)
		}
		defer output.Body.Close()

		// 验证内容
		body := make([]byte, len(content))
		_, err = output.Body.Read(body)
		if err != nil {
			t.Fatalf("读取对象内容失败: %v", err)
		}

		if string(body) != content {
			t.Errorf("对象内容不匹配，期望: %s, 实际: %s", content, string(body))
		}

		client.AddTestCase("小文件上传下载成功")
		t.Logf("成功上传并下载对象: %s", objectKey)
	})

	// 测试获取对象元数据
	t.Run("ShouldGetObjectMetadataSuccessfully_GivenValidObject", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("metadata-test.txt")
		content := "Test content for metadata."

		// 先上传对象
		err := client.TestClient().PutObject(bucket, objectKey, content, nil)
		if err != nil {
			t.Fatalf("上传对象失败: %v", err)
		}

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

		// 获取对象元数据
		input := &obs.GetObjectMetadataInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObjectMetadata(input)
		if err != nil {
			t.Fatalf("获取对象元数据失败: %v", err)
		}

		if output.ContentLength != int64(len(content)) {
			t.Errorf("对象长度不匹配，期望: %d, 实际: %d", len(content), output.ContentLength)
		}

		client.AddTestCase("获取对象元数据成功")
		t.Logf("对象元数据: Content-Type=%s, Content-Length=%d", output.ContentType, output.ContentLength)
	})

	// 测试列出对象
	t.Run("ShouldListObjectsSuccessfully_GivenValidBucket", func(t *testing.T) {
		// 创建多个测试对象
		objects := []string{"list-test-1.txt", "list-test-2.txt", "list-test-3.txt"}

		for _, objName := range objects {
			objectKey := client.GetTestObjectKey(objName)
			content := fmt.Sprintf("Test content for %s", objName)

			err := client.TestClient().PutObject(bucket, objectKey, content, nil)
			if err != nil {
				t.Fatalf("上传对象 %s 失败: %v", objName, err)
			}

			client.AddCleanup(func(t *testing.T) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    objectKey,
				}
				_, err := client.TestClient().DeleteObject(deleteInput)
				if err != nil {
					t.Logf("删除对象 %s 失败: %v", objName, err)
				}
			})
		}

		// 列出对象
		input := &obs.ListObjectsInput{
			Bucket: bucket,
		}

		output, err := client.TestClient().ListObjects(input)
		if err != nil {
			t.Fatalf("列出对象失败: %v", err)
		}

		if len(output.Contents) < 3 {
			t.Errorf("期望至少 3 个对象，实际找到 %d 个", len(output.Contents))
		}

		client.AddTestCase("列出对象成功")
		t.Logf("成功列出 %d 个对象", len(output.Contents))
	})
}

// TestBasicOperations_Concurrent 测试并发操作
func TestBasicOperations_Concurrent(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldHandleConcurrentOperations_GivenMultipleRequests", func(t *testing.T) {
		const concurrency = 10
		var wg sync.WaitGroup
		errChan := make(chan error, concurrency)

		// 并发上传多个对象
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				objectKey := client.GetTestObjectKey(fmt.Sprintf("concurrent-%d.txt", id))
				content := fmt.Sprintf("Concurrent test content %d", id)

				err := client.TestClient().PutObject(bucket, objectKey, content, nil)
				if err != nil {
					errChan <- fmt.Errorf("并发 %d 上传失败: %v", id, err)
					return
				}

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
			}(i)
		}

		wg.Wait()
		close(errChan)

		// 收集错误
		for err := range errChan {
			t.Error(err)
		}

		client.AddTestCase("并发操作成功")
		t.Logf("成功完成 %d 个并发操作", concurrency)
	})
}

// TestBasicOperations_WithTimeout 测试带超时的操作
func TestBasicOperations_WithTimeout(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldCompleteWithinTimeout_GivenValidOperation", func(t *testing.T) {
		objectKey := client.GetTestObjectKey("timeout-test.txt")
		content := "Test content with timeout."

		// 创建带超时的context
		ctx, cancel := client.WithContext(30 * time.Second)
		defer cancel()

		// 使用context执行操作
		done := make(chan error, 1)
		go func() {
			err := client.TestClient().PutObject(bucket, objectKey, content, nil)
			done <- err
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("上传对象失败: %v", err)
			}

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

			client.AddTestCase("超时测试成功")
			t.Log("操作在超时时间内完成")

		case <-ctx.Done():
			t.Error("操作超时")
		}
	})
}
