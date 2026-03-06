# 子任务 1.3：测试计划

## 测试目标
验证清单功能的请求转换和参数映射正确性。

## 测试用例

### 1. trans() 方法测试
```go
func TestSetBucketInventoryInput_ShouldReturnCorrectParams_GivenValidInput(t *testing.T) {
    // 验证参数映射正确
}

func TestSetBucketInventoryInput_ShouldSerializeXML_GivenCompleteConfig(t *testing.T) {
    // 验证 XML 序列化正确
}
```

### 2. 删除操作测试
```go
func TestDeleteBucketInventoryInput_ShouldIncludeId_GivenInventoryId(t *testing.T) {
    // 验证删除操作包含 ID 参数
}
```

### 3. 参数验证测试
```go
func TestSetBucketInventoryInput_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
    // 验证空桶名称返回错误
}
```

## 测试工具

- testify: 断言库
- encoding/xml: XML 验证

## 验收标准

- [ ] trans() 方法返回正确的参数映射
- [ ] XML 序列化输出符合 API 规范
- [ ] 参数验证逻辑完整
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/trait_bucket_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
