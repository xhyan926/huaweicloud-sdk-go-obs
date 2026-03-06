# 子任务 1.2：测试计划

## 测试目标
验证清单相关的常量和类型定义正确性。

## 测试用例

### 1. 常量值测试
```go
func TestSubResourceInventory_ShouldHaveCorrectValue(t *testing.T) {
    // 验证常量值为 "inventory"
}
```

### 2. 频率类型测试
```go
func TestInventoryFrequencyDaily_ShouldHaveCorrectValue(t *testing.T) {
    // 验证 Daily 频率值正确
}

func TestInventoryFrequencyWeekly_ShouldHaveCorrectValue(t *testing.T) {
    // 验证 Weekly 频率值正确
}
```

## 测试工具

- testify: 断言库

## 验收标准

- [ ] 所有常量值正确
- [ ] 所有类型定义完整
- [ ] 代码通过 go vet 检查
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/type_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 修复发现的问题
