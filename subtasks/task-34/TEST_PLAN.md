# 子任务 9.3：测试计划

## 重要提示
**必须使用 `/go-sdk-ut` skill 编写测试**

## 测试目标
确保 WORM 策略功能的完整性和正确性。

## 测试场景

### 1. 成功场景
- [ ] SetBucketObjectLock 成功设置
- [ ] GetBucketObjectLock 成功获取
- [ ] DeleteBucketObjectLock 成功删除

### 2. 错误场景
- [ ] nil 输入返回错误
- [ ] 空桶名称返回错误

### 3. 边界条件
- [ ] 处理各种响应状态
- [ ] 处理保留天数和年数

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
