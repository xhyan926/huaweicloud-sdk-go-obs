# 子任务 8.3 测试计划：在线解压策略单元测试

## 测试目标
验证在线解压策略功能的正确性和稳定性，确保代码质量达标。

## 测试用例（BDD 风格）

### 1. 数据模型测试
- **测试名称**：Should_validate_decompression_configuration_structure
- **测试内容**：
  - 验证 DecompressionConfiguration 结构字段
  - 测试 XML 标签正确性
  - 验证必填字段
- **预期结果**：结构定义正确

### 2. 序列化测试
- **测试名称**：Should_serialize_deconfiguration_to_xml
- **测试内容**：
  - 序列化配置到 XML
  - 验证 XML 格式
  - 测试特殊字符处理
- **预期结果**：序列化正确

- **测试名称**：Should_deserialize_xml_to_configuration
- **测试内容**：
  - 从 XML 反序列化配置
  - 验证数据完整性
  - 测试错误 XML 处理
- **预期结果**：反序列化正确

### 3. Trait 层测试
- **测试名称**：Should_set_bucket_decompression_policy
- **测试内容**：
  - 设置桶解压策略
  - 验证参数传递
  - 检查返回值
- **预期结果**：设置成功

- **测试名称**：Should_get_bucket_decompression_policy
- **测试内容**：
  - 获取桶解压策略
  - 验证返回数据
  - 处理空策略情况
- **预期结果**：获取正确

- **测试名称**：Should_delete_bucket_decompression_policy
- **测试内容**：
  - 删除桶解压策略
  - 验证删除操作
  - 检查错误处理
- **预期结果**：删除成功

### 4. 客户端方法测试
- **测试名称**：Should_call_client_methods_with_correct_parameters
- **测试内容**：
  - 客户端方法调用
  - 参数验证
  - 选项处理
- **预期结果**：调用正常

### 5. 错误处理测试
- **测试名称**：Should_return_error_when_invalid_bucket_name
- **测试内容**：
  - 测试无效桶名
  - 验证错误消息
  - 检查错误码
- **预期结果**：返回正确错误

- **测试名称**：Should_return_error_when_invalid_configuration
- **测试内容**：
  - 测试无效配置
  - 验证验证逻辑
  - 检查错误类型
- **预期结果**：返回正确错误

## 测试执行
```bash
# 运行单元测试
go test -tags unit ./obs -v -run "TestDecompression"

# 运行内部测试
go test -tags internal ./obs -v -run "TestDecompression"

# 检查测试覆盖率
go test -tags unit ./obs -cover -run "TestDecompression"
```

## 测试覆盖率要求
- 数据模型：95%
- 序列化：95%
- 业务逻辑：90%
- 客户端接口：85%
- 错误处理：85%
- 总体覆盖率：>85%
