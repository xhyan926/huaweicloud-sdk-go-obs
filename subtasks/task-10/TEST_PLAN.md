# 子任务 2.5：测试计划

## 测试目标
验证示例代码的可运行性和文档的完整性。

## 测试步骤

### 1. 代码编译测试
```bash
cd examples
go build -o post_upload_sample post_upload_sample.go
```
预期：编译成功，无错误

### 2. 环境变量测试
```bash
# 不设置环境变量运行
./post_upload_sample
```
预期：输出错误提示需要设置环境变量

### 3. 完整流程测试
```bash
# 设置环境变量
export OBS_AK="your-access-key"
export OBS_SK="your-secret-key"
export OBS_ENDPOINT="https://obs.cn-north-1.myhuaweicloud.com"

# 运行示例
./post_upload_sample
```
预期：成功生成 Policy 和 HTML 表单

### 4. HTML 输出验证
检查输出的 HTML 是否：
- [ ] 是有效的 HTML5
- [ ] 包含所有必需字段
- [ ] CSS 样式正确
- [ ] 表单 action 正确

### 5. 文档完整性检查
检查示例代码是否有：
- [ ] 详细的文件头注释
- [ ] 函数文档注释
- [ ] 使用说明
- [ ] 错误处理

## 验证标准

- [ ] 示例代码可编译
- [ ] 环境变量检查正常
- [ ] Policy 生成成功
- [ ] HTML 输出有效
- [ ] 文档完整清晰

## 执行步骤

1. 创建示例文件
2. 尝试编译代码
3. 使用测试环境变量运行
4. 验证输出内容
5. 检查 HTML 有效性（可使用在线验证工具）
6. 完善文档和注释
