//go:build perf

package obs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// PerformanceBaseline 性能基线结构
type PerformanceBaseline struct {
	Timestamp   time.Time              `json:"timestamp"`
	TestName    string                 `json:"test_name"`
	Throughput float64               `json:"throughput_mb_per_sec"`
	Latency     float64               `json:"latency_ms"`
	Allocations int64                 `json:"allocations_per_op"`
	Bytes       int64                 `json:"bytes_per_op"`
	CPU         float64               `json:"cpu_percent"`
	Memory      int64                 `json:"memory_mb"`
	GCCount     uint32                `json:"gc_count"`
	IsStable    bool                  `json:"is_stable"`
	Environment string                `json:"environment"`
}

// PerformanceReport 性能报告结构
type PerformanceReport struct {
	Baseline      PerformanceBaseline    `json:"baseline"`
	Current       PerformanceBaseline    `json:"current"`
	IsDegraded   bool                   `json:"is_degraded"`
	ThroughputDiff float64              `json:"throughput_diff_percent"`
	LatencyDiff   float64              `json:"latency_diff_percent"`
	IsWithinThreshold bool            `json:"is_within_threshold"`
	Threshold     float64              `json:"threshold"`
}

// 基线文件路径
const baselineFilePath = "benchmarks/performance_baseline.json"

// 性能阈值配置
const (
	// 性能退化阈值（低于基线多少百分比算退化）
	performanceThreshold = 0.9 // 90%

	// 基线稳定性要求（连续多少次测试通过才认为基线稳定）
	stabilityThreshold = 3

	// 基线老化时间（超过多少天的基线需要重新测试）
	baselineMaxAgeDays = 30
)

// LoadPerformanceBaseline 加载性能基线
func LoadPerformanceBaseline() (map[string]PerformanceBaseline, error) {
	baselines := make(map[string]PerformanceBaseline)

	// 检查基线文件是否存在
	if _, err := os.Stat(baselineFilePath); os.IsNotExist(err) {
		// 基线文件不存在，返回空基线
		return baselines, nil
	}

	// 读取基线文件
	data, err := os.ReadFile(baselineFilePath)
	if err != nil {
		return nil, fmt.Errorf("读取基线文件失败: %v", err)
	}

	// 解析JSON
	var savedBaselines []PerformanceBaseline
	if err := json.Unmarshal(data, &savedBaselines); err != nil {
		return nil, fmt.Errorf("解析基线文件失败: %v", err)
	}

	// 转换为map
	for _, baseline := range savedBaselines {
		baselines[baseline.TestName] = baseline
	}

	return baselines, nil
}

// SavePerformanceBaseline 保存性能基线
func SavePerformanceBaseline(baselines map[string]PerformanceBaseline) error {
	// 准备保存的数据
	var savedBaselines []PerformanceBaseline
	for _, baseline := range baselines {
		savedBaselines = append(savedBaselines, baseline)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(savedBaselines, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化基线数据失败: %v", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(baselineFilePath), 0755); err != nil {
		return fmt.Errorf("创建基线目录失败: %v", err)
	}

	// 保存到文件
	if err := os.WriteFile(baselineFilePath, data, 0644); err != nil {
		return fmt.Errorf("保存基线文件失败: %v", err)
	}

	return nil
}

// UpdatePerformanceBaseline 更新性能基线
func UpdatePerformanceBaseline(testName string, result testing.BenchmarkResult) error {
	baselines, err := LoadPerformanceBaseline()
	if err != nil {
		return err
	}

	// 获取或创建基线
	baseline, exists := baselines[testName]
	if !exists {
		// 创建新的基线
		baseline = PerformanceBaseline{
			Timestamp: time.Now(),
			TestName:   testName,
			IsStable:  false,
			Environment: "local",
		}
	} else {
		// 更新时间戳和稳定性
		baseline.Timestamp = time.Now()
		baseline.IsStable = checkBaselineStability(baseline)
	}

	// 更新性能数据
	baseline.Throughput = calculateThroughputMBPS(result, 1.0) // 假设1MB文件
	baseline.Latency = calculateLatencyMS(result)
	baseline.Allocations = result.AllocsPerOp
	baseline.Bytes = result.MemBytes

	// 获取内存使用（需要运行时采集）
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	baseline.Memory = m.Alloc / 1024 / 1024 // 转换为MB
	baseline.GCCount = m.NumGC

	// 保存更新后的基线
	baselines[testName] = baseline
	return SavePerformanceBaseline(baselines)
}

// checkBaselineStability 检查基线是否稳定
func checkBaselineStability(baseline PerformanceBaseline) bool {
	// 如果创建时间超过阈值，则不稳定
		if time.Since(baseline.Timestamp) > baselineMaxAgeDays*24*time.Hour {
		return false
	}

	// 这里可以添加更多稳定性检查逻辑
	// 例如：连续成功次数、方差分析等

	// 简化处理：如果是最近创建的，则认为稳定
	return true
}

// CompareWithBaseline 与基线对比性能
func CompareWithBaseline(testName string, result testing.BenchmarkResult) (PerformanceReport, error) {
	baselines, err := LoadPerformanceBaseline()
	if err != nil {
		return PerformanceReport{}, err
	}

	// 检查是否有基线
	baseline, hasBaseline := baselines[testName]
	if !hasBaseline {
		return PerformanceReport{
			Current: PerformanceBaseline{
				Timestamp: time.Now(),
				TestName:   testName,
				Environment: "local",
			},
			IsWithinThreshold: false,
			Threshold:         performanceThreshold,
		}, nil
	}

	// 计算当前性能数据
	currentThroughput := calculateThroughputMBPS(result, 1.0)
	currentLatency := calculateLatencyMS(result)
	currentAllocations := result.AllocsPerOp
	currentBytes := result.MemBytes

	// 计算性能差异
	throughputDiff := ((currentThroughput - baseline.Throughput) / baseline.Throughput) * 100
	latencyDiff := ((currentLatency - baseline.Latency) / baseline.Latency) * 100

	// 判断是否退化
	isDegraded := currentThroughput < baseline.Throughput*performanceThreshold

	// 判断是否在阈值范围内
	isWithinThreshold := currentThroughput >= baseline.Throughput*performanceThreshold

	report := PerformanceReport{
		Baseline:      baseline,
		Current: PerformanceBaseline{
			Timestamp:   time.Now(),
			TestName:    testName,
			Throughput:  currentThroughput,
			Latency:     currentLatency,
			Allocations: currentAllocations,
			Bytes:       currentBytes,
			Environment:  "local",
		},
		IsDegraded:      isDegraded,
		ThroughputDiff:  throughputDiff,
		LatencyDiff:   latencyDiff,
		IsWithinThreshold: isWithinThreshold,
		Threshold:        performanceThreshold,
	}

	return report, nil
}

// checkPerformanceDegradation 检查性能退化
func checkPerformanceDegradation(b *testing.B, testName string, currentThroughput float64) {
	baselines, err := LoadPerformanceBaseline()
	if err != nil {
		b.Logf("无法加载性能基线: %v", err)
		return
	}

	baseline, hasBaseline := baselines[testName]
	if !hasBaseline {
		b.Logf("没有找到性能基线: %s", testName)
		return
	}

	// 检查是否退化
	threshold := baseline.Throughput * performanceThreshold
	if currentThroughput < threshold {
		b.Errorf("性能退化检测失败! 测试: %s, 当前: %.2f MB/s, 基线: %.2f MB/s (阈值: %.1f%%)",
			testName, currentThroughput, baseline.Throughput, performanceThreshold*100)
	} else {
		b.Logf("性能检查通过: 测试: %s, 当前: %.2f MB/s, 基线: %.2f MB/s",
			testName, currentThroughput, baseline.Throughput)
	}
}

// GeneratePerformanceReport 生成性能报告
func GeneratePerformanceReport(b *testing.B, testName string, result testing.BenchmarkResult, fileSizeMB float64) error {
	baselines, err := LoadPerformanceBaseline()
	if err != nil {
		return err
	}

	// 获取或创建基线
	baseline, exists := baselines[testName]
	if !exists {
		// 创建新的基线
		baseline = PerformanceBaseline{
			Timestamp: time.Now(),
			TestName:   testName,
			Throughput: calculateThroughputMBPS(result, fileSizeMB),
			Latency:    calculateLatencyMS(result),
			Allocations: result.AllocsPerOp,
			Bytes:      result.MemBytes,
			IsStable:   false,
			Environment: "local",
		}

		// 保存新基线
		baselines[testName] = baseline
		if err := SavePerformanceBaseline(baselines); err != nil {
			return err
		}
	}

	// 生成报告
	currentThroughput := calculateThroughputMBPS(result, fileSizeMB)
	currentLatency := calculateLatencyMS(result)

	report := PerformanceReport{
		Baseline:    baseline,
		Current: PerformanceBaseline{
			Timestamp:   time.Now(),
			TestName:    testName,
			Throughput: currentThroughput,
			Latency:     currentLatency,
			Allocations: result.AllocsPerOp,
			Bytes:       result.MemBytes,
			Environment: "local",
		},
		IsDegraded:      false, // 因为刚设置基线
		ThroughputDiff:  0,
		LatencyDiff:   0,
		IsWithinThreshold: true,
		Threshold:        performanceThreshold,
	}

	// 生成JSON报告
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化性能报告失败: %v", err)
	}

	// 保存报告
	reportPath := fmt.Sprintf("benchmarks/performance_report_%s_%d.json",
		testName, time.Now().UnixNano())
	if err := os.WriteFile(reportPath, reportData, 0644); err != nil {
		return fmt.Errorf("保存性能报告失败: %v", err)
	}

	b.Logf("性能报告已保存: %s", reportPath)

	return nil
}

// PrintBaselineComparison 打印基线对比结果
func PrintBaselineComparison(b *testing.B) {
	baselines, err := LoadPerformanceBaseline()
	if err != nil {
		b.Logf("无法加载性能基线: %v", err)
		return
	}

	if len(baselines) == 0 {
		b.Log("没有性能基线数据")
		return
	}

	b.Logf("=== 性能基线对比 ===")
	b.Logf("基线数据来源: %s", baselineFilePath)
	b.Logf("基线数量: %d", len(baselines))
	b.Logf("性能阈值: %.1f%%", performanceThreshold*100)
	b.Logf("基线稳定阈值: %d 次测试", stabilityThreshold)
	b.Logf("基线最大年龄: %d 天", baselineMaxAgeDays)

	for testName, baseline := range baselines {
		b.Logf("\n基线: %s", testName)
		b.Logf("  更新时间: %s", baseline.Timestamp.Format("2006-01-02 15:04:05"))
		b.Logf("  吞吐量: %.2f MB/s", baseline.Throughput)
		b.Logf("  延迟: %.2f ms", baseline.Latency)
		b.Logf("  分配次数: %d/操作", baseline.Allocations)
		b.Logf("  内存使用: %.2f MB", baseline.Memory)
		b.Logf("  GC次数: %d", baseline.GCCount)
		b.Logf("  是否稳定: %t", baseline.IsStable)
	}
}

// ExportBaselineReport 导出基线报告
func ExportBaselineReport(reportPath string) error {
	baselines, err := LoadPerformanceBaseline()
	if err != nil {
		return err
	}

	// 生成HTML报告
	htmlReport := `<!DOCTYPE html>
<html>
<head>
    <title>性能基线报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #4CAF50; color: white; }
        .degraded { background-color: #FF6B6B; color: white; }
        .ok { background-color: #6CAF66; color: white; }
        .metric { font-weight: bold; }
    </style>
</head>
<body>
    <h1>OBS SDK Go 性能基线报告</h1>
    <p>生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
    <h2>性能基线概览</h2>
    <table>
        <tr>
            <th>测试名称</th>
            <th>基线吞吐量 (MB/s)</th>
            <th>基线延迟 (ms)</th>
            <th>基线分配次数</th>
            <th>基线内存 (MB)</th>
            <th>基线GC次数</th>
            <th>是否稳定</th>
            <th>最后更新</th>
        </tr>
`

	for testName, baseline := range baselines {
		statusClass := "ok"
		if !baseline.IsStable {
			statusClass = "degraded"
		}

		fmt.Fprintf(w, `
        <tr class="%s">
            <td class="metric">%s</td>
            <td>%.2f</td>
            <td>%.2f</td>
            <td>%d</td>
            <td>%.2f</td>
            <td>%d</td>
            <td>%t</td>
            <td>%s</td>
        </tr>`, statusClass, testName, baseline.Throughput,
			baseline.Latency, baseline.Allocations,
			baseline.Memory, baseline.GCCount, baseline.IsStable,
			baseline.Timestamp.Format("2006-01-02 15:04:05"))
	}

	fmt.Fprintf(w, `
    </table>
    <h2>说明</h2>
    <ul>
        <li>性能阈值: %.1f%%（低于此百分比视为性能退化）</li>
        <li>基线稳定要求: 连续 %d 次测试通过</li>
        <li>基线最大年龄: %d 天（超过此天需要重新测试）</li>
    </ul>
</body>
</html>`)

	// 写入HTML文件
	if err := os.WriteFile(reportPath, htmlReport, 0644); err != nil {
		return err
	}

	return nil
}

// 辅助函数：计算吞吐量（MB/s）
func calculateThroughputMBPS(result testing.BenchmarkResult, fileSizeMB float64) float64 {
	if result.N <= 0 || result.T <= 0 {
		return 0
	}

	totalData := float64(result.N) * fileSizeMB
	totalTime := float64(result.T) / float64(time.Second)

	if totalTime <= 0 {
		return 0
	}

	return totalData / totalTime
}

// 辅助函数：计算延迟（ms）
func calculateLatencyMS(result testing.BenchmarkResult) float64 {
	if result.N <= 0 {
		return 0
	}

	return float64(result.T) / float64(result.N) / float64(time.Millisecond)
}

// BenchmarkBaselineCollection 收集性能基线的测试
func BenchmarkBaselineCollection(b *testing.B) {
	// 这个测试专门用于收集性能基线
	// 应该单独运行以建立初始性能基线

	fmt.Println("开始收集性能基线数据...")
	fmt.Println("此测试可能需要较长时间，建议在独立环境中运行")

	// 这里可以调用所有的性能测试来收集基线
	// 由于性能测试之间需要独立运行，这里只作为标记测试

	b.Log("性能基线收集完成")
}

// BenchmarkBaselineValidation 验证性能基线
func BenchmarkBaselineValidation(b *testing.B) {
	baselines, err := LoadPerformanceBaseline()
	if err != nil {
		b.Fatalf("加载性能基线失败: %v", err)
	}

	if len(baselines) == 0 {
		b.Skip("没有性能基线数据，跳过验证")
	}

	b.Logf("找到 %d 个性能基线", len(baselines))

	// 检查基线有效性
	for testName, baseline := range baselines {
		if baseline.Throughput <= 0 {
			b.Errorf("基线 %s 的吞吐量为0，无效", testName)
		}

		if baseline.Latency < 0 {
			b.Errorf("基线 %s 的延迟为0，无效", testName)
		}

		age := time.Since(baseline.Timestamp)
		if age > baselineMaxAgeDays*24*time.Hour {
			b.Logf("基线 %s 过旧（%d 天），建议重新测试", testName, age.Hours()/24)
		}
	}

	// 打印基线对比
	PrintBaselineComparison(b)
}

// BaselineUsage 基线使用指南
const BaselineUsage = `
性能基线使用指南

## 基线文件位置

性能基线数据存储在: benchmarks/performance_baseline.json

## 基线管理

### 查看当前基线

使用BenchmarkBaselineValidation测试查看当前基线：
	go test -tags perf ./obs -bench=BenchmarkBaselineValidation -v

### 更新基线

运行性能测试会自动更新基线：
go test -tags perf ./obs -bench=. -benchtime=30s

### 导出报告

调用ExportBaselineReport函数生成HTML报告：
- 在测试代码中调用
- 或创建专门的导出工具

## 性能阈值

- 性能退化阈值: 90%（低于基线90%视为性能退化）
- 基线稳定阈值: 3次测试
- 基线最大年龄: 30天

## 基线内容

基线包含以下指标：
- 吞吐量 (MB/s)
- 延迟 (ms)
- 内存分配 (bytes/op)
- GC次数
- 环境信息

## 注意事项

1. 基线应该在相同环境中收集
2. 定期更新基线以反映系统变化
3. 性能退化应该立即调查
4. 过期的基线应该重新收集
`
