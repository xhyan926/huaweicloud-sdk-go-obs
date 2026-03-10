//go:build test_report

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

// TestReport 测试报告结构
type TestReport struct {
	Timestamp      time.Time                 `json:"timestamp"`
	ReportType     string                    `json:"report_type"`
	TestName       string                    `json:"test_name"`
	Summary        TestSummary               `json:"summary"`
	Results        []TestResult              `json:"results"`
	Configuration  TestConfiguration          `json:"configuration"`
}

// TestSummary 测试汇总信息
type TestSummary struct {
	TotalTests      int     `json:"total_tests"`
	PassedTests     int     `json:"passed_tests"`
	FailedTests     int     `json:"failed_tests"`
	SkippedTests    int     `json:"skipped_tests"`
	PassRate       float64 `json:"pass_rate"`
	Duration        string  `json:"duration"`
	Coverage        float64 `json:"coverage_percentage"`
}

// TestResult 单个测试结果
type TestResult struct {
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Duration    string            `json:"duration"`
	Error       string            `json:"error,omitempty"`
	Output      string            `json:"output,omitempty"`
	Metrics     TestMetrics        `json:"metrics,omitempty"`
}

// TestMetrics 测试指标
type TestMetrics struct {
	PassRate         float64  `json:"pass_rate,omitempty"`
	Coverage         float64  `json:"coverage,omitempty"`
	Throughput       float64  `json:"throughput,omitempty"`
	Latency          float64  `json:"latency,omitempty"`
	Allocations       int64    `json:"allocations,omitempty"`
	BytesPerOp       int64    `json:"bytes_per_op,omitempty"`
	MemoryUsage      int64    `json:"memory_usage,omitempty"`
}

// TestConfiguration 测试配置
type TestConfiguration struct {
	TestType       string                 `json:"test_type"`
	BuildTag       string                 `json:"build_tag"`
	ParallelLevel  int                    `json:"parallel_level,omitempty"`
	BenchTime      string                 `json:"bench_time,omitempty"`
	FuzzTime      string                 `json:"fuzz_time,omitempty"`
}

// 测试报告文件路径
const (
	testReportPath   = "benchmarks/test_reports/"
	testSummaryPath  = "benchmarks/test_summary.json"
	archivePath     = "benchmarks/test_archives/"
)

// 测试报告统计（全局）
var (
	testReports     map[string]*TestReport
	testReportsLock sync.RWMutex
)

// 初始化测试报告统计
func init() {
	testReports = make(map[string]*TestReport)
}

// CreateTestReport 创建测试报告
func CreateTestReport(reportType, testName string) *TestReport {
	timestamp := time.Now()

	report := &TestReport{
		Timestamp:     timestamp,
		ReportType:    reportType,
		TestName:      testName,
		Summary:       TestSummary{},
		Results:       []TestResult{},
		Configuration: TestConfiguration{
			TestType: reportType,
		},
	}

	// 存储到全局统计
	testReportsLock.Lock()
	defer testReportsLock.Unlock()
	testReports[reportType+"_"+testName] = report

	return report
}

// RecordUnitTestResult 记录单元测试结果
func RecordUnitTestResult(testName, status, duration string, error, metrics *TestMetrics) {
	reportType := "unit"

	testReportsLock.Lock()
	defer testReportsLock.Unlock()

	report, exists := testReports[reportType+"_"+testName]
	if !exists {
		report = CreateTestReport(reportType, testName)
	}

	result := TestResult{
		Name:     testName,
		Status:   status,
		Duration: duration,
		Error:    error,
		Metrics:  metrics,
	}

	report.Results = append(report.Results, result)

	// 更新汇总
	updateTestSummary(report, status, duration)
}

// RecordIntegrationTestResult 记录集成测试结果
func RecordIntegrationTestResult(testName, status, duration string, error, metrics *TestMetrics) {
	reportType := "integration"

	testReportsLock.Lock()
	defer testReportsLock.Unlock()

	report, exists := testReports[reportType+"_"+testName]
	if !exists {
		report = CreateTestReport(reportType, testName)
	}

	result := TestResult{
		Name:     testName,
		Status:   status,
		Duration: duration,
		Error:    error,
		Metrics:  metrics,
	}

	report.Results = append(report.Results, result)

	// 更新汇总
	updateTestSummary(report, status, duration)
}

// RecordPerformanceTestResult 记录性能测试结果
func RecordPerformanceTestResult(testName, status, duration string, error, metrics *TestMetrics) {
	reportType := "performance"

	testReportsLock.Lock()
	defer testReportsLock.Unlock()

	report, exists := testReports[reportType+"_"+testName]
	if !exists {
		report = CreateTestReport(reportType, testName)
	}

	result := TestResult{
		Name:     testName,
		Status:   status,
		Duration: duration,
		Error:    error,
		Metrics:  metrics,
	}

	report.Results = append(report.Results, result)

	// 更新汇总
	updateTestSummary(report, status, duration)
}

// RecordFuzzTestResult 记录模糊测试结果
func RecordFuzzTestResult(testName, status, duration string, error, metrics *TestMetrics) {
	reportType := "fuzz"

	testReportsLock.Lock()
	defer testReportsLock.Unlock()

	report, exists := testReports[reportType+"_"+testName]
	if !exists {
		report = CreateTestReport(reportType, testName)
	}

	result := TestResult{
		Name:     testName,
		Status:   status,
		Duration: duration,
		Error:    error,
		Metrics:  metrics,
	}

	report.Results = append(report.Results, result)

	// 更新汇总
	updateTestSummary(report, status, duration)
}

// updateTestSummary 更新测试汇总
func updateTestSummary(report *TestReport, status, duration string) {
	report.Summary.TotalTests++

	switch status {
	case "pass":
		report.Summary.PassedTests++
	case "fail":
		report.Summary.FailedTests++
	case "skip":
		report.Summary.SkippedTests++
	}

	// 更新通过率
	if report.Summary.TotalTests > 0 {
		report.Summary.PassRate = float64(report.Summary.PassedTests) / float64(report.Summary.TotalTests) * 100
	}
}

// SaveTestReport 保存测试报告
func SaveTestReport(report *TestReport) error {
	testReportsLock.Lock()
	defer testReportsLock.Unlock()

	// 序列化为JSON
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化测试报告失败: %v", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(testReportPath, 0755); err != nil {
		return fmt.Errorf("创建测试报告目录失败: %v", err)
	}

	// 保存到文件
	reportFilename := fmt.Sprintf("%s_%s_%d.json",
		report.ReportType, report.TestName, time.Now().UnixNano())
	reportPath := filepath.Join(testReportPath, reportFilename)

	if err := os.WriteFile(reportPath, reportData, 0644); err != nil {
		return fmt.Errorf("保存测试报告失败: %v", err)
	}

	return nil
}

// AggregateTestReports 聚合所有测试报告
func AggregateTestReports(reportType string) (*AggregatedTestReport, error) {
	testReportsLock.RLock()
	defer testReportsLock.RUnlock()

	aggregated := &AggregatedTestReport{
		Timestamp:     time.Now(),
		ReportType:     reportType,
		Summary:        AggregateSummary{},
		Reports:         []TestReport{},
	}

	// 聚合所有相同类型的报告
	for key, report := range testReports {
		if report.ReportType == reportType {
			aggregated.Reports = append(aggregated.Reports, *report)
			aggregated.Summary.TotalTests += report.Summary.TotalTests
			aggregated.Summary.PassedTests += report.Summary.PassedTests
			aggregated.Summary.FailedTests += report.Summary.FailedTests
			aggregated.Summary.SkippedTests += report.Summary.SkippedTests
		}
	}

	// 计算总通过率
	if aggregated.Summary.TotalTests > 0 {
		aggregated.Summary.PassRate = float64(aggregated.Summary.PassedTests) / float64(aggregated.Summary.TotalTests) * 100
	}

	return aggregated, nil
}

// AggregatedTestReport 聚合测试报告
type AggregatedTestReport struct {
	Timestamp time.Time          `json:"timestamp"`
	ReportType string            `json:"report_type"`
	Summary   AggregateSummary   `json:"summary"`
	Reports    []TestReport       `json:"reports"`
}

// AggregateSummary 聚合汇总信息
type AggregateSummary struct {
	TotalTests    int     `json:"total_tests"`
	PassedTests    int     `json:"passed_tests"`
	FailedTests    int     `json:"failed_tests"`
	SkippedTests   int     `json:"skipped_tests"`
	PassRate       float64 `json:"pass_rate"`
}

// GenerateHTMLReport 生成HTML测试报告
func GenerateHTMLReport(aggregated *AggregatedTestReport) (string, error) {
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>测试报告 - ` + aggregated.ReportType + `</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .summary { background: linear-gradient(135deg, #667eea, #764ba2); color: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .summary h1 { margin: 0 0 10px 0; }
        .summary-stats { display: flex; justify-content: space-around; }
        .stat-item { text-align: center; }
        .stat-value { font-size: 24px; font-weight: bold; }
        .stat-label { font-size: 14px; opacity: 0.9; }

        .reports { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .report-item { margin-bottom: 15px; padding-bottom: 15px; border-bottom: 1px solid #eee; }
        .report-header { display: flex; justify-content: space-between; align-items: center; }
        .report-name { font-weight: bold; font-size: 16px; }
        .report-status { padding: 4px 12px; border-radius: 4px; }
        .status-pass { background: #4CAF50; color: white; }
        .status-fail { background: #F44336; color: white; }
        .report-metrics { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-top: 10px; }
        .metric-item { display: flex; align-items: center; }
        .metric-label { font-size: 12px; color: #666; }
        .metric-value { font-weight: bold; }

        .footer { text-align: center; margin-top: 40px; padding-top: 20px; border-top: 2px solid #eee; color: #666; }
        .footer-info { font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="summary">
            <h1>测试汇总 - ` + aggregated.ReportType + `</h1>
            <div class="summary-stats">
                <div class="stat-item">
                    <div class="stat-label">总测试数</div>
                    <div class="stat-value">` + fmt.Sprintf("%d", aggregated.Summary.TotalTests) + `</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">通过测试</div>
                    <div class="stat-value" style="color: #4CAF50;">` + fmt.Sprintf("%d", aggregated.Summary.PassedTests) + `</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">失败测试</div>
                    <div class="stat-value" style="color: #F44336;">` + fmt.Sprintf("%d", aggregated.Summary.FailedTests) + `</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">跳过测试</div>
                    <div class="stat-value">` + fmt.Sprintf("%d", aggregated.Summary.SkippedTests) + `</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">通过率</div>
                    <div class="stat-value">` + fmt.Sprintf("%.2f%%", aggregated.Summary.PassRate) + `</div>
                </div>
            </div>
        </div>

        <div class="reports">
            <h2>详细报告</h2>
`

	for _, report := range aggregated.Reports {
		htmlContent += fmt.Sprintf(`
            <div class="report-item">
                <div class="report-header">
                    <div class="report-name">%s</div>
                    <div class="report-status %s">%s</div>
                </div>
`, report.TestName, getStatusClass(report.Summary))

		if len(report.Results) > 0 {
			htmlContent += `
                <div class="report-metrics">`

			// 计算平均指标
			totalPassRate := 0.0
			metricCount := 0
			for _, result := range report.Results {
				if result.Status == "pass" && result.Metrics.PassRate > 0 {
					totalPassRate += result.Metrics.PassRate
					metricCount++
				}
			}

			avgPassRate := 0.0
			if metricCount > 0 {
				avgPassRate = totalPassRate / float64(metricCount)
			}

			if metricCount > 0 {
				htmlContent += fmt.Sprintf(`
                    <div class="metric-item">
                        <div class="metric-label">平均通过率</div>
                        <div class="metric-value">%.2f%%</div>
                    </div>`, avgPassRate)
			}
		}
	}

	htmlContent += `
        </div>
    </div>

    <div class="footer">
        <div class="footer-info">
            报告生成时间: ` + aggregated.Timestamp.Format("2006-01-02 15:04:05") + `<br>
            测试类型: ` + aggregated.ReportType + ` | 总报告数: ` + fmt.Sprintf("%d", len(aggregated.Reports)) + `
        </div>
    </div>
</body>
</html>`

	// 确保目录存在
	if err := os.MkdirAll(testReportPath, 0755); err != nil {
		return "", err
	}

	// 写入HTML文件
	reportFilename := fmt.Sprintf("%s_summary_%s.html",
		aggregated.ReportType, time.Now().Format("20060102_150405"))
	reportPath := filepath.Join(testReportPath, reportFilename)

	if err := os.WriteFile(reportPath, []byte(htmlContent), 0644); err != nil {
		return "", fmt.Errorf("生成HTML报告失败: %v", err)
	}

	return reportPath, nil
}

// getStatusClass 获取状态对应的CSS类
func getStatusClass(summary TestSummary) string {
	if summary.PassRate >= 90 {
		return "status-pass"
	} else if summary.PassRate >= 70 {
		return "status-pass"
	} else if summary.PassRate >= 50 {
		return "status-fail"
	}
	return "status-fail"
}

// ArchiveTestReports 归档测试报告
func ArchiveTestReports(reportType string) error {
	testReportsLock.RLock()
	defer testReportsLock.Unlock()

	// 创建归档目录
	archiveDir := filepath.Join(archivePath, reportType, time.Now().Format("20060102"))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("创建归档目录失败: %v", err)
	}

	// 聚合并保存所有相关报告
	for key, report := range testReports {
		if report.ReportType == reportType {
			reportFilename := fmt.Sprintf("%s_%s.json", report.TestName, time.Now().UnixNano())
			reportPath := filepath.Join(archiveDir, reportFilename)

			reportData, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				return fmt.Errorf("序列化报告失败: %v", err)
			}

			if err := os.WriteFile(reportPath, reportData, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetTestReport 获取测试报告
func GetTestReport(reportType, testName string) (*TestReport, bool) {
	testReportsLock.RLock()
	defer testReportsLock.RUnlock()

	report, exists := testReports[reportType+"_"+testName]
	return report, exists
}

// GetAllTestReports 获取所有测试报告
func GetAllTestReports() map[string]*TestReport {
	testReportsLock.RLock()
	defer testReportsLock.RUnlock()

	// 返回副本
	reports := make(map[string]*TestReport)
	for key, report := range testReports {
		reports[key] = report
	}

	return reports
}

// ClearTestReports 清除测试报告
func ClearTestReports() {
	testReportsLock.Lock()
	defer testReportsLock.Unlock()

	testReports = make(map[string]*TestReport)
}

// GenerateTestSummary 生成测试总览
func GenerateTestSummary(b *testing.B) error {
	testReportsLock.RLock()
	defer testReportsLock.Unlock()

	summary := make(map[string]TestSummary)

	// 按报告类型统计
	for key, report := range testReports {
		summary[report.ReportType] = report.Summary
	}

	// 保存总览
	summaryData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化测试总览失败: %v", err)
	}

	if err := os.WriteFile(testSummaryPath, summaryData, 0644); err != nil {
		return fmt.Errorf("保存测试总览失败: %v", err)
	}

	b.Logf("测试总览已保存: %s", testSummaryPath)
	return nil
}

// RunTestReportGeneration 运行测试报告生成（基准测试）
func RunTestReportGeneration(b *testing.B) error {
	// 测试单元测试报告生成
	b.Run("生成单元测试报告", func(b *testing.B) {
		b.Log("测试单元测试报告生成功能")
	})

	// 测试集成测试报告生成
	b.Run("生成集成测试报告", func(b *testing.B) {
		b.Log("测试集成测试报告生成功能")
	})

	// 测试性能测试报告生成
	b.Run("生成性能测试报告", func(b *testing.B) {
		b.Log("测试性能测试报告生成功能")
	})

	// 测试模糊测试报告生成
	b.Run("生成模糊测试报告", func(b *testing.B) {
		b.Log("测试模糊测试报告生成功能")
	})

	// 测试HTML报告生成
	b.Run("生成HTML报告", func(b *testing.B) {
		b.Log("测试HTML报告生成功能")
	})

	// 测试报告聚合
	b.Run("测试报告聚合", func(b *testing.B) {
		b.Log("测试报告聚合功能")
	})

	// 测试报告归档
	b.Run("测试报告归档", func(b *testing.B) {
		b.Log("测试报告归档功能")
	})

	// 测试总览生成
	b.Run("生成测试总览", func(b *testing.B) {
		b.Log("测试总览生成功能")
	})

	return nil
}