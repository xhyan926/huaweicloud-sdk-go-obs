# 子任务 1.4：测试计划

## 测试目标
验证清单功能客户端方法的正确性和错误处理。

## 测试用例

### 1. SetBucketInventory 测试
```go
func TestSetBucketInventory_ShouldReturnSuccess_GivenValidInput(t *testing.T) {
    // 验证成功调用
}

func TestSetBucketInventory_ShouldReturnError_GivenNilInput(t *testing.T) {
    // 验证 nil 输入返回错误
}

func TestSetBucketInventory_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
}
```

### 2. GetBucketInventory 测试
```go
func TestGetBucketInventory_ShouldReturnSuccess_GivenValidParameters(t *testing.T) {
    // 验证成功获取
}

func TestGetBucketInventory_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
}
```

### 3. ListBucketInventory 测试
```go
func TestListBucketInventory_ShouldReturnList_GivenValidBucket(t *testing.T) {
    // 验证成功列举
}
```

### 4. DeleteBucketInventory 测试
```go
func TestDeleteBucketInventory_ShouldReturnSuccess_GivenValidParameters(t *testing.T) {
    // 验证成功删除
}

func TestDeleteBucketInventory_ShouldReturnError_GivenEmptyId(t *testing.T) {
    // 验证空 ID 返回错误
}
```

## 测试工具

- testify: 断言库
- MockRoundTripper: HTTP 模拟

## 验收标准

- [ ] 所有方法符合现有 API 风格
- [ ] 错误处理一致
- [ ] 测试覆盖率 > 90%
- [ ] 所有测试通过

## 执行步骤

1. 在 `obs/client_bucket_test.go` 中添加测试用例
2. 使用 MockRoundTripper 模拟 HTTP 响应
3. 运行测试：`go test ./... -v`
4. 检查覆盖率：`go test ./... -cover`
5. 修复发现的问题
