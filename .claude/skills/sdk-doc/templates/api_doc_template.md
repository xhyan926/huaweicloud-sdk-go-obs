# {{ FEATURE_NAME }} API 接口文档

本文档包含 {{ SDK_NAME }} 中 {{ FEATURE_NAME }} 相关的所有 API 接口说明。

## 目录

- [{{ API_1_NAME }}](#{{ API_1_ID }})
- [{{ API_2_NAME }}](#{{ API_2_ID }})

---

## {{ API_1_NAME }}

{{ API_1_DESCRIPTION }}

### 方法签名

```go
func {{ API_1_SIGNATURE }}
```

### 参数说明

{{ API_1_INPUT_TYPE }}

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
{{ API_1_PARAMS_TABLE }}

### 返回值

{{ API_1_OUTPUT_TYPE }}

| 字段 | 类型 | 说明 |
|------|------|------|
{{ API_1_FIELDS_TABLE }}

### 使用示例

```go
package main

import (
    "fmt"
    {{ IMPORT_PATH }}
)

func main() {
    client, err := {{ IMPORT_PATH }}.New("your-key", "your-secret", "endpoint")
    if err != nil {
        panic(err)
    }

    input := &{{ IMPORT_PATH }}.{{ API_1_INPUT_TYPE }}{
{{ API_1_EXAMPLE_INPUT }}
    }

    output, err := client.{{ API_1_METHOD_NAME }}(input)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("结果: %v\n", output)
}
```

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
{{ API_1_ERROR_CODES_TABLE }}

### 注意事项

{{ API_1_NOTES }}

---

## {{ API_2_NAME }}

{{ API_2_DESCRIPTION }}

### 方法签名

```go
func {{ API_2_SIGNATURE }}
```

### 参数说明

{{ API_2_INPUT_TYPE }}

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
{{ API_2_PARAMS_TABLE }}

### 返回值

{{ API_2_OUTPUT_TYPE }}

| 字段 | 类型 | 说明 |
|------|------|------|
{{ API_2_FIELDS_TABLE }}

### 使用示例

```go
{{ API_2_EXAMPLE }}
```

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
{{ API_2_ERROR_CODES_TABLE }}

### 注意事项

{{ API_2_NOTES }}

---

## 常量定义

```go
{{ CONSTANTS_CODE }}
```

## 使用场景

### 场景 1：{{ SCENARIO_1_NAME }}

{{ SCENARIO_1_DESCRIPTION }}

```go
{{ SCENARIO_1_CODE }}
```

### 场景 2：{{ SCENARIO_2_NAME }}

{{ SCENARIO_2_DESCRIPTION }}

```go
{{ SCENARIO_2_CODE }}
```

## 相关文档

- [API 参考文档]({{ API_REFERENCE_URL }})
- [示例代码](../../examples/{{ FEATURE_DIR }}/)

---

**文档版本**: 1.0
**更新日期**: {{ UPDATE_DATE }}
