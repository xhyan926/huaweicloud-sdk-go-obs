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
var uploadPerformanceBaseline = map[string]float64{
	"BenchmarkLight_PutObject_1MB":     150.0, // MB/s
	"BenchmarkDeep_PutObject_100MB":    50.0,  // MB/s
	"BenchmarkConcurrent_PutObject_100": 500.0, // ops/s
}

// BenchmarkLight_PutObject_1MB 轻量级小文件上传性能测试
func BenchmarkLight_PutObject_1MB(b *testing.B) {
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
	content := bytes.Repeat([]byte("A"), 1024*1024)
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("light-upload-%d", time.Now().UnixNano())
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

	// 检查性能退化
	checkPerformanceDegradation(b, "BenchmarkLight_PutObject_1MB",
		calculateThroughput(b.Result, 1.0))
}

// BenchmarkLight_PutObject_1MB_Metadata 带元数据的轻量级上传测试
func BenchmarkLight_PutObject_1MB_Metadata(b *testing.B) {
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

	content := bytes.Repeat([]byte("B"), 1024*1024)
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("light-upload-meta-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
				Metadata: map[string]string{
					"author":     "test-user",
					"upload-time": time.Now().Format(time.RFC3339),
					"test":       "metadata-test",
				},
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("带元数据上传失败: %v", err)
			}

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	checkPerformanceDegradation(b, "BenchmarkLight_PutObject_1MB_Metadata",
		calculateThroughput(b.Result, 1.0))
}

// BenchmarkLight_PutObject_10MB 中等大小文件上传测试
func BenchmarkLight_PutObject_10MB(b *testing.B) {
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

	content := bytes.Repeat([]byte("C"), 10*1024*1024)
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("light-upload-10mb-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("10MB文件上传失败: %v", err)
			}

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	checkPerformanceDegradation(b, "BenchmarkLight_PutObject_10MB",
		calculateThroughput(b.Result, 10.0))
}

// BenchmarkDeep_PutObject_100MB 深度大文件上传性能测试
func BenchmarkDeep_PutObject_100MB(b *testing.B) {
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

	// 测试数据：100MB文件
	content := bytes.Repeat([]byte("D"), 100*1024*1024)
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("deep-upload-100mb-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("100MB文件上传失败: %v", err)
			}

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	checkPerformanceDegradation(b, "BenchmarkDeep_PutObject_100MB",
		calculateThroughput(b.Result, 100.0))
}

// BenchmarkDeep_PutObject_1GB 超大文件上传性能测试
func BenchmarkDeep_PutObject_1GB(b *testing.B) {
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

	// 测试数据：1GB文件
	content := bytes.Repeat([]byte("E"), 1024*1024*1024)
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("deep-upload-1gb-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("1GB文件上传失败: %v", err)
			}

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	checkPerformanceDegradation(b, "BenchmarkDeep_PutObject_1GB",
		calculateThroughput(b.Result, 1024.0))
}

// BenchmarkConcurrent_PutObject 不同并发级别上传性能测试
func BenchmarkConcurrent_PutObject(b *testing.B) {
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

	content := bytes.Repeat([]byte("F"), 5*1024*1024)
	bucket := cfg.GetTestBucket()

	// 不同的并发级别
	concurrencyLevels := []int{1, 10, 50, 100, 500}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					testKey := fmt.Sprintf("concurrent-upload-%d-%d", concurrency, time.Now().UnixNano())
					input := &obs.PutObjectInput{
						Bucket: bucket,
						Key:    testKey,
						Body:   bytes.NewReader(content),
					}

					err := obsClient.PutObject(input)
					if err != nil {
						b.Errorf("并发上传失败: %v", err)
					}

					go func(key string) {
						deleteInput := &obs.DeleteObjectInput{
							Bucket: bucket,
							Key:    key,
						}
						obsClient.DeleteObject(deleteInput)
					}(testKey)
				}
			})

			checkPerformanceDegradation(b,
				fmt.Sprintf("BenchmarkConcurrent_PutObject_Concurrency-%d", concurrency),
				calculateThroughput(b.Result, 5.0))
		})
	}
}

// BenchmarkMemory_PutObject 内存使用性能测试
func BenchmarkMemory_PutObject(b *testing.B) {
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

	content := bytes.Repeat([]byte("G"), 10*1024*1024)
	bucket := cfg.GetTestBucket()

	var m runtime.MemStats
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("memory-upload-%d", time.Now().UnixNano())

			// 记录内存使用前
			runtime.ReadMemStats(&m)
			memBefore := m.Alloc

			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("内存监控上传失败: %v", err)
			}

			// 记录内存使用后
			runtime.ReadMemStats(&m)
			memAfter := m.Alloc
			memDelta := memAfter - memBefore

			// 记录内存使用
			b.ReportMetric(float64(memDelta)/1024/1024, "Alloc_MB")

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	// 输出内存统计信息
	runtime.ReadMemStats(&m)
	b.Logf("内存使用统计: Alloc=%d MB, HeapAlloc=%d MB, Sys=%d MB, NumGC=%d",
		m.Alloc/1024/1024, m.HeapAlloc/1024/1024, m.Sys/1024/1024, m.NumGC)
}

// BenchmarkUploadResourceMonitoring 上传资源监控测试
func BenchmarkUploadResourceMonitoring(b *testing.B) {
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

	content := bytes.Repeat([]byte("H"), 5*1024*1024)
	bucket := cfg.GetTestBucket()

	var m runtime.MemStats
	var totalGC uint32
	startTime := time.Now()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("resource-upload-%d", time.Now().UnixNano())

			// 记录资源使用
			runtime.ReadMemStats(&m)
			totalGC += m.NumGC - totalGC

			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("资源监控上传失败: %v", err)
			}

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
	runtime.ReadMemStats(&m)

	b.Logf("=== 上传资源使用统计 ===")
	b.Logf("总操作数: %d", b.N)
	b.Logf("总时间: %v", elapsed)
	b.Logf("平均操作时间: %.2f ms", float64(elapsed)/float64(b.N)/float64(time.Millisecond))
	b.Logf("最大内存: %d MB", m.Alloc/1024/1024)
	b.Logf("堆内存: %d MB", m.HeapAlloc/1024/1024)
	b.Logf("系统内存: %d MB", m.Sys/1024/1024)
	b.Logf("GC次数: %d", totalGC)
}

// BenchmarkUploadConnectionPooling 连接池性能测试
func BenchmarkUploadConnectionPooling(b *testing.B) {
	cfg := config.LoadTestConfig()
	obsClient, err := obs.New(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Endpoint,
		obs.WithSecurityToken(cfg.SecurityToken),
		obs.WithConnectTimeout(30*time.Second),
		obs.WithSocketTimeout(30*time.Second),
	)
	if err != nil {
		b.Fatalf("创建客户端失败: %v", err)
	}
	defer obsClient.Close()

	content := bytes.Repeat([]byte("I"), 2*1024*1024)
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("pooling-upload-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("连接池上传失败: %v", err)
			}

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	checkPerformanceDegradation(b, "BenchmarkUploadConnectionPooling",
		calculateThroughput(b.Result, 2.0))
}

// BenchmarkUploadWithContentType 带内容类型的上传性能测试
func BenchmarkUploadWithContentType(b *testing.B) {
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

	content := bytes.Repeat([]byte("J"), 2*1024*1024)
	bucket := cfg.GetTestBucket()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := fmt.Sprintf("content-type-upload-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket:      bucket,
				Key:         testKey,
				Body:        bytes.NewReader(content),
				ContentType: "application/octet-stream",
			}

			err := obsClient.PutObject(input)
			if err != nil {
				b.Errorf("带内容类型上传失败: %v", err)
			}

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)
		}
	})

	checkPerformanceDegradation(b, "BenchmarkUploadWithContentType",
		calculateThroughput(b.Result, 2.0))
}

// BenchmarkUploadDifferentSizes 不同文件大小上传性能测试
func BenchmarkUploadDifferentSizes(b *testing.B) {
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

	// 测试不同文件大小
	fileSizes := []struct {
		name  string
		size  int64
		sizeMB float64
	}{
		{"100KB", 100 * 1024, 0.1},
		{"1MB", 1024 * 1024, 1.0},
		{"10MB", 10 * 1024 * 1024, 10.0},
		{"50MB", 50 * 1024 * 1024, 50.0},
	}

	for _, test := range fileSizes {
		b.Run(test.name, func(b *testing.B) {
			content := bytes.Repeat([]byte("K"), int(test.size))

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					testKey := fmt.Sprintf("size-upload-%s-%d", test.name, time.Now().UnixNano())
					input := &obs.PutObjectInput{
						Bucket: bucket,
						Key:    testKey,
						Body:   bytes.NewReader(content),
					}

					err := obsClient.PutObject(input)
					if err != nil {
						b.Errorf("%s文件上传失败: %v", test.name, err)
					}

					go func(key string) {
						deleteInput := &obs.DeleteObjectInput{
							Bucket: bucket,
							Key:    key,
						}
						obsClient.DeleteObject(deleteInput)
					}(testKey)
				}
			})

			checkPerformanceDegradation(b,
				fmt.Sprintf("BenchmarkUploadDifferentSizes_%s", test.name),
				calculateThroughput(b.Result, test.sizeMB))
		})
	}
}

// BenchmarkUploadSequential 顺序上传性能测试
func BenchmarkUploadSequential(b *testing.B) {
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

	content := bytes.Repeat([]byte("L"), 1*1024*1024)
	bucket := cfg.GetTestBucket()

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			testKey := fmt.Sprintf("sequential-upload-%d", time.Now().UnixNano())
			input := &obs.PutObjectInput{
				Bucket: bucket,
				Key:    testKey,
				Body:   bytes.NewReader(content),
			}

			err := obsClient.PutObject(input)
			if err != nil {
				errors <- fmt.Errorf("顺序上传 %d 失败: %v", index, err)
				return
			}

			go func(key string) {
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    key,
				}
				obsClient.DeleteObject(deleteInput)
			}(testKey)

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	// 收集错误
	for err := range errors {
		if err != nil {
			b.Error(err)
		}
	}

	checkPerformanceDegradation(b, "BenchmarkUploadSequential",
		calculateThroughput(b.Result, 1.0))
}
