# 子任务 3.1：测试计划

## 测试目标
验证新增参数的正确映射和向后兼容性。

## 测试用例

### 1. 新参数映射测试
```go
func TestCreateBucket_ShouldSetBucketType_GivenBucketType(t *testing.T) {
    // 验证桶类型正确映射到 HTTP 头
}

func TestCreateBucket_ShouldSetKmsKeyId_GivenKmsKeyId(t *testing.T) {
    // 验证 KMS 密钥 ID 正确映射
}

func TestCreateBucket_ShouldSetDataEncryption_GivenDataEncryption(t *testing.T) {
    // 验证数据加密算法正确映射
}
```

### 2. 向后兼容性测试
```go
func TestCreateBucket_ShouldWork_GivenWithoutNewParameters(t *testing.T) {
    // 验证不设置新参数时仍然可以创建桶
}

func TestCreateBucket_ShouldSendDefaultHeaders_GivenValidInput(t *testing.T) {
    // 验证默认行为不受影响
}
```

### 3. 组合参数测试
```go
func TestCreateBucket_ShouldHandleAllParameters_GivenCompleteInput(t *testing.T) {
    // 验证所有参数一起使用时正常工作
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

1. 在 `obs/client_bucket_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
