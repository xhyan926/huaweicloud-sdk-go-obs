# 子任务 2.2：测试计划

## 测试目标
验证 Policy 构建和验证逻辑的正确性。

## 测试用例

### 1. Policy 构建测试
```go
func TestBuildPostPolicyJSON_ShouldGenerateValidJSON_GivenValidPolicy(t *testing.T) {
    // 验证生成有效的 JSON
}

func TestBuildPostPolicyJSON_ShouldIncludeAllConditions_GivenMultipleConditions(t *testing.T) {
    // 验证所有条件都包含在 JSON 中
}
```

### 2. 过期时间测试
```go
func TestBuildPostPolicyExpiration_ShouldGenerateCorrectFormat_GivenSeconds(t *testing.T) {
    // 验证过期时间格式正确
}

func TestBuildPostPolicyExpiration_ShouldUseUTC_GivenLocalTime(t *testing.T) {
    // 验证使用 UTC 时区
}
```

### 3. 验证逻辑测试
```go
func TestValidatePostPolicy_ShouldReturnNil_GivenValidPolicy(t *testing.T) {
    // 验证有效策略通过验证
}

func TestValidatePostPolicy_ShouldReturnError_GivenNilPolicy(t *testing.T) {
    // 验证 nil 策略返回错误
}

func TestValidatePostPolicy_ShouldReturnError_GivenEmptyExpiration(t *testing.T) {
    // 验证空过期时间返回错误
}

func TestValidatePostPolicy_ShouldReturnError_GivenEmptyConditions(t *testing.T) {
    // 验证空条件列表返回错误
}
```

### 4. 条件创建测试
```go
func TestCreateBucketCondition_ShouldCreateCorrectCondition_GivenBucket(t *testing.T) {
    // 验证桶条件创建正确
}

func TestCreateKeyCondition_ShouldCreateCorrectCondition_GivenKey(t *testing.T) {
    // 验证键条件创建正确
}

func TestCreateContentLengthRangeCondition_ShouldCreateCorrectCondition_GivenRange(t *testing.T) {
    // 验证范围条件创建正确
}
```

## 测试工具

- testify: 断言库
- encoding/json: JSON 验证

## 验收标准

- [ ] JSON 生成符合 AWS S3 POST 规范
- [ ] 验证逻辑能检测无效策略
- [ ] 过期时间处理正确
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/post_policy_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
