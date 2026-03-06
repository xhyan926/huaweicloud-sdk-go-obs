# 子任务 3.3：列举对象参数补充

## 目标
补充 ListObjects 的缺失参数。

## 范围
- 在 `ListObjectsInput` 中添加 EncodingType 字段
- 在 `trait_object.go` 的 `trans()` 方法中添加参数映射
- 处理 URL 编码响应

## 依赖
- 前置子任务：无
- 阻塞：task-14

## 实施步骤
1. 在 `ListObjectsInput` 中添加 EncodingType 字段
2. 在 `trait_object.go` 的 `trans()` 方法中添加参数映射逻辑
3. 更新 `ListObjectsOutput` 以处理编码响应
4. 确保与现有功能兼容

## 验收标准
- [ ] EncodingType 参数正确传递到查询字符串
- [ ] 响应正确处理编码
- [ ] 向后兼容
- [ ] 代码通过 go vet 检查

## 状态
pending
