# 子任务 7.1：测试计划

## 测试目标
验证 DIS 策略数据结构的正确性。

## 测试用例

### 1. 结构体完整性测试
```go
func TestSetDisPolicyInput_ShouldHaveRequiredFields_GivenValidInput(t *testing.T) {
    // 验证必需字段存在
}

func TestGetDisPolicyOutput_ShouldParseResponse_GivenValidXML(t *testing.T) {
    // 验证响应解析
}
```

### 2. 常量测试
```go
func TestSubResourceDisPolicy_ShouldHaveCorrectValue(t *testing.T) {
    // 验证常量值正确
}
```

## 测试工具

- testify: 断言库
- encoding/xml: XML 验证

## 验收标准

- [ ] 结构体定义完整
- [ ] 常量定义正确
- [ ] 测试覆盖率 > 80%

## 执行步骤

1. 在相关测试文件中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
