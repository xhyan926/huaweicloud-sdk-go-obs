# OBS SDK Go 性能测试指南

## 快速开始

### 1. 环境准备

确保已安装Go 1.18+和必要的性能分析工具：

```bash
go version  # 需要 >= 1.18

# 安装性能分析工具
go install golang.org/x/perf/cmd/benchstat@latest
```

### 2. 运行性能测试

```bash
# 进入项目目录
cd /path/to/huaweicloud-sdk-go-obs

# 运行轻量级性能测试
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem

# 运行深度性能测试
go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem

# 运行所有性能测试
go test -tags perf ./obs -bench=. -benchtime=1s
```

## 性能测试类型

### 轻量级测试（定期运行）

**特点**：
- 小文件（1MB）
- 低并发（10）
- 短时长（1s）
- 用于日常性能监控

```bash
# 运行轻量级测试
make test-perf-light
# 或
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem
```

### 深度测试（手动运行）

**特点**：
- 大文件（100MB-1GB）
- 高并发（100-1000）
- 长时长（30s）
- 用于详细性能分析

```bash
# 运行深度测试
make test-perf-deep
# 或
go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem
```

## 使用方法

### 基本性能测试

```go
//go:build perf

package obs

import (
	"testing"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/config"
)

// BenchmarkLight_PutObject_SmallFile 小文件上传性能测试
func BenchmarkLight_PutObject_SmallFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := make([]byte, 1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := "small-file-" + time.Now().Format("20060102150405")
			obsClient.PutObject(bucket, testKey, string(content), nil)
		}
	})
}

// BenchmarkDeep_GetObject_LargeFile 大文件下载性能测试
func BenchmarkDeep_GetObject_LargeFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 准备测试数据：100MB文件
	content := make([]byte, 100*1024*1024) // 100MB
	bucket := cfg.GetTestBucket()

	// 先上传一个测试文件
	testKey := "large-file-test.dat"
	obsClient.PutObject(bucket, testKey, string(content), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := &obs.GetObjectInput{
			Bucket: bucket,
			Key:    testKey,
		}
		obsClient.GetObject(input)
	}
}
```

### 高级性能测试

#### 1. 并发性能测试

```go
// BenchmarkConcurrent_PutObject_ConcurrentUpload 并发上传性能测试
func BenchmarkConcurrent_PutObject_ConcurrentUpload(b *testing.B) {
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := make([]byte, 1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	// 不同的并发级别
	concurrencyLevels := []int{10, 50, 100, 500}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					testKey := fmt.Sprintf("concurrent-%d", time.Now().UnixNano())
					obsClient.PutObject(bucket, testKey, string(content), nil)
				}
			})
		})
	}
}
```

#### 2. 内存性能测试

```go
// BenchmarkMemory_PutObject_MemoryUsage 内存使用性能测试
func BenchmarkMemory_PutObject_MemoryUsage(b *testing.B) {
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := make([]byte, 1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	var m runtime.MemStats
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := "memory-" + time.Now().Format("20060102150405")
			obsClient.PutObject(bucket, testKey, string(content), nil)

			// 记录内存使用
			runtime.ReadMemStats(&m)
			b.Logf("Alloc: %d MB", m.Alloc/1024/1024)
		}
	})
}
```

#### 3. 资源监控测试

```go
// BenchmarkWithResourceMonitoring 带资源监控的性能测试
func BenchmarkWithResourceMonitoring(b *testing.B) {
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := make([]byte, 1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	var m runtime.MemStats
	startTime := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testKey := "resource-" + time.Now().Format("20060102150405")
		obsClient.PutObject(bucket, testKey, string(content), nil)

		// 记录资源使用
		runtime.ReadMemStats(&m)
	}

	// 分析资源使用
	elapsed := time.Since(startTime)
	b.Logf("总时间: %v", elapsed)
	b.Logf("最大内存: %d MB", m.Alloc/1024/1024)
	b.Logf("GC次数: %d", m.NumGC)
}
```

## 运行指南

### 基本命令

```bash
# 运行轻量级测试
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s

# 运行深度测试
go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s

# 运行所有性能测试
go test -tags perf ./obs -bench=.

# 保存测试结果
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -o bench.out

# 生成HTML报告
go tool benchmark -html bench.out > report.html
```

### 使用性能分析工具

```bash
# 使用benchstat对比性能
go test -tags perf ./obs -bench=BenchmarkLight > bench1.out
go test -tags perf ./obs -bench=BenchmarkLight > bench2.out
benchstat bench1.out bench2.out

# 持续性能监控
while true; do
    go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s
    sleep 3600
done
```

## 性能指标分析

### 关键指标

1. **ns/op** - 每次操作的纳秒数
2. **MB/s** - 每秒传输的兆字节数
3. **allocs/op** - 每次操作的内存分配次数
4. **bytes/op** - 每次操作分配的字节数

### 性能基线

| 测试名称 | 基准值 | 阈值 |
|---------|--------|------|
| BenchmarkLight_PutObject_SmallFile | 150 MB/s | 90% |
| BenchmarkLight_GetObject_SmallFile | 200 MB/s | 90% |
| BenchmarkDeep_PutObject_LargeFile | 50 MB/s | 90% |
| BenchmarkDeep_GetObject_LargeFile | 100 MB/s | 90% |

### 性能退化检测

```bash
# 使用benchstat检测性能退化
benchstat baseline.out current.out

# 输出示例
name             old time/op  new time/op  delta
BenchmarkLight_PutObject_SmallFile-8  45.2ms ± 2%  52.1ms ± 3%  +15.25%  (p=0.008 n=5+5)
```

## 故障排除

### 常见问题

1. **性能测试不稳定**
   ```bash
   # 增加测试时长
   go test -tags perf ./obs -bench=. -benchtime=10s

   # 使用 -count 参数多次运行
   go test -tags perf ./obs -bench=BenchmarkLight -count=5
   ```

2. **内存使用过高**
   ```bash
   # 检查内存分配
   go test -tags perf ./obs -bench=. -benchmem

   # 使用内存分析工具
   go test -tags perf ./obs -bench=. -memprofile=mem.out
   go tool pprof mem.out
   ```

3. **网络影响测试结果**
   ```bash
   # 使用Mock服务器减少网络影响
   export OBS_MOCK_ENABLED=true
   go test -tags perf ./obs -bench=.
   ```

4. **测试数据不一致**
   ```bash
   # 使用固定的测试数据
   const fixedContent = string(bytes.Repeat([]byte('x'), 1024*1024))
   ```

### 调试技巧

1. **分析单个测试**
   ```bash
   # 运行特定测试
   go test -tags perf ./obs -bench=BenchmarkLight_PutObject_SmallFile

   # 使用 -v 查看详细信息
   go test -tags perf ./obs -bench=. -v
   ```

2. **CPU性能分析**
   ```bash
   # 生成CPU profile
   go test -tags perf ./obs -bench=BenchmarkLight -cpuprofile=cpu.out
   go tool pprof cpu.out

   # 生成火焰图
   go tool pprof -http=:8080 cpu.out
   ```

3. **内存性能分析**
   ```bash
   # 生成内存profile
   go test -tags perf ./obs -bench=BenchmarkLight -memprofile=mem.out
   go tool pprof mem.out
   ```

## 最佳实践

### 1. 建立性能基线

```bash
# 第一次运行时建立基线
go test -tags perf ./obs -bench=. > baseline.out
```

### 2. 定期性能测试

```bash
# 每天自动运行性能测试
0 2 * * * cd /path/to/project && go test -tags perf ./obs -bench=BenchmarkLight > $(date +%Y%m%d)-bench.out
```

### 3. 性能告警

```bash
# 使用benchstat设置告警
if benchstat baseline.out current.out | grep -q "+"; then
    echo "性能退化警告!"
    # 发送告警邮件或通知
fi
```

### 4. 测试环境标准化

```bash
# 使用相同的硬件和网络环境
# 在Docker容器中运行测试
# 记录测试环境信息
```

## 下一步

1. 阅读完整文档：[skill.md](skill.md)
2. 查看更多示例：[benchmark_test.go.tmpl](./templates/benchmark_test.go.tmpl)
3. 使用技能生成测试：`/go-sdk-perf`