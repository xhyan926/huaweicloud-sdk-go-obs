# OBS SDK Go 模糊测试指南

## 快速开始

### 1. 确认Go版本

模糊测试需要Go 1.18+：

```bash
go version  # 需要 >= 1.18
```

### 2. 运行模糊测试

```bash
# 进入项目目录
cd /path/to/huaweicloud-sdk-go-obs

# 运行所有模糊测试
go test -tags fuzz ./obs -fuzz=.

# 运行特定函数
go test -tags fuzz ./obs -fuzz=FuzzTransToXml

# 设置测试时长（30秒）
go test -tags fuzz ./obs -fuzz=. -fuzztime=30s

# 使用多个工作线程
go test -tags fuzz ./obs -fuzz=. -parallel=4
```

## 使用方法

### 基本模糊测试

```go
//go:build fuzz

package obs

import (
	"testing"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

// FuzzXmlParsing 测试XML解析
func FuzzXmlParsing(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		input interface{}
	}{
		{&obs.CreateBucketInput{Bucket: "test-bucket"}},
		{&obs.PutObjectInput{Bucket: "test", Key: "test-key"}},
		{&obs.GetObjectInput{Bucket: "test", Key: "test-key"}},
	}

	for _, seed := range seeds {
		f.Add(seed.input)
	}

	// 模糊测试
	f.Fuzz(func(t *testing.T, input interface{}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 10000 {
			t.Skip("跳过长时间运行")
		}

		// 执行XML转换
		_, err := TransToXml(input)
		if err != nil {
			t.Logf("XML转换失败: %v, input: %v", err, input)
			return
		}
	})
}

// FuzzUrlHandler 测试URL处理
func FuzzUrlHandler(f *testing.F) {
	// 添加种子URL
	seeds := []string{
		"https://obs.cn-north-4.myhuaweicloud.com",
		"http://localhost:8080",
		"https://example.com/bucket/object",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	// 模糊测试
	f.Fuzz(func(t *testing.T, urlStr string) {
		// 限制URL长度
		if len(urlStr) > 2048 {
			t.Skip("URL过长，跳过")
		}

		// 执行URL解析
		_, err := parseObsUrl(urlStr)
		if err != nil {
			t.Logf("URL解析失败: %v, url: %s", err, urlStr)
			return
		}
	})
}
```

### 高级用法

#### 1. 使用种子文件

```bash
# 创建种子目录
mkdir -p ./corpus/FuzzFunction

# 将测试用例保存为种子文件
echo "test data" > ./corpus/FuzzFunction/seed1.txt

# 使用种子运行测试
go test -tags fuzz ./obs -fuzz=FuzzFunction -fuzzinput=./corpus/
```

#### 2. 监控资源使用

```go
func FuzzWithMonitoring(f *testing.F) {
	f.Fuzz(func(t *testing.T, input interface{}) {
		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 100*1024*1024 { // 100MB限制
			t.Skip("内存使用过高")
		}

		// 执行测试...
	})
}
```

#### 3. 自定义超时控制

```go
func FuzzWithTimeout(f *testing.F) {
	f.Fuzz(func(t *testing.T, input string) {
		// 防止输入过大
		if len(input) > 1024*1024 {
			t.Skip("输入过大")
		}

		// 执行测试...
	})
}
```

## 运行指南

### 命令行参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-fuzz=` | 指定要运行的模糊测试函数 | `-fuzz=FuzzXmlParsing` |
| `-fuzztime=` | 设置测试时长 | `-fuzztime=30s` |
| `-parallel=` | 设置工作线程数 | `-parallel=4` |
| `-count=` | 设置测试次数 | `-count=10` |

### 运行模式

1. **发现崩溃模式**
   ```bash
   go test -tags fuzz ./obs -fuzz=.
   ```

2. **持续测试模式**
   ```bash
   # 循环运行，直到按Ctrl+C停止
   while true; do go test -tags fuzz ./obs -fuzz=. -fuzztime=30s; done
   ```

3. **种子文件模式**
   ```bash
   # 使用现有种子文件
   go test -tags fuzz ./obs -fuzz=FuzzFunction -fuzzinput=./seeds/
   ```

### 输出分析

1. **正常输出**
   ```
   fuzz: elapsed 0s, 0/0 inputs
   ```
   表示正在运行，尚未发现崩溃

2. **发现崩溃**
   ```
   --- FAIL: FuzzXmlParsing (0.00s)
       --- FAIL: FuzzXmlParsing (0.00s)
           xml: syntax error
           --- FAIL: FuzzXmlParsing (0.00s)
               --- FAIL: FuzzXmlParsing (0.00s)
                   xml: syntax error: line 1: element "input" not closed
   ```
   会生成崩溃文件用于分析

3. **超时输出**
   ```
   --- PASS: FuzzXmlParsing (30.00s)
   ```
   测试完成未发现崩溃

## 目标函数分析

### XML解析函数

**关键函数**：
- `TransToXml` - XML序列化
- `XmlToTrans` - XML反序列化

**测试重点**：
- 恶意XML注入
- 格式错误处理
- 大型XML文档

### URL解析函数

**关键函数**：
- `New`函数中的endpoint处理
- URL解析和构造
- 路径规范化

**测试重点**：
- 恶意URL
- 超长URL
- 格式错误的URL

### 签名验证函数

**关键函数**：
- `calculateSignatureV2`
- `calculateSignatureV4`
- `SignatureObs`

**测试重点**：
- 空字符串输入
- 超长字符串
- 特殊字符处理

## 故障排除

### 常见问题

1. **模糊测试没有运行**
   ```bash
   # 检查build tag
   #go:build fuzz
   ```

2. **测试运行缓慢**
   ```bash
   # 增加工作线程
   go test -tags fuzz ./obs -fuzz=. -parallel=4

   # 或减少测试时长
   go test -tags fuzz ./obs -fuzz=. -fuzztime=10s
   ```

3. **内存不足**
   ```bash
   # 限制输入大小
   go test -tags fuzz ./obs -fuzz=. -fuzztime=30s
   ```

4. **测试超时**
   ```bash
   # 增加测试时长
   go test -tags fuzz ./obs -fuzz=. -fuzztime=60s
   ```

### 调试技巧

1. **查看崩溃文件**
   ```bash
   ls fuzz/crash/
   cat fuzz/crash/testdata/*
   ```

2. **使用最小化输入**
   ```bash
   # 限制输入大小快速定位问题
   go test -tags fuzz ./obs -fuzz=FuzzFunction -fuzztime=5s
   ```

3. **分析种子文件**
   ```bash
   # 查看生成的种子
   ls fuzz/corpus/FuzzFunction/
   ```

## 最佳实践

### 1. 输入限制

```go
f.Fuzz(func(t *testing.T, input string) {
    // 防止输入过大
    if len(input) > 1024*1024 {
        t.Skip("输入过大，跳过")
    }

    // 防止测试运行时间过长
    if len(f.Fuzzing()) > 10000 {
        t.Skip("跳过长时间运行")
    }

    // 执行测试...
})
```

### 2. 错误处理

```go
f.Fuzz(func(t *testing.T, input interface{}) {
    result, err := processInput(input)
    if err != nil {
        // 记录错误但不崩溃
        t.Logf("处理失败: %v", err)
        return
    }

    // 验证结果
    validateResult(result)
})
```

### 3. 性能监控

```go
func FuzzWithPerfCheck(f *testing.F) {
    f.Fuzz(func(t *testing.T, input string) {
        // 内存监控
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        if m.Alloc > 100*1024*1024 {
            t.Skip("内存使用过高")
            return
        }

        // 执行测试...
    })
}
```

## 下一步

1. 阅读完整文档：[skill.md](skill.md)
2. 查看更多示例：[fuzz_test.go.tmpl](./templates/fuzz_test.go.tmpl)
3. 使用技能生成测试：`/go-sdk-fuzz`