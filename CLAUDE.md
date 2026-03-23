# 华为云 OBS Go SDK 开发规范

本文档定义了华为云 OBS Go SDK 开发的项目特定规范，供 Claude Code 参考使用。

## 快速参考

- **API 开发指南**: `docs/API_DEVELOPMENT_GUIDE.md` - 新 API 功能开发的完整规范
- **测试规范**: 见下方测试部分
- **代码审查规范**: 见下方审查检查项

---

## 一、新 API 功能开发流程

当用户请求开发新 API 功能时，必须遵循以下流程：

### 1.1 开发前检查（必须）

在编写任何代码之前，先完成：

```bash
# 1. 查阅官方 API 文档
# 2. 创建 API 规范记录文件（临时）
# 3. 确认以下信息：
#    - 数据格式：XML 或 JSON
#    - Content-Type：application/xml 或 application/json
#    - 子资源名称：query 参数（如 ?replication、?disPolicy）
#    - 所有字段类型和必需性
```

### 1.2 文件修改顺序

按以下顺序修改文件（避免遗漏）：

| 顺序 | 文件 | 作用 | 必须修改 |
|------|------|------|----------|
| 1 | `obs/type.go` | 子资源常量定义 | ✅ 是 |
| 2 | `obs/const.go` | 允许的资源参数列表 | ✅ 是（签名计算必需） |
| 3 | `obs/model_base.go` | 数据结构定义 | ✅ 是 |
| 4 | `obs/model_bucket.go` | Input/Output 模型 | ✅ 是 |
| 5 | `obs/convert.go` | 序列化函数（仅 JSON API） | JSON API 需要 |
| 6 | `obs/trait_bucket.go` | trans() 参数转换方法 | ✅ 是 |
| 7 | `obs/client_bucket.go` | API 方法实现 | ✅ 是 |

### 1.3 关键规则（避免常见错误）

#### 规则 1: 子资源必须添加到 allowedResourceParameterNames

```go
// obs/const.go
allowedResourceParameterNames = map[string]bool{
    // ... 其他参数 ...
    "replication": true,  // ⚠️ 使用小写！
    "dispolicy":   true,  // ⚠️ 不是 disPolicy
}
```

**后果**：如果遗漏，签名计算会跳过该参数，导致 API 调用失败。

#### 规则 2: JSON API 必须序列化，不能直接赋值结构体

```go
// ❌ 错误实现
func (input SetBucketDisPolicyInput) trans(isObs bool) (...) {
    data = input.DisPolicyConfiguration  // 错误！直接赋值结构体
    return
}

// ✅ 正确实现
func (input SetBucketDisPolicyInput) trans(isObs bool) (...) {
    data, err = convertDisPolicyToJSON(input.DisPolicyConfiguration)  // 正确！序列化
    return
}
```

#### 规则 3: 数据结构必须直接映射 API，不要过度抽象

```go
// ❌ 错误：创建不必要的中间结构
type DisEvent struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled"`  // API 中没有这个字段！
}

// ✅ 正确：直接映射 API 结构
type DisPolicyRule struct {
    Events []string `json:"events"`  // 直接字符串数组
}
```

### 1.4 trans() 方法模板

**XML API 模板**：
```go
func (input SetBucketXXXInput) trans(isObs bool) (...) {
    params = map[string]string{string(SubResourceXXX): ""}
    reader, md5, _ := ConvertRequestToIoReaderV2(input.XXXConfiguration, false)
    data = reader
    headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
    return
}
```

**JSON API 模板**：
```go
func (input SetBucketXXXInput) trans(isObs bool) (...) {
    params = map[string]string{string(SubResourceXXX): ""}
    headers = make(map[string][]string, 1)
    headers[HEADER_CONTENT_TYPE] = []string{mimeTypes["json"]}
    data, err = convertXXXToJSON(input.XXXConfiguration)  // ⚠️ 必须序列化
    return
}
```

---

## 二、测试规范

### 2.1 单元测试

**文件位置**: 与业务代码同目录，命名为 `*_test.go`

**命名规范** (BDD 风格):
```go
// 格式：Test{功能}_Should{期望}_When{条件}
func TestSetBucketReplication_ShouldReturnSuccess_WhenValidInput() { }
func TestSetBucketReplication_ShouldHandleError_WhenBucketNotFound() { }
```

**测试工具**:
- `github.com/stretchr/testify/assert` - 断言库
- `github.com/stretchr/testify/require` - 必须通过的断言

**最小覆盖**:
- [ ] 正常流程测试
- [ ] 参数验证测试
- [ ] 序列化/反序列化测试

### 2.2 集成测试

**文件位置**: `obs/test/integration/`

**配置方式**:
1. 认证信息（ak, sk, endpoint, bucket）必须通过环境变量设置
2. 特性配置（DIS stream、复制目标等）可通过配置文件管理

**环境变量**:
```bash
export OBS_TEST_AK=your_access_key
export OBS_TEST_SK=your_secret_key
export OBS_TEST_ENDPOINT=your_endpoint
export OBS_TEST_BUCKET=your_test_bucket
```

**配置文件** (可选): `obs/test/integration/test.config.json`
```json
{
  "dis": {
    "stream": "your-dis-stream",
    "project": "your-project-id",
    "agency": "your-iam-agency"
  },
  "replication": {
    "destBucket": "your-dest-bucket",
    "location": "cn-south-1"
  }
}
```

**运行集成测试**:
```bash
export OBS_TEST_AK=xxx
export OBS_TEST_SK=xxx
export OBS_TEST_ENDPOINT=xxx
export OBS_TEST_BUCKET=xxx
go test ./obs/test/integration/ -v -tags=integration
```

### 2.3 仅支持 OBS 签名的功能

某些功能仅支持 OBS 签名，不支持 AWS 签名。集成测试中应验证：

```go
func TestFeature_ShouldSupportOnlyOBSSignature(t *testing.T) {
    // 使用 OBS 签名应该成功
    client := createClient(t, obs.SignatureObs)
    // ... 验证成功 ...

    // 使用 AWS 签名应该失败或不支持
    // client := createClient(t, obs.SignatureV4)
    // ... 验证失败 ...
}
```

---

## 三、代码审查检查项

在提交代码或审查代码时，使用此检查清单：

### 3.1 API 实现检查

- [ ] `obs/type.go` - 子资源常量已定义
- [ ] `obs/const.go` - 子资源已添加到 `allowedResourceParameterNames`（使用小写）
- [ ] `obs/model_base.go` - 数据结构定义正确，与 API 规范一致
- [ ] `obs/model_bucket.go` - Input/Output 已定义
- [ ] `obs/convert.go` - JSON 序列化函数已添加（如适用）
- [ ] `obs/trait_bucket.go` - trans() 方法正确实现
  - [ ] 子资源参数已设置
  - [ ] JSON API: Content-Type header 已设置
  - [ ] JSON API: 使用 convertXXXToJSON() 序列化（不是直接赋值）
  - [ ] XML API: 使用 ConvertRequestToIoReaderV2() 序列化
- [ ] `obs/client_bucket.go` - API 方法已实现

### 3.2 常见错误检查

- [ ] trans() 方法中 `data` 不是直接赋值结构体（JSON API）
- [ ] 数据结构没有过度抽象（没有不必要的中间结构）
- [ ] 字段类型与 API 规范一致（字符串数组不是对象数组）
- [ ] 子资源名称使用小写（`dispolicy` 不是 `disPolicy`）

### 3.3 测试检查

- [ ] 单元测试已编写并通过
- [ ] 集成测试已编写（使用正确的 build tag）
- [ ] 配置文件不包含敏感信息（ak, sk）
- [ ] 敏感信息通过环境变量获取

---

## 四、文档规范

### 4.1 API 文档位置

- **总索引**: `docs/README.md`
- **功能文档**: `docs/{feature}/README.md`

### 4.2 文档内容

每个功能的文档应包含：
- [ ] 功能概述
- [ ] API 方法列表
- [ ] 使用示例
- [ ] 数据结构说明
- [ ] 限制和注意事项
- [ ] 签名支持说明

---

## 五、常见问题

### Q: API 调用返回签名错误？
**A**: 检查 `obs/const.go` 的 `allowedResourceParameterNames` 是否包含子资源名称（使用小写）。

### Q: JSON API 返回请求体格式错误？
**A**: 检查 trans() 方法是否使用 `convertXXXToJSON()` 序列化，而不是直接赋值结构体。

### Q: 数据结构反序列化失败？
**A**: 检查结构体定义是否与 API 规范一致，特别注意：
- 字符串数组 vs 对象数组
- 字段名称大小写
- 嵌套结构层级

---

## 六、相关文档

- [API 开发指南](docs/API_DEVELOPMENT_GUIDE.md) - 完整的 API 开发规范
- [集成测试说明](obs/test/integration/README.md) - 集成测试配置和运行
- [华为云 OBS 官方文档](https://support.huaweicloud.com/obs/)

---

## 版本信息

- 文档版本: 1.0
- 创建日期: 2026-03-23
- SDK 版本: 3.26.0+
