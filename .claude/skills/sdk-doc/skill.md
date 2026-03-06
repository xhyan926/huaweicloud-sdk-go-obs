# SDK API 文档编写指南

## 技能用途

本技能为 SDK（Software Development Kit）项目提供 API 接口文档编写的完整指导。它包含了文档结构设计、内容规范、格式标准、最佳实践和模板，帮助开发者编写高质量、易维护、用户友好的 API 文档。

## 适用场景

- 为 SDK 项目创建 API 接口文档
- 设计文档目录结构和组织方式
- 编写单个 API 接口的详细文档
- 创建文档总索引和导航
- 审查和改进现有 SDK 文档
- 建立文档编写规范和流程

## 前置条件

在使用此技能前，确保：

1. **项目环境正确**：已了解 SDK 的基本结构和功能
2. **API 接口已知**：清楚需要文档化的 API 接口及其参数
3. **代码已实现**：对应的 API 方法已经实现并可用
4. **测试已通过**：API 接口功能已验证

## 输入格式

用户提供以下信息之一：

```
"为桶清单功能生成 API 接口文档"
"为 PutObject 接口编写文档"
"设计 SDK 的文档目录结构"
"创建文档总索引"
```

## 输出格式

技能生成以下文档：

### 1. 文档目录结构
```
docs/
├── README.md                 # 文档总索引
├── bucket/                   # 按特性组织的子目录
│   └── README.md            # 特性详细文档
├── object/
├── multipart/
└── ...
```

### 2. 总索引文档模板
包含导航、使用基础、规范说明等

### 3. API 接口文档模板
每个接口包含完整的信息：方法签名、参数、返回值、示例、错误码等

## 文档结构设计

### 目录组织原则

**按特性/功能模块组织**

```
docs/
├── README.md                 # 总索引
├── bucket/                   # 桶操作
├── object/                   # 对象操作
├── multipart/                # 分块上传
├── lifecycle/                # 生命周期管理
├── encryption/               # 加密
└── [其他特性]/
```

**优点**：
- 清晰的功能分类
- 易于维护和扩展
- 符合用户使用习惯

### 目录命名规范

| 特性类型 | 目录命名 | 示例 |
|---------|---------|------|
| 桶操作 | bucket/ | CreateBucket, DeleteBucket |
| 对象操作 | object/ | PutObject, GetObject |
| 分块上传 | multipart/ | InitiateMultipartUpload |
| 生命周期 | lifecycle/ | SetBucketLifecycle |
| 加密 | encryption/ | SetBucketEncryption |
| 权限 | acl/ | SetBucketAcl, SetObjectAcl |

## 文档内容规范

### 1. 总索引文档 (docs/README.md)

#### 必需章节

```markdown
# SDK API 接口文档

## 目录结构
(展示 docs/ 目录结构)

## 快速导航
(按特性分类的 API 列表)

## SDK 使用基础
- 创建客户端
- 使用扩展选项
- 错误处理

## API 命名规范
(说明 SDK 中的命名规则)

## 文档规范
(说明文档的组织方式和格式)

## 版本信息
- 文档版本
- SDK 版本
- 更新日期

## 相关资源
(链接到官方文档、示例代码等)

## 反馈与支持
(问题反馈渠道)
```

#### SDK 使用基础模板

```go
// 创建客户端
sdkClient, err := sdk.New("your-key", "your-secret", "endpoint")

// 使用扩展选项
output, err := sdkClient.Method(input, sdk.WithOption(value))

// 错误处理
if err != nil {
    if sdkError, ok := err.(sdk.SdkError); ok {
        fmt.Printf("错误码: %s\n", sdkError.Code)
        fmt.Printf("错误信息: %s\n", sdkError.Message)
    }
}
```

### 2. API 接口文档 (docs/[feature]/README.md)

#### 单个 API 接口文档结构

```markdown
### InterfaceName

简要描述接口功能（1-2 句话）

#### 方法签名

\`\`\`go
func MethodSignature(params...) (output, error)
\`\`\`

#### 参数说明

**InputStructure**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| param1 | Type | Yes | 参数说明 |
| param2 | Type | No | 可选参数说明 |

**NestedStructure**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ... | ... | ... | ... |

#### 返回值

**OutputStructure**

| 字段 | 类型 | 说明 |
|------|------|------|
| field1 | Type | 字段说明 |
| field2 | Type | 字段说明 |

#### 使用示例

\`\`\`go
// 完整可运行的代码示例
package main

import (
    "fmt"
    sdk "path/to/sdk"
)

func main() {
    client, _ := sdk.New("key", "secret", "endpoint")
    input := &sdk.MethodInput{
        // 设置参数
    }

    output, err := client.Method(input)
    if err != nil {
        // 错误处理
        return
    }
    // 使用输出
}
\`\`\`

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| ErrorCode1 | 400 | 错误说明 |
| ErrorCode2 | 404 | 错误说明 |

#### 注意事项

1. 重要提示1
2. 重要提示2
```

#### 常量定义章节

```markdown
## 常量定义

### 常量类型说明

\`\`\`go
type ConstantType string

const (
    ConstantValue1 ConstantType = "value1"  // 说明
    ConstantValue2 ConstantType = "value2"  // 说明
)
\`\`\`
```

#### 使用场景章节

```markdown
## 使用场景

### 场景 1：场景名称

简要描述场景用途。

\`\`\`go
input := &sdk.MethodInput{
    // 场景特定的参数配置
}
\`\`\`
```

## 文档格式规范

### Markdown 规范

#### 标题层级

```markdown
# 一级标题（文档标题）
## 二级标题（主要章节）
### 三级标题（子章节）
#### 四级标题（详细说明）
```

**使用原则**：
- 每个一级标题只出现一次
- 二级标题按功能划分
- 三级标题用于具体条目
- 避免超过四级标题

#### 代码块

```markdown
\`\`\`go
// 代码块必须指定语言类型
func Example() {
    // 代码内容
}
\`\`\`
```

**要求**：
- 必须指定语言类型（go、java、python 等）
- 代码必须可运行
- 包含必要的 import 语句
- 有错误处理

#### 表格

```markdown
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | Type | Yes | 描述 |
```

**要求**：
- 使用规范的 Markdown 表格语法
- 表头清晰
- 对齐一致
- 避免空行

#### 链接

```markdown
- 内部链接: [文本](./relative/path.md)
- 外部链接: [文本](https://example.com)
- 锚点链接: [文本](#section-name)
```

**要求**：
- 内部链接使用相对路径
- 锚点名称使用小写和连字符
- 确保链接有效

### 代码示例规范

#### 完整性要求

```go
// ✅ 正确：完整可运行
package main

import (
    "fmt"
    sdk "github.com/example/sdk"
)

func main() {
    client, err := sdk.New("key", "secret", "endpoint")
    if err != nil {
        panic(err)
    }

    input := &sdk.MethodInput{
        Param: "value",
    }

    output, err := client.Method(input)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("结果: %v\n", output)
}
```

```go
// ❌ 错误：缺少关键部分
func main() {
    output, _ := client.Method(input)
}
```

#### 示例命名

```
示例文件命名: [feature]_sample.go
示例函数命名: Example[Feature][Method]()
```

### 参数说明规范

#### 表格格式

| 字段 | 类型 | 说明 |
|------|------|------|
| 参数名称 | 参数类型 | 参数说明 |

**必填标识**：
- Yes / No
- 是 / 否
- Required / Optional

#### 类型说明

| 类型 | 示例 |
|------|------|
| string | "example" |
| int | 123 |
| bool | true/false |
| []Type | []string{"a", "b"} |
| map[K]V | map[string]int{"key": 1} |

## 最佳实践

### 1. 以用户为中心

**原则**：
- 从用户视角编写文档
- 假设用户是第一次使用
- 提供清晰的入门指南

**示例**：

```markdown
✅ 好的做法：
首先，创建一个客户端实例...

❌ 差的做法：
客户端初始化需要三个参数...
```

### 2. 保持简洁准确

**原则**：
- 避免冗余信息
- 描述准确清晰
- 突出重点

**示例**：

```markdown
✅ 好的做法：
Bucket - 存储对象的容器

❌ 差的做法：
Bucket 是 OBS 中用来存储对象的基本单元，它类似于文件系统中的目录，但在 OBS 中...
```

### 3. 提供完整示例

**原则**：
- 示例代码完整可运行
- 包含错误处理
- 展示常见用法
- 添加必要的注释

**示例**：

```go
// ✅ 好的做法
input := &sdk.SetBucketInventoryInput{
    Bucket: "my-bucket",
    InventoryConfiguration: sdk.InventoryConfiguration{
        Id:        "inventory-1",  // 清单规则 ID
        IsEnabled: true,           // 启用清单
        Destination: sdk.InventoryDestination{
            Format: "CSV",         // CSV 格式
            Bucket: "report-bucket",
            Prefix: "reports/",
        },
        Schedule: sdk.InventorySchedule{
            Frequency: sdk.InventoryFrequencyDaily,  // 每日生成
        },
    },
}

output, err := client.SetBucketInventory(input)
if err != nil {
    return err  // 处理错误
}
```

### 4. 清晰的错误处理

**原则**：
- 列出所有可能的错误码
- 说明错误原因
- 提供解决建议

**示例**：

```markdown
| 错误码 | HTTP 状态码 | 说明 | 解决方法 |
|--------|-------------|------|---------|
| NoSuchBucket | 404 | 桶不存在 | 确认桶名正确，桶已创建 |
| AccessDenied | 403 | 权限不足 | 检查 AK/SK 是否正确 |
```

### 5. 维护文档一致性

**原则**：
- 使用统一的术语
- 保持格式一致
- 同步更新文档和代码
- 定期审查和更新

**检查清单**：

- [ ] 所有接口文档结构一致
- [ ] 参数命名风格统一
- [ ] 示例代码风格一致
- [ ] 错误码说明格式一致

## 文档模板

### 总索引文档模板

```markdown
# [SDK Name] API 接口文档

欢迎使用 [SDK Name] API 接口文档。本文档提供了 SDK 中所有 API 接口的详细说明、使用示例和最佳实践。

## 目录结构

本文档按照功能特性进行组织：

```
docs/
├── README.md
├── [feature1]/
│   └── README.md
├── [feature2]/
│   └── README.md
└── ...
```

## 快速导航

### [特性 1]

- [API 1](./[feature1]/README.md#api1)
- [API 2](./[feature1]/README.md#api2)

## SDK 使用基础

### 创建客户端

\`\`\`go
import sdk "path/to/sdk"

client, err := sdk.New("your-key", "your-secret", "endpoint")
\`\`\`

### 错误处理

\`\`\`go
output, err := client.Method(input)
if err != nil {
    if sdkError, ok := err.(sdk.SdkError); ok {
        fmt.Printf("错误码: %s\n", sdkError.Code)
    }
}
\`\`\`

## API 命名规范

| 操作类型 | 前缀 | 示例 |
|---------|------|------|
| 创建 | Set/Create | SetBucket, CreateObject |
| 获取 | Get | GetBucket, GetObject |
| 列举 | List | ListObjects, ListBuckets |
| 删除 | Delete | DeleteBucket, DeleteObject |

## 版本信息

- 文档版本: 1.0
- SDK 版本: X.Y.Z+
- 更新日期: YYYY-MM-DD

## 相关资源

- [官方文档](https://example.com)
- [示例代码](../examples/)
- [更新日志](../CHANGELOG.md)
```

### API 接口文档模板

```markdown
# [特性名称] API 接口文档

本文档包含 [SDK Name] 中 [特性名称] 相关的所有 API 接口说明。

## 目录

- [接口 1](#interface1)
- [接口 2](#interface2)

---

## 接口 1

简要描述接口功能。

### 方法签名

\`\`\`go
func (client Client) Method(input *InputType, options ...OptionType) (*OutputType, error)
\`\`\`

### 参数说明

**InputType**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| param1 | string | 是 | 参数1说明 |
| param2 | int | 否 | 参数2说明，默认值0 |

### 返回值

**OutputType**

| 字段 | 类型 | 说明 |
|------|------|------|
| field1 | string | 字段1说明 |
| field2 | int | 字段2说明 |

### 使用示例

\`\`\`go
package main

import (
    "fmt"
    sdk "path/to/sdk"
)

func main() {
    client, err := sdk.New("key", "secret", "endpoint")
    if err != nil {
        panic(err)
    }

    input := &sdk.MethodInput{
        Param1: "value",
        Param2: 123,
    }

    output, err := client.Method(input)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("结果: %v\n", output)
}
\`\`\`

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidArgument | 400 | 参数错误 |
| AccessDenied | 403 | 权限不足 |
| NotFound | 404 | 资源不存在 |

### 注意事项

1. 注意点1
2. 注意点2

---

## 常量定义

\`\`\`go
const (
    Constant1 = "value1"  // 说明
    Constant2 = "value2"  // 说明
)
\`\`\`

## 使用场景

### 场景 1：场景名称

\`\`\`go
input := &sdk.MethodInput{
    // 场景特定配置
}
\`\`\`

## 相关文档

- [API 参考文档](https://example.com/api)
- [示例代码](../../examples/)
```

## 质量检查清单

### 内容完整性

- [ ] 包含所有必需章节
- [ ] 方法签名准确
- [ ] 参数说明完整
- [ ] 返回值说明完整
- [ ] 示例代码可运行
- [ ] 错误码列表完整

### 格式规范

- [ ] Markdown 格式正确
- [ ] 代码块带语言标识
- [ ] 表格格式正确
- [ ] 链接有效
- [ ] 标题层级合理

### 可读性

- [ ] 描述清晰简洁
- [ ] 示例易于理解
- [ ] 注释适当
- [ ] 无错别字
- [ ] 术语一致

### 维护性

- [ ] 版本信息更新
- [ ] 同步代码变更
- [ ] 链接有效
- [ ] 示例代码最新

## 工作流程

### 第一阶段：需求分析

1. 确定需要文档化的 API 接口范围
2. 分析功能特性，确定目录结构
3. 收集必要的接口信息（签名、参数、返回值）

### 第二阶段：目录设计

1. 设计 docs/ 目录结构
2. 按特性创建子目录
3. 创建总索引文档框架

### 第三阶段：文档编写

1. 为每个特性编写详细文档
2. 为每个 API 接口编写完整文档
3. 添加使用示例和场景

### 第四阶段：质量检查

1. 使用质量检查清单验证
2. 确保示例代码可运行
3. 检查所有链接有效

### 第五阶段：验收和发布

1. 生成验收报告
2. 更新项目文档索引
3. 同步更新 README

## 注意事项

### 常见错误

1. **过度复杂**：避免不必要的复杂结构
2. **信息过时**：及时同步代码变更
3. **示例错误**：确保示例代码可运行
4. **链接失效**：定期检查内部链接
5. **术语混乱**：保持术语一致性

### 安全注意事项

- 不要在文档中硬编码真实凭证
- 使用占位符替代敏感信息
- 提醒用户注意凭证安全

## 参考资源

- [Markdown 语法指南](https://www.markdownguide.org/)
- [技术文档最佳实践](https://developers.google.com/tech-writing)
- [API 文档设计指南](https://stoplight.io/blog/api-documentation-guide/)
