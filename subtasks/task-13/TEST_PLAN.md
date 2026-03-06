# 子任务 3.3：测试计划

## 测试目标
验证 EncodingType 参数的正确处理和响应编码。

## 测试用例

### 1. 参数映射测试
```go
func TestListObjects_ShouldAddEncodingType_GivenEncodingType(t *testing.T) {
    // 验证 EncodingType 参数添加到查询字符串
}

func TestListObjects_ShouldNotAddEncodingType_GivenEmptyEncodingType(t *testing.T) {
    // 验证空 EncodingType 不添加参数
}
```

### 2. 响应处理测试
```go
func TestListObjects_ShouldHandleURLEncodedResponse_GivenEncodingTypeURL(t *testing.T) {
    // 验证 URL 编码响应正确处理
}

func TestListObjects_ShouldReturnEncodingType_GivenResponse(t *testing.T) {
    // 验证响应中的 EncodingType 字段正确返回
}
```

### 3. 向后兼容性测试
```go
func TestListObjects_ShouldWork_GivenWithoutEncodingType(t *testing.T) {
    // 验证不设置 EncodingType 时正常工作
}

func TestListObjects_ShouldHandleExistingParameters_GivenValidInput(t *testing.T) {
    // 验证现有参数功能不受影响
}
```

## 测试工具

- testify: 断言库
- MockRoundTripper: HTTP 模拟
- net/url: URL 编码验证

## 验收标准

- [ ] EncodingType 参数正确传递到查询字符串
- [ ] 响应正确处理编码
- [ ] 向后兼容性保持
- [ ] 测试覆盖率 > 90%

## 执行步骤

1. 在 `obs/client_object_test.go` 中添加测试用例
2. 运行测试：`go test ./... -v`
3. 检查覆盖率：`go test ./... -cover`
4. 修复发现的问题
