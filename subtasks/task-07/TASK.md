# 子任务 2.2：POST 策略构建和验证

## 目标
实现 Policy 的构建、验证和 JSON 生成逻辑。

## 范围
- 创建 `obs/post_policy.go` 新文件
- 实现 Policy JSON 生成
- 实现策略验证
- 实现过期时间处理

## 依赖
- 前置子任务：task-06
- 阻塞：task-08

## 实施步骤
1. 创建 `obs/post_policy.go` 新文件
2. 实现 Policy JSON 生成逻辑
3. 添加策略验证方法
4. 实现过期时间计算和格式化
5. 确保生成的 JSON 符合 AWS S3 POST 规范

## 验收标准
- [ ] 生成的 JSON 符合 AWS S3 POST 规范
- [ ] 验证逻辑能检测无效策略
- [ ] 过期时间处理正确
- [ ] 代码通过 go vet 检查

## 状态
pending
