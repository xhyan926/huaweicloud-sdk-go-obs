# 子任务 6.2：归档直读实现

## 目标
实现 Set/Get/DeleteDirectcoldaccess 方法。

## 范围
- Trait 层的 trans() 实现
- 客户端方法实现
- 参数验证

## 依赖
- 前置子任务：task-23
- 阻塞：task-25

## 实施步骤
1. 在 `obs/trait_bucket.go` 和 `obs/client_bucket.go` 中添加方法
2. 实现参数验证
3. 确保符合现有 API 风格

## 验收标准
- [ ] 方法符合现有 API 风格
- [ ] 错误处理一致

## 状态
pending
