//go:build perf

package obs

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/config"
)

// 性能基线配置
var downloadPerformanceBaseline = map[string]float64{
	"BenchmarkLight_GetObject_1MB":      200.0, // MB/s
	"BenchmarkDeep_GetObject_100MB":     100.0, // MB/s
	"BenchmarkConcurrent_GetObject_100": 800.0, // ops/s
}

// BenchmarkLight_GetObject_1MB 轻量级小文件下载性能测试
func BenchmarkLight_GetObject_1MB(b *testing.B) {
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
	objectKey := fmt.Sprintf("light-download-1mb-%d", time.Now().UnixNano())
	content := bytes.Repeat([]byte("A"), 1024*1024) // 1MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
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

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				b.Errorf("下载失败: %v", err)
				return
			}

			// 读取内容
			io.Copy(io.Discard, output.Body)
			output.Body.Close()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkLight_GetObject_1MB",
		calculateThroughput(b.Result, 1.0))
}

// BenchmarkLight_GetObject_10MB 中等大小文件下载性能测试
func BenchmarkLight_GetObject_10MB(b *testing.B) {
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
	objectKey := fmt.Sprintf("light-download-10mb-%d", time.Now().UnixNano())
	content := bytes.Repeat([]byte("B"), 10*1024*1024) // 10MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
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
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				b.Errorf("下载失败: %v", err)
				return
			}

			io.Copy(io.Discard, output.Body)
			output.Body.Close()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkLight_GetObject_10MB",
		calculateThroughput(b.Result, 10.0))
}

// BenchmarkLight_GetObjectWithRange 范围下载性能测试
func BenchmarkLight_GetObjectWithRange(b *testing.B) {
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
	objectKey := fmt.Sprintf("range-download-%d", time.Now().UnixNano())
	content := bytes.Repeat([]byte("C"), 5*1024*1024) // 5MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
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
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
				Range:  "0-1023999", // 前1MB
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				b.Errorf("范围下载失败: %v", err)
				return
			}

			io.Copy(io.Discard, output.Body)
			output.Body.Close()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkLight_GetObjectWithRange",
		calculateThroughput(b.Result, 1.0))
}

// BenchmarkDeep_GetObject_100MB 深度大文件下载性能测试
func BenchmarkDeep_GetObject_100MB(b *testing.B) {
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
	objectKey := fmt.Sprintf("deep-download-100mb-%d", time.Now().UnixNano())
	content := bytes.Repeat([]byte("D"), 100*1024*1024) // 100MB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
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
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				b.Errorf("下载失败: %v", err)
				return
			}

			io.Copy(io.Discard, output.Body)
			output.Body.Close()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkDeep_GetObject_100MB",
		calculateThroughput(b.Result, 100.0))
}

// BenchmarkDeep_GetObject_1GB 超大文件下载性能测试
func BenchmarkDeep_GetObject_1GB(b *testing.B) {
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
	objectKey := fmt.Sprintf("deep-download-1gb-%d", time.Now().UnixNano())
	content := bytes.Repeat([]byte("E"), 1024*1024*1024) // 1GB

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
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
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				b.Errorf("下载失败: %v", err)
				return
			}

			io.Copy(io.Discard, output.Body)
			output.Body.Close()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkDeep_GetObject_1GB",
		calculateThroughput(b.Result, 1024.0))
}

// BenchmarkConcurrent_GetObject 不同并发级别下载性能测试
func BenchmarkConcurrent_GetObject(b *testing.B) {
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

	// 准备多个测试对象
	objectKeys := make([]string, 10)
	for i := 0; i < 10; i++ {
		objectKeys[i] = fmt.Sprintf("concurrent-download-%d-%d", i, time.Now().UnixNano())

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKeys[i],
			Body:   bytes.NewReader(content),
		}

		if _, err := obsClient.PutObject(putInput); err != nil {
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

	// 不同的并发级别
	concurrencyLevels := []int{10, 50, 100, 200}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)

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
						b.Errorf("并发下载失败: %v", err)
						return
					}

					io.Copy(io.Discard, output.Body)
					output.Body.Close()
				}
			})

			checkPerformanceDegradation(b,
				fmt.Sprintf("BenchmarkConcurrent_GetObject_Concurrency-%d", concurrency),
				calculateThroughput(b.Result, 2.0))
		})
	}
}

// BenchmarkMemory_GetObject 内存使用下载性能测试
func BenchmarkMemory_GetObject(b *testing.B) {
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

	content := bytes.Repeat([]byte("G"), 10*1024*1024) // 10MB
	bucket := cfg.GetTestBucket()
	objectKey := fmt.Sprintf("memory-download-%d", time.Now().UnixNano())

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
		b.Fatalf("准备测试对象失败: %v", err)
	}

	defer func() {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		obsClient.DeleteObject(deleteInput)
	}()

	var m runtime.MemStats
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runtime.ReadMemStats(&m)
			memBefore := m.Alloc

			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				b.Errorf("下载失败: %v", err)
				return
			}

			io.Copy(io.Discard, output.Body)
			output.Body.Close()

			// 记录内存使用
			runtime.ReadMemStats(&m)
			memAfter := m.Alloc
			memDelta := memAfter - memBefore

			b.ReportMetric(float64(memDelta)/1024/1024, "Alloc_MB")
		}
	})

	runtime.ReadMemStats(&m)
	b.Logf("内存使用统计: Alloc=%d MB, HeapAlloc=%d MB, Sys=%d MB, NumGC=%d",
		m.Alloc/1024/1024, m.HeapAlloc/1024/1024,
		m.Sys/1024/1024, m.NumGC)
}

// BenchmarkDownloadResourceMonitoring 下载资源监控测试
func BenchmarkDownloadResourceMonitoring(b *testing.B) {
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

	content := bytes.Repeat([]byte("H"), 5*1024*1024) // 5MB
	bucket := cfg.GetTestBucket()
	objectKey := fmt.Sprintf("resource-download-%d", time.Now().UnixNano())

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
		b.Fatalf("准备测试对象失败: %v", err)
	}

	defer func() {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    objectKey,
		}
		obsClient.DeleteObject(deleteInput)
	}()

	var m runtime.MemStats
	var totalGC uint32
	startTime := time.Now()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				b.Errorf("下载失败: %v", err)
				return
			}

			io.Copy(io.Discard, output.Body)
			output.Body.Close()

			// 记录资源使用
			runtime.ReadMemStats(&m)
			totalGC += m.NumGC - totalGC
		}
	})

	// 分析资源使用
	elapsed := time.Since(startTime)
	runtime.ReadMemStats(&m)

	b.Logf("=== 下载资源使用统计 ===")
	b.Logf("总操作数: %d", b.N)
	b.Logf("总时间: %v", elapsed)
	b.Logf("平均操作时间: %.2f ms", float64(elapsed)/float64(b.N)/float64(time.Millisecond))
	b.Logf("最大内存: %d MB", m.Alloc/1024/1024)
	b.Logf("堆内存: %d MB", m.HeapAlloc/1024/1024)
	b.Logf("系统内存: %d MB", m.Sys/1024/1024)
	b.Logf("GC次数: %d", totalGC)
}

// BenchmarkDownloadSequential 顺序下载性能测试
func BenchmarkDownloadSequential(b *testing.B) {
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

	content := bytes.Repeat([]byte("I"), 1*1024*1024) // 1MB
	bucket := cfg.GetTestBucket()

	// 准备测试对象
	objectKeys := make([]string, 20)
	for i := 0; i < 20; i++ {
		objectKeys[i] = fmt.Sprintf("sequential-download-%d-%d", i, time.Now().UnixNano())

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKeys[i],
			Body:   bytes.NewReader(content),
		}

		if _, err := obsClient.PutObject(putInput); err != nil {
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

	var wg sync.WaitGroup
	errors := make(chan error, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			keyIndex := index % len(objectKeys)
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKeys[keyIndex],
			}

			output, err := obsClient.GetObject(input)
			if err != nil {
				errors <- fmt.Errorf("顺序下载 %d 失败: %v", index, err)
				return
			}

			io.Copy(io.Discard, output.Body)
			output.Body.Close()
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

	checkPerformanceDegradation(b, "BenchmarkDownloadSequential",
		calculateThroughput(b.Result, 1.0))
}

// BenchmarkDownloadMetadata 对象元数据获取性能测试
func BenchmarkDownloadMetadata(b *testing.B) {
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

	content := bytes.Repeat([]byte("J"), 2*1024*1024) // 2MB
	bucket := cfg.GetTestBucket()
	objectKey := fmt.Sprintf("metadata-download-%d", time.Now().UnixNano())

	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    objectKey,
		Body:   bytes.NewReader(content),
	}
	if _, err := obsClient.PutObject(putInput); err != nil {
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
			input := &obs.GetObjectMetadataInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			_, err := obsClient.GetObjectMetadata(input)
			if err != nil {
				b.Errorf("获取元数据失败: %v", err)
			}
		}
	})
}

// BenchmarkDownloadMixedSizes 不同文件大小下载性能测试
func BenchmarkDownloadMixedSizes(b *testing.B) {
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
		objectKey := fmt.Sprintf("mixed-download-%s-%d", test.name, time.Now().UnixNano())
		content := bytes.Repeat([]byte("K"), int(test.size))

		putInput := &obs.PutObjectInput{
			Bucket: bucket,
			Key:    objectKey,
			Body:   bytes.NewReader(content),
		}

		if _, err := obsClient.PutObject(putInput); err != nil {
			b.Fatalf("准备测试对象失败: %v", err)
		}

		defer func(key string) {
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    key,
			}
			obsClient.DeleteObject(deleteInput)
		}(objectKey)

		b.Run(test.name, func(b *testing.B) {
			input := &obs.GetObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					output, err := obsClient.GetObject(input)
					if err != nil {
						b.Errorf("%s文件下载失败: %v", test.name, err)
						return
					}

					io.Copy(io.Discard, output.Body)
					output.Body.Close()
				}
			})

			checkPerformanceDegradation(b,
				fmt.Sprintf("BenchmarkDownloadMixedSizes_%s", test.name),
				calculateThroughput(b.Result, test.sizeMB))
		})
	}
}
