//go:build fuzz

package obs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

// FuzzingReport 模糊测试报告结构
type FuzzingReport struct {
	Timestamp      time.Time              `json:"timestamp"`
	TestName      string                 `json:"test_name"`
	TotalRuns     int64                  `json:"total_runs"`
	CrashCount    int                    `json:"crash_count"`
	HangCount     int                    `json:"hang_count"`
	MemoryPeak    int64                  `json:"memory_peak_bytes"`
	CoveragePct   float64               `json:"coverage_percentage"`
	InputSize     int64                  `json:"max_input_size"`
	Duration      time.Duration           `json:"duration_seconds"`
	Crashes       []CrashDetail          `json:"crashes"`
	Configuration FuzzingConfiguration   `json:"configuration"`
}

// CrashDetail 崩溃详情
type CrashDetail struct {
	Input          string    `json:"input"`
	StackTrace      string    `json:"stack_trace"`
	Reproducible    bool      `json:"reproducible"`
	Timestamp       time.Time `json:"timestamp"`
	InputSize       int       `json:"input_size"`
	CrashType      string    `json:"crash_type"`
}

// FuzzingConfiguration 模糊测试配置
type FuzzingConfiguration struct {
	MaxInputSize      int           `json:"max_input_size_bytes"`
	MaxDuration      time.Duration `json:"max_duration_seconds"`
	Workers          int           `json:"workers"`
	MemoryLimit      int64         `json:"memory_limit_bytes"`
	TimeoutThreshold time.Duration `json:"timeout_threshold_seconds"`
}

// FuzzingBaseline 模糊测试基线
type FuzzingBaseline struct {
	TestName      string    `json:"test_name"`
	LastRunTime  time.Time `json:"last_run_time"`
	TotalRuns     int64     `json:"total_runs"`
	CrashCount   int       `json:"crash_count"`
	IsStable      bool      `json:"is_stable"`
	Threshold     int       `json:"crash_threshold"`
}

// 模糊测试配置文件路径
const (
	fuzzConfigPath     = "benchmarks/fuzzing_config.json"
	fuzzBaselinePath   = "benchmarks/fuzzing_baseline.json"
	fuzzReportPath    = "benchmarks/fuzzing_reports/"
	fuzzCorpusPath    = "fuzzing/corpus/"
)

// 默认模糊测试配置
const (
	defaultMaxInputSize      = 1024 * 100      // 100KB
	defaultMaxDuration       = 30 * time.Minute   // 30分钟
	defaultWorkers          = 4                 // 4个工作线程
	defaultMemoryLimit      = 2 * 1024 * 1024 * 1024 // 2GB
	defaultTimeoutThreshold = 10 * time.Second  // 10秒超时
	crashStabilityThreshold = 5                  // 5次运行无崩溃认为稳定
)

// FuzzingStats 模糊测试统计（全局）
var (
	fuzzStats     map[string]*FuzzingReport
	fuzzStatsLock sync.RWMutex
)

// 初始化模糊测试统计
func init() {
	fuzzStats = make(map[string]*FuzzingReport)
}

// LoadFuzzingConfiguration 加载模糊测试配置
func LoadFuzzingConfiguration() (*FuzzingConfiguration, error) {
	config := &FuzzingConfiguration{
		MaxInputSize:      defaultMaxInputSize,
		MaxDuration:       defaultMaxDuration,
		Workers:          defaultWorkers,
		MemoryLimit:      defaultMemoryLimit,
		TimeoutThreshold: defaultTimeoutThreshold,
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(fuzzConfigPath); os.IsNotExist(err) {
		// 配置文件不存在，返回默认配置
		return config, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(fuzzConfigPath)
	if err != nil {
		return nil, fmt.Errorf("读取模糊测试配置文件失败: %v", err)
	}

	// 解析JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析模糊测试配置文件失败: %v", err)
	}

	return config, nil
}

// SaveFuzzingConfiguration 保存模糊测试配置
func SaveFuzzingConfiguration(config *FuzzingConfiguration) error {
	// 序列化为JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化模糊测试配置失败: %v", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fuzzConfigPath), 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 保存到文件
	if err := os.WriteFile(fuzzConfigPath, data, 0644); err != nil {
		return fmt.Errorf("保存模糊测试配置文件失败: %v", err)
	}

	return nil
}

// LoadFuzzingBaseline 加载模糊测试基线
func LoadFuzzingBaseline() (map[string]FuzzingBaseline, error) {
	baselines := make(map[string]FuzzingBaseline)

	// 检查基线文件是否存在
	if _, err := os.Stat(fuzzBaselinePath); os.IsNotExist(err) {
		// 基线文件不存在，返回空基线
		return baselines, nil
	}

	// 读取基线文件
	data, err := os.ReadFile(fuzzBaselinePath)
	if err != nil {
		return nil, fmt.Errorf("读取模糊测试基线文件失败: %v", err)
	}

	// 解析JSON
	var savedBaselines []FuzzingBaseline
	if err := json.Unmarshal(data, &savedBaselines); err != nil {
		return nil, fmt.Errorf("解析模糊测试基线文件失败: %v", err)
	}

	// 转换为map
	for _, baseline := range savedBaselines {
		baselines[baseline.TestName] = baseline
	}

	return baselines, nil
}

// SaveFuzzingBaseline 保存模糊测试基线
func SaveFuzzingBaseline(baselines map[string]FuzzingBaseline) error {
	// 准备保存的数据
	var savedBaselines []FuzzingBaseline
	for _, baseline := range baselines {
		savedBaselines = append(savedBaselines, baseline)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(savedBaselines, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化模糊测试基线数据失败: %v", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fuzzBaselinePath), 0755); err != nil {
		return fmt.Errorf("创建基线目录失败: %v", err)
	}

	// 保存到文件
	if err := os.WriteFile(fuzzBaselinePath, data, 0644); err != nil {
		return fmt.Errorf("保存模糊测试基线文件失败: %v", err)
	}

	return nil
}

// UpdateFuzzingBaseline 更新模糊测试基线
func UpdateFuzzingBaseline(testName string, totalRuns int64, crashCount int) error {
	baselines, err := LoadFuzzingBaseline()
	if err != nil {
		return err
	}

	// 获取或创建基线
	baseline, exists := baselines[testName]
	if !exists {
		// 创建新的基线
		baseline = FuzzingBaseline{
			TestName:      testName,
			LastRunTime:  time.Now(),
			TotalRuns:     totalRuns,
			CrashCount:   crashCount,
			IsStable:      false,
			Threshold:     crashStabilityThreshold,
		}
	} else {
		// 更新现有基线
		baseline.LastRunTime = time.Now()
		baseline.TotalRuns += totalRuns
		baseline.CrashCount = crashCount
		baseline.IsStable = crashCount == 0 // 只有当崩溃数为0时才认为稳定
	}

	// 保存更新后的基线
	baselines[testName] = baseline
	return SaveFuzzingBaseline(baselines)
}

// SaveFuzzingReport 保存模糊测试报告
func SaveFuzzingReport(report *FuzzingReport) error {
	// 生成报告文件名
	reportFilename := fmt.Sprintf("%s_%d.json", report.TestName, time.Now().UnixNano())
	reportPath := filepath.Join(fuzzReportPath, reportFilename)

	// 序列化为JSON
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化模糊测试报告失败: %v", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(reportPath), 0755); err != nil {
		return fmt.Errorf("创建报告目录失败: %v", err)
	}

	// 保存到文件
	if err := os.WriteFile(reportPath, reportData, 0644); err != nil {
		return fmt.Errorf("保存模糊测试报告失败: %v", err)
	}

	return nil
}

// GenerateFuzzingHTMLReport 生成模糊测试HTML报告
func GenerateFuzzingHTMLReport(reports []*FuzzingReport) error {
	// 生成HTML报告
	htmlReport := `<!DOCTYPE html>
<html>
<head>
    <title>模糊测试报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #4CAF50; color: white; }
        .crash { background-color: #FF6B6B; color: white; }
        .stable { background-color: #6CAF66; color: white; }
        .warning { background-color: #FFB74D; color: white; }
        .metric { font-weight: bold; }
    </style>
</head>
<body>
    <h1>OBS SDK Go 模糊测试报告</h1>
    <p>生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
    <h2>测试概览</h2>
    <table>
        <tr>
            <th>测试名称</th>
            <th>总运行次数</th>
            <th>崩溃次数</th>
            <th>崩溃率</th>
            <th>内存峰值</th>
            <th>状态</th>
        </tr>
`

	totalRuns := int64(0)
	totalCrashes := 0

	for _, report := range reports {
		crashRate := float64(0)
		if report.TotalRuns > 0 {
			crashRate = float64(report.CrashCount) / float64(report.TotalRuns) * 100
		}

		statusClass := "stable"
		if report.CrashCount > 0 {
			statusClass = "crash"
		}

		htmlReport += fmt.Sprintf(`
        <tr class="%s">
            <td class="metric">%s</td>
            <td>%d</td>
            <td>%d</td>
            <td>%.2f%%</td>
            <td>%d MB</td>
            <td>%s</td>
        </tr>`, statusClass, report.TestName, report.TotalRuns, report.CrashCount, crashRate, report.MemoryPeak/(1024*1024), getStatusText(report.CrashCount))

		totalRuns += report.TotalRuns
		totalCrashes += report.CrashCount
	}

	htmlReport += `
    </table>
    <h2>统计汇总</h2>
    <ul>
        <li>总测试数: <strong>` + fmt.Sprintf("%d", len(reports)) + `</strong></li>
        <li>总运行次数: <strong>` + fmt.Sprintf("%d", totalRuns) + `</strong></li>
        <li>总崩溃次数: <strong>` + fmt.Sprintf("%d", totalCrashes) + `</strong></li>
        <li>平均崩溃率: <strong>` + fmt.Sprintf("%.4f%%", float64(totalCrashes)/float64(totalRuns)*100) + `</strong></li>
    </ul>

    <h2>崩溃详情</h2>
`

	hasCrashes := false
	for _, report := range reports {
		if len(report.Crashes) > 0 {
			hasCrashes = true
			htmlReport += `<h3>` + report.TestName + `</h3>
    <table>
        <tr>
            <th>输入大小</th>
            <th>崩溃类型</th>
            <th>崩溃时间</th>
            <th>可重现</th>
        </tr>
`

			for _, crash := range report.Crashes {
				htmlReport += fmt.Sprintf(`
        <tr>
            <td>%d bytes</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
        </tr>`, crash.InputSize, crash.CrashType, crash.Timestamp.Format("2006-01-02 15:04:05"), getYesNo(crash.Reproducible))
			}

			htmlReport += `
    </table>`
		}
	}

	if !hasCrashes {
		htmlReport += `<p>🎉 未发现崩溃！所有测试通过。</p>`
	}

	htmlReport += `
    <h2>配置说明</h2>
    <ul>
        <li>最大输入大小: <strong>` + fmt.Sprintf("%d KB", defaultMaxInputSize/1024) + `</strong></li>
        <li>最大测试时长: <strong>` + fmt.Sprintf("%.0f 分钟", defaultMaxDuration.Minutes()) + `</strong></li>
        <li>工作线程数: <strong>` + fmt.Sprintf("%d", defaultWorkers) + `</strong></li>
        <li>内存限制: <strong>` + fmt.Sprintf("%.2f GB", float64(defaultMemoryLimit)/(1024*1024*1024)) + `</strong></li>
        <li>超时阈值: <strong>` + fmt.Sprintf("%.0f 秒", defaultTimeoutThreshold.Seconds()) + `</strong></li>
    </ul>
</body>
</html>`

	// 写入HTML文件
	reportPath := filepath.Join(fuzzReportPath, "fuzzing_summary.html")
	if err := os.WriteFile(reportPath, []byte(htmlReport), 0644); err != nil {
		return fmt.Errorf("保存模糊测试HTML报告失败: %v", err)
	}

	return nil
}

// RecordFuzzingCrash 记录模糊测试崩溃
func RecordFuzzingCrash(testName, input, stackTrace, crashType string, reproducible bool, inputSize int) error {
	fuzzStatsLock.Lock()
	defer fuzzStatsLock.Unlock()

	// 获取或创建报告
	report, exists := fuzzStats[testName]
	if !exists {
		report = &FuzzingReport{
			Timestamp:    time.Now(),
			TestName:     testName,
			Crashes:     []CrashDetail{},
			Configuration: FuzzingConfiguration{
				MaxInputSize:      defaultMaxInputSize,
				MaxDuration:       defaultMaxDuration,
				Workers:          defaultWorkers,
				MemoryLimit:      defaultMemoryLimit,
				TimeoutThreshold: defaultTimeoutThreshold,
			},
		}
		fuzzStats[testName] = report
	}

	// 添加崩溃详情
	crashDetail := CrashDetail{
		Input:       input,
		StackTrace:   stackTrace,
		Reproducible: reproducible,
		Timestamp:    time.Now(),
		InputSize:    inputSize,
		CrashType:   crashType,
	}

	report.Crashes = append(report.Crashes, crashDetail)
	report.CrashCount++

	return nil
}

// UpdateFuzzingStats 更新模糊测试统计
func UpdateFuzzingStats(testName string, totalRuns int64, memoryPeak int64, duration time.Duration) {
	fuzzStatsLock.Lock()
	defer fuzzStatsLock.Unlock()

	// 获取或创建报告
	report, exists := fuzzStats[testName]
	if !exists {
		report = &FuzzingReport{
			Timestamp:    time.Now(),
			TestName:     testName,
			Crashes:     []CrashDetail{},
			Configuration: FuzzingConfiguration{
				MaxInputSize:      defaultMaxInputSize,
				MaxDuration:       defaultMaxDuration,
				Workers:          defaultWorkers,
				MemoryLimit:      defaultMemoryLimit,
				TimeoutThreshold: defaultTimeoutThreshold,
			},
		}
		fuzzStats[testName] = report
	}

	// 更新统计
	report.TotalRuns = totalRuns
	report.MemoryPeak = memoryPeak
	report.Duration = duration
	report.Timestamp = time.Now()

	return nil
}

// GetFuzzingStats 获取模糊测试统计
func GetFuzzingStats(testName string) (*FuzzingReport, bool) {
	fuzzStatsLock.RLock()
	defer fuzzStatsLock.RUnlock()

	report, exists := fuzzStats[testName]
	return report, exists
}

// ClearFuzzingStats 清除模糊测试统计
func ClearFuzzingStats() {
	fuzzStatsLock.Lock()
	defer fuzzStatsLock.Unlock()

	fuzzStats = make(map[string]*FuzzingReport)
}

// CompareWithBaseline 与基线对比模糊测试结果
func CompareWithBaseline(testName string, currentCrashCount int) (string, error) {
	baselines, err := LoadFuzzingBaseline()
	if err != nil {
		return "", err
	}

	// 检查是否有基线
	baseline, hasBaseline := baselines[testName]
	if !hasBaseline {
		return fmt.Sprintf("没有找到 %s 的基线数据", testName), nil
	}

	// 比较崩溃次数
	if currentCrashCount > baseline.CrashCount {
		return fmt.Sprintf("⚠️ 警告: 崩溃次数增加 (之前: %d, 当前: %d)", baseline.CrashCount, currentCrashCount), nil
	} else if currentCrashCount < baseline.CrashCount {
		return fmt.Sprintf("✅ 改进: 崩溃次数减少 (之前: %d, 当前: %d)", baseline.CrashCount, currentCrashCount), nil
	} else {
		return fmt.Sprintf("✅ 稳定: 崩溃次数保持不变 (%d)", currentCrashCount), nil
	}
}

// RunFuzzingBaselineValidation 验证模糊测试基线
func RunFuzzingBaselineValidation(b *testing.B) {
	baselines, err := LoadFuzzingBaseline()
	if err != nil {
		b.Fatalf("加载模糊测试基线失败: %v", err)
	}

	if len(baselines) == 0 {
		b.Skip("没有模糊测试基线数据")
	}

	b.Logf("找到 %d 个模糊测试基线", len(baselines))

	// 检查基线有效性
	for testName, baseline := range baselines {
		if baseline.TotalRuns < 0 {
			b.Errorf("基线 %s 的总运行次数无效", testName)
		}

		if baseline.CrashCount < 0 {
			b.Errorf("基线 %s 的崩溃次数无效", testName)
		}

		age := time.Since(baseline.LastRunTime)
		if age > 7*24*time.Hour { // 超过7天
			b.Logf("基线 %s 过旧（%.1f 天），建议重新测试", testName, age.Hours()/24)
		}

		if baseline.IsStable {
			b.Logf("基线 %s 稳定 (无崩溃)", testName)
		} else {
			b.Logf("基线 %s 不稳定 (崩溃次数: %d)", testName, baseline.CrashCount)
		}
	}

	// 打印基线对比
	PrintFuzzingBaselineComparison(b)
}

// PrintFuzzingBaselineComparison 打印基线对比结果
func PrintFuzzingBaselineComparison(b *testing.B) {
	baselines, err := LoadFuzzingBaseline()
	if err != nil {
		b.Logf("无法加载模糊测试基线: %v", err)
		return
	}

	if len(baselines) == 0 {
		b.Log("没有模糊测试基线数据")
		return
	}

	b.Logf("=== 模糊测试基线对比 ===")
	b.Logf("基线数据来源: %s", fuzzBaselinePath)
	b.Logf("基线数量: %d", len(baselines))
	b.Logf("稳定性阈值: %d 次运行无崩溃", crashStabilityThreshold)

	for testName, baseline := range baselines {
		b.Logf("\n基线: %s", testName)
		b.Logf("  最后运行时间: %s", baseline.LastRunTime.Format("2006-01-02 15:04:05"))
		b.Logf("  总运行次数: %d", baseline.TotalRuns)
		b.Logf("  崩溃次数: %d", baseline.CrashCount)
		b.Logf("  是否稳定: %t", baseline.IsStable)
	}
}

// PrintFuzzingConfiguration 打印模糊测试配置
func PrintFuzzingConfiguration(b *testing.B) {
	config, err := LoadFuzzingConfiguration()
	if err != nil {
		b.Fatalf("加载模糊测试配置失败: %v", err)
	}

	b.Logf("=== 模糊测试配置 ===")
	b.Logf("最大输入大小: %d bytes (%.2f KB)", config.MaxInputSize, float64(config.MaxInputSize)/1024)
	b.Logf("最大测试时长: %.0f 分钟", config.MaxDuration.Minutes())
	b.Logf("工作线程数: %d", config.Workers)
	b.Logf("内存限制: %.2f GB", float64(config.MemoryLimit)/(1024*1024*1024))
	b.Logf("超时阈值: %.0f 秒", config.TimeoutThreshold.Seconds())
}

// MonitorFuzzingMemory 监控模糊测试内存使用
func MonitorFuzzingMemory(testName string) (int64, func() int64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	peak := m.Alloc

	return peak, func() int64 {
		runtime.ReadMemStats(&m)
		if m.Alloc > peak {
			peak = m.Alloc
		}
		return peak
	}
}

// SaveFuzzingCorpus 保存模糊测试语料库
func SaveFuzzingCorpus(testName string, inputs []string) error {
	if len(inputs) == 0 {
		return nil
	}

	// 创建语料库目录
	corpusDir := filepath.Join(fuzzCorpusPath, testName)
	if err := os.MkdirAll(corpusDir, 0755); err != nil {
		return fmt.Errorf("创建语料库目录失败: %v", err)
	}

	// 保存每个输入
	for i, input := range inputs {
		filename := fmt.Sprintf("input_%d.txt", i)
		filePath := filepath.Join(corpusDir, filename)
		if err := os.WriteFile(filePath, []byte(input), 0644); err != nil {
			return fmt.Errorf("保存语料库文件失败: %v", err)
		}
	}

	return nil
}

// LoadFuzzingCorpus 加载模糊测试语料库
func LoadFuzzingCorpus(testName string) ([]string, error) {
	corpusDir := filepath.Join(fuzzCorpusPath, testName)

	// 检查目录是否存在
	if _, err := os.Stat(corpusDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// 读取目录中的文件
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		return nil, fmt.Errorf("读取语料库目录失败: %v", err)
	}

	var inputs []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(corpusDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("读取语料库文件失败: %v", err)
		}

		inputs = append(inputs, string(data))
	}

	return inputs, nil
}

// ===== 辅助函数 =====

// getStatusText 根据崩溃次数获取状态文本
func getStatusText(crashCount int) string {
	if crashCount == 0 {
		return "✅ 稳定"
	} else if crashCount < 5 {
		return "⚠️ 轻微"
	} else if crashCount < 10 {
		return "🔴 严重"
	} else {
		return "❌ 危险"
	}
}

// getYesNo 将布尔值转换为是/否文本
func getYesNo(value bool) string {
	if value {
		return "是"
	}
	return "否"
}