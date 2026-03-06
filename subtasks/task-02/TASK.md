# 子任务 1.2：桶清单常量和类型定义

## 目标
添加清单功能相关的常量和类型定义。

## 范围
- 在 `obs/type.go` 中添加 `SubResourceInventory` 常量
- 定义 Frequency 类型枚举
- 定义 OptionalFields 类型

## 依赖
- 前置子任务：task-01
- 阻塞：task-03

## 实施步骤
1. 在 `obs/type.go` 中添加 `SubResourceInventory` 常量
2. 定义 Frequency 类型（Daily/Weekly）
3. 定义 OptionalFields 类型
4. 确保常量命名符合现有规范

## 验收标准
- [ ] 常量命名符合现有规范
- [ ] 类型定义完整覆盖 API 需求
- [ ] 代码通过 go vet 检查

## 状态
pending
