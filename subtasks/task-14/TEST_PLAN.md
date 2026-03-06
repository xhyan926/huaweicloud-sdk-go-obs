# 子任务 3.4：测试计划

## 重要提示
**必须使用 `/go-sdk-ut` skill 编写测试**

## 测试目标
确保所有新增参数的正确性和向后兼容性。

## 测试场景

### 1. CreateBucket 参数场景
- [ ] 设置 BucketType 参数
- [ ] 设置 KMS 密钥 ID 参数
- [ ] 设置 KMS 密钥项目 ID 参数
- [ ] 设置数据加密算法参数
- [ ] 组合使用多个新参数
- [ ] 不设置新参数（向后兼容）

### 2. PutObject 参数场景
- [ ] 设置 Expires 参数
- [ ] 设置对象 WORM 模式参数
- [ ] 设置 WORM 保留时间参数
- [ ] 设置数据加密算法参数
- [ ] 设置 KMS 密钥 ID 参数
- [ ] 组合使用 WORM 和加密参数
- [ ] 不设置新参数（向后兼容）

### 3. ListObjects 参数场景
- [ ] 设置 EncodingType 为 "url"
- [ ] 设置 EncodingType 为空
- [ ] 验证查询字符串正确
- [ ] 验证响应处理正确
- [ ] 不设置 EncodingType（向后兼容）

### 4. 向后兼容性验证
- [ ] 运行所有现有测试套件
- [ ] 验证没有测试失败
- [ ] 检查测试覆盖率
- [ ] 确保功能不受影响

## 测试工具

- **testify**: 断言库
- **httptest**: HTTP 服务器模拟
- **gomonkey**: Mock 工具

## 验收标准

- [ ] 所有新测试通过
- [ ] 现有测试不受影响
- [ ] 测试覆盖率 > 90%
- [ ] 已使用 `/go-sdk-ut` skill

## 执行步骤

1. 调用 `/go-sdk-ut` skill
2. 根据指导编写测试用例
3. 运行测试：`go test ./... -v`
4. 检查覆盖率：`go test ./... -coverprofile=coverage.out`
5. 生成覆盖率报告：`go tool cover -html=coverage.out`
6. 运行完整测试套件验证向后兼容性
7. 修复发现的问题
8. 确保所有测试通过
