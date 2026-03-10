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
var multipartPerformanceBaseline = map[string]float64{
	"BenchmarkLight_MultipartUpload_10MB":  100.0,  // MB/s
	"BenchmarkDeep_MultipartUpload_100MB":  50.0,  // MB/s
	"BenchmarkConcurrent_MultipartUpload_50": 300.0, // ops/s
}

// BenchmarkLight_MultipartUpload_10MB 轻量级分块上传性能测试
func BenchmarkLight_MultipartUpload_10MB(b *testing.B) {
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
	objectKey := fmt.Sprintf("light-multipart-10mb-%d", time.Now().UnixNano())
	fileSize := 10 * 1024 * 1024 // 10MB
	partSize := 5 * 1024 * 1024 // 5MB分块

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := objectKey + "-" + fmt.Sprintf("%d", time.Now().UnixNano())

			// 初始化分块上传
			initInput := &obs.InitiateMultipartUploadInput{
				Bucket: bucket,
				Key:    testKey,
			}

			initOutput, err := obsClient.InitiateMultipartUpload(initInput)
			if err != nil {
				b.Errorf("初始化分块上传失败: %v", err)
				return
			}

			// 上传分块
			parts := make([]obs.Part, 2)
			for i := 0; i < 2; i++ {
				partContent := bytes.Repeat([]byte("A"), partSize)
				partKey := fmt.Sprintf("part-%d-%d", i, time.Now().UnixNano())

				// 创建临时分块对象
				putInput := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    partKey,
					Body:   bytes.NewReader(partContent),
				}

				if _, err := obsClient.PutObject(putInput); err != nil {
					b.Errorf("创建临时分块失败: %v", err)
					break
				}

				// 上传分块
				partInput := &obs.UploadPartInput{
					Bucket:     bucket,
					Key:        testKey,
					PartNumber: int32(i + 1),
					UploadId:   initOutput.UploadId,
					SourceFile: partKey,
				}

				partOutput, err := obsClient.UploadPart(partInput)
				if err != nil {
					b.Errorf("上传分块失败: %v", err)
					break
				}

				parts[i] = obs.Part{
					PartNumber: partOutput.PartNumber,
					ETag:       partOutput.ETag,
				}

				// 清理临时分块对象
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    partKey,
				}
				obsClient.DeleteObject(deleteInput)
			}

			// 完成分块上传
			completeInput := &obs.CompleteMultipartUploadInput{
				Bucket:   bucket,
				Key:      testKey,
				UploadId: initOutput.UploadId,
				Parts:    parts,
			}

			_, err = obsClient.CompleteMultipartUpload(completeInput)
			if err != nil {
				b.Errorf("完成分块上传失败: %v", err)
				return
			}

			// 清理对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    testKey,
			}
			go func() {
				obsClient.DeleteObject(deleteInput)
			}()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkLight_MultipartUpload_10MB",
		calculateThroughput(b.Result, 10.0))
}

// BenchmarkDeep_MultipartUpload_100MB 深度大文件分块上传性能测试
func BenchmarkDeep_MultipartUpload_100MB(b *testing.B) {
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
	objectKey := fmt.Sprintf("deep-multipart-100mb-%d", time.Now().UnixNano())
	fileSize := 100 * 1024 * 1024 // 100MB
	partSize := 10 * 1024 * 1024 // 10MB分块

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testKey := objectKey + "-" + fmt.Sprintf("%d", time.Now().UnixNano())

			// 初始化分块上传
			initInput := &obs.InitiateMultipartUploadInput{
				Bucket: bucket,
				Key:    testKey,
			}

			initOutput, err := obsClient.InitiateMultipartUpload(initInput)
			if err != nil {
				b.Errorf("初始化分块上传失败: %v", err)
				return
			}

			// 上传10个分块
			parts := make([]obs.Part, 10)
			for i := 0; i < 10; i++ {
				partContent := bytes.Repeat([]byte("B"), partSize)
				partKey := fmt.Sprintf("part-%d-%d", i, time.Now().UnixNano())

				// 创建临时分块对象
				putInput := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    partKey,
					Body:   bytes.NewReader(partContent),
				}

				if _, err := obsClient.PutObject(putInput); err != nil {
					b.Errorf("创建临时分块失败: %v", err)
					break
				}

				// 上传分块
				partInput := &obs.UploadPartInput{
					Bucket:     bucket,
					Key:        testKey,
					PartNumber: int32(i + 1),
					UploadId:   initOutput.UploadId,
					SourceFile: partKey,
				}

				partOutput, err := obsClient.UploadPart(partInput)
				if err != nil {
					b.Errorf("上传分块失败: %v", err)
					break
				}

				parts[i] = obs.Part{
					PartNumber: partOutput.PartNumber,
					ETag:       partOutput.ETag,
				}

				// 清理临时分块对象
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    partKey,
				}
				obsClient.DeleteObject(deleteInput)
			}

			// 完成分块上传
			completeInput := &obs.CompleteMultipartUploadInput{
				Bucket:   bucket,
				Key:      testKey,
				UploadId: initOutput.UploadId,
				Parts:    parts,
			}

			_, err = obsClient.CompleteMultipartUpload(completeInput)
			if err != nil {
				b.Errorf("完成分块上传失败: %v", err)
				return
			}

			// 清理对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    testKey,
			}
			go func() {
				obsClient.DeleteObject(deleteInput)
			}()
		}
	})

	checkPerformanceDegradation(b, "BenchmarkDeep_MultipartUpload_100MB",
		calculateThroughput(b.Result, 100.0))
}

// BenchmarkConcurrent_MultipartUpload 不同分块数量的性能测试
func BenchmarkConcurrent_MultipartUpload(b *testing.B) {
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
	objectKey := fmt.Sprintf("concurrent-multipart-%d", time.Now().UnixNano())
	fileSize := 20 * 1024 * 1024 // 20MB
	partSize := 5 * 1024 * 1024 // 5MB分块

	// 不同的分块数量
	partCounts := []int{2, 4, 8, 16}

	for _, partCount := range partCounts {
		b.Run(fmt.Sprintf("Parts-%d", partCount), func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					testKey := objectKey + "-" + fmt.Sprintf("%d", time.Now().UnixNano())

					// 初始化分块上传
					initInput := &obs.InitiateMultipartUploadInput{
						Bucket: bucket,
						Key:    testKey,
					}

					initOutput, err := obsClient.InitiateMultipartUpload(initInput)
					if err != nil {
						b.Errorf("初始化分块上传失败: %v", err)
						return
					}

					// 上传指定数量的分块
					parts := make([]obs.Part, partCount)
					for i := 0; i < partCount; i++ {
						partContent := bytes.Repeat([]byte("C"), partSize)
						partKey := fmt.Sprintf("part-%d-%d", i, time.Now().UnixNano())

						// 创建临时分块对象
						putInput := &obs.PutObjectInput{
							Bucket: bucket,
							Key:    partKey,
							Body:   bytes.NewReader(partContent),
						}

						if _, err := obsClient.PutObject(putInput); err != nil {
							break
						}

						// 上传分块
						partInput := &obs.UploadPartInput{
							Bucket:     bucket,
							Key:        testKey,
							PartNumber: int32(i + 1),
							UploadId:   initOutput.UploadId,
							SourceFile: partKey,
						}

						partOutput, err := obsClient.UploadPart(partInput)
						if err != nil {
							break
						}

						parts[i] = obs.Part{
							PartNumber: partOutput.PartNumber,
							ETag:       partOutput.ETag,
						}

						// 清理临时分块对象
						deleteInput := &obs.DeleteObjectInput{
							Bucket: bucket,
							Key:    partKey,
						}
						obsClient.DeleteObject(deleteInput)
					}

					// 完成分块上传
					completeInput := &obs.CompleteMultipartUploadInput{
						Bucket:   bucket,
						Key:      testKey,
						UploadId: initOutput.UploadId,
						Parts:    parts,
					}

					_, err = obsClient.CompleteMultipartUpload(completeInput)
					if err != nil {
						return
					}

					// 清理对象
					deleteInput := &obs.DeleteObjectInput{
						Bucket: bucket,
						Key:    testKey,
					}
					go func() {
						obsClient.DeleteObject(deleteInput)
					}()
				}
			})

			checkPerformanceDegradation(b,
				fmt.Sprintf("BenchmarkConcurrent_MultipartUpload_Parts-%d", partCount),
				calculateThroughput(b.Result, 20.0))
		})
	}
}

// BenchmarkMultipart_Initialization 分块上传初始化性能测试
func BenchmarkMultipart_Initialization(b *testing.B) {
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
			objectKey := fmt.Sprintf("init-multipart-%d", time.Now().UnixNano())

			// 仅测试初始化性能
			initInput := &obs.InitiateMultipartUploadInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			_, err := obsClient.InitiateMultipartUpload(initInput)
			if err != nil {
				b.Errorf("初始化分块上传失败: %v", err)
			}
		}
	})
}

// BenchmarkMultipart_PartUpload 分块上传性能测试
func BenchmarkMultipart_PartUpload(b *testing.B) {
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
	objectKey := fmt.Sprintf("part-upload-%d", time.Now().UnixNano())
	partSize := 5 * 1024 * 1024 // 5MB分块

	// 初始化一次分块上传
	initInput := &obs.InitiateMultipartUploadInput{
		Bucket: bucket,
		Key:    objectKey,
	}

	initOutput, err := obsClient.InitiateMultipartUpload(initInput)
	if err != nil {
		b.Fatalf("初始化分块上传失败: %v", err)
	}

	// 清理函数
	defer func() {
		abortInput := &obs.AbortMultipartUploadInput{
			Bucket:   bucket,
			Key:      objectKey,
			UploadId: initOutput.UploadId,
		}
		obsClient.AbortMultipartUpload(abortInput)
	}()

	// 测试单个分块上传性能
	partContent := bytes.Repeat([]byte("D"), partSize)
	partKey := fmt.Sprintf("test-part-%d", time.Now().UnixNano())

	// 创建临时分块对象
	putInput := &obs.PutObjectInput{
		Bucket: bucket,
		Key:    partKey,
		Body:   bytes.NewReader(partContent),
	}

	if _, err := obsClient.PutObject(putInput); err != nil {
		b.Fatalf("创建临时分块失败: %v", err)
	}

	// 清理临时分块对象
	defer func() {
		deleteInput := &obs.DeleteObjectInput{
			Bucket: bucket,
			Key:    partKey,
		}
		obsClient.DeleteObject(deleteInput)
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			partInput := &obs.UploadPartInput{
				Bucket:     bucket,
				Key:        objectKey,
				PartNumber: 1,
				UploadId:   initOutput.UploadId,
				SourceFile: partKey,
			}

			_, err := obsClient.UploadPart(partInput)
			if err != nil {
				b.Errorf("上传分块失败: %v", err)
			}
		}
	})

	checkPerformanceDegradation(b, "BenchmarkMultipart_PartUpload",
		calculateThroughput(b.Result, 5.0))
}

// BenchmarkMultipart_Completion 分块上传完成性能测试
func BenchmarkMultipart_Completion(b *testing.B) {
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
			objectKey := fmt.Sprintf("complete-multipart-%d", time.Now().UnixNano())

			// 初始化分块上传
			initInput := &obs.InitiateMultipartUploadInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			initOutput, err := obsClient.InitiateMultipartUpload(initInput)
			if err != nil {
				b.Errorf("初始化分块上传失败: %v", err)
				return
			}

			// 准备分块
			parts := make([]obs.Part, 5)
			for i := 0; i < 5; i++ {
				partContent := bytes.Repeat([]byte("E"), 2*1024*1024)
				partKey := fmt.Sprintf("part-%d-%d", i, time.Now().UnixNano())

				// 创建临时分块对象
				putInput := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    partKey,
					Body:   bytes.NewReader(partContent),
				}

				if _, err := obsClient.PutObject(putInput); err != nil {
					break
				}

				// 上传分块
				partInput := &obs.UploadPartInput{
					Bucket:     bucket,
					Key:        objectKey,
					PartNumber: int32(i + 1),
					UploadId:   initOutput.UploadId,
					SourceFile: partKey,
				}

				partOutput, err := obsClient.UploadPart(partInput)
				if err != nil {
					break
				}

				parts[i] = obs.Part{
					PartNumber: partOutput.PartNumber,
					ETag:       partOutput.ETag,
				}

				// 清理临时分块对象
				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    partKey,
				}
				obsClient.DeleteObject(deleteInput)
			}

			// 测试完成性能
			completeInput := &obs.CompleteMultipartUploadInput{
				Bucket:   bucket,
				Key:      objectKey,
				UploadId: initOutput.UploadId,
				Parts:    parts,
			}

			_, err = obsClient.CompleteMultipartUpload(completeInput)
			if err != nil {
				b.Errorf("完成分块上传失败: %v", err)
				return
			}

			// 清理对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			go func() {
				obsClient.DeleteObject(deleteInput)
			}()
		}
	})
}

// BenchmarkMultipart_MemoryUsage 内存使用性能测试
func BenchmarkMultipart_MemoryUsage(b *testing.B) {
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
	objectKey := fmt.Sprintf("memory-multipart-%d", time.Now().UnixNano())
	partSize := 2 * 1024 * 1024 // 2MB分块

	var m runtime.MemStats

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runtime.ReadMemStats(&m)
			memBefore := m.Alloc

			// 初始化分块上传
			initInput := &obs.InitiateMultipartUploadInput{
				Bucket: bucket,
				Key:    objectKey,
			}

			initOutput, err := obsClient.InitiateMultipartUpload(initInput)
			if err != nil {
				b.Errorf("初始化分块上传失败: %v", err)
				return
			}

			// 上传3个分块
			parts := make([]obs.Part, 3)
			for i := 0; i < 3; i++ {
				partContent := bytes.Repeat([]byte("F"), partSize)
				partKey := fmt.Sprintf("part-%d-%d", i, time.Now().UnixNano())

				putInput := &obs.PutObjectInput{
					Bucket: bucket,
					Key:    partKey,
					Body:   bytes.NewReader(partContent),
				}

				if _, err := obsClient.PutObject(putInput); err != nil {
					break
				}

				partInput := &obs.UploadPartInput{
					Bucket:     bucket,
					Key:        objectKey,
					PartNumber: int32(i + 1),
					UploadId:   initOutput.UploadId,
					SourceFile: partKey,
				}

				partOutput, err := obsClient.UploadPart(partInput)
				if err != nil {
					break
				}

				parts[i] = obs.Part{
					PartNumber: partOutput.PartNumber,
					ETag:       partOutput.ETag,
				}

				deleteInput := &obs.DeleteObjectInput{
					Bucket: bucket,
					Key:    partKey,
				}
				obsClient.DeleteObject(deleteInput)
			}

			// 完成分块上传
			completeInput := &obs.CompleteMultipartUploadInput{
				Bucket:   bucket,
				Key:      objectKey,
				UploadId: initOutput.UploadId,
				Parts:    parts,
			}

			obsClient.CompleteMultipartUpload(completeInput)

			// 记录内存使用
			runtime.ReadMemStats(&m)
			memAfter := m.Alloc
			memDelta := memAfter - memBefore

			b.ReportMetric(float64(memDelta)/1024/1024, "Alloc_MB")

			// 清理对象
			deleteInput := &obs.DeleteObjectInput{
				Bucket: bucket,
				Key:    objectKey,
			}
			go func() {
				obsClient.DeleteObject(deleteInput)
			}()
		}
	})

	runtime.ReadMemStats(&m)
	b.Logf("内存使用统计: Alloc=%d MB, HeapAlloc=%d MB, Sys=%d MB, NumGC=%d",
		m.Alloc/1024/1024, m.HeapAlloc/1024/1024,
		m.Sys/1024/1024, m.NumGC)
}
