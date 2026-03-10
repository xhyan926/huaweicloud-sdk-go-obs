# OBS SDK Go Makefile

# 项目信息
PROJECT_NAME := huaweicloud-sdk-go-obs
VERSION := $(shell grep "VERSION" VERSION | cut -d'=' -f2)
BUILD_TIME := $(shell date +"%Y%m%d_%H%M%S")
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Go 相关设置
GO := go
GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0

# 测试相关设置
TEST_TIMEOUT := 30s
TEST_COVERAGE := ./coverage.out
TEST_HTML_REPORT := ./coverage.html

# 默认目标
.PHONY: all
all: clean deps build test

# 清理
.PHONY: clean
clean:
	@echo "Cleaning..."
	@$(GO) clean -i ./...
	@rm -f $(TEST_COVERAGE) $(TEST_HTML_REPORT)
	@rm -rf bin/
	@rm -rf dist/

# 下载依赖
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy

# 构建项目
.PHONY: build
build:
	@echo "Building..."
	@mkdir -p bin
	@$(GO) build -v -ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GitBranch=$(GIT_BRANCH)'" -o bin/$(PROJECT_NAME) ./cmd/...

# 运行单元测试
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@$(GO) test -tags unit ./obs -v -run TestUnit

# 运行集成测试
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@if [ -z "$(OBS_TEST_AK)" ]; then \
		echo "Error: OBS_TEST_AK not set"; \
		exit 1; \
	fi
	@if [ -z "$(OBS_TEST_SK)" ]; then \
		echo "Error: OBS_TEST_SK not set"; \
		exit 1; \
	fi
	@if [ -z "$(OBS_TEST_ENDPOINT)" ]; then \
		echo "Error: OBS_TEST_ENDPOINT not set"; \
		exit 1; \
	fi
	@if [ -z "$(OBS_TEST_BUCKET)" ]; then \
		echo "Error: OBS_TEST_BUCKET not set"; \
		exit 1; \
	fi
	@$(GO) test -tags integration ./obs/test/integration -v

# 运行轻量级性能测试
.PHONY: test-perf-light
test-perf-light:
	@echo "Running light performance tests..."
	@$(GO) test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem -v

# 运行深度性能测试
.PHONY: test-perf-deep
test-perf-deep:
	@echo "Running deep performance tests..."
	@if [ -z "$(OBS_TEST_AK)" ]; then \
		echo "Error: OBS_TEST_AK not set"; \
		exit 1; \
	fi
	@if [ -z "$(OBS_TEST_SK)" ]; then \
		echo "Error: OBS_TEST_SK not set"; \
		exit 1; \
	fi
	@$(GO) test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem -v

# 运行模糊测试
.PHONY: test-fuzz
test-fuzz:
	@echo "Running fuzz tests..."
	@$(GO) test -tags fuzz ./obs -fuzz=. -fuzztime=30s

# 运行所有测试
.PHONY: test-all
test-all: test-unit test-integration test-perf-light

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GO) test -tags unit ./obs -coverprofile=$(TEST_COVERAGE) -v
	@$(GO) tool cover -html=$(TEST_COVERAGE) -o $(TEST_HTML_REPORT)
	@echo "Coverage report generated: $(TEST_HTML_REPORT)"

# 运行测试并生成报告
.PHONY: test-report
test-report:
	@echo "Generating test report..."
	@mkdir -p reports
	@$(GO) run ./cmd/test-report/main.go > reports/test-report-$(BUILD_TIME).json
	@echo "Test report generated: reports/test-report-$(BUILD_TIME).json"

# 生成测试报告
	.PHONY: test-report-generate
	test-report-generate:
		@echo "Generating test report..."
		@$(GO) test -tags test_report ./obs -bench=RunTestReportGeneration -v
		@echo "Test report generated successfully"

# 聚合测试报告
	.PHONY: test-report-aggregate
	test-report-aggregate:
		@echo "Aggregating test reports..."
		@$(GO) test -tags test_report ./obs -bench=AggregateTestReports -v
		@echo "Test reports aggregated successfully"

# 生成HTML测试总览
	.PHONY: test-html-summary
	test-html-summary:
		@echo "Generating HTML test summary..."
		@$(GO) test -tags test_report ./obs -bench=GenerateTestSummary -v
		@echo "HTML test summary generated successfully"

# 清除测试报告
	.PHONY: test-report-clean
	test-report-clean:
		@echo "Cleaning test reports..."
		@rm -f benchmarks/test_reports/*.json
		@rm -f benchmarks/test_summary.json
		@rm -f benchmarks/test_archives/**/*.json
		@echo "Test reports cleaned successfully"

# 安装到本地
.PHONY: install
install:
	@echo "Installing to GOPATH..."
	@$(GO) install ./...

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "Vetting code..."
	@$(GO) vet ./...

# 交叉构建
.PHONY: build-cross
build-cross:
	@echo "Cross building..."
	@$(GO) build -v -ldflags "-X 'main.Version=$(VERSION)'" -o bin/$(PROJECT_NAME)-linux-amd64 ./cmd/... -ldflags "-extldflags -static"
	@GOOS=darwin GOARCH=amd64 $(GO) build -v -ldflags "-X 'main.Version=$(VERSION)'" -o bin/$(PROJECT_NAME)-darwin-amd64 ./cmd/...
	@GOOS=windows GOARCH=amd64 $(GO) build -v -ldflags "-X 'main.Version=$(VERSION)'" -o bin/$(PROJECT_NAME)-windows-amd64.exe ./cmd/...

# 生成API文档
.PHONY: docs
docs:
	@echo "Generating API documentation..."
	@mkdir -p docs/api
	@$(GO) doc -all > docs/api/README.md

# 发布
.PHONY: release
release: build build-cross test-all test-coverage
	@echo "Preparing release..."
	@mkdir -p release
	@cp bin/* release/
	@cp README.md release/
	@cp LICENSE release/
	@tar -czf release/$(PROJECT_NAME)-$(VERSION)-src.tar.gz --exclude=bin --exclude=dist .
	@echo "Release prepared: release/$(PROJECT_NAME)-$(VERSION)-src.tar.gz"

# 开发环境设置
.PHONY: setup-dev
setup-dev:
	@echo "Setting up development environment..."
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@$(GO) install golang.org/x/tools/cmd/cover@latest
	@$(GO) install golang.org/x/perf/cmd/benchstat@latest
	@echo "Development environment setup complete"

# 显示帮助
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Clean, deps, build, test"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  build         - Build the project"
	@echo "  test-unit     - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-perf-light - Run light performance tests"
	@echo "  test-perf-deep - Run deep performance tests"
	@echo "  test-fuzz     - Run fuzz tests"
	@echo "  test-all      - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-report   - Generate test report"
	@echo "  install       - Install to GOPATH"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  build-cross   - Cross build"
	@echo "  docs          - Generate API documentation"
	@echo "  release       - Prepare release"
	@echo "  setup-dev     - Setup development environment"
	@echo "  help          - Show this help"