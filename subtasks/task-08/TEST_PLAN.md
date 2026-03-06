# 子任务 2.3：测试计划

## 测试目标
验证签名计算和 Token 生成的正确性。

## 测试用例

### 1. 签名计算测试
```go
func TestCalculatePostPolicySignature_ShouldCalculateCorrectSignature_GivenValidInput(t *testing.T) {
    // 验证签名计算正确
    // 使用已知输入和预期输出
}

func TestCalculatePostPolicySignature_ShouldReturnError_GivenInvalidInput(t *testing.T) {
    // 验证无效输入返回错误
}
```

### 2. Token 生成测试
```go
func TestBuildPostPolicyToken_ShouldGenerateCorrectFormat_GivenValidComponents(t *testing.T) {
    // 验证 Token 格式: ak:signature:policy
}

func TestBuildPostPolicyToken_ShouldIncludeAllComponents_GivenValidInput(t *testing.T) {
    // 验证所有组件都包含在 Token 中
}
```

### 3. CreatePostPolicy 方法测试
```go
func TestCreatePostPolicy_ShouldReturnPolicy_GivenValidInput(t *testing.T) {
    // 验证成功创建 Policy
}

func TestCreatePostPolicy_ShouldIncludeDefaultConditions_GivenNoConditions(t *testing.T) {
    // 验证自动添加桶和键条件
}

func TestCreatePostPolicy_ShouldReturnError_GivenNilInput(t *testing.T) {
    // 验证 nil 输入返回错误
}

func TestCreatePostPolicy_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
}
```

### 4. 完整流程测试
```go
func TestCreatePostPolicy_ShouldGenerateValidToken_GivenCompleteInput(t *testing.T) {
    // 验证完整的 Token 生成流程
    // 验证 Token 可以被解析和使用
}
```

## 测试工具

- testify: 断言库
- crypto/hmac: HMAC 验证

## 验收标准

- [ ] 签名计算正确
- [ ] Token 格式符合规范
- [ ] Base64 编码正确
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/post_policy_test.go` 和 `obs/client_object_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
