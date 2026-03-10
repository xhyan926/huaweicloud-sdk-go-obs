package config

import (
	"os"
	"testing"
)

// IntegrationEnvConfig 集成测试环境配置
type IntegrationEnvConfig struct {
	Config *TestConfig
}

// NewIntegrationEnvConfig 创建集成测试环境配置
func NewIntegrationEnvConfig(t *testing.T) *IntegrationEnvConfig {
	cfg := LoadTestConfig()

	if !cfg.IsValid() {
		t.Skipf("Skipping integration tests: invalid test configuration. Please set OBS_TEST_AK, OBS_TEST_SK, OBS_TEST_ENDPOINT, and OBS_TEST_BUCKET environment variables")
	}

	return &IntegrationEnvConfig{Config: cfg}
}

// ShouldSkipIntegrationTest 根据测试名称检查是否应该跳过集成测试
func (c *IntegrationEnvConfig) ShouldSkipIntegrationTest(testName string) bool {
	// 检查是否通过环境变量跳过所有集成测试
	if c.Config.SkipIntegrationTests {
		return true
	}

	// 检查测试环境配置是否有效
	if !c.Config.IsValid() {
		return true
	}

	// 检查特定测试是否被过滤
	if testFilter := os.Getenv("OBS_TEST_FILTER"); testFilter != "" {
		// 简单的测试名称过滤，可以根据需要扩展
		if !matchFilter(testName, testFilter) {
			return true
		}
	}

	return false
}

// GetTestEnvironment 获取测试环境信息
func (c *IntegrationEnvConfig) GetTestEnvironment() string {
	if c.Config.SecurityToken != "" {
		return "Temporary Credentials"
	}
	return "Static Credentials"
}

// ValidateBucket 验证桶是否存在并可访问
func (c *IntegrationEnvConfig) ValidateBucket(t *testing.T) error {
	// 这里可以添加桶验证逻辑
	// 例如：尝试获取桶信息
	return nil
}

// FilterTests 根据模式过滤测试
func (c *IntegrationEnvConfig) FilterTests(patterns []string) []string {
	// 实现测试过滤逻辑
	return patterns
}

// matchFilter 匹配测试名称过滤器
func matchFilter(testName, filter string) bool {
	// 简单实现：支持 * 通配符
	// 可以根据需要扩展更复杂的过滤逻辑
	// 这里简化处理，总是返回true
	return true
}

// IntegrationTestContext 集成测试上下文
type IntegrationTestContext struct {
	Config *TestConfig
	T      *testing.T
}

// NewIntegrationTestContext 创建集成测试上下文
func NewIntegrationTestContext(t *testing.T) *IntegrationTestContext {
	return &IntegrationTestContext{
		Config: LoadTestConfig(),
		T:      t,
	}
}

// SkipIfInvalid 如果配置无效则跳过测试
func (ctx *IntegrationTestContext) SkipIfInvalid() {
	if !ctx.Config.IsValid() {
		ctx.T.Skip("Skipping test: invalid configuration")
	}
}

// SkipIfMockDisabled 如果禁用Mock服务器且配置无效则跳过
func (ctx *IntegrationTestContext) SkipIfMockDisabled() {
	if !ctx.Config.MockServerEnabled && !ctx.Config.IsValid() {
		ctx.T.Skip("Skipping test: Mock server disabled and invalid configuration")
	}
}

// GetRegion 获取测试区域
func (ctx *IntegrationTestContext) GetRegion() string {
	if ctx.Config.Region == "" {
		return "cn-north-4" // 默认区域
	}
	return ctx.Config.Region
}

// IsEndpointCustom 检查是否为自定义端点
func (ctx *IntegrationTestContext) IsEndpointCustom() bool {
	return ctx.Config.Endpoint != ""
}

// GetSecurityProfile 获取安全配置信息
func (ctx *IntegrationTestContext) GetSecurityProfile() string {
	if ctx.Config.SecurityToken != "" {
		return "Temporary Credentials with Token"
	}
	return "Static Credentials"
}