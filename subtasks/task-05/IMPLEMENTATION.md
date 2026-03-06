# 子任务 1.5：实施计划

## 详细实施步骤

### 1. 使用 go-sdk-ut skill

**必须调用** `/go-sdk-ut` skill 来编写测试

### 2. 测试文件位置
- **目标文件**: `obs/client_bucket_test.go`
- **追加位置**: 在现有测试之后

### 3. 测试用例结构

#### 3.1 SetBucketInventory 测试

```go
func TestSetBucketInventory_ShouldSetInventory_GivenValidInput(t *testing.T) {
    // 准备测试数据
    input := &SetBucketInventoryInput{
        Bucket: "test-bucket",
        InventoryConfiguration: InventoryConfiguration{
            Id:          "test-inventory",
            IsEnabled:    true,
            Destination:  InventoryDestination{Format: "CSV", Bucket: InventoryBucket{Name: "dest-bucket"}, Prefix: "inventory/"},
            Schedule:     InventorySchedule{Frequency: string(InventoryFrequencyDaily)},
        },
    }

    // 创建模拟响应
    expectedResp := `<BucketInfo><Id>test-bucket</Id></BucketInfo>`
    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(204, expectedResp, nil)
        },
    }

    // 创建客户端
    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))

    // 执行测试
    output, err := client.SetBucketInventory(input)

    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

#### 3.2 GetBucketInventory 测试

```go
func TestGetBucketInventory_ShouldGetInventory_GivenValidId(t *testing.T) {
    // 准备测试数据
    bucket := "test-bucket"
    id := "test-inventory"
    expectedResp := `<?xml version="1.0" encoding="UTF-8"?>
        <InventoryConfiguration>
            <Id>test-inventory</Id>
            <IsEnabled>true</IsEnabled>
            <Destination>
                <Format>CSV</Format>
                <Bucket><Name>dest-bucket</Name></Bucket>
                <Prefix>inventory/</Prefix>
            </Destination>
            <Schedule>
                <Frequency>Daily</Frequency>
            </Schedule>
        </InventoryConfiguration>`

    // 创建模拟响应
    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(200, expectedResp, nil)
        },
    }

    // 创建客户端
    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))

    // 执行测试
    output, err := client.GetBucketInventory(bucket, id)

    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.Equal(t, "test-inventory", output.Id)
    assert.True(t, output.IsEnabled)
}
```

#### 3.3 ListBucketInventory 测试

```go
func TestListBucketInventory_ShouldListAllInventories_GivenValidBucket(t *testing.T) {
    // 准备测试数据
    bucket := "test-bucket"
    expectedResp := `<?xml version="1.0" encoding="UTF-8"?>
        <ListInventoryConfigurationsResult>
            <InventoryConfiguration>
                <Id>inventory-1</Id>
                <IsEnabled>true</IsEnabled>
            </InventoryConfiguration>
            <InventoryConfiguration>
                <Id>inventory-2</Id>
                <IsEnabled>false</IsEnabled>
            </InventoryConfiguration>
        </ListInventoryConfigurationsResult>`

    // 创建模拟响应
    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(200, expectedResp, nil)
        },
    }

    // 创建客户端
    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))

    // 执行测试
    output, err := client.ListBucketInventory(bucket)

    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.Len(t, output.InventoryConfigurationList, 2)
}
```

#### 3.4 DeleteBucketInventory 测试

```go
func TestDeleteBucketInventory_ShouldDeleteInventory_GivenValidId(t *testing.T) {
    // 准备测试数据
    bucket := "test-bucket"
    id := "test-inventory"

    // 创建模拟响应
    mockTransport := &MockRoundTripper{
        ResponseFunc: func(req *http.Request) *http.Response {
            return CreateTestHTTPResponse(204, "", nil)
        },
    }

    // 创建客户端
    client := New(ak, sk, endpoint, WithHttpTransport(mockTransport))

    // 执行测试
    output, err := client.DeleteBucketInventory(bucket, id)

    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

### 4. 时间估算
- SetBucketInventory 测试：30 分钟
- GetBucketInventory 测试：30 分钟
- ListBucketInventory 测试：30 分钟
- DeleteBucketInventory 测试：30 分钟
- 边界条件和错误测试：30 分钟
- **总计**: 约 2.5 小时（0.31 天）

## 技术要点

### BDD 风格命名
- 格式: Test<功能>_Should<预期结果>_When<条件>
- 描述测试的意图
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
