# 子任务 4.2：跨区域复制常量和类型

## 目标
添加跨区域复制功能相关的常量和类型定义。

## 范围
- 在 `obs/type.go` 中添加 `SubResourceReplication` 常量
- 定义复制状态枚举
- 定义存储类类型枚举

## 依赖
- 前置子任务：task-15
- 阻塞：task-17

## 实施步骤
1. 在 `obs/type.go` 中添加 `SubResourceReplication` 常量
2. 定义 ReplicationStatus 类型（Enabled/Disabled）
3. 定义相关枚举类型
4. 确保常量命名符合现有规范

## 验收标准
- [ ] 常量命名符合现有规范
- [ ] 类型定义完整覆盖 API 需求
- [ ] 代码通过 go vet 检查

## 状态
pending
