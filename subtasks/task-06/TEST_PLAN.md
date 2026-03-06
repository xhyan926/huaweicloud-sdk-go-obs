# 子任务 2.1：测试计划

## 测试目标
验证 POST 策略数据结构的正确性和 JSON 序列化功能。

## 测试用例

### 1. 结构体完整性测试
```go
func TestPostPolicy_ShouldHaveRequiredFields_GivenValidInput(t *testing.T) {
    // 验证必需字段存在
}

func TestPostPolicyCondition_ShouldSupportMultipleTypes_GivenDifferentConditions(t *testing.T) {
    // 验证条件支持不同类型
}
```

### 2. JSON 序列化测试
```go
func TestPostPolicy_ShouldSerializeToJSON_GivenValidPolicy(t *testing.T) {
    // 验证 JSON 序列化正确
}

func TestPostPolicy_ShouldIncludeConditions_GivenMultipleConditions(t *testing.T) {
    // 验证条件列表正确序列化
}
```

### 3. 条件类型测试
```go
func TestPostPolicyCondition_ShouldSupportEquals_GivenStringValue(t *testing.T) {
    // 验证 equals 条件
}

func TestPostPolicyCondition_ShouldSupportStartsWith_GivenPrefix(t *testing.T) {
    // 验证 starts-with 条件
}

func TestPostPolicyCondition_ShouldSupportRange_GivenMinMax(t *testing.T) {
    // 验证 content-length-range 条件
}
```

## 测试工具

- testify: 断言库
- encoding/json: JSON 验证

## 验收标准

- [ ] 所有结构体定义通过 go vet 检查
- [ ] JSON 序列化输出符合规范
- [ ] 所有条件类型支持完整
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/model_object_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
