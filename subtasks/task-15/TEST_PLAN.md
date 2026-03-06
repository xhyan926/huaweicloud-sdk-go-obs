# 子任务 4.1：测试计划

## 测试目标
验证跨区域复制数据结构的正确性和 XML 序列化功能。

## 测试用例

### 1. 结构体完整性测试
```go
func TestReplicationRule_ShouldHaveRequiredFields_GivenValidInput(t *testing.T) {
    // 验证必需字段存在且类型正确
}

func TestReplicationDestination_ShouldHaveBucket_GivenValidInput(t *testing.T) {
    // 验证目标配置包含桶名称
}
```

### 2. XML 序列化测试
```go
func TestSetBucketReplicationInput_ShouldSerializeToXML_GivenValidConfig(t *testing.T) {
    // 验证配置可以正确序列化为 XML
}

func TestReplicationConfiguration_ShouldIncludeMultipleRules_GivenMultipleRules(t *testing.T) {
    // 验证多个规则正确序列化
}
```

### 3. 字段类型测试
```go
func TestReplicationRule_ShouldAcceptValidStatus_GivenEnabled(t *testing.T) {
    // 验证状态字段接受有效值
}
```

## 测试工具

- testify: 断言库
- encoding/xml: XML 验证

## 验收标准

- [ ] 所有结构体定义通过 go vet 检查
- [ ] XML 序列化输出符合 API 规范
- [ ] 必选字段和可选字段验证通过
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/model_bucket_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
