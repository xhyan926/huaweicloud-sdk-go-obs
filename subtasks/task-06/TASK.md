# 子任务 2.1：POST 策略数据模型定义

## 目标
创建 POST 上传策略的数据结构。

## 范围
- 在 `obs/model_object.go` 中添加策略结构体
- 定义 CreatePostPolicyInput 和 Output
- 定义 PostPolicyCondition 结构体
- 定义 PostPolicyRule 结构体
- 支持 JSON 序列化

## 依赖
- 前置子任务：无
- 阻塞：task-07

## 实施步骤
1. 在 `obs/model_object.go` 中添加 POST 策略相关结构体
2. 定义 Input/Output 结构体
3. 定义条件和规则结构体
4. 添加 JSON 标签映射
5. 确保结构体支持所有常见策略条件

## 验收标准
- [ ] 所有结构体通过 `go vet` 检查
- [ ] JSON 序列化正确
- [ ] 支持所有常见策略条件
- [ ] 字段类型定义正确

## 状态
pending
