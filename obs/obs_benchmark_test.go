//go:build perf

package obs

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/config"
)

// 性能基准配置
var performanceBaseline = map[string]float64{
	"BenchmarkLight_PutObject_SmallFile":     150.0, // MB/s
	"BenchmarkLight_GetObject_SmallFile":     200.0, // MB/s
	"BenchmarkLight_DeleteObject_SmallFile":   1000.0, // ops/s
	"BenchmarkLight_ListObjects_SmallBucket": 500.0,  // ops/s
}

// 性能退化阈值
const performanceThreshold = 0.9 // 90%

// BenchmarkLight_PutObject_SmallFile 小文件上传性能测试
func BenchmarkLight_PutObject_SmallFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := bytes.Repeat([]byte("test"), 256*1024) // 1MB
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("perf-upload-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("上传失败: %v", err)
			}

			// 异步清理
			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})
}

// BenchmarkLight_GetObject_SmallFile 小文件下载性能测试
func BenchmarkLight_GetObject_SmallFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	// 准备测试对象
	bucket := cfg.GetTestBucket()
	objectKey := fmt.Sprintf("perf-download-%d", time.Now().UnixNano())
	content := bytes.Repeat([]byte("test"), 256*1024) // 1MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if err := obsClient.PutObject(putInput); err != nil {
		b.Fatalf("准备测试对象失败: %v", err)
	}
	defer func() {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		obsClient.DeleteObject(deleteInput)
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getInput := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(getInput)
			if err != nil {
				b.Errorf("下载失败: %v", err)
				continue
			}

			// 读取并关闭响应
			buf := make([]byte, len(content))
			_, _ = output.Body.Read(buf)
			output.Body.Close()
		}
	})
}

// BenchmarkLight_DeleteObject_SmallFile 小文件删除性能测试
func BenchmarkLight_DeleteObject_SmallFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	// 准备测试对象
	bucket := cfg.GetTestBucket()
	testKeys := make([]string, 100) // 预创建100个对象
	for i := 0; i < 100; i++ {
		testKeys[i] = fmt.Sprintf("perf-delete-%d-%d", i, time.Now().UnixNano())
		content := bytes.Repeat([]byte("test"), 1024) // 1KB
		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    testKeys[i],
			Body:   bytes.NewReader(content),
		}
		if err := obsClient.PutObject(putInput); err != nil {
			b.Fatalf("准备测试对象失败: %v", err)
		}
	}

	// 清理函数
	defer func() {
		var wg sync.WaitGroup
		for _, key := range testKeys {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    k,
				}
				obsClient.DeleteObject(deleteInput)
			}(key)
		}
		wg.Wait()
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		keyIndex := 0
		for pb.Next() {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    testKeys[keyIndex%len(testKeys)],
			}

			err := obsClient.DeleteObject(deleteInput)
			if err != nil {
				b.Errorf("删除失败: %v", err)
			}

			keyIndex++
		}
	})
}

// BenchmarkLight_ListObjects_SmallBucket 小桶列表性能测试
func BenchmarkLight_ListObjects_SmallBucket(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			listInput := &obs.ListObjectsInput{
				Bucket: bucket,
				MaxKeys: 100,
			}

			_, err := obsClient.ListObjects(listInput)
			if err != nil {
				b.Errorf("列表失败: %v", err)
			}
		}
	})
}

// BenchmarkLight_HeadObject_SmallFile 小文件元数据获取性能测试
func BenchmarkLight_HeadObject_SmallFile(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	// 准备测试对象
	bucket := cfg.GetTestBucket()
	objectKey := fmt.Sprintf("perf-head-%d", time.Now().UnixNano())
	content := bytes.Repeat([]byte("test"), 1024) // 1KB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if err := obsClient.PutObject(putInput); err != nil {
		b.Fatalf("准备测试对象失败: %v", err)
	}
	defer func() {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		obsClient.DeleteObject(deleteInput)
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			headInput := &obs.GetObjectMetadataInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			_, err := obsClient.GetObjectMetadata(headInput)
			if err != nil {
				b.Errorf("获取元数据失败: %v", err)
			}
		}
	})
}

// BenchmarkConcurrent_PutObject_ConcurrentUpload 并发上传性能测试
func BenchmarkConcurrent_PutObject_ConcurrentUpload(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := bytes.Repeat([]byte("test"), 256*1024) // 1MB
	bucket := cfg.GetTestBucket()

	// 不同的并发级别
	concurrencyLevels := []int{1, 10, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					testKey := fmt.Sprintf("perf-concurrent-%d", time.Now().UnixNano())
					input := &obs.PutObjectInput{
						Bucket: bucket,
						Key:    testKey,
						Body:   bytes.NewReader(content),
					}

					err := obsClient.PutObject(input)
					if err != nil {
						b.Errorf("上传失败: %v", err)
					}

					// 异步清理
					go func(key string) {
						deleteInput := &obs.DeleteObjectInput{
							Bucket: bucket,
							Key:    key,
						}
						obsClient.DeleteObject(deleteInput)
					}(testKey)
				}
			})
		})
	}
}

// BenchmarkMemory_PutObject_MemoryUsage 内存使用性能测试
func BenchmarkMemory_PutObject_MemoryUsage(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := bytes.Repeat([]byte("test"), 256*1024) // 1MB
	bucket := cfg.GetTestBucket()

	var m runtime.MemStats
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("perf-memory-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("上传失败: %v", err)
			}

			// 记录内存使用
			runtime.ReadMemStats(&m)
			b.ReportMetric(float64(m.Alloc)/1024/1024, "Alloc_MB")

			// 异步清理
			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})
}

// BenchmarkWithResourceMonitoring 带资源监控的性能测试
func BenchmarkWithResourceMonitoring(b *testing.B) {
	// 创建测试客户端
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	// 测试数据：1MB文件
	content := bytes.Repeat([]byte("test"), 256*1024) // 1MB
	bucket := cfg.GetTestBucket()

	var m runtime.MemStats
	startTime := time.Now()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("perf-resource-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("上传失败: %v", err)
			}

			// 记录资源使用
			runtime.ReadMemStats(&m)

			// 异步清理
			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	// 分析资源使用
	elapsed := time.Since(startTime)
	b.Logf("=== 资源使用统计 ===")
	b.Logf("总时间: %v", elapsed)
	b.Logf("最大内存: %d MB", m.Alloc/1024/1024)
	b.Logf("GC次数: %d", m.NumGC)
}

// checkPerformanceDegradation 检查性能退化
func checkPerformanceDegradation(b *testing.B, name string, currentValue float64) {
	if baseline, exists := performanceBaseline[name]; exists {
		threshold := baseline * performanceThreshold
		if currentValue < threshold {
			b.Errorf("性能退化! 当前: %.2f, 基线: %.2f (阈值: %.1f%%)",
				currentValue, baseline, performanceThreshold*100)
		} else {
			b.Logf("性能检查通过: %.2f >= %.2f", currentValue, threshold)
		}
	}
}

// calculateThroughput 计算吞吐量 (MB/s)
func calculateThroughput(result testing.BenchmarkResult, fileSizeMB float64) float64 {
	if result.N <= 0 || result.T <= 0 {
		return 0
	}

	// 计算总数据量 (MB)
	totalData := float64(result.N) * fileSizeMB

	// 计算总时间 (秒)
	totalTime := float64(result.T) / float64(time.Second)

	if totalTime <= 0 {
		return 0
	}

	return totalData / totalTime
}

// calculateLatency 计算延迟 (ms)
func calculateLatency(result testing.BenchmarkResult) float64 {
	if result.N <= 0 {
		return 0
	}

	// 计算平均延迟 (纳秒)
	avgLatencyNs := float64(result.T) / float64(result.N)

	// 转换为毫秒
	return avgLatencyNs / float64(time.Millisecond)
}

// generatePerfReport 生成性能报告
func generatePerfReport(b *testing.B, testName string, result testing.BenchmarkResult, fileSizeMB float64) {
	throughput := calculateThroughput(result, fileSizeMB)
	latency := calculateLatency(result)

	// 检查性能退化
	checkPerformanceDegradation(b, testName, throughput)

	b.Logf("=== 性能报告 ===")
	b.Logf("测试名称: %s", testName)
	b.Logf("操作次数: %d", result.N)
	b.Logf("总时间: %v", result.T)
	b.Logf("吞吐量: %.2f MB/s", throughput)
	b.Logf("平均延迟: %.2f ms", latency)
	b.Logf("内存分配: %d bytes/op", result.MemBytes)
	b.Logf("内存分配次数: %d allocs/op", result.AllocsPerOp)
}
