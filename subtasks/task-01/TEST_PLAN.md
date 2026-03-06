# 子任务 1.1：测试计划

## 测试目标
验证所有清单配置数据结构的正确性和 XML 序列化功能。

## 测试用例

### 1. 结构体完整性测试
```go
func TestInventoryConfiguration_ShouldHaveRequiredFields_GivenValidInput(t *testing.T) {
    // 验证必需字段存在且类型正确
}

func TestInventoryConfiguration_ShouldAllowOptionalFields_GivenFilter(t *testing.T) {
    // 验证可选字段可以设置
}
```

### 2. XML 序列化测试
```go
func TestInventoryConfiguration_ShouldSerializeToXML_GivenValidConfig(t *testing.T) {
    // 验证配置可以正确序列化为 XML
}

func TestInventoryConfiguration_ShouldIncludeNestedStructures_GivenCompleteConfig(t *testing.T) {
    // 验证嵌套结构正确序列化
}
```

### 3. 字段类型测试
```go
func TestInventoryDestination_ShouldAcceptValidFormat_GivenCSV(t *testing.T) {
    // 验证格式字段接受有效值
}

func TestInventorySchedule_ShouldAcceptValidFrequency_GivenDaily(t *testing.T) {
    // 验证频率字段接受有效值
}
```

## 测试工具

- testify: 断言库
- encoding/xml: XML 序列化验证

## 验收标准

- [ ] 所有结构体定义通过 go vet 检查
- [ ] XML 序列化输出符合 API 规范
- [ ] 必选字段和可选字段验证通过
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 创建 `obs/model_bucket_test.go`（如果不存在）
2. 添加上述测试用例
3. 运行测试：`go test ./... -v`
4. 检查覆盖率：`go test ./... -cover`
5. 修复发现的问题
