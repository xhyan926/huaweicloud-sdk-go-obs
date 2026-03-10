//go:build perf

package obs

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/config"
)

// 性能基线配置
var concurrentPerformanceBaseline = map[string]float64{
	"BenchmarkConcurrent_MixedOperations_10":    1000.0, // ops/s
	"BenchmarkConcurrent_StressTest_100":      800.0,  // ops/s
	"BenchmarkConcurrent_ConnectionPool_50":   2000.0, // ops/s
}

// BenchmarkConcurrent_MixedOperations 混合操作并发性能测试
func BenchmarkConcurrent_MixedOperations(b *testing.B) {
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

	content := bytes.Repeat([]byte("A"), 1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			operationType := time.Now().UnixNano() % 3
			objectKey := fmt.Sprintf("mixed-op-%d-%d", operationType, time.Now().UnixNano())

			switch operationType {
			case 0: // 上传
				input := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    objectKey,
					Body:   bytes.NewReader(content),
				}

				if err := obsClient.PutObject(input); err != nil {
					b.Errorf("并发上传失败: %v", err)
				}

				// 清理
				go func(key string) {
					deleteInput := &obs.DeleteObjectInput{
						Bucket: bucket,
						Key:    key,
					}
					obsClient.DeleteObject(deleteInput)
				}(objectKey)

			case 1: // 下载
				putInput := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    objectKey,
					Body:   bytes.NewReader(content),
				}

				// 先上传
				if err := obsClient.PutObject(putInput); err != nil {
					break
				}

				getInput := &obs.GetObjectInput{
					Bucket: bucket,
					Key:    objectKey,
				}

				output, err := obsClient.GetObject(getInput)
				if err != nil {
					b.Errorf("并发下载失败: %v", err)
					return
				}

				output.Body.Close()

				// 清理
				go func(key string) {
					deleteInput := &obs.DeleteObjectInput{
						Bucket: bucket,
						Key:    key,
					}
					obsClient.DeleteObject(deleteInput)
				}(objectKey)

			case 2: // 元数据查询
				getInput := &obs.GetObjectMetadataInput{
					Bucket: bucket,
					Key:    objectKey,
				}

				_, err := obsClient.GetObjectMetadata(getInput)
				if err != nil {
					// 对象不存在是预期的
				}
			}
		}
	})

	checkPerformanceDegradation(b, "BenchmarkConcurrent_MixedOperations_10",
		calculateThroughput(b.Result, 1.0))
}

// BenchmarkConcurrent_StressTest 高并发压力测试
func BenchmarkConcurrent_StressTest(b *testing.B) {
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

	content := bytes.Repeat([]byte("B"), 512*1024) // 512KB
	bucket := cfg.GetTestBucket()

	// 高并发水平
	concurrencyLevels := []int{50, 100, 200, 500}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					objectKey := fmt.Sprintf("stress-test-%d-%d", concurrency, time.Now().UnixNano())

					input := &obs.PutObjectInput{
						Bucket: bucket,
						Key:    objectKey,
						Body:   bytes.NewReader(content),
					}

					if err := obsClient.PutObject(input); err != nil {
						return
					}

					// 异步清理
					go func(key string) {
						deleteInput := &obs.DeleteObjectInput{
							Bucket: bucket,
							Key:    key,
						}
						obsClient.DeleteObject(deleteInput)
					}(objectKey)
				}
			})

			checkPerformanceDegradation(b,
				fmt.Sprintf("BenchmarkConcurrent_StressTest_Concurrency-%d", concurrency),
				calculateThroughput(b.Result, 0.5))
		})
	}
}

// BenchmarkConcurrent_ConnectionPool 连接池性能测试
func BenchmarkConcurrent_ConnectionPool(b *testing.B) {
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
		obs.WithConnectTimeout(10*time.Second),
		obs.WithSocketTimeout(10*time.Second),
		obs.WithMaxRedirectCount(5),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	content := bytes.Repeat([]byte("C"), 2*1024*1024) // 2MB
	bucket := cfg.GetTestBucket()

	// 准备测试对象
	objectKeys := make([]string, 20)
	for i := 0; i < 20; i++ {
		objectKeys[i] = fmt.Sprintf("connpool-%d-%d", i, time.Now().UnixNano())

		input := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKeys[i],
			Body:   bytes.NewReader(content),
		}

		if _, err := obsClient.PutObject(input); err != nil {
			b.Fatalf("准备测试对象 %d 失败: %v", i, err)
		}

		// 清理函数
		defer func(key string) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    key,
			}
			obsClient.DeleteObject(deleteInput)
		}(objectKeys[i])
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			keyIndex := int(time.Now().UnixNano()) % len(objectKeys)
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKeys[keyIndex],
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				return
			}

			output.Body.Close()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkConcurrent_ConnectionPool_20",
		calculateThroughput(b.Result, 2.0))
}

// BenchmarkConcurrent_ResourceRace 资源竞争性能测试
func BenchmarkConcurrent_ResourceRace(b *testing.B) {
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

	content := bytes.Repeat([]byte("D"), 1*1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	var counter int64
	var mutex sync.Mutex

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			objectKey := fmt.Sprintf("race-test-%d", time.Now().UnixNano())

			// 使用互斥锁模拟资源竞争
			mutex.Lock()
			counter++
			testCounter := counter
			mutex.Unlock()

			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    fmt.Sprintf("race-test-%d-%d", testCounter%10, time.Now().UnixNano()),
				Body:   bytes.NewReader(content),
			}

			if err := obsClient.PutObject(input); err != nil {
				return
			}

			// 异步清理
			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(objectKey)
		}
	})

	checkPerformanceDegradation(b, "BenchmarkConcurrent_ResourceRace",
		calculateThroughput(b.Result, 1.0))
}

// BenchmarkConcurrent_ListObjects 并发列表操作性能测试
func BenchmarkConcurrent_ListObjects(b *testing.B) {
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

	content := bytes.Repeat([]byte("E"), 1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	// 准备测试对象
	objectKeys := make([]string, 100)
	for i := 0; i < 100; i++ {
		objectKeys[i] = fmt.Sprintf("list-test-%d", i)

		input := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKeys[i],
			Body:   bytes.NewReader(content),
		}

		if _, err := obsClient.PutObject(input); err != nil {
			b.Fatalf("准备测试对象 %d 失败: %v", i, err)
		}

		// 清理函数
		defer func(key string) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    key,
			}
			obsClient.DeleteObject(deleteInput)
		}(objectKeys[i])
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			input := &obs.ListObjectsInput{
				Bucket: bucket,
				MaxKeys: 100,
			}

			_, err := obsClient.ListObjects(input)
			if err != nil {
				b.Errorf("列表操作失败: %v", err)
			}
		}
	})
}

// BenchmarkConcurrent_MemoryUsage 并发操作内存使用测试
func BenchmarkConcurrent_MemoryUsage(b *testing.B) {
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

	content := bytes.Repeat([]byte("F"), 2*1024*1024) // 2MB
	bucket := cfg.GetTestBucket()

	var m runtime.MemStats

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runtime.ReadMemStats(&m)
			memBefore := m.Alloc

			objectKey := fmt.Sprintf("memory-concurrent-%d", time.Now().UnixNano())

			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    objectKey,
				Body:   bytes.NewReader(content),
			}

			if err := obsClient.PutObject(input); err != nil {
				return
			}

			runtime.ReadMemStats(&m)
			memAfter := m.Alloc
			memDelta := memAfter - memBefore

			b.ReportMetric(float64(memDelta)/1024/1024, "Alloc_MB")

			// 清理
			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(objectKey)
		}
	})

	runtime.ReadMemStats(&m)
	b.Logf("并发操作内存使用: Alloc=%d MB, HeapAlloc=%d MB, Sys=%d MB, NumGC=%d",
		m.Alloc/1024/1024, m.HeapAlloc/1024/1024,
		m.Sys/1024/1024, m.NumGC)
}

// BenchmarkConcurrent_ErrorHandling 并发错误处理性能测试
func BenchmarkConcurrent_ErrorHandling(b *testing.B) {
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
			// 测试访问不存在的对象（预期错误）
			objectKey := fmt.Sprintf("error-handling-%d", time.Now().UnixNano())

			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			_, err := obsClient.GetObject(input)
			if err == nil {
				// 成功是意外的，但继续测试
			}
		}
	})
}

// BenchmarkConcurrent_Latency 并发操作延迟测试
func BenchmarkConcurrent_Latency(b *testing.B) {
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
		obs.WithConnectTimeout(5*time.Second),
		obs.WithSocketTimeout(5*time.Second),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	content := bytes.Repeat([]byte("G"), 1*1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	// 准备测试对象
	objectKey := fmt.Sprintf("latency-test-%d", time.Now().UnixNano())

	input := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}

	if _, err := obsClient.PutObject(input); err != nil {
		b.Fatalf("准备测试对象失败: %v", err)
	}

	// 清理函数
	defer func() {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		obsClient.DeleteObject(deleteInput)
	}()

	var latencies []time.Duration

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()

			getInput := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(getInput)
			if err != nil {
				return
			}

			output.Body.Close()

			latency := time.Since(start)

			// 记录延迟
			latencyLock := sync.Mutex{}
			latencyLock.Lock()
			latencies = append(latencies, latency)
			if len(latencies) > 1000 {
				latencies = latencies[len(latencies)-100:]
			}
			latencyLock.Unlock()
		}
	})

	// 计算延迟统计
	if len(latencies) > 0 {
		var sum time.Duration
		min := latencies[0]
		max := latencies[0]

		for _, latency := range latencies {
			sum += latency
			if latency < min {
				min = latency
			}
			if latency > max {
				max = latency
			}
		}

		avg := sum / time.Duration(len(latencies))

		b.Logf("=== 并发操作延迟统计 ===")
		b.Logf("总请求数: %d", len(latencies))
		b.Logf("平均延迟: %v", avg)
		b.Logf("最小延迟: %v", min)
		b.Logf("最大延迟: %v", max)
	}
}

// BenchmarkConcurrent_ThroughputScaling 吞吐量扩展性能测试
func BenchmarkConcurrent_ThroughputScaling(b *testing.B) {
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

	content := bytes.Repeat([]byte("H"), 1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	// 不同的并发水平
	concurrencyLevels := []int{1, 5, 10, 25, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					objectKey := fmt.Sprintf("scaling-%d-%d", concurrency, time.Now().UnixNano())

					input := &obs.PutObjectInput{
						Bucket: bucket,
						Key:    objectKey,
						Body:   bytes.NewReader(content),
					}

					if err := obsClient.PutObject(input); err != nil {
						return
					}

					// 清理
					go func(key string) {
						deleteInput := &obs.DeleteObjectInput{
							Bucket: bucket,
							Key:    key,
						}
						obsClient.DeleteObject(deleteInput)
					}(objectKey)
				}
			})

			checkPerformanceDegradation(b,
				fmt.Sprintf("BenchmarkConcurrent_ThroughputScaling_Concurrency-%d", concurrency),
				calculateThroughput(b.Result, 1.0))
		})
	}
}
