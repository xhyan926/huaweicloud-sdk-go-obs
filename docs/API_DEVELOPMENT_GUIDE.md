# OBS SDK Go API 开发规范

本文档定义了华为云 OBS Go SDK 新功能开发的规范和检查项，用于避免常见错误。

## 目录

- [开发前准备](#一开发前准备)
- [数据结构设计](#二数据结构设计)
- [trans() 方法实现](#三trans-方法实现)
- [签名计算](#四签名计算)
- [完整开发流程](#五完整开发流程)
- [快速检查清单](#六快速检查清单)

---

## 一、开发前准备

### 1.1 官方文档验证（强制）

在实现任何新 API 之前，必须：

- [ ] **阅读官方 API 文档**，记录请求/响应格式
- [ ] **确认数据格式**：XML 或 JSON
- [ ] **记录请求示例**：完整的请求体示例
- [ ] **记录响应示例**：完整的响应体示例
- [ ] **确认字段类型**：字符串/数组/对象/嵌套结构
- [ ] **确认必需/可选字段**：标注哪些字段是必填的
- [ ] **确认 Content-Type**：`application/xml` 或 `application/json`
- [ ] **确认子资源名称**：query 参数名称（如 `?replication`、`?disPolicy`）

### 1.2 API 实现文件映射

| 文件 | 作用 | 必须修改 |
|------|------|----------|
| `obs/type.go` | 子资源常量定义 | ✅ 是 |
| `obs/const.go` | 允许的资源参数列表 | ✅ 是 |
| `obs/model_base.go` | 数据结构定义 | ✅ 是 |
| `obs/model_bucket.go` | Input/Output 模型 | ✅ 是 |
| `obs/convert.go` | 序列化/反序列化函数 | JSON API 需要 |
| `obs/trait_bucket.go` | trans() 参数转换方法 | ✅ 是 |
| `obs/client_bucket.go` | API 方法实现 | ✅ 是 |

---

## 二、数据结构设计

### 2.1 XML API 结构设计

```go
// 示例：跨区域复制（XML）
type ReplicationRule struct {
    ID     string                  `xml:"id"`
    Status string                  `xml:"status"`  // Enabled/Disabled
    Prefix *ReplicationPrefix       `xml:"prefix"`
    Destination ReplicationDestination `xml:"destination"`
}

type ReplicationConfiguration struct {
    XMLName xml.Name               `xml:"ReplicationConfiguration"`
    Rules   []ReplicationRule      `xml:"Rule"`
}
```

**XML 检查项**：
- [ ] 结构体字段使用 `xml` 标签
- [ ] 数组标签使用单数形式（如 `xml:"Rule"`）
- [ ] 可选字段使用指针类型
- [ ] 使用 `ConvertRequestToIoReaderV2` 进行序列化

### 2.2 JSON API 结构设计

```go
// 示例：DIS 事件通知（JSON）
type DisPolicyRule struct {
    ID      string   `json:"id"`      // 规则 ID
    Stream  string   `json:"stream"`  // DIS 通道名称
    Project string   `json:"project"` // 项目 ID
    Events  []string `json:"events"`  // 事件列表（字符串数组）
    Prefix  string   `json:"prefix,omitempty"`
    Suffix  string   `json:"suffix,omitempty"`
    Agency  string   `json:"agency"`  // IAM 委托名
}

type DisPolicyConfiguration struct {
    Rules []DisPolicyRule `json:"rules"` // 规则数组
}
```

**JSON 检查项**：
- [ ] 结构体字段使用 `json` 标签
- [ ] 可选字段添加 `omitempty`
- [ ] **数组类型直接用 `[]string` 或 `[]struct`，不要创建中间包装结构**
- [ ] 嵌套结构体直接定义，不要过度抽象

### 2.3 常见错误模式

❌ **错误：过度抽象的中间结构**
```go
type DisEvent struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled"`  // API 中没有此字段！
}

type DisPolicyConfiguration struct {
    AgencyName string     `json:"agency_name"`  // ❌ 应该在每条规则内
    Events     []DisEvent `json:"events"`        // ❌ 应该是 []string
}
```

✅ **正确：直接映射 API 结构**
```go
type DisPolicyRule struct {
    ID     string   `json:"id"`
    Events []string `json:"events"`  // 直接字符串数组
    Agency string   `json:"agency"`  // agency 在规则内
}

type DisPolicyConfiguration struct {
    Rules []DisPolicyRule `json:"rules"`
}
```

---

## 三、trans() 方法实现

### 3.1 XML API trans 方法模板

```go
func (input SetBucketXXXInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
    // 1. 设置子资源参数
    params = map[string]string{string(SubResourceXXX): ""}

    // 2. 序列化为 XML
    reader, md5, convertErr := ConvertRequestToIoReaderV2(input.XXXConfiguration, false)
    if convertErr != nil {
        return nil, nil, nil, convertErr
    }

    // 3. 验证大小（如果需要）
    readerLen, err := GetReaderLen(reader)
    if err != nil {
        return nil, nil, nil, err
    }
    err = validateLength(int(readerLen), 0, MAX_SIZE, XML_SIZE)
    if err != nil {
        return nil, nil, nil, err
    }

    // 4. 设置 data 和 headers
    data = reader
    headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
    return
}
```

### 3.2 JSON API trans 方法模板

```go
func (input SetBucketXXXInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
    // 1. 设置子资源参数
    params = map[string]string{string(SubResourceXXX): ""}

    // 2. 设置 Content-Type header
    headers = make(map[string][]string, 1)
    headers[HEADER_CONTENT_TYPE] = []string{mimeTypes["json"]}

    // 3. 序列化为 JSON 字符串（⚠️ 关键：不能直接赋值结构体）
    data, err = convertXXXToJSON(input.XXXConfiguration)
    if err != nil {
        return nil, nil, nil, err
    }

    return
}
```

### 3.3 trans 方法对比

| 检查项 | XML API | JSON API |
|--------|---------|----------|
| 子资源参数 | ✅ `params[SubResourceXXX] = ""` | ✅ `params[SubResourceXXX] = ""` |
| Content-Type | ❌ 不需要 | ✅ `headers[HEADER_CONTENT_TYPE]` |
| 数据序列化 | ✅ `ConvertRequestToIoReaderV2` | ✅ `convertXXXToJSON` |
| MD5/SHA256 | ✅ 设置 `HEADER_MD5_CAMEL` | ❌ 不需要 |
| data 类型 | ✅ `io.Reader` | ✅ `string` |

---

## 四、签名计算

### 4.1 必须添加到 allowedResourceParameterNames

在 `obs/const.go` 中：

```go
allowedResourceParameterNames = map[string]bool{
    // ... 其他参数 ...
    "replication": true,  // 跨区域复制
    "dispolicy":   true,  // DIS 事件通知
    // ⚠️ 注意：使用小写，不是 disPolicy
}
```

### 4.2 签名计算原理

OBS SDK 的签名计算包含：
1. HTTP Method（GET/PUT/DELETE等）
2. URI 路径
3. **Query 参数**（包括子资源）← 必须在允许列表中
4. Headers
5. Request Body 的哈希

如果子资源不在 `allowedResourceParameterNames` 中，签名计算时会跳过该参数，导致签名验证失败。

---

## 五、完整开发流程

### Step 1: API 规范分析（30分钟）

1. 访问华为云官方文档
2. 记录 API 规范到临时文件
3. 确认所有字段类型

### Step 2: 代码实现（2-3小时）

按顺序修改：
1. `obs/type.go` - 添加 SubResource 常量
2. `obs/const.go` - 添加到 allowedResourceParameterNames
3. `obs/model_base.go` - 定义数据结构
4. `obs/model_bucket.go` - 定义 Input/Output
5. `obs/convert.go` - 添加序列化函数（JSON API）
6. `obs/trait_bucket.go` - 实现 trans() 方法
7. `obs/client_bucket.go` - 实现 API 方法

### Step 3: 单元测试（1-2小时）

1. 测试数据结构序列化/反序列化
2. 测试 trans() 方法输出
3. 测试参数验证

### Step 4: 集成测试（1小时）

1. 测试与真实 API 通信
2. 验证签名计算
3. 验证响应解析

### Step 5: 文档生成（30分钟）

1. 更新 docs/README.md
2. 创建功能文档目录
3. 添加使用示例

---

## 六、快速检查清单

### 编译检查
```bash
go build ./obs/...
```

### 单元测试
```bash
go test ./obs -v -run TestXXX
```

### 提交前检查

**文件检查**：
- [ ] `obs/type.go` - 子资源常量已定义
- [ ] `obs/const.go` - 子资源已添加到 `allowedResourceParameterNames`（小写）
- [ ] `obs/model_base.go` - 数据结构定义正确
- [ ] `obs/model_bucket.go` - Input/Output 已定义
- [ ] `obs/convert.go` - JSON 序列化函数已添加（JSON API）
- [ ] `obs/trait_bucket.go` - trans() 方法正确实现
- [ ] `obs/client_bucket.go` - API 方法已实现

**trans() 方法特别检查**：
- [ ] 子资源参数设置：`params[SubResourceXXX] = ""`
- [ ] JSON API：Content-Type header 已设置
- [ ] JSON API：使用 `convertXXXToJSON()` 序列化（不是直接赋值结构体）
- [ ] XML API：使用 `ConvertRequestToIoReaderV2()` 序列化
- [ ] data 变量类型正确（JSON 用 string，XML 用 io.Reader）

---

## 版本信息

- 文档版本: 1.0
- 创建日期: 2026-03-23
- 适用 SDK 版本: 3.26.0+
