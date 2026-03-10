# OBS SDK Go 性能测试使用指南

## 概述

本指南说明如何使用OBS SDK Go的性能测试系统来评估和监控SDK性能表现。

## 性能测试文件

### 上传性能测试 (`upload_benchmark_test.go`)

| 测试名称 | 文件大小 | 类型 | 说明 |
|---------|---------|------|------|
| `BenchmarkLight_PutObject_1MB` | 1MB | 轻量级 | 基本小文件上传 |
| `BenchmarkLight_PutObject_1MB_Metadata` | 1MB | 轻量级 | 带元数据上传 |
| `BenchmarkLight_PutObject_10MB` | 10MB | 轻量级 | 中等大小文件上传 |
| `BenchmarkDeep_PutObject_100MB` | 100MB | 深度 | 大文件上传 |
| `BenchmarkDeep_PutObject_1GB` | 1GB | 深度 | 超大文件上传 |
| `BenchmarkConcurrent_PutObject` | 5MB | 并发 | 不同并发级别上传 (1,10,50,100,500) |
| `BenchmarkMemory_PutObject` | 10MB | 内存 | 内存使用监控 |
| `BenchmarkUploadResourceMonitoring` | 5MB | 资源 | 资源使用监控 |
| `BenchmarkUploadConnectionPooling` | 2MB | 连接池 | 连接池性能 |
| `BenchmarkUploadWithContentType` | 2MB | 内容类型 | 不同内容类型上传 |
| `BenchmarkUploadDifferentSizes` | 多种 | 文件大小 | 不同文件大小上传 (100KB-50MB) |
| `BenchmarkUploadSequential` | 1MB | 顺序 | 顺序上传测试 |

### 下载性能测试 (`download_benchmark_test.go`)

| 测试名称 | 文件大小 | 类型 | 说明 |
|---------|---------|------|------|
| `BenchmarkLight_GetObject_1MB` | 1MB | 轻量级 | 基本小文件下载 |
| `BenchmarkLight_GetObject_10MB` | 10MB | 轻量级 | 中等大小文件下载 |
| `BenchmarkLight_GetObjectWithRange` | 5MB | 范围 | 范围下载 (1MB) |
| `BenchmarkDeep_GetObject_100MB` | 100MB | 深度 | 大文件下载 |
| `BenchmarkDeep_GetObject_1GB` | 1GB | 深度 | 超大文件下载 |
| `BenchmarkConcurrent_GetObject` | 2MB | 并发 | 不同并发级别下载 (10,50,100,200) |
| `BenchmarkMemory_GetObject` | 10MB | 内存 | 内存使用监控 |
| `BenchmarkDownloadResourceMonitoring` | 5MB | 资源 | 资源使用监控 |
| `BenchmarkDownloadSequential` | 1MB | 顺序 | 顺序下载测试 (20个对象) |
| `BenchmarkDownloadMetadata` | 2MB | 元数据 | 对象元数据获取 |
| `BenchmarkDownloadMixedSizes` | 多种 | 文件大小 | 不同文件大小下载 (100KB-50MB) |

### 分块上传性能测试 (`multipart_benchmark_test.go`)

| 测试名称 | 文件大小 | 分块数 | 说明 |
|---------|---------|--------|------|
| `BenchmarkLight_MultipartUpload_10MB` | 10MB | 2 | 轻量级分块上传 |
| `BenchmarkDeep_MultipartUpload_100MB` | 100MB | 10 | 深度分块上传 |
| `BenchmarkConcurrent_MultipartUpload` | 20MB | 多种 | 不同分块数量 (2,4,8,16) |
| `BenchmarkMultipart_Initialization` | - | - | 分块上传初始化性能 |
| `BenchmarkMultipart_PartUpload` | 5MB | - | 单个分块上传性能 |
| `BenchmarkMultipart_Completion` | 10MB | 5 | 分块上传完成性能 |
| `BenchmarkMultipart_MemoryUsage` | 6MB | 3 | 分块上传内存使用 |

### 并发性能测试 (`concurrent_benchmark_test.go`)

| 测试名称 | 并发数 | 类型 | 说明 |
|---------|-------|------|------|
| `BenchmarkConcurrent_MixedOperations_10` | 10 | 混合操作 | 上传、下载、元数据混合操作 |
| `BenchmarkConcurrent_StressTest_100` | 100 | 压力测试 | 高并发压力测试 |
| `BenchmarkConcurrent_ConnectionPool_20` | 20 | 连接池 | 连接池性能 |
| `BenchmarkConcurrent_ResourceRace` | - | 资源竞争 | 内存竞争测试 |
| `BenchmarkConcurrent_ListObjects` | 100 | 列表 | 并发列表操作 |
| `BenchmarkConcurrent_MemoryUsage` | - | 内存 | 并发内存使用 |
| `BenchmarkConcurrent_ErrorHandling` | - | 错误 | 并发错误处理 |
| `BenchmarkConcurrent_Latency` | - | 延迟 | 并发操作延迟 |
| `BenchmarkConcurrent_ThroughputScaling` | - | 吞吐量 | 吞吐量扩展测试 |

### 性能基线管理 (`performance_baseline_test.go`)

包含完整的性能基线管理功能：
- 基线数据加载和保存
- 性能对比和退化检测
- HTML报告生成
- 基线稳定性验证

## 运行性能测试

### 基本命令

```bash
# 进入项目目录
cd /path/to/huaweicloud-sdk-go-obs

# 运行轻量级性能测试
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem

# 运行深度性能测试
go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem

# 运行所有性能测试
go test -tags perf ./obs -bench=. -benchtime=1s

# 运行并发性能测试
go test -tags perf ./obs -bench=BenchmarkConcurrent -benchtime=10s -benchmem
```

### 使用不同并发级别

```bash
# 设置并发级别
go test -tags perf ./obs -bench=BenchmarkConcurrent -benchtime=10s -parallel=50

# 运行特定并发级别
go test -tags perf ./obs -bench=BenchmarkConcurrent/Concurrency-100 -benchtime=10s
```

### 保存测试结果

```bash
# 保存测试结果到文件
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -o light-bench.out

# 生成HTML报告
go test -tags perf ./obs -bench=. > bench.out
```

## 性能基线管理

### 查看当前性能基线

```bash
# 运行基线验证测试
go test -tags perf ./obs -bench=BenchmarkBaselineValidation -v
```

### 更新性能基线

运行任何性能测试会自动更新对应的基线数据：
```bash
# 运行性能测试（会自动更新基线）
go test -tags perf ./obs -bench=BenchmarkLight_PutObject_1MB -benchtime=1s
```

### 生成性能报告

```bash
# 在测试代码中调用报告生成功能
# 或创建专门的报告生成工具
```

### 导出性能报告

```bash
# 导出HTML报告
go test -tags perf ./obs -bench=BenchmarkBaselineExport
```

## 性能指标说明

### 关键指标

1. **ns/op** - 每次操作的纳秒数
   - 数值越小越好
   - 反映单个操作的执行时间

2. **MB/s** - 每秒传输的兆字节数
   - 数值越大越好
   - 衡量数据传输效率

3. **allocs/op** - 每次操作的内存分配次数
   - 数值越小越好
   - 反映内存分配效率

4. **bytes/op** - 每次操作分配的字节数
   - 数值越小越好
   - 反映内存使用量

### 性能基线

性能基线是性能测试的参考标准，包含以下内容：

- 轻量级测试基线：小文件、低并发、短时长
- 深度测试基线：大文件、高并发、长时长
- 并发测试基线：不同并发级别

### 性能退化检测

- 阈值：基线值的 90%
- 检测：当前性能低于基线 90% 时触发告警
- 处理：性能退化应立即调查和修复

### 资源监控

- **CPU使用率**：CPU占用百分比
- **内存使用**：堆内存、系统内存
- **GC次数**：垃圾回收次数
- **网络带宽**：实际传输速率

## 测试环境配置

### 环境变量

```bash
# 必需配置
export OBS_TEST_AK="your-access-key"
export OBS_TEST_SK="your-secret-key"
export OBS_TEST_ENDPOINT="https://obs.cn-north-4.myhuaweicloud.com"
export OBS_TEST_BUCKET="your-test-bucket"

# 可选配置
export OBS_TEST_REGION="cn-north-4"
export OBS_TEST_TOKEN="your-temporary-token"
```

### 系统要求

- Go版本：1.18+
- 操作系统：Linux/macOS/Windows
- 网络：稳定的网络连接
- 硬件：充足的CPU和内存

### 测试环境注意事项

1. **网络条件**：网络条件会影响测试结果
2. **系统负载**：关闭其他程序以获得稳定结果
3. **磁盘IO**：确保磁盘性能不会成为瓶颈
4. **内存充足**：确保有足够的内存避免GC影响

## 性能测试最佳实践

### 1. 定期运行

```bash
# 每天自动运行
0 2 * * * cd /path/to/project && go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s

# 每周深度测试
0 2 * * 0 cd /path/to/project && go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s
```

### 3. 性能数据收集

收集和存储性能数据用于趋势分析：
- 存储每次测试结果
- 绘制性能趋势图
- 分析性能变化模式
- 识别性能瓶颈

### 4. 性能优化建议

根据性能测试结果优化代码：

1. **减少内存分配**：
   - 重用缓冲区
   - 避免不必要的拷贝
   - 使用对象池

2. **提高并发效率**：
   - 优化锁的使用
   - 减少上下文切换
   - 合理设置并发级别

3. **优化网络IO**：
   - 使用连接池
   - 批量操作
   - 减少网络往返

4. **减少GC压力**：
   - 避免大对象分配
   - 使用sync.Pool
   - 减少临时对象

## 性能问题诊断

### 常见性能问题

1. **内存泄漏**
   ```bash
   # 检查内存泄漏
   go test -tags perf ./obs -bench=. -benchtime=10s -memprofile=mem.out
   go tool pprof mem.out
   ```

2. **CPU热点**
   ```bash
   # 分析CPU性能
   go test -tags perf ./obs -bench=. -benchtime=10s -cpuprofile=cpu.out
   go tool pprof cpu.out
   ```

3. **锁竞争**
   ```bash
   # 检测锁竞争
   go test -tags perf ./obs -bench=. -race
   ```

4. **网络瓶颈**
   - 检查网络延迟
   - 分析DNS解析时间
   - 检查连接建立时间

## 性能测试流程

### 1. 初始基线建立

1. 在稳定环境中运行所有性能测试
2. 收集基线数据
3. 验证基线稳定性（多次运行确认）
4. 保存初始基线

### 2. 持续监控

1. 定期运行轻量级测试
2. 监控性能变化趋势
3. 检测性能退化
4. 及时发现性能问题

### 3. 性能优化

1. 分析性能瓶颈
2. 实施优化措施
3. 重新测试验证效果
4. 更新基线数据

### 4. 回归测试

1. 确保优化不影响功能
2. 运行完整的测试套件
3. 验证性能提升效果

## 性能测试工具

### benchstat

Go官方性能对比工具：
```bash
# 安装
go install golang.org/x/perf/cmd/benchstat@latest

# 对比性能
benchstat baseline.out current.out
```

### go-tc

代码覆盖和性能分析：
```bash
# 安装
go install github.com/golang/tools/cmd/go-tc@latest

# 运行性能测试
go-tc -benchmem -benchtime=5s ./obs/
```

### flamegraph

火焰图可视化：
```bash
# 安装
go install github.com/uber/go-torch@latest

# 生成火焰图
go test -tags perf ./obs -bench=. -cpuprofile=cpu.out
go-torch -http=:8080 cpu.out
```

## 性能报告解读

### 报告结构

1. **测试概览**
   - 测试名称和版本
   - 测试环境信息
   - 测试时间戳

2. **性能数据**
   - 基线性能数据
   - 当前性能数据
   - 性能差异百分比
   - 是否在阈值范围内

3. **资源使用**
   - 内存使用情况
   - CPU使用率
   - GC次数
   - 网络带宽

4. **性能趋势**
   - 历史性能趋势
   - 性能变化分析
   - 性能建议

### 性能评级

| 评级 | 标准 |
|------|------|
| 优秀 | 性能 > 基线 110% |
| 良好 | 基线 90% ≤ 性能 ≤ 基线 110% |
| 一般 | 基线 80% ≤ 性能 < 基线 90% |
| 较差 | 性能 < 基线 80% |
| 需要优化 | 性能 < 基线 60% |

## 下一步

1. 定期审查性能基线
3. 监控性能趋势
4. 持续优化性能
5. 完善性能测试覆盖

更多详细信息请参考：
- Go官方性能测试文档
- OBS SDK开发指南
- 性能测试最佳实践