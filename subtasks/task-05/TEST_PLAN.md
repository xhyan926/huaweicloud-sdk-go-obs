# 子任务 1.5：测试计划

## 重要提示
**必须使用 `/go-sdk-ut` skill 编写测试**

## 测试目标
确保清单功能的完整性和正确性。

## 测试场景

### 1. 成功场景
- [ ] SetBucketInventory 成功设置清单
- [ ] GetBucketInventory 成功获取清单
- [ ] ListBucketInventory 成功列举清单
- [ ] DeleteBucketInventory 成功删除清单

### 2. 错误场景
- [ ] SetBucketInventory 输入为 nil 返回错误
- [ ] SetBucketInventory 桶名称为空返回错误
- [ ] GetBucketInventory 桶名称为空返回错误
- [ ] GetBucketInventory ID 为空返回错误
- [ ] DeleteBucketInventory 桶名称为空返回错误
- [ ] DeleteBucketInventory ID 为空返回错误

### 3. 边界条件
- [ ] 清单配置包含所有可选字段
- [ ] 清单配置只包含必选字段
- [ ] 清单频率为 Daily
- [ ] 清单频率为 Weekly
- [ ] 多个清单配置的列举

## 测试工具

- **testify**: 断言库
- **httptest**: HTTP 服务器模拟
- **gomonkey**: Mock 工具

## 验收标准

- [ ] 测试覆盖率 > 90%
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
