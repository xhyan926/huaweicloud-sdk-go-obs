# 子任务 4.4：实施计划

## 详细实施步骤

### 1. 使用 go-sdk-ut skill

**必须调用** `/go-sdk-ut` skill 来编写测试

### 2. 测试文件位置
- **目标文件**: `obs/client_bucket_test.go`

### 3. 测试用例结构

#### 3.1 SetBucketReplication 测试

```go
func TestSetBucketReplication_ShouldSetReplication_GivenValidInput(t *testing.T) {
    input := &SetBucketReplicationInput{
        Bucket: "src-bucket",
        ReplicationConfiguration: ReplicationConfiguration{
            Role: "arn:aws:iam::123456789012:role/some-role",
            Rules: []ReplicationRule{
                {
                    ID:     "rule-1",
                    Status: string(ReplicationStatusEnabled),
                    Prefix: "prefix/",
                    Destination: ReplicationDestination{
                        Bucket: "dest-bucket",
                    },
                },
            },
        },
    }

    expectedResp := `<BucketInfo><Id>src-bucket</Id></BucketInfo>`
    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(204, expectedResp, nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.SetBucketReplication(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

#### 3.2 GetBucketReplication 测试

```go
func TestGetBucketReplication_ShouldGetReplication_GivenValidBucket(t *testing.T) {
    bucket := "test-bucket"
    expectedResp := `<?xml version="1.0" encoding="UTF-8"?>
        <ReplicationConfiguration>
            <Role>arn:aws:iam::123456789012:role/some-role</Role>
            <Rule>
                <ID>rule-1</ID>
                <Status>Enabled</Status>
                <Prefix>prefix/</Prefix>
                <Destination><Bucket>dest-bucket</Bucket></Destination>
            </Rule>
        </ReplicationConfiguration>`

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(200, expectedResp, nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.GetBucketReplication(bucket)

    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.Equal(t, "arn:aws:iam::123456789012:role/some-role", output.Role)
    assert.Len(t, output.Rules, 1)
}
```

#### 3.3 DeleteBucketReplication 测试

```go
func TestDeleteBucketReplication_ShouldDeleteReplication_GivenValidBucket(t *testing.T) {
    bucket := "test-bucket"

    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(204, "", nil)
        },
    }

    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))
    output, err := client.DeleteBucketReplication(bucket)

    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

#### 3.4 边界条件测试

```go
func TestSetBucketReplication_ShouldHandleMultipleRules_GivenMultipleRules(t *testing.T) {
    // 验证多个复制规则
}

func TestSetBucketReplication_ShouldHandleEmptyPrefix_GivenNoPrefix(t *testing.T) {
    // 验证空前缀
}
```

### 4. 时间估算
- SetBucketReplication 测试：30 分钟
- GetBucketReplication 测试：30 分钟
- DeleteBucketReplication 测试：20 分钟
- 边界条件和错误测试：20 分钟
- 测试覆盖率优化：20 分钟
- **总计**: 约 2 小时（0.25 天）

## 技术要点

### BDD 风格命名
- 格式: Test<功能>_Should<预期结果>_When<条件>
- 描述测试意图
- 易于理解

### MockRoundTripper 使用
- 模拟 HTTP 响应
- 验证请求参数
- 控制测试场景

### testify 断言
- assert.NoError(t, err)
- assert.NotNil(t, output)
- assert.Equal(t, expected, actual)
- assert.Len(t, list, expectedLen)

### 测试覆盖率
- 目标: > 80%
- 使用 `go test -cover` 检查
- 添加边界条件测试提高覆盖率
