# 子任务 2.4：测试计划

## 重要提示
**必须使用 `/go-sdk-ut` skill 编写测试**

## 测试目标
确保 POST 策略功能的完整性和正确性。

## 测试场景

### 1. Policy 生成场景
- [ ] 生成简单 Policy（仅桶和键条件）
- [ ] 生成复杂 Policy（多个条件）
- [ ] 生成带 content-length-range 的 Policy
- [ ] 生成带 content-type 条件的 Policy

### 2. 签名计算场景
- [ ] 计算简单 Policy 的签名
- [ ] 计算复杂 Policy 的签名
- [ ] 验证签名一致性
- [ ] 处理编码错误

### 3. Token 生成场景
- [ ] 生成完整 Token
- [ ] Token 格式验证（ak:signature:policy）
- [ ] Token 组件完整性

### 4. 客户端方法场景
- [ ] 成功创建 Policy
- [ ] 处理无效输入
- [ ] 处理空桶名称
- [ ] 处理空键名
- [ ] 添加默认条件

### 5. 边界条件
- [ ] 零过期时间
- [ ] 非常长过期时间
- [ ] 空条件列表
- [ ] 单一条件
- [ ] 多个条件

## 测试工具

- **testify**: 断言库
- **httptest**: HTTP 服务器模拟
- **gomonkey**: Mock 工具

## 验收标准

- [ ] 测试覆盖率 > 85%
- [ ] 所有测试通过
- [ ] 符合 BDD 命名规范
- [ ] 已使用 `/go-sdk-ut` skill

## 执行步骤

1. 调用 `/go-sdk-ut` skill
2. 根据指导编写测试用例
3. 运行测试：`go test ./... -v`
4. 检查覆盖率：`go test ./... -coverprofile=coverage.out`
5. 生成覆盖率报告：`go tool cover -html=coverage.out`
6. 修复发现的问题
7. 确保所有测试通过
