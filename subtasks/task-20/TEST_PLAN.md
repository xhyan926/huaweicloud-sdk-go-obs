# 子任务 5.2：测试计划

## 测试目标
验证存量信息相关的常量定义正确性。

## 测试用例

### 1. 常量值测试
```go
func TestSubResourceStorageInfo_ShouldHaveCorrectValue(t *testing.T) {
    // 验证常量值为 "storageinfo"
}
```

## 测试工具

- testify: 断言库

## 验收标准

- [ ] 常量值正确
- [ ] 代码通过 go vet 检查
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/type_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 修复发现的问题
