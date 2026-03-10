# go-sdk-perf

## 技能概述

OBS SDK Go性能测试编写指南。本技能指导用户如何为华为云OBS SDK编写轻量级和深度性能测试，用于评估SDK的性能表现、建立性能基线、检测性能退化。

## 使用场景

- 评估SDK性能表现（并发、带宽、资源占用）
- 建立性能基线
- 检测性能退化
- 资源监控和性能调优

## 核心功能

### 1. 性能测试定义

区分两种性能测试类型：
- **轻量级测试**：小文件（1MB）、低并发（10）、短时长（1s）
- **深度测试**：大文件（100MB-1GB）、高并发（100-1000）、长时长（30s）

### 2. 性能基线建立

指导如何建立和对比性能基线：
- 初始基准值设定
- 定期性能对比
- 性能趋势分析

### 3. 性能退化检测

提供系统性的性能退化检测方法：
- 自动对比历史数据
- 生成性能报告
- 告警阈值设置

### 4. 资源监控方法

监控关键资源指标：
- CPU使用率
- 内存占用
- 带宽使用
- 垃圾回收情况

## 使用方法

### 基本用法

```bash
/go-sdk-perf
```

### 带参数使用

```bash
/go-sdk-perf --type=light
/go-sdk-perf --type=deep
/go-sdk-perf --concurrency=100
/go-sdk-perf --duration=30s
/go-sdk-perf --output=./perf-report/
```

### 技能输出

1. **性能测试文件**：
   ```
   obs/*_benchmark_test.go
   ```

2. **配置模板**：
   - 轻量级测试配置
   - 深度测试配置
   - 资源监控设置

3. **基线数据**：
   - 初始性能基线
   - 历史性能数据
   - 性能对比分析

4. **监控工具**：
   - 资源监控脚本
   - 性能报告生成器

## 测试类型说明

### 轻量级测试（BenchmarkLight）

**特点**：
- 小文件（1MB）
- 低并发（10）
- 短时长（1s）
- 定期运行（CI/CD）

**运行命令**：
```bash
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem
```

### 深度测试（BenchmarkDeep）

**特点**：
- 大文件（100MB-1GB）
- 高并发（100-1000）
- 长时长（30s）
- 定期手动运行

**运行命令**：
```bash
go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem
```

## 输出示例

### 生成的性能测试文件示例

```go
//go:build perf

package obs

import (
	"bytes"
	"testing"
	"sync"
	"time"
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
	content := bytes.Repeat([]byte("test"), 256*1024) // 1MB
	bucket := cfg.GetTestBucket()
	objectKey := "small-file-test.dat"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := objectKey + "-" + time.Now().Format("20060102150405")
			obsClient.PutObject(bucket, testKey, content, nil)
		}
	})
}

// BenchmarkLight_GetObject_SmallFile 小文件下载性能测试
func BenchmarkLight_GetObject_SmallFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 准备测试对象
	bucket := cfg.GetTestBucket()
	objectKey := "small-file-download.dat"
	content := bytes.Repeat([]byte("test"), 256*1024) // 1MB

	// 先上传一个测试文件
	obsClient.PutObject(bucket, objectKey, content, nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			obsClient.GetObject(input)
		}
	})
}

// BenchmarkDeep_PutObject_LargeFile 大文件上传性能测试
func BenchmarkDeep_PutObject_LargeFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 测试数据：100MB文件
	content := bytes.Repeat([]byte("large-file-content"), 100*1024*1024) // 100MB
	bucket := cfg.GetTestBucket()
	objectKey := "large-file-upload.dat"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; i < b.N; i++ {
			testKey := objectKey + "-" + time.Now().Format("20060102150405")
			obsClient.PutObject(bucket, testKey, content, nil)
		}
	})
}

// BenchmarkDeep_GetObject_LargeFile 大文件下载性能测试
func BenchmarkDeep_GetObject_LargeFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 准备测试对象
	bucket := cfg.GetTestBucket()
	objectKey := "large-file-download.dat"
	content := bytes.Repeat([]byte("large-file-content"), 100*1024*1024) // 100MB

	// 先上传一个测试文件
	obsClient.PutObject(bucket, objectKey, content, nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; i < b.N; i++ {
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			obsClient.GetObject(input)
		}
	})
}

// BenchmarkConcurrent_PutObject_ConcurrentUpload 并发上传性能测试
func BenchmarkConcurrent_PutObject_ConcurrentUpload(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := bytes.Repeat([]byte("concurrent-test"), 256*1024) // 1MB
	bucket := cfg.GetTestBucket()

	// 不同的并发级别
	concurrencyLevels := []int{10, 50, 100, 500}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					testKey := fmt.Sprintf("concurrent-%d-%s", concurrency, time.Now().Format("20060102150405"))
					obsClient.PutObject(bucket, testKey, content, nil)
				}
			})
		})
	}
}

// BenchmarkMemory_PutObject_MemoryUsage 内存使用性能测试
func BenchmarkMemory_PutObject_MemoryUsage(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, _ := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := bytes.Repeat([]byte("memory-test"), 256*1024) // 1MB
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var m runtime.MemStats
		for pb.Next() {
			testKey := "memory-" + time.Now().Format("20060102150405")
			obsClient.PutObject(bucket, testKey, content, nil)

			// 记录内存使用
			runtime.ReadMemStats(&m)
			b.Logf("Alloc: %d MB", m.Alloc/1024/1024)
		}
	})
}
```

### 性能报告示例

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "benchmark_name": "BenchmarkLight_PutObject_SmallFile",
  "results": [
    {
      "concurrency": 1,
      "ns_per_op": 45023456,
      "mb_per_s": 22.5,
      "allocs_per_op": 25,
      "bytes_per_op": 1024
    },
    {
      "concurrency": 10,
      "ns_per_op": 5678901,
      "mb_per_s": 178.3,
      "allocs_per_op": 28,
      "bytes_per_op": 1024
    },
    {
      "concurrency": 100,
      "ns_per_op": 890123,
      "mb_per_s": 1135.2,
      "allocs_per_op": 30,
      "bytes_per_op": 1024
    }
  ],
  "performance_baseline": {
    "expected_mb_per_s": 150,
    "threshold": 0.9,
    "status": "ok"
  },
  "resource_usage": {
    "avg_cpu_percent": 45.2,
    "max_memory_mb": 512,
    "garbage_collections": 25
  }
}
```

## 运行性能测试

### 基本运行

```bash
# 运行轻量级性能测试
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem

# 运行深度性能测试
go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem

# 运行所有性能测试
go test -tags perf ./obs -bench=. -benchtime=1s
```

### 使用不同的并发级别

```bash
# 设置并行度
go test -tags perf ./obs -bench=BenchmarkConcurrent -benchtime=10s -parallel=100

# 运行特定并发级别
go test -tags perf ./obs -bench=BenchmarkConcurrent/PutObject_Concurrency-10 -benchtime=10s
```

### 保存测试结果

```bash
# 保存测试结果到文件
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -o benchmark.out

# 生成HTML报告
go tool benchmark -html benchmark.out > benchmark.html
```

## 性能监控

### 内存监控

```go
func BenchmarkMemory_Upload(b *testing.B) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    // 执行测试...

    // 记录内存使用
    runtime.ReadMemStats(&m)
    b.Logf("内存使用: %d MB", m.Alloc/1024/1024)
}
```

### CPU监控

```go
func BenchmarkCPU_Upload(b *testing.B) {
    startTime := time.Now()

    // 执行测试...

    elapsed := time.Since(startTime)
    b.Logf("CPU时间: %v", elapsed)
}
```

### 垃圾回收监控

```go
func BenchmarkGC_Upload(b *testing.B) {
    var m runtime.MemStats

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 执行测试...

        // 记录垃圾回收次数
        runtime.ReadMemStats(&m)
        if m.NumGC > 0 {
            b.Logf("GC次数: %d", m.NumGC)
        }
    }
}
```

## 最佳实践

### 1. 性能基线设置

```go
var performanceBaseline = map[string]float64{
    "BenchmarkLight_PutObject_SmallFile": 150.0, // MB/s
    "BenchmarkLight_GetObject_SmallFile": 200.0, // MB/s
    "BenchmarkDeep_PutObject_LargeFile":  50.0,  // MB/s
}

func checkPerformance(b *testing.B, name string, current float64) {
    baseline := performanceBaseline[name]
    threshold := baseline * 0.9 // 90%阈值

    if current < threshold {
        b.Errorf("性能退化! 当前: %.2f MB/s, 基线: %.2f MB/s", current, baseline)
    }
}
```

### 2. 性能退化检测

```go
func BenchmarkWithRegressionCheck(b *testing.B) {
    // 运行测试
    b.Run("Upload", func(b *testing.B) {
        result := testing.Benchmark(func(b *testing.B) {
            // 上传测试代码
        })

        // 检查性能退化
        checkPerformance(b, "Upload", calculateThroughput(result))
    })
}
```

### 3. 资源监控

```go
func BenchmarkWithResourceMonitoring(b *testing.B) {
    var m runtime.MemStats
    startTime := time.Now()

    // 执行测试
    b.ResetTimer()
    b.N = 1000

    for i := 0; i < b.N; i++ {
        // 测试代码
        runtime.ReadMemStats(&m)
    }

    // 分析资源使用
    elapsed := time.Since(startTime)
    b.Logf("总时间: %v", elapsed)
    b.Logf("最大内存: %d MB", m.Alloc/1024/1024)
    b.Logf("GC次数: %d", m.NumGC)
}
```

## 故障排除

### 常见问题

1. **性能测试不稳定**
   ```bash
   # 增加测试时长
   go test -tags perf ./obs -bench=. -benchtime=10s
   ```

2. **内存使用过高**
   ```bash
   # 添加内存监控
   go test -tags perf ./obs -bench=. -benchmem
   ```

3. **测试运行缓慢**
   ```bash
   # 减少测试数据大小
   # 或者使用更快的存储
   ```

4. **并发问题**
   ```bash
   # 调整并行度
   go test -tags perf ./obs -bench=. -parallel=4
   ```

### 调试技巧

1. **分析单个测试**
   ```bash
   go test -tags perf ./obs -bench=BenchmarkLight -count=1
   ```

2. **查看详细输出**
   ```bash
   go test -tags perf ./obs -bench=. -v
   ```

3. **生成HTML报告**
   ```bash
   go test -tags perf ./obs -bench=. > bench.out
   go tool benchmark -html bench.out > report.html
   ```

## 注意事项

1. **环境变量**：确保设置了正确的环境变量
2. **网络条件**：网络会影响测试结果
3. **硬件资源**：确保有足够的硬件资源
4. **测试数据**：使用一致的数据进行测试

## 技能版本

- 版本：1.1.0
- 最后更新：2024-01-09
- 兼容性：OBS SDK v3.x
- Go版本：1.18+
- **支持技能调用**: 可被go-sdk-dev-task技能调用和协调

## 技能调用接口

本技能可以被go-sdk-dev-task技能调用，以实现开发任务的自动化性能测试。

### 调用方式

go-sdk-dev-task可以通过以下方式调用本技能：

```bash
# 基本调用（使用默认参数）
/go-sdk-dev-task --type=optimization

# 指定性能测试类型
/go-sdk-dev-task --perf-types=light,deep

# 指定测试时长
/go-sdk-dev-task --benchtime=30s
```

### 技能协调机制

1. **性能策略一致性**
   - 根据任务类型选择合适的性能测试策略
   - 轻量级：小文件、低并发、短时长
   - 深度：大文件、高并发、长时长
   - 避免重复测试相同功能

2. **基线管理协调**
   - 统一性能基线存储格式
   - 避免基线数据冲突
   - 协调基线更新时机
   - 确保基线数据可追溯

3. **资源管理协调**
   - 使用统一的测试资源目录
   - 避免测试资源冲突
   - 协调测试数据清理
   - 确保测试环境清洁

4. **进度同步**
   - 向go-sdk-dev-task报告性能测试生成进度
   - 汇总性能测试结果
   - 识别性能瓶颈
   - 协调测试执行时机

### 技能调用建议

go-sdk-dev-task在以下情况会调用本技能：

1. **新功能开发**: 生成基础性能测试
2. **Bug修复**: 生成性能回归测试
3. **性能优化**: 生成优化前后对比测试
4. **代码重构**: 确保无性能退化
5. **集成测试前**: 验证性能基线