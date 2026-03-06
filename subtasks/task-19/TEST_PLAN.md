# 子任务 5.1：测试计划

## 测试目标
验证存量信息数据结构的正确性和 XML 解析功能。

## 测试用例

### 1. 结构体完整性测试
```go
func TestStorageInfo_ShouldHaveRequiredFields_GivenValidInput(t *testing.T) {
    // 验证必需字段存在且类型正确
}
```

### 2. XML 解析测试
```go
func TestGetBucketStorageInfoOutput_ShouldParseXML_GivenValidResponse(t *testing.T) {
    // 验证可以正确解析 XML 响应
}

func TestGetBucketStorageInfoOutput_ShouldHandleLargeNumbers_GivenBigStorage(t *testing.T) {
    // 验证大数值正确解析
}
```

### 3. 字段类型测试
```go
func TestStorageInfo_ShouldAcceptInt64Values_GivenLargeStorage(t *testing.T) {
    // 验证 int64 类型支持大数值
}
```

## 测试工具

- testify: 断言库
- encoding/xml: XML 验证

## 验收标准

- [ ] 所有结构体定义通过 go vet 检查
- [ ] XML 解析输出符合 API 规范
- [ ] 大数值处理正确
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/model_bucket_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
