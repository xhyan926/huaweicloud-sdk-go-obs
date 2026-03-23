# DIS 事件通知 API 接口文档

本文档包含华为云 OBS Go SDK 中 DIS 事件通知（DIS Event Notification）相关的所有 API 接口说明。

## 目录

- [设置 DIS 事件通知策略](#设置-dis-事件通知策略-setbucketdispolicy)
- [获取 DIS 事件通知配置](#获取-dis-事件通知配置-getbucketdispolicy)
- [删除 DIS 事件通知策略](#删除-dis-事件通知策略-deletebucketdispolicy)
- [数据结构定义](#数据结构定义)
- [常量定义](#常量定义)
- [使用场景](#使用场景)

---

## 设置 DIS 事件通知策略

为指定桶设置 DIS 事件通知策略，将对象操作事件发送到 DIS 数据接入服务。

### 方法签名

```go
func (obsClient ObsClient) SetBucketDisPolicy(input *SetBucketDisPolicyInput, extensions ...extensionOptions) (output *BaseModel, err error)
```

### 参数说明

**SetBucketDisPolicyInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名 |
| DisPolicyConfiguration | DisPolicyConfiguration | 是 | DIS 策略配置 |

**DisPolicyConfiguration**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Rules | []DisPolicyRule | 是 | 规则数组，1-10 条 |

**DisPolicyRule**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ID | string | 是 | 规则 ID |
| Stream | string | 是 | DIS 通道名称 |
| Project | string | 是 | DIS 项目 ID |
| Events | []string | 是 | 事件列表（字符串数组） |
| Prefix | string | 否 | 对象名前缀 |
| Suffix | string | 否 | 对象名后缀 |
| Agency | string | 是 | IAM 委托名 |

### 返回值

**BaseModel**

| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID |
| Headers | map[string][]string | 响应头 |

### 使用示例

```go
package main

import (
	"fmt"
	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func main() {
	// 创建 OBS 客户端
	ak := "your-access-key"
	sk := "your-secret-key"
	endpoint := "https://obs.cn-north-4.myhuaweicloud.com"

	client, err := obs.New(ak, sk, endpoint, obs.WithSignature(obs.SignatureObs))
	if err != nil {
		panic(err)
	}

	// 设置 DIS 事件通知策略
	input := &obs.SetBucketDisPolicyInput{
		Bucket: "my-bucket",
		DisPolicyConfiguration: obs.DisPolicyConfiguration{
			Rules: []obs.DisPolicyRule{
				{
					ID:      "rule-1",
					Stream:  "my-dis-stream",
					Project: "my-project-id",
					Events:  []string{"ObjectCreated:*", "ObjectRemoved:*"},
					Prefix:  "images/",
					Suffix:  ".jpg",
					Agency:  "my-dis-agency",
				},
			},
		},
	}

	output, err := client.SetBucketDisPolicy(input)
	if err != nil {
		fmt.Printf("设置 DIS 事件通知策略失败: %v\n", err)
		return
	}

	fmt.Printf("设置成功，RequestId: %s\n", output.RequestId)
}
```

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| MalformedJSON | 400 | JSON 格式错误 |
| InvalidRequest | 400 | 请求参数无效 |
| NoSuchBucket | 404 | 桶不存在 |
| Unauthorized | 403 | 未授权 |
| AgencyNotFound | 400 | 委托不存在 |

### 注意事项

1. DIS 事件通知功能仅支持 OBS 签名，不支持 AWS 签名
2. 同一桶内最多可配置 10 条事件规则
3. 需要先创建委托并授予 OBS 访问 DIS 的权限
4. 事件列表为字符串数组，不是对象数组

---

## 获取 DIS 事件通知配置

获取指定桶的 DIS 事件通知配置。

### 方法签名

```go
func (obsClient ObsClient) GetBucketDisPolicy(bucketName string, extensions ...extensionOptions) (output *GetBucketDisPolicyOutput, err error)
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

### 返回值

**GetBucketDisPolicyOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| BaseModel | - | 基础响应模型 |
| DisPolicyConfiguration | string | DIS 策略配置（JSON 字符串） |

### 使用示例

```go
// 获取 DIS 事件通知配置
output, err := client.GetBucketDisPolicy("my-bucket")
if err != nil {
	fmt.Printf("获取 DIS 事件通知配置失败: %v\n", err)
	return
}

// 解析 JSON 配置
var config obs.DisPolicyConfiguration
if err := json.Unmarshal([]byte(output.DisPolicyConfiguration), &config); err != nil {
	fmt.Printf("解析配置失败: %v\n", err)
	return
}

// 遍历规则
for _, rule := range config.Rules {
	fmt.Printf("规则 ID: %s\n", rule.ID)
	fmt.Printf("DIS 通道: %s\n", rule.Stream)
	fmt.Printf("事件列表: %v\n", rule.Events)
}
```

---

## 删除 DIS 事件通知策略

删除指定桶的 DIS 事件通知策略。

### 方法签名

```go
func (obsClient ObsClient) DeleteBucketDisPolicy(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error)
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

### 返回值

**BaseModel**

| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID |
| Headers | map[string][]string | 响应头 |

### 使用示例

```go
// 删除 DIS 事件通知策略
output, err := client.DeleteBucketDisPolicy("my-bucket")
if err != nil {
	fmt.Printf("删除 DIS 事件通知策略失败: %v\n", err)
	return
}

fmt.Printf("删除成功，RequestId: %s\n", output.RequestId)
```

---

## 数据结构定义

```go
// DIS 策略规则
type DisPolicyRule struct {
	ID      string   `json:"id"`      // 规则 ID
	Stream  string   `json:"stream"`  // DIS 通道名称
	Project string   `json:"project"` // 项目 ID
	Events  []string `json:"events"`  // 事件列表（字符串数组）
	Prefix  string   `json:"prefix,omitempty"`  // 对象名前缀
	Suffix  string   `json:"suffix,omitempty"`  // 对象名后缀
	Agency  string   `json:"agency"`  // IAM 委托名
}

// DIS 策略配置
type DisPolicyConfiguration struct {
	Rules []DisPolicyRule `json:"rules"` // 规则数组，1-10 条
}
```

---

## 常量定义

```go
// 子资源类型
const (
	SubResourceDisPolicy SubResourceType = "disPolicy"
)
```

---

## 支持的事件类型

| 事件名称 | 说明 |
|---------|------|
| ObjectCreated:* | 所有对象创建事件 |
| ObjectCreated:Put | PUT 对象事件 |
| ObjectCreated:Post | POST 对象事件 |
| ObjectCreated:Copy | 复制对象事件 |
| ObjectCreated:CompleteMultipartUpload | 完成分块上传事件 |
| ObjectRemoved:* | 所有对象删除事件 |
| ObjectRemoved:Delete | DELETE 对象事件 |
| ObjectRemoved:DeleteMarkerCreated | 删除标记创建事件 |

---

## 使用场景

### 场景 1：监听所有对象创建和删除事件

```go
input := &obs.SetBucketDisPolicyInput{
	Bucket: "my-bucket",
	DisPolicyConfiguration: obs.DisPolicyConfiguration{
		Rules: []obs.DisPolicyRule{
			{
				ID:      "rule-all-events",
				Stream:  "my-dis-stream",
				Project: "my-project-id",
				Events:  []string{"ObjectCreated:*", "ObjectRemoved:*"},
				Agency:  "my-dis-agency",
			},
		},
	},
}
```

### 场景 2：按文件类型过滤事件

只监听图片文件的创建事件。

```go
input := &obs.SetBucketDisPolicyInput{
	Bucket: "my-bucket",
	DisPolicyConfiguration: obs.DisPolicyConfiguration{
		Rules: []obs.DisPolicyRule{
			{
				ID:      "rule-images",
				Stream:  "my-dis-stream",
				Project: "my-project-id",
				Events:  []string{"ObjectCreated:*"},
				Prefix:  "images/",
				Suffix:  ".jpg",
				Agency:  "my-dis-agency",
			},
		},
	},
}
```

### 场景 3：多规则配置

配置多条规则处理不同的对象前缀。

```go
input := &obs.SetBucketDisPolicyInput{
	Bucket: "my-bucket",
	DisPolicyConfiguration: obs.DisPolicyConfiguration{
		Rules: []obs.DisPolicyRule{
			{
				ID:      "rule-images",
				Stream:  "my-dis-stream",
				Project: "my-project-id",
				Events:  []string{"ObjectCreated:*"},
				Prefix:  "images/",
				Agency:  "my-dis-agency",
			},
			{
				ID:      "rule-videos",
				Stream:  "my-dis-stream",
				Project: "my-project-id",
				Events:  []string{"ObjectRemoved:*"},
				Prefix:  "videos/",
				Agency:  "my-dis-agency",
			},
		},
	},
}
```

### 场景 4：选择性事件类型

只监听特定类型的对象操作。

```go
input := &obs.SetBucketDisPolicyInput{
	Bucket: "my-bucket",
	DisPolicyConfiguration: obs.DisPolicyConfiguration{
		Rules: []obs.DisPolicyRule{
			{
				ID:      "rule-specific",
				Stream:  "my-dis-stream",
				Project: "my-project-id",
				Events:  []string{
					"ObjectCreated:Put",
					"ObjectCreated:Post",
					"ObjectRemoved:Delete",
				},
				Agency: "my-dis-agency",
			},
		},
	},
}
```

---

## 委托配置说明

在使用 DIS 事件通知功能前，需要先创建委托并授予相应权限：

1. 登录华为云统一身份认证服务
2. 创建委托，委托名称为 `dis_access_agency`
3. 授予 OBS 服务访问 DIS 服务的权限
4. 在配置 DIS 策略时使用该委托名称

### 委托权限要求

委托需要具有以下权限：
- DIS 的操作权限（DIS OperateAccess）
- OBS 的读写权限

### 创建委托示例

在华为云控制台创建委托：
- 委托名称：`dis_access_agency`
- 委托类型：服务委托
- 被委托方：DIS 服务
- 权限配置：授予 OBS 访问 DIS 的权限

---

## 相关文档

- [华为云 OBS DIS 事件通知文档](https://support.huaweicloud.com/api-obs/obs_04_0139.html)
- [DIS 服务文档](https://support.huaweicloud.com/dis/)
- [API 参考文档](https://support.huaweicloud.com/api-obs/obs_04_0139.html)
- [示例代码](../../examples/)
