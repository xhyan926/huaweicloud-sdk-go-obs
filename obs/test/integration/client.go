//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/config"
)

// CleanupFunction 清理函数类型
type CleanupFunction func(t *testing.T)

// TestClient 集成测试客户端
type TestClient struct {
	ObsClient    *obs.ObsClient
	Config       *config.TestConfig
	TestPrefix   string
	CleanupFuncs []CleanupFunction
	TestCases    []string
}

// NewTestClient 创建集成测试客户端
func NewTestClient(t *testing.T) *TestClient {
	envConfig := config.NewIntegrationEnvConfig(t)
	cfg := envConfig.Config

	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
		obs.WithRegion(cfg.Region),
	)

	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	return &TestClient{
		ObsClient:    obsClient,
		Config:       cfg,
		TestPrefix:   cfg.TestPrefix,
		CleanupFuncs: make([]CleanupFunction, 0),
		TestCases:    make([]string, 0),
	}
}

// TestClient 获取测试客户端
func (c *TestClient) TestClient() *obs.ObsClient {
	return c.ObsClient
}

// AddCleanup 添加清理函数
func (c *TestClient) AddCleanup(f CleanupFunction) {
	c.CleanupFuncs = append(c.CleanupFuncs, f)
}

// AddTestCase 添加测试用例记录
func (c *TestClient) AddTestCase(testCase string) {
	c.TestCases = append(c.TestCases, testCase)
}

// GetTestBucket 获取测试桶名称
func (c *TestClient) GetTestBucket() string {
	return c.Config.GetTestBucket()
}

// GetTestObjectKey 获取测试对象键
func (c *TestClient) GetTestObjectKey(key string) string {
	return c.Config.GetTestObjectKey(key)
}

// CleanTestObjectKey 清理测试对象键
func (c *TestClient) CleanTestObjectKey(key string) string {
	return c.Config.CleanTestObjectKey(key)
}

// Cleanup 清理测试资源
func (c *TestClient) Cleanup(t *testing.T) {
	for i := len(c.CleanupFuncs) - 1; i >= 0; i-- {
		cleanupFunc := c.CleanupFuncs[i]
		cleanupFunc(t)
	}
	c.CleanupFuncs = make([]CleanupFunction, 0)
}

// GenerateUniqueKey 生成唯一的测试键
func (c *TestClient) GenerateUniqueKey(suffix string) string {
	return fmt.Sprintf("%s%s-%d-%s",
		c.TestPrefix,
		"test",
		time.Now().UnixNano(),
		suffix)
}

// WithContext 带超时的context
func (c *TestClient) WithContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// WaitForObjectExists 等待对象存在
func (c *TestClient) WaitForObjectExists(bucket, key string, maxAttempts int) error {
	input := &obs.GetObjectMetadataInput{
		Bucket: bucket,
		Key:    key,
	}

	for i := 0; i < maxAttempts; i++ {
		_, err := c.ObsClient.GetObjectMetadata(input)
		if err == nil {
			return nil
		}

		// 如果是404错误，继续等待
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			time.Sleep(1 * time.Second)
			continue
		}

		// 其他错误直接返回
		return err
	}

	return fmt.Errorf("object %s/%s does not exist after %d attempts", bucket, key, maxAttempts)
}

// WaitForObjectNotExists 等待对象不存在
func (c *TestClient) WaitForObjectNotExists(bucket, key string, maxAttempts int) error {
	input := &obs.GetObjectMetadataInput{
		Bucket: bucket,
		Key:    key,
	}

	for i := 0; i < maxAttempts; i++ {
		_, err := c.ObsClient.GetObjectMetadata(input)
		if err == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// 如果是404错误，对象不存在，返回成功
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			return nil
		}

		// 其他错误继续等待
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("object %s/%s still exists after %d attempts", bucket, key, maxAttempts)
}

// RegisterCleanupFunc 注册一个简单的清理函数
func (c *TestClient) RegisterCleanupFunc(action string, cleanupFunc func() error) {
	c.CleanupFuncs = append(c.CleanupFuncs, func(t *testing.T) {
		t.Logf("Cleaning up: %s", action)
		err := cleanupFunc()
		if err != nil {
			t.Logf("Cleanup error for %s: %v", action, err)
		}
	})
}

// PrintTestCases 打印测试用例记录
func (c *TestClient) PrintTestCases() {
	if len(c.TestCases) > 0 {
		fmt.Println("\n=== Test Cases Executed ===")
		for i, testCase := range c.TestCases {
			fmt.Printf("%d. %s\n", i+1, testCase)
		}
		fmt.Println("============================")
	}
}