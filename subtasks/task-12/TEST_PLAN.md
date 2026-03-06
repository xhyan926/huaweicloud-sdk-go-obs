# 子任务 3.2：测试计划

## 测试目标
验证 PUT 上传新增参数的正确映射和向后兼容性。

## 测试用例

### 1. 新参数映射测试
```go
func TestPutObject_ShouldSetExpires_GivenExpiresValue(t *testing.T) {
    // 验证过期时间正确映射到 HTTP 头
}

func TestPutObject_ShouldSetObjectLock_GivenLockParameters(t *testing.T) {
    // 验证 WORM 参数正确映射
}

func TestPutObject_ShouldSetDataEncryption_GivenEncryptionValue(t *testing.T) {
    // 验证数据加密算法正确映射
}
```

### 2. 向后兼容性测试
```go
func TestPutObject_ShouldWork_GivenWithoutNewParameters(t *testing.T) {
    // 验证不设置新参数时仍然可以上传对象
}

func TestPutObject_ShouldHandleExistingParameters_GivenValidInput(t *testing.T) {
    // 验证现有参数功能不受影响
}
```

### 3. 组合参数测试
```go
func TestPutObject_ShouldHandleEncryptionAndLock_GivenCompleteInput(t *testing.T) {
    // 验证加密和 WORM 参数一起使用
}
```

## 测试工具

- testify: 断言库
- MockRoundTripper: HTTP 模拟

## 验收标准

- [ ] 所有新参数正确映射到 HTTP 头部
- [ ] 现有功能不受影响
- [ ] 向后兼容性保持
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/client_object_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
