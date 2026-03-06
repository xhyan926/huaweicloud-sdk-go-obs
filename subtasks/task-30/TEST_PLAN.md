# 子任务 8.2：测试计划

## 测试目标
验证在线解压方法的正确性。

## 测试用例

### 1. SetZipPolicy 测试
```go
func TestSetZipPolicy_ShouldReturnSuccess_GivenValidInput(t *testing.T) {
    // 验证成功调用
}

func TestSetZipPolicy_ShouldReturnError_GivenNilInput(t *testing.T) {
    // 验证 nil 输入返回错误
}
```

### 2. GetZipPolicy 测试
```go
func TestGetZipPolicy_ShouldReturnSuccess_GivenValidBucket(t *testing.T) {
    // 验证成功获取
}
```

### 3. DeleteZipPolicy 测试
```go
func TestDeleteZipPolicy_ShouldReturnSuccess_GivenValidBucket(t *testing.T) {
    // 验证成功删除
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
