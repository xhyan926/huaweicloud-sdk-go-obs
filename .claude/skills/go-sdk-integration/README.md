# OBS SDK Go 集成测试指南

## 快速开始

### 1. 安装依赖

确保已安装Go 1.19+：

```bash
go version
```

### 2. 环境配置

设置必需的环境变量：

```bash
export OBS_TEST_AK="your-access-key"
export OBS_TEST_SK="your-secret-key"
export OBS_TEST_ENDPOINT="https://obs.cn-north-4.myhuaweicloud.com"
export OBS_TEST_BUCKET="your-test-bucket"
```

### 3. 运行集成测试

```bash
# 进入项目目录
cd /path/to/huaweicloud-sdk-go-obs

# 运行所有集成测试
go test -tags integration ./obs/test/integration -v

# 运行特定测试
go test -tags integration ./obs/test/integration/e2e -run TestBasicObjectOperations -v
```

## 使用集成测试客户端

### 基本用法

```go
//go:build integration

package e2e

import (
    "testing"
    "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
    "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

func TestExample(t *testing.T) {
    // 创建测试客户端
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    // 获取测试桶
    bucket := client.GetTestBucket()

    // 上传对象
    objectKey := client.GetTestObjectKey("test.txt")
    content := "test content"
    err := client.TestClient().PutObject(bucket, objectKey, content, nil)
    if err != nil {
        t.Fatalf("上传失败: %v", err)
    }

    // 清理资源（自动执行）
}
```

### 高级用法

#### 1. 使用Mock服务器

```go
func TestWithMockServer(t *testing.T) {
    // 创建Mock服务器
    mockServer := integration.NewMockServer()
    mockServer.Start(":8080")
    defer mockServer.Stop()

    // 创建指向Mock服务器的客户端
    mockEndpoint := "http://localhost:8080"
    client := integration.NewTestClientWithEndpoint(t, mockEndpoint)
    defer client.Cleanup(t)

    // 执行测试...
}
```

#### 2. 注册自定义清理函数

```go
func TestWithCustomCleanup(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    // 注册自定义清理函数
    client.AddCleanup(func(t *testing.T) {
        // 清理特定资源
        t.Log("执行自定义清理")
    })

    // 执行测试...
}
```

#### 3. 使用带超时的Context

```go
func TestWithTimeout(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    // 创建带超时的context
    ctx, cancel := client.WithContext(30 * time.Second)
    defer cancel()

    // 使用context...
}
```

## 测试场景示例

### 认证测试

```go
func TestAuthentication(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    // 测试静态凭证
    t.Run("Static Credentials", func(t *testing.T) {
        input := &obs.GetBucketLocationInput{
            Bucket: client.GetTestBucket(),
        }

        _, err := client.TestClient().GetBucketLocation(input)
        if err != nil {
            t.Fatalf("认证失败: %v", err)
        }
    })

    // 测试临时凭证
    if os.Getenv("OBS_TEST_TOKEN") != "" {
        t.Run("Temporary Credentials", func(t *testing.T) {
            input := &obs.GetBucketLocationInput{
                Bucket: client.GetTestBucket(),
            }

            _, err := client.TestClient().GetBucketLocation(input)
            if err != nil {
                t.Fatalf("认证失败: %v", err)
            }
        })
    }
}
```

### 对象操作测试

```go
func TestObjectOperations(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    bucket := client.GetTestBucket()
    objectKey := client.GetTestObjectKey("test-object.txt")
    content := "Test content for object operations"

    // 上传对象
    t.Run("Upload Object", func(t *testing.T) {
        err := client.TestClient().PutObject(bucket, objectKey, content, nil)
        if err != nil {
            t.Fatalf("上传失败: %v", err)
        }

        // 注册清理函数
        client.AddCleanup(func(t *testing.T) {
            err := client.TestClient().DeleteObject(bucket, objectKey)
            if err != nil {
                t.Logf("删除失败: %v", err)
            }
        })
    })

    // 获取对象
    t.Run("Get Object", func(t *testing.T) {
        input := &obs.GetObjectInput{
            Bucket: bucket,
            Key:    objectKey,
        }

        output, err := client.TestClient().GetObject(input)
        if err != nil {
            t.Fatalf("获取失败: %v", err)
        }
        defer output.Body.Close()

        // 验证内容...
    })
}
```

### 分块上传测试

```go
func TestMultipartUpload(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    bucket := client.GetTestBucket()
    objectKey := client.GetTestObjectKey("multipart-test.txt")

    // 初始化分块上传
    initReq := &obs.InitiateMultipartUploadInput{
        Bucket: bucket,
        Key:    objectKey,
    }

    initOutput, err := client.TestClient().InitiateMultipartUpload(initReq)
    if err != nil {
        t.Fatalf("初始化分块上传失败: %v", err)
    }

    // 记录清理函数
    client.AddCleanup(func(t *testing.T) {
        abortReq := &obs.AbortMultipartUploadInput{
            Bucket:   bucket,
            Key:      objectKey,
            UploadId: initOutput.UploadId,
        }
        client.TestClient().AbortMultipartUpload(abortReq)
    })

    // 上传分块...
}
```

## 环境变量说明

| 变量名 | 必需 | 说明 |
|--------|------|------|
| OBS_TEST_AK | 是 | 访问密钥 Access Key |
| OBS_TEST_SK | 是 | 访问密钥 Secret Key |
| OBS_TEST_ENDPOINT | 是 | OBS服务端点 |
| OBS_TEST_BUCKET | 是 | 测试桶名称 |
| OBS_TEST_REGION | 否 | 区域 |
| OBS_TEST_TOKEN | 否 | 临时安全令牌 |
| OBS_MOCK_ENABLED | 否 | 是否启用Mock服务器 |
| OBS_MOCK_PORT | 否 | Mock服务器端口 |
| OBS_SKIP_INTEGRATION_TESTS | 否 | 是否跳过集成测试 |

## 命令行工具

### 运行所有集成测试

```bash
make test-integration
```

### 运行特定测试文件

```bash
go test -tags integration ./obs/test/integration/e2e -v
```

### 运行带过滤的测试

```bash
# 只运行桶相关测试
go test -tags integration ./obs/test/integration/e2e -run "Bucket" -v

# 只运行对象相关测试
go test -tags integration ./obs/test/integration/e2e -run "Object" -v
```

## Mock服务器使用

### 启动Mock服务器

```go
func TestWithMock(t *testing.T) {
    mockServer := integration.NewMockServer()
    mockServer.Start(":8080")
    defer mockServer.Stop()

    client := integration.NewTestClientWithEndpoint(t, "http://localhost:8080")
    defer client.Cleanup(t)

    // 执行测试...
}
```

### 查看Mock服务器统计

```go
func TestWithMock(t *testing.T) {
    mockServer := integration.NewMockServer()
    mockServer.Start(":8080")
    defer mockServer.Stop()

    // 执行测试...

    // 查看统计信息
    stats := mockServer.GetStats()
    t.Logf("Mock服务器统计: %+v", stats)
}
```

## 故障排除

### 常见错误

1. **认证失败**
   ```
   Error: OBS_TEST_AK not set
   ```
   解决：检查环境变量是否正确设置

2. **连接超时**
   ```
   context deadline exceeded
   ```
   解决：检查网络连接，或增加超时时间

3. **Mock服务器启动失败**
   ```
   address already in use
   ```
   解决：更换端口，或关闭占用端口的进程

### 调试技巧

1. **启用详细日志**
   ```bash
   go test -tags integration ./test/integration -v
   ```

2. **使用Mock服务器**
   ```bash
   export OBS_MOCK_ENABLED=true
   export OBS_MOCK_PORT=8080
   ```

3. **逐步测试**
   ```bash
   go test -tags integration ./test/integration -run "TestSpecificTest" -v
   ```

## 最佳实践

1. **测试隔离**：每个测试用例应该独立运行
2. **资源清理**：使用`defer client.Cleanup(t)`确保资源清理
3. **Mock服务器**：开发环境优先使用Mock服务器
4. **错误处理**：提供详细的错误信息
5. **超时设置**：设置合理的超时时间

## 下一步

1. 阅读完整文档：[skill.md](skill.md)
2. 查看示例代码：[example_test.go](../obs/test/integration/e2e/example_test.go)
3. 使用技能生成测试：`/go-sdk-integration`