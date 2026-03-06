# 子任务 5.4：实施计划

## 详细实施步骤

### 1. 使用 go-sdk-ut skill

**必须调用** `/go-sdk-ut` skill 来编写测试

### 2. 测试文件位置
- **目标文件**: `obs/client_bucket_test.go`

### 3. 测试用例结构

#### 3.1 成功场景测试

```go
func TestGetBucketStorageInfo_ShouldReturnStorageInfo_GivenValidBucket(t *testing.T) {
    bucket := "test-bucket"
    expectedResp := `<?xml version="1.0" encoding="UTF-8"?>
        <StorageInfo>
            <Size>1234567890</Size>
            <ObjectNumber>100</ObjectNumber>
        </StorageInfo>`

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(200, expectedResp, nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.GetBucketStorageInfo(bucket)

    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.Equal(t, int64(1234567890), output.Size)
    assert.Equal(t, int64(100), output.ObjectNumber)
}
```

#### 3.2 错误场景测试

```go
func TestGetBucketStorageInfo_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    bucket := ""

    client := New(ak, sk, endpoint)
    output, err := client.GetBucketStorageInfo(bucket)

    assert.Error(t, err)
    assert.Nil(t, output)
    assert.Contains(t, err.Error(), "bucketName is empty")
}
```

#### 3.3 边界条件测试

```go
func TestGetBucketStorageInfo_ShouldHandleZeroStorage_GivenEmptyBucket(t *testing.T) {
    // 验证空桶的存量信息
}

func TestGetBucketStorageInfo_ShouldHandleLargeNumbers_GivenBigStorage(t *testing.T) {
    // 验证大数值处理
}
```

### 4. 时间估算
- 成功场景测试：15 分钟
- 错误场景测试：15 分钟
- 边界条件测试：15 分钟
- 测试覆盖率优化：15 分钟
- **总计**: 约 1 小时（0.125 天）

## 技术要点

### BDD 风格命名
- 格式: Test<功能>_Should<预期结果>_When<条件>
- 描述测试意图
- 易于维护

### MockRoundTripper 使用
- 模拟 HTTP 响应
- 验证请求参数
- 控制测试场景

### testify 断言
- assert.NoError(t, err)
- assert.NotNil(t, output)
- assert.Equal(t, expected, actual)
- assert.Contains(t, str, substr)

### 测试覆盖率
- 目标: > 80%
- 使用 `go test -cover` 检查
- 添加边界条件测试提高覆盖率
