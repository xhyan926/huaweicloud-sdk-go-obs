# 子任务 3.4：实施计划

## 详细实施步骤

### 1. 使用 go-sdk-ut skill

**必须调用** `/go-sdk-ut` skill 来编写测试

### 2. 测试文件位置
- **目标文件**: `obs/client_bucket_test.go`
- **目标文件**: `obs/client_object_test.go`

### 3. CreateBucket 参数测试

```go
func TestCreateBucket_ShouldSetBucketType_GivenBucketTypeParameter(t *testing.T) {
    input := &CreateBucketInput{
        Bucket:     "test-bucket",
        BucketType: "POSIX",
    }

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            // 验证请求头包含 BucketType
            assert.Equal(t, "POSIX", req.Header.Get("x-obs-bucket-type"))
            return CreateTestHTTPResponse(200, "", nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.CreateBucket(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
}

func TestCreateBucket_ShouldSetKmsKeyId_GivenKmsKeyId(t *testing.T) {
    input := &CreateBucketInput{
        Bucket:     "test-bucket",
        SseKmsKeyId: "kms-key-id-123",
    }

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            assert.Equal(t, "kms-key-id-123", req.Header.Get("x-obs-server-side-encryption-kms-key-id"))
            return CreateTestHTTPResponse(200, "", nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.CreateBucket(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

### 4. PutObject 参数测试

```go
func TestPutObject_ShouldSetExpires_GivenExpiresParameter(t *testing.T) {
    input := &PutObjectInput{
        Bucket:  "test-bucket",
        Key:     "test.txt",
        Expires: 30, // 30 天后过期
    }

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            assert.Equal(t, "30", req.Header.Get("x-obs-expires"))
            return CreateTestHTTPResponse(200, "", nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.PutObject(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
}

func TestPutObject_ShouldSetObjectLock_GivenLockParameters(t *testing.T) {
    input := &PutObjectInput{
        Bucket:                   "test-bucket",
        Key:                      "test.txt",
        ObjectLockMode:            "COMPLIANCE",
        ObjectLockRetainUntilDate: "2026-03-05T10:30:00.000Z",
    }

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            assert.Equal(t, "COMPLIANCE", req.Header.Get("x-obs-object-lock-mode"))
            assert.Equal(t, "2026-03-05T10:30:00.000Z", req.Header.Get("x-obs-object-lock-retain-until-date"))
            return CreateTestHTTPResponse(200, "", nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.PutObject(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

### 5. ListObjects 参数测试

```go
func TestListObjects_ShouldAddEncodingType_GivenEncodingTypeURL(t *testing.T) {
    input := &ListObjectsInput{
        Bucket:      "test-bucket",
        EncodingType: "url",
    }

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            // 验证查询字符串包含 encoding-type=url
            assert.Contains(t, req.URL.RawQuery, "encoding-type=url")
            return CreateTestHTTPResponse(200, mockListResponse, nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.ListObjects(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.Equal(t, "url", output.EncodingType)
}
```

### 6. 向后兼容性测试

```go
func TestParameterAdditions_ShouldNotBreakExistingFeatures(t *testing.T) {
    // 运行所有现有的测试
    // 确保没有任何测试因新增参数而失败
}
```

### 7. 时间估算
- CreateBucket 测试：30 分钟
- PutObject 测试：30 分钟
- ListObjects 测试：20 分钟
- 向后兼容性测试：20 分钟
- 测试覆盖率优化：20 分钟
- **总计**: 约 2 小时（0.25 天）

## 技术要点

### MockRoundTripper 验证
- 在响应函数中验证请求参数
- 检查 HTTP 头
- 检查查询字符串
- 确保参数正确传递

### 测试覆盖
- 每个新参数至少一个测试
- 向后兼容性测试
- 边界条件测试
- 错误场景测试

### 现有测试验证
- 运行完整测试套件
- 确保现有测试继续通过
- 检查测试覆盖率
