# 子任务 8.3：在线解压单元测试

## 目标
为在线解压功能编写完整的单元测试。

## 范围
- 单元测试
- Mock 测试
- 集成测试

## 依赖
- 前置子任务：task-30
- 必须使用 `/go-sdk-ut` skill 编写测试

## 实施步骤
1. 使用 `/go-sdk-ut` skill 编写测试
2. 在测试文件中添加测试用例
3. 使用 MockRoundTripper 模拟 HTTP 响应
4. 编写 BDD 风格的测试用例

## 验收标准
- [ ] 测试覆盖率 > 80%
- [ ] 所有测试通过
- [ ] 符合 BDD 命名规范
- [ ] 使用 testify、httptest、gomonkey

## 状态
pending
