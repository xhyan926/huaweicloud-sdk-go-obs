# 子任务 2.4：POST 策略单元测试

## 目标
为 POST 策略功能编写完整的单元测试。

## 范围
- Policy 生成测试
- 签名计算测试
- 验证逻辑测试
- 集成测试
- 边界条件测试

## 依赖
- 前置子任务：task-08
- 必须使用 `/go-sdk-ut` skill 编写测试

## 实施步骤
1. 使用 `/go-sdk-ut` skill 编写测试
2. 在 `obs/post_policy_test.go` 中添加测试用例
3. 测试各种策略条件
4. 测试签名计算准确性
5. 测试边界条件

## 验收标准
- [ ] 测试覆盖率 > 85%
- [ ] 所有测试通过
- [ ] 符合 BDD 命名规范
- [ ] 使用 testify、httptest、gomonkey

## 状态
pending
