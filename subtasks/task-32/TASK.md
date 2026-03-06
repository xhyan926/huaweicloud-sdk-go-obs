# 子任务 9.1：WORM 策略数据模型和常量

## 目标
创建桶级 WORM 策略的数据结构和常量。

## 范围
- 在 `obs/model_bucket.go` 中添加结构体
- 在 `obs/type.go` 中添加常量

## 依赖
- 前置子任务：无
- 阻塞：task-33

## 实施步骤
1. 添加 Set/GetBucketObjectLock 的 Input/Output 结构体
2. 添加 `SubResourceObjectLock` 常量
3. 添加 XML 标签映射

## 验收标准
- [ ] 结构体和常量定义完整
- [ ] 代码通过 go vet 检查

## 状态
pending
