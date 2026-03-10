package config

import (
	"os"
	"strconv"
)

// TestConfig 测试配置管理
type TestConfig struct {
	// 基础配置
	AccessKey       string
	SecretKey       string
	SecurityToken   string
	Endpoint        string
	Region          string

	// 测试资源
	TestBucket      string
	TestObject      string
	TestPrefix      string

	// 性能测试配置
	PerfLargeFileSize   int64  // 大文件大小（字节）
	PerfConcurrency     int    // 并发数
	PerfTestDuration    int    // 测试时长（秒）

	// Mock服务器配置
	MockServerEnabled bool
	MockServerPort    int

	// 跳过测试标记
	SkipIntegrationTests bool
	SkipFuzzTests       bool
	SkipPerfTests       bool

	// 日志配置
	LogLevel string
	LogPath  string
}

// LoadTestConfig 从环境变量加载测试配置
func LoadTestConfig() *TestConfig {
	cfg := &TestConfig{
		AccessKey:        getEnvString("OBS_TEST_AK", ""),
		SecretKey:        getEnvString("OBS_TEST_SK", ""),
		SecurityToken:    getEnvString("OBS_TEST_TOKEN", ""),
		Endpoint:         getEnvString("OBS_TEST_ENDPOINT", ""),
		Region:           getEnvString("OBS_TEST_REGION", ""),
		TestBucket:       getEnvString("OBS_TEST_BUCKET", ""),
		TestObject:       getEnvString("OBS_TEST_OBJECT", ""),
		TestPrefix:       "test-" + getEnvString("OBS_TEST_PREFIX", "default") + "-",

		PerfLargeFileSize: getEnvInt64("OBS_PERF_LARGE_FILE_SIZE", 100*1024*1024), // 默认100MB
		PerfConcurrency:  getEnvInt("OBS_PERF_CONCURRENCY", 100),
		PerfTestDuration: getEnvInt("OBS_PERF_TEST_DURATION", 30),

		MockServerEnabled: getEnvBool("OBS_MOCK_ENABLED", false),
		MockServerPort:    getEnvInt("OBS_MOCK_PORT", 8080),

		SkipIntegrationTests: getEnvBool("OBS_SKIP_INTEGRATION_TESTS", false),
		SkipFuzzTests:       getEnvBool("OBS_SKIP_FUZZ_TESTS", false),
		SkipPerfTests:       getEnvBool("OBS_SKIP_PERF_TESTS", false),

		LogLevel: getEnvString("OBS_LOG_LEVEL", "INFO"),
		LogPath:  getEnvString("OBS_LOG_PATH", ""),
	}

	return cfg
}

// IsValid 检查配置是否有效
func (c *TestConfig) IsValid() bool {
	if c.SkipIntegrationTests {
		return true
	}
	return c.AccessKey != "" && c.SecretKey != "" && c.Endpoint != "" && c.TestBucket != ""
}

// ShouldSkipIntegrationTest 检查是否应该跳过集成测试
func (c *TestConfig) ShouldSkipIntegrationTest() bool {
	return c.SkipIntegrationTests || !c.IsValid()
}

// GetTestBucket 获取测试桶名称（带前缀）
func (c *TestConfig) GetTestBucket() string {
	return c.TestBucket
}

// GetTestObjectKey 获取测试对象键（带前缀）
func (c *TestConfig) GetTestObjectKey(key string) string {
	return c.TestPrefix + key
}

// CleanTestObjectKey 清理测试对象键（移除前缀）
func (c *TestConfig) CleanTestObjectKey(key string) string {
	if len(key) > len(c.TestPrefix) && key[:len(c.TestPrefix)] == c.TestPrefix {
		return key[len(c.TestPrefix):]
	}
	return key
}

// getEnvString 获取环境变量，支持默认值
func getEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// getEnvInt 获取环境变量（整数），支持默认值
func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 32); err == nil {
			return int(i)
		}
	}
	return defaultValue
}

// getEnvInt64 获取环境变量（int64），支持默认值
func getEnvInt64(key string, defaultValue int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

// getEnvBool 获取环境变量（布尔值），支持默认值
func getEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultValue
}