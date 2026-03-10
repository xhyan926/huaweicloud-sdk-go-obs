# 子任务 9.3 测试计划：WORM 策略单元测试

## 测试目标
验证 WORM 策略功能的正确性和稳定性，确保代码质量达标。

## 测试用例（BDD 风格）

### 1. 数据模型测试
- **测试名称**：Should_validate_worm_configuration_structure
- **测试内容**：
  - 验证 BucketWormConfiguration 结构字段
  - 测试 XML 标签正确性
  - 验证必填字段
  - 测试状态枚举
- **预期结果**：结构定义正确

### 2. 序列化测试
- **测试名称**：Should_serialize_worm_configuration_to_xml
- **测试内容**：
  - 序列化配置到 XML
  - 验证 XML 格式
  - 测试版本信息
  - 测试特殊字符处理
- **预期结果**：序列化正确

- **测试名称**：Should_deserialize_xml_to_worm_configuration
- **测试内容**：
  - 从 XML 反序列化配置
  - 验证数据完整性
  - 测试版本处理
  - 测试错误 XML 处理
- **预期结果**：反序列化正确

### 3. Trait 层测试
- **测试名称**：Should_set_bucket_worm_policy
- **测试内容**：
  - 设置桶 WORM 策略
  - 验证参数传递
  - 检查返回值
  - 测试状态初始化
- **预期结果**：设置成功

- **测试名称**：Should_get_bucket_worm_policy
- **测试内容**：
  - 获取桶 WORM 策略
  - 验证返回数据
  - 处理空策略情况
  - 测试版本信息
- **预期结果**：获取正确

- **测试名称**：Should_extend_bucket_worm_policy
- **测试内容**：
  - 延长桶 WORM 策略期限
  - 验证延长逻辑
  - 保持状态不变
  - 更新版本信息
- **预期结果**：延长成功

### 4. WORM 约束测试
- **测试名称**：Should_enforce_compliance_mode_constraints
- **测试内容**：
  - 测试 COMPLIANCE 状态下的不可修改性
  - 验证参数锁定
  - 测试错误消息
- **预期结果**：约束生效

- **测试名称**：Should_handle_suspended_mode_correctly
- **测试内容**：
  - 测试 SUSPENDED 状态的特殊处理
  - 验证状态转换
  - 测试修改权限
- **预期结果**：状态处理正确

### 5. 客户端方法测试
- **测试名称**：Should_call_worm_client_methods_with_correct_parameters
- **测试内容**：
  - 客户端方法调用
  - 参数验证
  - 选项处理
  - 返回值检查
- **预期结果**：调用正常

### 6. 错误处理测试
- **测试名称**：Should_return_error_when_invalid_worm_configuration
- **测试内容**：
  - 测试无效配置
  - 验证错误消息
  - 检查错误码
- **预期结果**：返回正确错误

- **测试名称**：Should_return_error_when_violate_worm_constraints
- **测试内容**：
  - 测试违反约束的操作
  - 验证业务规则检查
  - 检查特定错误码
- **预期结果**：返回正确错误

## 测试执行
```bash
# 运行单元测试
go test -tags unit ./obs -v -run "TestWorm"

# 运行内部测试
go test -tags internal ./obs -v -run "TestWorm"

# 检查测试覆盖率
go test -tags unit ./obs -cover -run "TestWorm"
```

## 测试覆盖率要求
- 数据模型：95%
- 序列化：95%
- 业务逻辑：90%
- 客户端接口：85%
- WORM 约束处理：90%
- 错误处理：85%
- 总体覆盖率：>85%
