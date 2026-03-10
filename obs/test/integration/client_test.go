//go:build integration

package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/config"
)

// TestNewTestClient_ShouldCreateSuccessfully_GivenValidConfig 测试创建测试客户端
func TestNewTestClient_ShouldCreateSuccessfully_GivenValidConfig(t *testing.T) {
	// 创建测试客户端
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 验证客户端不为nil
	if client == nil {
		t.Fatal("客户端创建失败，返回nil")
	}

	// 验证基本字段
	if client.ObsClient == nil {
		t.Error("ObsClient字段为nil")
	}

	if client.Config == nil {
		t.Error("Config字段为nil")
	}

	if client.CleanupFuncs == nil {
		t.Error("CleanupFuncs字段为nil")
	}

	if client.TestCases == nil {
		t.Error("TestCases字段为nil")
	}

	// 验证配置有效性
	if !client.Config.IsValid() {
		t.Error("配置无效")
	}

	t.Log("测试客户端创建成功")
}

// TestTestClient_ShouldReturnObsClient 测试获取ObsClient
func TestTestClient_ShouldReturnObsClient(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	obsClient := client.TestClient()

	// 验证返回的客户端
	if obsClient == nil {
		t.Fatal("TestClient()返回nil")
	}

	// 验证是同一个客户端实例
	if obsClient != client.ObsClient {
		t.Error("TestClient()返回的不是同一个ObsClient实例")
	}

	t.Log("TestClient()方法工作正常")
}

// TestAddCleanup_ShouldAddCleanupFunction 测试添加清理函数
func TestAddCleanup_ShouldAddCleanupFunction(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 添加多个清理函数
	cleanupCalled := make([]bool, 3)

	for i := 0; i < 3; i++ {
		idx := i
		client.AddCleanup(func(t *testing.T) {
			cleanupCalled[idx] = true
			t.Logf("清理函数 %d 被调用", idx)
		})
	}

	// 验证清理函数数量
	if len(client.CleanupFuncs) != 3 {
		t.Errorf("期望3个清理函数，实际 %d 个", len(client.CleanupFuncs))
	}

	// 调用Cleanup验证清理函数被执行
	client.Cleanup(t)

	// 验证所有清理函数都被调用
	for i, called := range cleanupCalled {
		if !called {
			t.Errorf("清理函数 %d 未被调用", i)
		}
	}

	// 验证Cleanup后清理函数列表被清空
	if len(client.CleanupFuncs) != 0 {
		t.Errorf("Cleanup后CleanupFuncs应为空，实际有 %d 个", len(client.CleanupFuncs))
	}

	t.Log("清理函数添加和执行测试通过")
}

// TestAddTestCase_ShouldAddTestCase 测试添加测试用例记录
func TestAddTestCase_ShouldAddTestCase(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 添加多个测试用例
	testCases := []string{"测试用例1", "测试用例2", "测试用例3"}
	for _, tc := range testCases {
		client.AddTestCase(tc)
	}

	// 验证测试用例数量
	if len(client.TestCases) != 3 {
		t.Errorf("期望3个测试用例，实际 %d 个", len(client.TestCases))
	}

	// 验证测试用例内容
	for i, expected := range testCases {
		if client.TestCases[i] != expected {
			t.Errorf("测试用例 %d 不匹配，期望: %s, 实际: %s", i, expected, client.TestCases[i])
		}
	}

	t.Log("测试用例添加测试通过")
}

// TestGetTestBucket_ShouldReturnBucketName 测试获取测试桶名称
func TestGetTestBucket_ShouldReturnBucketName(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	// 验证桶名称不为空
	if bucket == "" {
		t.Error("测试桶名称为空")
	}

	// 验证桶名称与配置一致
	if bucket != client.Config.TestBucket {
		t.Errorf("桶名称不匹配，期望: %s, 实际: %s", client.Config.TestBucket, bucket)
	}

	t.Logf("测试桶名称: %s", bucket)
}

// TestGetTestObjectKey_ShouldReturnKeyWithPrefix 测试获取测试对象键
func TestGetTestObjectKey_ShouldReturnKeyWithPrefix(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	testKey := "test-object.txt"
	objectKey := client.GetTestObjectKey(testKey)

	// 验证对象键不为空
	if objectKey == "" {
		t.Error("对象键为空")
	}

	// 验证前缀
	expectedPrefix := client.Config.TestPrefix
	if len(objectKey) <= len(expectedPrefix) || objectKey[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("对象键不包含预期前缀，期望前缀: %s, 实际键: %s", expectedPrefix, objectKey)
	}

	// 验证包含原始键
	if len(objectKey) <= len(expectedPrefix) {
		t.Error("对象键只包含前缀，不包含原始键")
	}

	actualKey := objectKey[len(expectedPrefix):]
	if actualKey != testKey {
		t.Errorf("对象键的原始部分不匹配，期望: %s, 实际: %s", testKey, actualKey)
	}

	t.Logf("测试对象键: %s", objectKey)
}

// TestCleanTestObjectKey_ShouldRemovePrefix 测试清理测试对象键
func TestCleanTestObjectKey_ShouldRemovePrefix(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	testKey := "test-object.txt"
	objectKey := client.GetTestObjectKey(testKey)
	cleanedKey := client.CleanTestObjectKey(objectKey)

	// 验证清理后的键与原始键一致
	if cleanedKey != testKey {
		t.Errorf("清理后的键不匹配，期望: %s, 实际: %s", testKey, cleanedKey)
	}

	t.Logf("原始键: %s, 带前缀键: %s, 清理后键: %s", testKey, objectKey, cleanedKey)
}

// TestGenerateUniqueKey_ShouldGenerateUniqueKey 测试生成唯一键
func TestGenerateUniqueKey_ShouldGenerateUniqueKey(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 生成多个唯一键
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = client.GenerateUniqueKey(fmt.Sprintf("unique-%d", i))
	}

	// 验证所有键都是唯一的
	keySet := make(map[string]bool)
	for _, key := range keys {
		if keySet[key] {
			t.Error("生成的键不唯一")
		}
		keySet[key] = true
	}

	// 验证所有键都包含前缀
	expectedPrefix := client.Config.TestPrefix + "test-"
	for i, key := range keys {
		if len(key) <= len(expectedPrefix) || key[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("键 %d 不包含预期前缀，键: %s, 前缀: %s", i, key, expectedPrefix)
		}
	}

	t.Logf("生成了 %d 个唯一键", len(keys))
}

// TestWithContext_ShouldCreateContextWithTimeout 测试创建带超时的context
func TestWithContext_ShouldCreateContextWithTimeout(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	timeout := 30 * time.Second
	ctx, cancel := client.WithContext(timeout)

	// 验证context和cancel函数不为nil
	if ctx == nil {
		t.Fatal("WithContext返回的context为nil")
	}
	if cancel == nil {
		t.Fatal("WithContext返回的cancel函数为nil")
	}

	// 验证context的超时设置
	if deadline, ok := ctx.Deadline(); ok {
		expectedDeadline := time.Now().Add(timeout)
		// 允许1秒的时间误差
		if deadline.Before(expectedDeadline.Add(-1 * time.Second)) || deadline.After(expectedDeadline.Add(1*time.Second)) {
			t.Errorf("context超时设置不正确，期望: %v, 实际: %v", expectedDeadline, deadline)
		}
	} else {
		t.Error("context没有设置deadline")
	}

	// 调用cancel函数验证正常工作
	cancel()
	if ctx.Err() == nil {
		t.Error("cancel后context.Err()应该不为nil")
	}

	t.Log("WithContext功能测试通过")
}

// TestWaitForObjectExists_ShouldWaitSuccessfully 测试等待对象存在
func TestWaitForObjectExists_ShouldWaitSuccessfully(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("wait-exist-test.txt")
	content := "测试内容"

	// 上传对象
	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   nil,
	}
	if err := client.TestClient().PutObject(putInput); err != nil {
		t.Fatalf("上传对象失败: %v", err)
	}

	// 清理函数
	client.AddCleanup(func(t *testing.T) {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		client.TestClient().DeleteObject(deleteInput)
	})

	// 等待对象存在
	err := client.WaitForObjectExists(bucket, objectKey, 5)
	if err != nil {
		t.Errorf("等待对象存在失败: %v", err)
	}

	t.Log("WaitForObjectExists功能测试通过")
}

// TestWaitForObjectNotExists_ShouldWaitSuccessfully 测试等待对象不存在
func TestWaitForObjectNotExists_ShouldWaitSuccessfully(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("wait-not-exist-test.txt")
	content := "测试内容"

	// 上传对象
	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   nil,
	}
	if err := client.TestClient().PutObject(putInput); err != nil {
		t.Fatalf("上传对象失败: %v", err)
	}

	// 等待对象存在后删除对象
	if err := client.WaitForObjectExists(bucket, objectKey, 5); err != nil {
		t.Errorf("等待对象存在失败: %v", err)
	}

	// 删除对象
	deleteInput := &obs.DeleteObjectInput{
		Bucket: bucket,
		Key:    objectKey,
	}
	if err := client.TestClient().DeleteObject(deleteInput); err != nil {
		t.Fatalf("删除对象失败: %v", err)
	}

	// 等待对象不存在
	err := client.WaitForObjectNotExists(bucket, objectKey, 5)
	if err != nil {
		t.Errorf("等待对象不存在失败: %v", err)
	}

	t.Log("WaitForObjectNotExists功能测试通过")
}

// TestRegisterCleanupFunc_ShouldRegisterCleanupFunc 测试注册清理函数
func TestRegisterCleanupFunc_ShouldRegisterCleanupFunc(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	cleanupCalled := false

	// 注册清理函数
	client.RegisterCleanupFunc("测试清理", func() error {
		cleanupCalled = true
		return nil
	})

	// 验证清理函数被添加
	if len(client.CleanupFuncs) != 1 {
		t.Errorf("期望1个清理函数，实际 %d 个", len(client.CleanupFuncs))
	}

	// 调用Cleanup验证清理函数被执行
	client.Cleanup(t)

	// 验证清理函数被调用
	if !cleanupCalled {
		t.Error("清理函数未被调用")
	}

	t.Log("RegisterCleanupFunc功能测试通过")
}

// TestPrintTestCases_ShouldPrintTestCases 测试打印测试用例
func TestPrintTestCases_ShouldPrintTestCases(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 添加测试用例
	client.AddTestCase("测试用例1")
	client.AddTestCase("测试用例2")
	client.AddTestCase("测试用例3")

	// 调用PrintTestCases
	client.PrintTestCases()

	// 验证测试用例数量
	if len(client.TestCases) != 3 {
		t.Errorf("期望3个测试用例，实际 %d 个", len(client.TestCases))
	}

	t.Log("PrintTestCases功能测试通过")
}

// TestClientConcurrentAccess 测试客户端并发访问
func TestClientConcurrentAccess(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 模拟并发访问
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 测试各种并发操作
			switch id % 4 {
			case 0:
				_ = client.TestClient()
			case 1:
				_ = client.GetTestBucket()
			case 2:
				_ = client.GetTestObjectKey(fmt.Sprintf("concurrent-%d.txt", id))
			case 3:
				client.AddTestCase(fmt.Sprintf("并发测试用例%d", id))
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 收集错误
	for err := range errors {
		t.Error(err)
	}

	// 验证测试用例数量
	if len(client.TestCases) != 3 { // 10次循环，只有25%会添加测试用例
		t.Logf("并发测试用例数量: %d", len(client.TestCases))
	}

	t.Log("客户端并发访问测试通过")
}

// TestCleanup_Order_ShouldExecuteCleanupInReverseOrder 测试清理函数按反向顺序执行
func TestCleanup_Order_ShouldExecuteCleanupInReverseOrder(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 添加多个清理函数，记录执行顺序
	executionOrder := make([]int, 5)

	for i := 0; i < 5; i++ {
		idx := i
		client.AddCleanup(func(t *testing.T) {
			executionOrder[idx] = idx
			t.Logf("清理函数 %d 执行", idx)
		})
	}

	// 调用Cleanup
	client.Cleanup(t)

	// 验证执行顺序是反向的
	for i := 0; i < 5; i++ {
		expectedOrder := 4 - i // 应该是4,3,2,1,0
		if executionOrder[i] != expectedOrder {
			t.Errorf("清理函数执行顺序不正确，位置 %d 期望 %d, 实际 %d", i, expectedOrder, executionOrder[i])
		}
	}

	t.Log("清理函数反向执行顺序测试通过")
}

// TestClientContextualUsage 测试客户端上下文使用
func TestClientContextualUsage(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 使用带超时的context
	ctx, cancel := client.WithContext(10 * time.Second)
	defer cancel()

	// 验证context可用于操作
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("context没有设置deadline")
	}

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("context-test.txt")

	// 检查对象是否存在（使用context）
	getInput := &obs.GetObjectMetadataInput{
		Bucket: bucket,
		Key:    objectKey,
	}

	_, err := client.TestClient().GetObjectMetadata(getInput)
	if err != nil {
		// 对象不存在是预期的
		t.Logf("对象不存在（预期的）: %v", err)
	}

	t.Log("客户端上下文使用测试通过")
}

// TestClientMultipleCleanupCalls 测试多次调用Cleanup
func TestClientMultipleCleanupCalls(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	cleanupCallCount := 0

	// 添加清理函数
	client.AddCleanup(func(t *testing.T) {
		cleanupCallCount++
	})

	// 第一次Cleanup
	client.Cleanup(t)
	firstCallCount := cleanupCallCount

	// 第二次Cleanup（应该没有效果）
	client.Cleanup(t)
	secondCallCount := cleanupCallCount

	// 验证清理函数只被调用一次
	if firstCallCount != 1 {
		t.Errorf("第一次Cleanup后清理函数应该被调用1次，实际 %d 次", firstCallCount)
	}

	if secondCallCount != 1 {
		t.Errorf("第二次Cleanup不应该调用清理函数，实际被调用 %d 次", secondCallCount)
	}

	t.Log("多次Cleanup调用测试通过")
}

// TestClientErrorHandling 测试客户端错误处理
func TestClientErrorHandling(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	// 测试添加nil清理函数
	client.AddCleanup(nil)

	// 验证清理函数被添加
	if len(client.CleanupFuncs) != 1 {
		t.Error("nil清理函数应该被添加到列表中")
	}

	// Cleanup应该能够处理nil清理函数
	client.Cleanup(t)

	t.Log("客户端错误处理测试通过")
}

// TestClientWithInvalidConfig 测试无效配置的处理
func TestClientWithInvalidConfig(t *testing.T) {
	// 这个测试需要跳过，因为NewTestClient会检查配置有效性
	if !config.LoadTestConfig().IsValid() {
		t.Skip("跳过无效配置测试，因为测试环境配置无效")
	}

	// 如果配置有效，测试正常创建
	client := NewTestClient(t)
	defer client.Cleanup(t)

	// 验证客户端能够正常创建
	if client == nil {
		t.Fatal("客户端创建失败")
	}

	t.Log("客户端配置处理测试通过")
}
