# Extension Rules

约束 OBS SDK 的扩展选项系统，包括函数式选项、header 扩展、进度回调和自定义头。

---

## EXT-01: 可选参数通过 extensionOptions 可变参数传递

**条款**: 禁止通过修改方法签名来添加可选参数（必须通过 `...extensionOptions` 可变参数传递）。

### 正确示例

```go
// 方法签名中的 extensions 参数
func (obsClient ObsClient) PutObject(input *PutObjectInput, extensions ...extensionOptions) (output *PutObjectOutput, err error) {
    // extensions 在 doAction 中被解析
}

// 调用方使用
output, err := client.PutObject(input,
    obs.WithProgress(myListener),
    obs.WithTrafficLimitHeader(838860800),
)
```

### 错误示例

```go
// 错误：在方法签名中添加专门的参数
func (obsClient ObsClient) PutObject(input *PutObjectInput, progressListener ProgressListener, trafficLimit int64) { ... }
```

### 验证清单
- [ ] 是否不存在因新增可选功能而修改方法签名的情况
- [ ] 是否不存在可选参数未通过 `...extensionOptions` 传递的情况
- [ ] 是否不存在调用方无法不传 extensions 的情况

---

## EXT-02: 新扩展实现对应接口

**条款**: 禁止定义不实现 `extensionHeaders` 或 `extensionProgressListener` 的扩展类型（如自定义 struct），禁止扩展函数返回非标准类型。

### 正确示例

```go
// extensionHeaders 实现
type extensionHeaders func(headers map[string][]string, isObs bool) error

func WithReqPaymentHeader(requester PayerType) extensionHeaders {
    return setHeaderPrefix(REQUEST_PAYER, string(requester))
}

// extensionProgressListener 实现
type extensionProgressListener func() ProgressListener

func WithProgress(progressListener ProgressListener) extensionProgressListener {
    return func() ProgressListener {
        return progressListener
    }
}
```

### 错误示例

```go
// 错误：自定义扩展类型，不实现标准接口
type myExtension struct {
    Key   string
    Value string
}

func WithMyExt(key, value string) myExtension {
    return myExtension{Key: key, Value: value}
}
```

### 验证清单
- [ ] 是否不存在新扩展返回非 `extensionHeaders` 或 `extensionProgressListener` 类型的情况
- [ ] 是否不存在 `extensionHeaders` 未接受 `(headers map[string][]string, isObs bool)` 参数的情况
- [ ] 是否不存在扩展函数未返回 `error`（用于校验失败）的情况

---

## EXT-03: 使用 setHeaderPrefix 自动处理 OBS/S3 双前缀

**条款**: 禁止手动拼接 `x-amz-`/`x-obs-` 前缀（必须使用 `setHeaderPrefix` 辅助函数），禁止绕过 `setHeaderPrefix` 直接操作 headers map。

### 正确示例

```go
// extension.go — 使用 setHeaderPrefix
func setHeaderPrefix(key string, value string) extensionHeaders {
    return func(headers map[string][]string, isObs bool) error {
        if strings.TrimSpace(value) == "" {
            return fmt.Errorf("set header %s with empty value", key)
        }
        setHeaders(headers, key, []string{value}, isObs)
        return nil
    }
}

func WithTrafficLimitHeader(trafficLimit int64) extensionHeaders {
    return setHeaderPrefix(TRAFFIC_LIMIT, strconv.FormatInt(trafficLimit, 10))
}
```

### 错误示例

```go
// 错误：手动处理前缀
func WithTrafficLimitHeader(trafficLimit int64) extensionHeaders {
    return func(headers map[string][]string, isObs bool) error {
        key := "x-amz-" + TRAFFIC_LIMIT  // 错误：手动拼接前缀
        if isObs {
            key = "x-obs-" + TRAFFIC_LIMIT
        }
        headers[key] = []string{strconv.FormatInt(trafficLimit, 10)}
        return nil
    }
}
```

### 验证清单
- [ ] 是否不存在新的 header 扩展未使用 `setHeaderPrefix` 辅助函数的情况
- [ ] 是否不存在手动拼接 `x-amz-` / `x-obs-` 前缀的情况
- [ ] 是否不存在空值校验未通过 `setHeaderPrefix` 统一处理的情况

---

## EXT-04: WithCustomHeader 直接使用原始 key

**条款**: 禁止在 `WithCustomHeader` 中自动添加 `x-amz-`/`x-obs-` 前缀（与 `setHeaderPrefix` 不同，此函数必须直接使用原始 key）。

### 正确示例

```go
// extension.go — WithCustomHeader 不添加前缀
func WithCustomHeader(key string, value string) extensionHeaders {
    return func(headers map[string][]string, isObs bool) error {
        if strings.TrimSpace(value) == "" {
            return fmt.Errorf("set header %s with empty value", key)
        }
        headers[key] = []string{value}
        return nil
    }
}

// 使用方指定完整的 header key
output, err := client.GetObject(input,
    obs.WithCustomHeader("X-Custom-Header", "value"),
)
```

### 错误示例

```go
// 错误：WithCustomHeader 中自动添加前缀
func WithCustomHeader(key string, value string) extensionHeaders {
    return func(headers map[string][]string, isObs bool) error {
        prefix := "x-amz-"
        if isObs {
            prefix = "x-obs-"
        }
        headers[prefix+key] = []string{value}  // 错误：不应添加前缀
        return nil
    }
}
```

### 验证清单
- [ ] 是否不存在 `WithCustomHeader` 自动添加 `x-amz-`/`x-obs-` 前缀的情况
- [ ] 是否不存在 `WithCustomHeader` 调用 `setHeaderPrefix` 或 `setHeaders` 的情况
- [ ] 是否不存在 `WithCustomHeader` 缺少空值校验的情况
