# go-sdk-fuzz

## 技能概述

OBS SDK Go模糊测试编写指南。本技能指导用户如何为华为云OBS SDK识别关键输入解析函数并编写高质量的模糊测试，用于发现潜在的漏洞和边界条件问题。

## 使用场景

- 识别关键的输入解析函数
- 编写模糊测试发现潜在bug
- 测试XML解析、URL解析、签名验证等敏感功能
- 验证输入处理的鲁棒性

## 核心功能

### 1. 关键函数识别

自动识别需要模糊测试的关键函数：
- **XML解析函数**：`TransToXml`、`XmlToTrans`
- **URL解析函数**：endpoint处理、URL构造函数
- **签名验证函数**：`authV2_test.go`中的签名相关函数
- **响应解析函数**：各种API响应的解析函数

### 2. Fuzzing配置

配置模糊测试参数：
- 测试时长（默认30s）
- 工作线程数量
- 最大输入大小限制
- 种子文件管理

### 3. 崩溃分析流程

提供系统性的崩溃分析方法：
- 崩溃日志分析
- 输入复现
- 根本原因定位
- 修复建议

### 4. Fuzzing最佳实践

指导如何编写有效的模糊测试：
- 超时设置
- 输入边界
- 资源限制
- 性能监控

## 目标函数清单

### XML解析相关
- `TransToXml` - 将结构体转换为XML
- `XmlToTrans` - 将XML解析为结构体
- 各种XML解析辅助函数

### URL解析相关
- `New`函数中的endpoint处理
- URL构造和解析函数
- 路径处理函数

### 签名验证相关
- `authV2.go`中的签名生成函数
- `authV4.go`中的签名验证函数
- `SignatureObs`相关函数

### 响应解析相关
- 错误响应解析函数
- 成功响应解析函数
- 自定义头部解析函数

## 使用方法

### 基本用法

```bash
/go-sdk-fuzz
```

### 带参数使用

```bash
/go-sdk-fuzz --targets=xml,url,auth
/go-sdk-fuzz --duration=60s
/go-sdk-fuzz --workers=4
/go-sdk-fuzz --seeds=./fuzz_seeds/
```

### 技能输出

1. **模糊测试文件**：
   ```
   obs/*_fuzz_test.go
   ```

2. **目标函数列表**：
   - 关键函数识别结果
   - 测试优先级排序

3. **配置指南**：
   - fuzzing配置建议
   - 超时和资源设置

4. **分析工具**：
   - 崩溃分析脚本
   - 输入生成器

## 输出示例

### 生成的模糊测试文件示例

```go
//go:build fuzz

package obs

import (
	"testing"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

// FuzzTransToXml 测试XML解析函数
func FuzzTransToXml(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		input interface{}
	}{
		{&obs.CreateBucketInput{Bucket: "test"}},
		{&obs.PutObjectInput{Bucket: "test", Key: "test"}},
		{&obs.GetObjectInput{Bucket: "test", Key: "test"}},
	}

	for _, seed := range seeds {
		f.Add(seed.input)
	}

	// Fuzzing测试
	f.Fuzz(func(t *testing.T, input interface{}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 10000 {
			t.Skip("跳过长时间运行")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			// 记录XML解析错误
			t.Logf("XML转换失败: %v, input: %v", err, input)
			return
		}

		// 验证生成的XML
		if len(xmlBytes) == 0 {
			t.Error("生成的XML为空")
		}
	})
}

// FuzzUrlParsing 测试URL解析函数
func FuzzUrlParsing(f *testing.F) {
	// 添加种子数据
	seeds := []string{
		"https://obs.cn-north-4.myhuaweicloud.com",
		"http://localhost:8080",
		"https://example.com/bucket/object",
		"https://example.com/bucket/",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	// Fuzzing测试
	f.Fuzz(func(t *testing.T, urlStr string) {
		// 防止URL过长
		if len(urlStr) > 2048 {
			t.Skip("URL过长，跳过")
		}

		// 执行URL解析
		_, err := parseObsUrl(urlStr)
		if err != nil {
			// 记录解析错误
			t.Logf("URL解析失败: %v, url: %s", err, urlStr)
			return
		}
	})
}

// FuzzAuthSignature 测试签名验证函数
func FuzzAuthSignature(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		accessKey string
		secretKey string
		method    string
		path      string
	}{
		{"test", "test", "GET", "/"},
		{"test", "test", "PUT", "/bucket/object"},
		{"test", "test", "POST", "/bucket/object?uploadId=test"},
	}

	for _, seed := range seeds {
		f.Add(seed.accessKey, seed.secretKey, seed.method, seed.path)
	}

	// Fuzzing测试
	f.Fuzz(func(t *testing.T, accessKey, secretKey, method, path string) {
		// 防止凭据过长
		if len(accessKey) > 100 || len(secretKey) > 100 {
			t.Skip("凭据过长，跳过")
		}

		// 执行签名计算
		signature := calculateSignatureV4(accessKey, secretKey, method, path)
		if signature == "" {
			t.Error("签名生成失败")
		}
	})
}
```

### 崩溃分析报告示例

```json
{
  "crash_reports": [
    {
      "function": "TransToXml",
      "input": " malicious XML input here",
      "error": "xml: syntax error",
      "stack_trace": "...",
      "recommendation": "Add XML validation before parsing"
    },
    {
      "function": "parseObsUrl",
      "input": "https://[invalid-url]",
      "error": "parse failed",
      "stack_trace": "...",
      "recommendation": "Add URL validation"
    }
  ],
  "total_executions": 100000,
  "crashes_found": 2,
  "timeout_issues": 0
}
```

## 运行模糊测试

### 基本运行

```bash
# 运行所有模糊测试
go test -tags fuzz ./obs -fuzz=.

# 运行特定函数的模糊测试
go test -tags fuzz ./obs -fuzz=FuzzTransToXml

# 设置测试时长
go test -tags fuzz ./obs -fuzz=. -fuzztime=60s
```

### 使用工作线程

```bash
# 使用多个工作线程（并行执行）
go test -tags fuzz ./obs -fuzz=. -parallel=4
```

### 使用种子文件

```bash
# 使用种子文件启动
go test -tags fuzz ./obs -fuzz=FuzzFunction -fuzzinput=./seeds/
```

## 最佳实践

### 1. 输入边界设置

```go
f.Fuzz(func(t *testing.T, input string) {
    // 防止输入过大
    if len(input) > 1024*1024 { // 1MB限制
        t.Skip("输入过大，跳过")
    }

    // 执行测试...
})
```

### 2. 超时控制

```go
f.Fuzz(func(t *testing.T, input interface{}) {
    // 防止测试运行时间过长
    if len(f.Fuzzing()) > 10000 {
        t.Skip("跳过长时间运行")
    }

    // 执行测试...
})
```

### 3. 资源监控

```go
f.Fuzz(func(t *testing.T, input interface{}) {
    // 检查内存使用
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    if m.Alloc > 100*1024*1024 { // 100MB限制
        t.Skip("内存使用过高，跳过")
    }

    // 执行测试...
})
```

### 4. 错误记录

```go
f.Fuzz(func(t *testing.T, input interface{}) {
    // 执行测试
    result, err := processInput(input)
    if err != nil {
        // 记录错误但不崩溃
        t.Logf("处理失败: %v, input: %v", err, input)
        return
    }

    // 验证结果
    if result == nil {
        t.Error("结果为nil")
    }
})
```

## 故障排除

### 常见问题

1. **模糊测试没有输出**
   ```
   检查build tag: //go:build fuzz
   ```

2. **测试运行缓慢**
   ```
   使用 -parallel 参数增加工作线程
   或者减少测试时长
   ```

3. **内存不足**
   ```
   添加输入大小限制
   增加跳过条件
   ```

4. **测试超时**
   ```
   添加超时检查
   使用 -fuzztime 参数设置更长的测试时间
   ```

### 调试技巧

1. **查看种子文件**
   ```bash
   go test -tags fuzz ./obs -fuzz=FuzzFunction -fuzzoutput=./corpus/
   ```

2. **分析崩溃文件**
   ```bash
   cat crash-*.txt
   ```

3. **使用小输入范围测试**
   ```bash
   # 限制输入长度
   go test -tags fuzz ./obs -fuzz=FuzzFunction -fuzztime=10s
   ```

## 注意事项

1. **性能影响**：模糊测试会消耗大量CPU和内存资源
2. **测试时间**：建议单独运行，不要与其他测试一起
3. **敏感数据**：避免在生产环境运行模糊测试
4. **覆盖率**：关注代码覆盖率，确保测试有效

## 技能版本

- 版本：1.1.0
- 最后更新：2024-01-09
- 兼容性：OBS SDK v3.x
- Go版本：1.18+ (支持模糊测试)
- **支持技能调用**: 可被go-sdk-dev-task技能调用和协调

## 技能调用接口

本技能可以被go-sdk-dev-task技能调用，以实现开发任务的自动化安全测试。

### 调用方式

go-sdk-dev-task可以通过以下方式调用本技能：

```bash
# 基本调用（使用默认参数）
/go-sdk-dev-task --type=security

# 指定模糊测试目标
/go-sdk-dev-task --fuzz-targets=xml,url,auth

# 指定测试时长
/go-sdk-dev-task --fuzz-duration=60s
```

### 技能协调机制

1. **安全优先策略**
   - 优先为安全关键函数生成模糊测试
   - 识别潜在的输入验证漏洞
   - 检测SQL注入、XSS、XML注入等攻击

2. **去重检查**
   - 检查目标函数是否已有模糊测试
   - 避免重复生成相同测试用例
   - 采用一致的测试命名规范

3. **资源管理**
   - 使用统一的测试语料库目录
   - 避免模糊测试资源冲突
   - 协调崩溃报告存储位置

4. **进度同步**
   - 向go-sdk-dev-task报告模糊测试生成进度
   - 汇总模糊测试结果
   - 识别待解决的安全问题