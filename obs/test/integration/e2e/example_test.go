//go:build integration

package e2e

import (
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// ExampleClientUsage 展示如何使用集成测试客户端
func ExampleClientUsage(t *testing.T) {
	// 创建集成测试客户端
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	// 添加清理函数
	client.AddCleanup(func(t *testing.T) {
		t.Log("执行自定义清理操作")
	})

	// 记录测试用例
	client.AddTestCase("创建和获取对象")
	client.AddTestCase("删除对象")

	// 获取测试桶
	bucket := client.GetTestBucket()
	t.Logf("使用测试桶: %s", bucket)

	// 创建测试对象
	objectKey := client.GetTestObjectKey("test-object.txt")
	content := "This is a test object content for E2E testing."

	// 上传对象
	err := client.TestClient().PutObject(bucket, objectKey, content, nil)
	if err != nil {
		t.Fatalf("上传对象失败: %v", err)
	}

	// 添加对象清理
	client.AddCleanup(func(t *testing.T) {
		err := client.TestClient().DeleteObject(bucket, objectKey)
		if err != nil {
			t.Logf("删除对象失败: %v", err)
		}
	})

	t.Logf("成功创建对象: %s", objectKey)
}

// TestBasicBucketOperations 基本桶操作测试
func TestBasicBucketOperations(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	// 测试创建桶
	t.Run("创建桶", func(t *testing.T) {
		input := &obs.CreateBucketInput{
			Bucket: bucket,
		}

		_, err := client.TestClient().CreateBucket(input)
		if err != nil {
			t.Fatalf("创建桶失败: %v", err)
		}

		client.AddTestCase("创建桶成功")
	})

	// 测试获取桶信息
	t.Run("获取桶信息", func(t *testing.T) {
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		_, err := client.TestClient().GetBucketLocation(input)
		if err != nil {
			t.Fatalf("获取桶信息失败: %v", err)
		}

		client.AddTestCase("获取桶信息成功")
	})

	// 测试删除桶
	t.Run("删除桶", func(t *testing.T) {
		input := &obs.DeleteBucketInput{
			Bucket: bucket,
		}

		_, err := client.TestClient().DeleteBucket(input)
		if err != nil {
			t.Fatalf("删除桶失败: %v", err)
		}

		client.AddTestCase("删除桶成功")
	})
}

// TestBasicObjectOperations 基本对象操作测试
func TestBasicObjectOperations(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("test-file.txt")
	content := "This is a test file content for object operations."

	// 测试上传对象
	t.Run("上传对象", func(t *testing.T) {
		err := client.TestClient().PutObject(bucket, objectKey, content, nil)
		if err != nil {
			t.Fatalf("上传对象失败: %v", err)
		}

		client.AddTestCase("上传对象成功")
	})

	// 测试获取对象
	t.Run("获取对象", func(t *testing.T) {
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		output, err := client.TestClient().GetObject(input)
		if err != nil {
			t.Fatalf("获取对象失败: %v", err)
		}
		defer output.Body.Close()

		body := make([]byte, len(content))
		_, err = output.Body.Read(body)
		if err != nil {
			t.Fatalf("读取对象内容失败: %v", err)
		}

		if string(body) != content {
			t.Fatalf("对象内容不匹配，期望: %s, 实际: %s", content, string(body))
		}

		client.AddTestCase("获取对象成功")
	})

	// 测试获取对象元数据
	t.Run("获取对象元数据", func(t *testing.T) {
		input := &obs.GetObjectMetadataInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		_, err := client.TestClient().GetObjectMetadata(input)
		if err != nil {
			t.Fatalf("获取对象元数据失败: %v", err)
		}

		client.AddTestCase("获取对象元数据成功")
	})

	// 测试删除对象
	t.Run("删除对象", func(t *testing.T) {
		input := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}

		_, err := client.TestClient().DeleteObject(input)
		if err != nil {
			t.Fatalf("删除对象失败: %v", err)
		}

		client.AddTestCase("删除对象成功")
	})

	// 打印测试用例记录
	client.PrintTestCases()
}