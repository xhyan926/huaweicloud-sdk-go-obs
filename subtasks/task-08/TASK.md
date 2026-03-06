# 子任务 2.3：POST 签名计算

## 目标
实现基于 Policy 的签名计算和 Token 生成。

## 范围
- 在 `obs/post_policy.go` 中实现签名计算
- 复用 `obs/auth.go` 的签名逻辑
- 实现 Policy Base64 编码
- 生成完整的 Token 格式（ak:signature:policy）

## 依赖
- 前置子任务：task-07
- 阻塞：task-09

## 实施步骤
1. 复用 `obs/auth.go` 中的签名逻辑
2. 实现 Policy 到签名的转换
3. 生成完整的 Token 格式
4. 在 `obs/client_object.go` 中添加 CreatePostPolicy 客户端方法

## 验收标准
- [ ] 签名计算正确
- [ ] Token 格式符合规范
- [ ] Base64 编码正确
- [ ] 与现有认证逻辑一致

## 状态
pending
