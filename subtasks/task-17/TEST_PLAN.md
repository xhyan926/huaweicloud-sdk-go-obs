# 子任务 4.3：测试计划

## 测试目标
验证跨区域复制方法的正确性和错误处理。

## 测试用例

### 1. SetBucketReplication 测试
```go
func TestSetBucketReplication_ShouldReturnSuccess_GivenValidInput(t *testing.T) {
    // 验证成功调用
}

func TestSetBucketReplication_ShouldReturnError_GivenNilInput(t *testing.T) {
    // 验证 nil 输入返回错误
}

func TestSetBucketReplication_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
}
```

### 2. GetBucketReplication 测试
```go
func TestGetBucketReplication_ShouldReturnSuccess_GivenValidBucket(t *testing.T) {
    // 验证成功获取
}

func TestGetBucketReplication_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
}
```

### 3. DeleteBucketReplication 测试
```go
func TestDeleteBucketReplication_ShouldReturnSuccess_GivenValidBucket(t *testing.T) {
    // 验证成功删除
}

func TestDeleteBucketReplication_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
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
