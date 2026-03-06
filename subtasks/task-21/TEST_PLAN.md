# 子任务 5.3：测试计划

## 测试目标
验证存量信息方法的正确性和错误处理。

## 测试用例

### 1. 成功场景测试
```go
func TestGetBucketStorageInfo_ShouldReturnSuccess_GivenValidBucket(t *testing.T) {
    // 验证成功获取
}
```

### 2. 错误场景测试
```go
func TestGetBucketStorageInfo_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
}
```

### 3. 响应解析测试
```go
func TestGetBucketStorageInfo_ShouldParseResponse_GivenLargeStorage(t *testing.T) {
    // 验证大数值正确解析
}
```

## 测试工具

- testify: 断言库
- MockRoundTripper: HTTP 模拟

## 验收标准

- [ ] 方法符合现有 API 风格
- [ ] 错误处理一致
- [ ] 测试覆盖率 > 80%
- [ ] 所有测试通过

## 执行步骤

1. 在 `obs/client_bucket_test.go` 中添加测试用例
2. 使用 MockRoundTripper 模拟 HTTP 响应
3. 运行测试：`go test ./... -v`
4. 检查覆盖率：`go test ./... -cover`
5. 修复发现的问题
