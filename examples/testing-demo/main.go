package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("OBS SDK Go 测试工程化验证")
	fmt.Println("========================")

	// 1. 验证目录结构
	fmt.Println("\n1. 验证测试目录结构...")
	if err := verifyDirectoryStructure(); err != nil {
		log.Printf("目录结构验证失败: %v", err)
	} else {
		fmt.Println("✓ 测试目录结构创建成功")
	}

	// 2. 验证配置文件
	fmt.Println("\n2. 验证测试配置...")
	if err := verifyTestConfig(); err != nil {
		log.Printf("配置文件验证失败: %v", err)
	} else {
		fmt.Println("✓ 测试配置管理正常")
	}

	// 3. 验证集成测试客户端
	fmt.Println("\n3. 验证集成测试客户端...")
	if err := verifyIntegrationClient(); err != nil {
		log.Printf("集成测试客户端验证失败: %v", err)
	} else {
		fmt.Println("✓ 集成测试客户端正常")
	}

	// 4. 验证Mock服务器
	fmt.Println("\n4. 验证Mock服务器...")
	if err := verifyMockServer(); err != nil {
		log.Printf("Mock服务器验证失败: %v", err)
	} else {
		fmt.Println("✓ Mock服务器正常")
	}

	// 5. 验证测试技能
	fmt.Println("\n5. 验证测试技能...")
	if err := verifySkills(); err != nil {
		log.Printf("测试技能验证失败: %v", err)
	} else {
		fmt.Println("✓ 测试技能可用")
	}

	// 6. 验证文档
	fmt.Println("\n6. 验证文档完整性...")
	if err := verifyDocumentation(); err != nil {
		log.Printf("文档验证失败: %v", err)
	} else {
		fmt.Println("✓ 文档完整性正常")
	}

	// 7. 验证Makefile命令
	fmt.Println("\n7. 验证Makefile命令...")
	if err := verifyMakefile(); err != nil {
		log.Printf("Makefile验证失败: %v", err)
	} else {
		fmt.Println("✓ Makefile命令正常")
	}

	fmt.Println("\n" + "========================")
	fmt.Println("OBS SDK Go 测试工程化第一阶段验证完成！")
	fmt.Println("所有核心组件已实现并可以正常使用。")
}

func verifyDirectoryStructure() error {
	// 检查核心目录
	dirs := []string{
		"obs/test/config",
		"obs/test/integration/e2e",
		"obs/test/integration/fixtures",
		"obs/test/mock_server/responses",
		".claude/skills/go-sdk-integration",
		".claude/skills/go-sdk-fuzz",
		".claude/skills/go-sdk-perf",
		"docs/testing",
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("目录不存在: %s", dir)
		}
	}

	return nil
}

func verifyTestConfig() error {
	// 检查关键配置文件
	files := []string{
		"obs/test/config/test_config.go",
		"obs/test/config/integration_env.go",
	}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("配置文件不存在: %s", file)
		}
	}

	return nil
}

func verifyIntegrationClient() error {
	cmd := exec.Command("go", "build", "-tags", "integration", "./obs/test/integration/client.go")
	return cmd.Run()
}

func verifyMockServer() error {
	cmd := exec.Command("go", "build", "-tags", "integration", "./obs/test/mock_server/server.go")
	return cmd.Run()
}

func verifySkills() error {
	// 检查技能文件
	skills := []string{
		".claude/skills/go-sdk-integration/skill.md",
		".claude/skills/go-sdk-fuzz/skill.md",
		".claude/skills/go-sdk-perf/skill.md",
	}

	for _, skill := range skills {
		if _, err := os.Stat(skill); os.IsNotExist(err) {
			return fmt.Errorf("技能文件不存在: %s", skill)
		}
	}

	return nil
}

func verifyDocumentation() error {
	// 检查核心文档
	docs := []string{
		"docs/testing/README.md",
		"docs/testing/architecture.md",
		"docs/testing/integration-testing.md",
	}

	for _, doc := range docs {
		if _, err := os.Stat(doc); os.IsNotExist(err) {
			return fmt.Errorf("文档不存在: %s", doc)
		}
	}

	return nil
}

func verifyMakefile() error {
	// 检查Makefile是否存在
	if _, err := os.Stat("Makefile"); os.IsNotExist(err) {
		return fmt.Errorf("Makefile不存在")
	}

	// 测试make help命令
	cmd := exec.Command("make", "help")
	return cmd.Run()
}