# OBS SDK Go 集成测试规范

## 概述

集成测试验证OBS SDK与真实OBS服务或Mock服务之间的端到端功能，确保各个组件能够正确协作。本规范定义了集成测试的架构、实现方式和最佳实践。

## 测试范围

### 核心功能测试

1. **认证测试**
   - 静态凭证认证
   - 临时凭证认证（Security Token）
   - 认证失败处理

2. **存储桶操作**
   - 创建桶
   - 获取桶信息
   - 列举桶
   - 删除桶

3. **对象操作**
   - 上传对象（PutObject）
   - 获取对象（GetObject）
   - 删除对象（DeleteObject）
   - 获取对象元数据（GetObjectMetadata）
   - 复制对象（CopyObject）

4. **高级功能**
   - 分块上传（Multipart Upload）
   - 断点续传（Resumable Upload）
   - 下载进度（Download Progress）
   - 上传回调（Upload Callback）

### 错误场景测试

1. **认证错误**
   - 无效的Access Key
   - 无效的Secret Key
   - 过期的Token

2. **网络错误**
   - 连接超时
   - 网络中断
   - 服务器错误

3. **业务错误**
   - 不存在的桶
   - 不存在的对象
   - 权限不足

## 测试环境配置

### 环境变量

| 变量名 | 必需 | 说明 | 默认值 |
|--------|------|------|--------|
| OBS_TEST_AK | 是 | 访问密钥 Access Key | - |
| OBS_TEST_SK | 是 | 访问密钥 Secret Key | - |
| OBS_TEST_ENDPOINT | 是 | OBS服务端点 | - |
| OBS_TEST_BUCKET | 是 | 测试桶名称 | - |
| OBS_TEST_REGION | 否 | 区域 | cn-north-4 |
| OBS_TEST_TOKEN | 否 | 临时安全令牌 | - |
| OBS_MOCK_ENABLED | 否 | 是否启用Mock服务器 | false |
| OBS_MOCK_PORT | 否 | Mock服务器端口 | 8080 |
| OBS_SKIP_INTEGRATION_TESTS | 否 | 是否跳过集成测试 | false |

### 配置文件

```go
// obs/test_config.go
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

    // Mock服务器配置
    MockServerEnabled bool
    MockServerPort    int

    // 跳过测试标记
    SkipIntegrationTests bool
}
```

## 测试客户端

### IntegrationClient设计

```go
// obs/test/integration/client.go
type TestClient struct {
    ObsClient    *obs.ObsClient
    Config       *config.TestConfig
    TestPrefix   string
    CleanupFuncs []CleanupFunction
    TestCases    []string
}
```

### 核心功能

1. **自动跳过机制**
   ```go
   func NewTestClient(t *testing.T) *TestClient {
       cfg := config.NewIntegrationEnvConfig(t)
       if cfg.ShouldSkipIntegrationTest() {
           t.Skip("Skipping integration tests")
       }
       // 创建客户端...
   }
   ```

2. **Mock服务器支持**
   ```go
   func (c *TestClient) WithMockServer() *TestClient {
       mockServer := NewMockServer()
       mockServer.Start(":8080")
       // 创建指向Mock的客户端...
       return c
   }
   ```

3. **资源清理管理**
   ```go
   func (c *TestClient) AddCleanup(f CleanupFunction) {
       c.CleanupFuncs = append(c.CleanupFuncs, f)
   }

   func (c *TestClient) Cleanup(t *testing.T) {
       for i := len(c.CleanupFuncs) - 1; i >= 0; i-- {
           c.CleanupFuncs[i](t)
       }
   }
   ```

4. **测试用例记录**
   ```go
   func (c *TestClient) AddTestCase(testCase string) {
       c.TestCases = append(c.TestCases, testCase)
   }
   ```

## 测试实现规范

### 1. 基本测试结构

```go
//go:build integration

package e2e

import (
    "testing"
    "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
    "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// TestBucketOperations_ShouldSucceed_GivenValidCredentials
func TestBucketOperations(t *testing.T) {
    // 创建测试客户端
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    bucket := client.GetTestBucket()

    // 测试创建桶
    t.Run("CreateBucket", func(t *testing.T) {
        input := &obs.CreateBucketInput{
            Bucket: bucket,
        }
        _, err := client.TestClient().CreateBucket(input)
        if err != nil {
            t.Fatalf("创建桶失败: %v", err)
        }

        client.AddTestCase("创建桶成功")
    })

    // 测试删除桶
    t.Run("DeleteBucket", func(t *testing.T) {
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
```

### 2. 对象操作测试

```go
// TestObjectOperations_ShouldSucceed_GivenValidBucket
func TestObjectOperations(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    bucket := client.GetTestBucket()
    objectKey := client.GetTestObjectKey("test-object.txt")
    content := "This is a test object content."

    // 测试上传对象
    t.Run("PutObject", func(t *testing.T) {
        err := client.TestClient().PutObject(bucket, objectKey, content, nil)
        if err != nil {
            t.Fatalf("上传对象失败: %v", err)
        }

        // 注册清理函数
        client.AddCleanup(func(t *testing.T) {
            err := client.TestClient().DeleteObject(bucket, objectKey)
            if err != nil {
                t.Logf("删除对象失败: %v", err)
            }
        })

        client.AddTestCase("上传对象成功")
    })

    // 测试获取对象
    t.Run("GetObject", func(t *testing.T) {
        input := &obs.GetObjectInput{
            Bucket: bucket,
            Key:    objectKey,
        }

        output, err := client.TestClient().GetObject(input)
        if err != nil {
            t.Fatalf("获取对象失败: %v", err)
        }
        defer output.Body.Close()

        // 读取并验证内容
        body := make([]byte, len(content))
        _, err = output.Body.Read(body)
        if err != nil {
            t.Fatalf("读取对象内容失败: %v", err)
        }

        if string(body) != content {
            t.Errorf("对象内容不匹配，期望: %s, 实际: %s", content, string(body))
        }

        client.AddTestCase("获取对象成功")
    })
}
```

### 3. 分块上传测试

```go
// TestMultipartUpload_ShouldSucceed_GivenLargeFile
func TestMultipartUpload(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    bucket := client.GetTestBucket()
    objectKey := client.GetTestObjectKey("multipart-test.dat")

    // 创建100MB的测试数据
    content := bytes.Repeat([]byte("test"), 25*1024*1024) // 100MB

    // 初始化分块上传
    t.Run("InitiateMultipartUpload", func(t *testing.T) {
        initReq := &obs.InitiateMultipartUploadInput{
            Bucket: bucket,
            Key:    objectKey,
        }

        initOutput, err := client.TestClient().InitiateMultipartUpload(initReq)
        if err != nil {
            t.Fatalf("初始化分块上传失败: %v", err)
        }

        // 注册清理函数
        client.AddCleanup(func(t *testing.T) {
            abortReq := &obs.AbortMultipartUploadInput{
                Bucket:   bucket,
                Key:      objectKey,
                UploadId: initOutput.UploadId,
            }
            client.TestClient().AbortMultipartUpload(abortReq)
        })

        client.AddTestCase("初始化分块上传成功")
    })
}
```

### 4. 错误场景测试

```go
// TestErrorScenarios_ShouldFail_GivenInvalidInput
func TestErrorScenarios(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    bucket := client.GetTestBucket()
    nonexistentObject := client.GetTestObjectKey("nonexistent.txt")

    // 测试获取不存在的对象
    t.Run("GetObject_Nonexistent", func(t *testing.T) {
        input := &obs.GetObjectInput{
            Bucket: bucket,
            Key:    nonexistentObject,
        }

        _, err := client.TestClient().GetObject(input)
        if err == nil {
            t.Error("预期获取不存在的对象会失败")
        }

        // 验证错误类型
        if obsError, ok := err.(obs.ObsError); ok {
            if obsError.StatusCode != 404 {
                t.Errorf("预期404错误，实际: %d", obsError.StatusCode)
            }
        } else {
            t.Error("预期ObsError类型")
        }

        client.AddTestCase("正确处理不存在的对象")
    })
}
```

### 5. Mock服务器测试

```go
// TestWithMockServer_ShouldSimulateOBSBehavior
func TestWithMockServer(t *testing.T) {
    // 创建Mock服务器
    mockServer := integration.NewMockServer()
    err := mockServer.Start(":8080")
    if err != nil {
        t.Fatalf("启动Mock服务器失败: %v", err)
    }
    defer mockServer.Stop()

    // 创建指向Mock服务器的客户端
    mockEndpoint := "http://localhost:8080"
    client := integration.NewTestClientWithEndpoint(t, mockEndpoint)
    defer client.Cleanup(t)

    bucket := client.GetTestBucket()
    objectKey := client.GetTestObjectKey("mock-test.txt")
    content := "Mock server test content."

    // 测试上传到Mock服务器
    t.Run("PutObject_ToMock", func(t *testing.T) {
        err := client.TestClient().PutObject(bucket, objectKey, content, nil)
        if err != nil {
            t.Fatalf("上传到Mock服务器失败: %v", err)
        }

        // 验证Mock服务器记录
        requests := mockServer.GetRequests()
        if len(requests) == 0 {
            t.Error("Mock服务器未记录任何请求")
        }

        // 查找上传请求
        found := false
        for _, req := range requests {
            if req.Method == "PUT" && req.URL == fmt.Sprintf("/%s/%s", bucket, objectKey) {
                found = true
                break
            }
        }

        if !found {
            t.Error("未找到上传请求记录")
        }

        client.AddTestCase("Mock服务器测试成功")
    })
}
```

## 测试命名规范

### BDD命名格式

使用 `Test{Function}_Should{Expected}_When{Condition}_Given{Precondition}` 格式：

```
TestBucketOperations_ShouldSucceed_GivenValidCredentials
TestObjectUpload_ShouldFail_GivenInvalidFile
TestAuthentication_ShouldConnect_GivenTempCredentials
```

### 测试场景命名

1. **成功场景**：`ShouldSucceed`、`ShouldCreateSuccessfully`
2. **失败场景**：`ShouldFail`、`ShouldReturnError`
3. **边界条件**：`ShouldHandleEdgeCase`、`ShouldValidateInput`
4. **性能测试**：`ShouldPerformWithinLimits`

## 资源清理规范

### 1. 清理函数设计

```go
// 清理函数类型
type CleanupFunction func(t *testing.T)

// 示例：清理测试对象
func cleanupObject(t *testing.T, client *integration.TestClient, bucket, key string) {
    err := client.TestClient().DeleteObject(bucket, key)
    if err != nil {
        t.Logf("清理对象失败: %v", err)
    }
}

// 注册清理函数
client.AddCleanup(func(t *testing.T) {
    cleanupObject(t, client, bucket, objectKey)
})
```

### 2. 清理时机

1. **测试完成时**：使用`defer client.Cleanup(t)`
2. **测试失败时**：自动执行清理
3. **测试异常时**：确保清理执行

### 3. 清理顺序

按照创建的反序执行清理：
1. 先删除对象
2. 再删除桶
3. 最后关闭连接

## 性能考虑

### 1. 测试优化

1. **复用连接**：使用HTTP keep-alive
2. **并发测试**：使用`RunParallel`
3. **批量操作**：减少API调用次数

### 2. 超时设置

```go
func (c *TestClient) WithContext(timeout time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), timeout)
}
```

### 3. 资源限制

- 单个测试最大执行时间：5分钟
- 内存使用限制：1GB
- 并发请求数：100

## 最佳实践

### 1. 测试独立性

- 每个测试用例独立运行
- 不依赖其他测试的状态
- 使用唯一的测试对象名

### 2. 错误处理

```go
// 好的实践
_, err := client.TestClient().GetObject(input)
if err != nil {
    if obsError, ok := err.(obs.ObsError); ok {
        // 处理特定错误
        if obsError.StatusCode == 404 {
            t.Log("对象不存在，这是预期的")
            return
        }
    }
    t.Fatalf("获取对象失败: %v", err)
}

// 不好的实践
_, err := client.TestClient().GetObject(input)
if err != nil {
    t.Fatal(err)
}
```

### 3. 测试数据管理

```go
// 使用测试前缀避免冲突
objectKey := client.GetTestObjectKey("unique-name")

// 使用固定测试数据
const testContent = "This is a fixed test content"
```

### 4. 日志记录

```go
// 记录测试步骤
t.Logf("开始测试: %s", testCase)
t.Logf("使用桶: %s", bucket)
t.Logf("创建对象: %s", objectKey)

// 记录测试结果
if err == nil {
    t.Logf("测试成功: %s", testCase)
} else {
    t.Logf("测试失败: %s, 错误: %v", testCase, err)
}
```

## 常见问题

### 1. 环境配置问题

**问题**：`Error: OBS_TEST_AK not set`

**解决方案**：
```bash
export OBS_TEST_AK="your-access-key"
export OBS_TEST_SK="your-secret-key"
export OBS_TEST_ENDPOINT="https://obs.cn-north-4.myhuaweicloud.com"
export OBS_TEST_BUCKET="your-test-bucket"
```

### 2. 网络连接问题

**问题**：`context deadline exceeded`

**解决方案**：
```go
// 增加超时时间
ctx, cancel := client.WithContext(30 * time.Second)
defer cancel()

// 使用带超时的context
input := &obs.GetObjectInput{
    Bucket: bucket,
    Key:    objectKey,
}
input.WithContext(ctx)
```

### 3. 权限问题

**问题**：`Access Denied`

**解决方案**：
- 检查Access Key和Secret Key是否正确
- 确保测试桶有正确的权限
- 检查区域是否匹配

## 测试报告

### 1. 测试用例统计

```go
// 打印测试用例记录
client.PrintTestCases()

// 输出示例：
// === Test Cases Executed ===
// 1. 创建桶成功
// 2. 上传对象成功
// 3. 获取对象成功
// 4. 删除对象成功
// 5. 删除桶成功
// ============================
```

### 2. 性能统计

```go
// 记录性能数据
startTime := time.Now()
// 执行测试...
elapsed := time.Since(startTime)
t.Logf("执行时间: %v", elapsed)
```

### 3. 错误统计

```go
// 统计错误类型
errorTypes := make(map[string]int)
for _, err := range errors {
    errorTypes[err.Type]++
}

t.Logf("错误统计: %v", errorTypes)
```

## 相关工具

### 1. Mock服务器工具

```go
// 启动Mock服务器
mockServer := integration.NewMockServer()
mockServer.Start(":8080")
defer mockServer.Stop()

// 查看统计信息
stats := mockServer.GetStats()
t.Logf("Mock服务器统计: %v", stats)
```

### 2. 测试数据生成器

```go
// 生成测试数据
func generateTestData(size int) []byte {
    return bytes.Repeat([]byte("test"), size)
}
```

### 3. 性能测试工具

```go
// 测量执行时间
func measurePerformance(t *testing.T, f func()) time.Duration {
    start := time.Now()
    f()
    return time.Since(start)
}
```

## 版本控制

- 版本：1.0.0
- 最后更新：2024-01-01
- 兼容性：OBS SDK v3.x

## 相关文档

- [测试架构设计](../architecture.md)
- [单元测试规范](unit-testing.md)
- [性能测试规范](performance-testing.md)
- [模糊测试规范](fuzz-testing.md)