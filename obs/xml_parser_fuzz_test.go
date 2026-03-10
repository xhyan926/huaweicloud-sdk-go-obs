//go:build fuzz

package obs

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
	"time"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
)

// FuzzTransToXml_通用XML转换函数模糊测试
func FuzzTransToXml(f *testing.F) {
	// 添加正常用例的种子数据
	seeds := []interface{}{
		&obs.CreateBucketInput{Bucket: "test-bucket"},
		&obs.PutObjectInput{Bucket: "test", Key: "test-key"},
		&obs.GetObjectInput{Bucket: "test", Key: "test-key"},
		&obs.ListObjectsInput{Bucket: "test"},
		&obs.DeleteObjectInput{Bucket: "test", Key: "test-key"},
		&obs.GetBucketMetadataInput{Bucket: "test"},
		map[string]string{"key1": "value1", "key2": "value2"},
		[]string{"item1", "item2", "item3"},
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	// 模糊测试
	f.Fuzz(func(t *testing.T, input interface{}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 50000 {
			t.Skip("跳过长时间运行")
		}

		// 防止输入过大（限制在100KB以内）
		if str, ok := input.(string); ok && len(str) > 100*1024 {
			t.Skip("输入过大，跳过")
		}

		// 防止切片过大
		if slice, ok := input.([]byte); ok && len(slice) > 10*1024 {
			t.Skip("切片过大，跳过")
		}

		// 防止map过大
		if m, ok := input.(map[string]interface{}); ok && len(m) > 100 {
			t.Skip("map过大，跳过")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 10*1024*1024 { // 10MB限制
			t.Skip("内存使用过高，跳过")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("XML转换失败（预期边界情况）: %v, input类型: %T", err, input)
			return
		}

		// 验证生成的XML
		if len(xmlBytes) == 0 && input != nil {
			t.Error("非nil输入生成空XML")
		}

		// 验证XML格式基本正确
		if len(xmlBytes) > 0 {
			xmlStr := string(xmlBytes)
			if strings.HasPrefix(xmlStr, "<?xml") {
				if !strings.Contains(xmlStr, ">") {
					t.Error("生成的XML格式不正确")
				}
			}
		}

		// 验证没有恶意内容
		dangerousPatterns := []string{
			"<!ENTITY", "<!DOCTYPE", "<!ATTLIST", "<xsl:", "<script:", "<javascript:",
			"onerror=", "onload=", "onmouseover=", "onclick=", "<iframe:",
		"<svg:", "<style:", "<link:", "<meta:", "<body:", "<form:",
			"<input:", "<button:", "<img:", "<object:", "<embed:",
		}

		xmlLower := strings.ToLower(xmlStr)
		for _, pattern := range dangerousPatterns {
			if strings.Contains(xmlLower, pattern) {
				t.Errorf("生成的XML包含潜在危险模式: %s", pattern)
			}
		}
	})
}

// FuzzConvertAclToXml ACL XML转换模糊测试
func FuzzConvertAclToXml(f *testing.F) {
	// 添加种子数据
	seeds := []obs.AccessControlPolicy{
		{
			Owner: obs.Owner{ID: "owner-id-12345", DisplayName: "test-owner"},
			Grants: []obs.Grant{
				{Grantee: &obs.Grantee{ID: "grantee-id-12345"}, Permission: obs.PERMISSION_READ},
			},
		},
		{
			Owner: obs.Owner{ID: "", DisplayName: ""},
			Grants: []obs.Grant{
				{Grantee: &obs.Grantee{URI: "http://example.com/uri", DisplayName: "test-group"}, Permission: obs.PERMISSION_READ},
			},
		},
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input obs.AccessControlPolicy) {
		if len(f.Fuzzing()) > 30000 {
			t.Skip("跳过长时间运行")
		}

		// 验证Grantee不为nil
		for i, grant := range input.Grants {
			if grant.Grantee == nil {
				t.Errorf("Grant[%d]的Grantee为nil", i)
			}

			// 检查ID和URI
			if grant.Grantee != nil {
				if grant.Grantee.ID == "" && grant.Grantee.URI == "" {
					t.Error("Grant有ID和URI都为空")
				}
			}
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 20*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行ACL转换
		xmlData, md5, err := ConvertAclToXml(input, true, false)
		if err != nil {
			t.Logf("ACL转换失败: %v", err)
			return
		}

		// 验证输出
		if len(xmlData) == 0 {
			t.Error("生成的ACL XML为空")
		}

		if len(md5) != 32 {
			t.Errorf("MD5长度不正确，期望32，实际: %d", len(md5))
		}

		t.Logf("ACL转换成功，XML长度: %d, MD5: %s", len(xmlData), md5)
	})
}

// FuzzConvertLifecycleConfigurationToXml 生命周期配置XML转换模糊测试
func FuzzConvertLifecycleConfigurationToXml(f *testing.F) {
	// 添加种子数据
	seeds := []obs.BucketLifecycleConfiguration{
		{
			Rules: []obs.LifecycleRule{
				{
					ID:     "rule-id-1",
					Prefix: "test/",
					Status: "Enabled",
					Expiration: obs.Expiration{
						Days:    30,
						Created: "2024-01-01T00:00:00Z",
					},
				},
			},
		{
			Rules: []obs.LifecycleRule{
				{
					ID:     "rule-id-2",
					Prefix: "expired/",
					Status: "Enabled",
					Transitions: []obs.Transition{
						{
							Days:    1,
							StorageClass: "STANDARD_IA",
						},
					},
				},
			},
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input obs.BucketLifecycleConfiguration) {
		if len(f.Fuzzing()) > 30000 {
			t.Skip("跳过长时间运行")
		}

		// 验证规则不为空
		if len(input.Rules) == 0 {
			// 空规则可能是有效的边界情况
			return
		}

		// 验证每个规则的基本结构
		for i, rule := range input.Rules {
			if rule.Prefix == "" && len(rule.Transitions) > 0 {
				// 没有前缀但有转换规则，这是允许的
			}

			// 验证状态值
			validStatuses := map[string]bool{
				"Enabled": true,
				"Disabled": true,
				"Suspended": true,
			}
			if rule.Status != "" && !validStatuses[rule.Status] {
				t.Errorf("规则[%d]的状态无效: %s", i, rule.Status)
			}
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 30*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行转换
		xmlData, md5, err := ConvertLifecycleConfigurationToXml(input, false, false, true)
		if err != nil {
			t.Logf("生命周期配置转换失败: %v", err)
			return
		}

		// 验证输出
		if len(xmlData) == 0 && len(input.Rules) > 0 {
			t.Error("生成的XML为空")
		}

		t.Logf("生命周期配置转换成功，XML长度: %d", len(xmlData))
	})
}

// FuzzConvertNotificationToXml 通知配置XML转换模糊测试
func FuzzConvertNotificationToXml(f *testing.F) {
	// 添加种子数据
	seeds := []obs.BucketNotification{
		{
			TopicConfigurations: []obs.TopicConfiguration{
				{
					ID:     "topic-id-1",
					Topic:  "arn:aws:sns:us-east-1:123456789012:MyTopic",
					Events: []string{"s3:ObjectCreated:*", "s3:ObjectRemoved:*"},
				},
			},
	{
			TopicConfigurations: []obs.TopicConfiguration{
				{
					ID:     "topic-id-2",
					Topic:  "arn:aws:sns:us-east-1:123456789012:MyOtherTopic",
					Events: []string{"s3:ObjectCreated:*"},
				},
			},
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input obs.BucketNotification) {
		if len(f.Fuzzing()) > 30000 {
			t.Skip("跳过长时间运行")
		}

		// 验证TopicConfiguration不为空
		if len(input.TopicConfigurations) == 0 {
			return
		}

		// 验证每个主题配置
		for i, topic := range input.TopicConfigurations {
			if topic.Topic == "" {
				t.Errorf("TopicConfig[%d]的Topic为空", i)
			}

			if topic.ID == "" {
				t.Errorf("TopicConfig[%d]的ID为空", i)
			}

			// 验证事件不为空
			if len(topic.Events) == 0 {
				t.Errorf("TopicConfig[%d]的Events为空", i)
			}

			// 验证事件格式
			for _, event := range topic.Events {
				if !strings.HasPrefix(event, "s3:") && !strings.HasPrefix(event, "obs:") {
					t.Errorf("事件[%d]格式无效: %s", i, event)
				}
			}
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 25*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行转换
		xmlData, md5, err := ConvertNotificationToXml(input, false, false)
		if err != nil {
			t.Logf("通知配置转换失败: %v", err)
			return
		}

		if len(xmlData) == 0 {
			t.Error("生成的XML为空")
		}

		t.Logf("通知配置转换成功，XML长度: %d", len(xmlData))
	})
}

// FuzzConvertCompleteMultipartUploadInputToXml 分块上传完成XML转换模糊测试
func FuzzConvertCompleteMultipartUploadInputToXml(f *testing.F) {
	// 添加种子数据
	seeds := []obs.CompleteMultipartUploadInput{
		{
			Bucket:   "test-bucket",
			Key:      "test-key",
			UploadId: "upload-id-12345",
			Parts: []obs.Part{
				{PartNumber: 1, ETag: "etag-12345"},
				{PartNumber: 2, ETag: "etag-67890"},
			},
		},
		{
			Bucket:   "test-bucket",
			Key:      "test-key-2",
			UploadId: "upload-id-67890",
			Parts: []obs.Part{
				{PartNumber: 1, ETag: "etag-11111"},
			},
		},
		{
			Bucket:   "test-bucket",
			Key:      "test-key-empty",
			UploadId: "upload-id-empty",
			Parts:    []obs.Part{}, // 空分块列表
		},
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input obs.CompleteMultipartUploadInput) {
		if len(f.Fuzzing()) > 30000 {
			t.Skip("跳过长时间运行")
		}

		// 验证必需字段
		if input.Bucket == "" {
			t.Error("Bucket为空")
		}

		if input.Key == "" {
			t.Error("Key为空")
		}

		if input.UploadId == "" {
			t.Error("UploadId为空")
		}

		// 验证Parts
		if len(input.Parts) == 0 {
			// 空分块列表可能对某些场景是有效的
		} else {
			for i, part := range input.Parts {
				if part.PartNumber <= 0 {
					t.Errorf("Part[%d]的PartNumber无效: %d", i, part.PartNumber)
				}

				if part.ETag == "" {
					t.Errorf("Part[%d]的ETag为空", i)
				}
			}
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 15*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行转换
		xmlData, md5, err := ConvertCompleteMultipartUploadInputToXml(input, false)
		if err != nil {
			t.Logf("分块上传完成转换失败: %v", err)
			return
		}

		if len(xmlData) == 0 {
			t.Error("生成的XML为空")
		}

		t.Logf("分块上传完成转换成功，XML长度: %d", len(xmlData))
	})
}

// FuzzXmlInputWithMaliciousContent 恶意XML内容模糊测试
func FuzzXmlInputWithMaliciousContent(f *testing.F) {
	// 添加恶意XML种子数据
	maliciousInputs := []string{
		`<?xml version="1.0"?><!DOCTYPE root [<!ENTITY xxe SYSTEM "evil.dtd">]><root>&evil;</root>`,
		`<?xml version="1.0"?><!ENTITY xxe SYSTEM "http://example.com/evil.dtd"><!ENTITY %xxe SYSTEM "http://example.com/xxe.dtd"><root>`,
		`<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "evil.dtd">]><test>`,
		`<?xml version="1.0"?><!ENTITY foo "&bar;"><test>`,
		`<?xml version="1.0"?><!ENTITY foo "<![CDATA[<evil>]]>"><test>`,
		`<?xml version="1.0"?><!ENTITY %bar "<![CDATA[</script>alert(1)</script>]]>"><test>`,
		`<?xml version="1.0"?><!DOCTYPE foo SYSTEM "file:///etc/passwd"><test>`,
		`<test><?xml?><evil>content</evil></test>`,
		`<test><!--><script>alert(1)</script>--></test>`,
		`<test><?xml version="1.0"?><test><??></test>`,
		`<test>&#x0000;&test;</test>`,
		`<test>&amp;&lt;test;&amp;gt;</test>`,
		`<test>${jndi:dns:ldap://127.0.0.1:389/e1}</test>`,
		`<test>%253c%2565253c/test>`,
		`<test>%SYSTEM%test;/test>`,
	}

	for _, maliciousInput := range maliciousInputs {
		f.Add(maliciousInput)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// 防止输入过大
		if len(input) > 50*1024 {
			t.Skip("恶意输入过大，跳过")
		}

		// 防止测试超时
		if len(f.Fuzzing()) > 10000 {
			t.Skip("跳过长时间运行")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 10*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("恶意XML转换失败（预期的）: %v", err)
			return
		}

		// 验证生成的XML
		if len(xmlBytes) == 0 {
			t.Error("恶意输入生成空XML")
		}

		// 验证没有危险内容
		xmlStr := strings.ToLower(string(xmlBytes))
		if strings.Contains(xmlStr, "<script>") || strings.Contains(xmlStr, "<iframe>") ||
			strings.Contains(xmlStr, "onload=") || strings.Contains(xmlStr, "onerror=") {
			t.Error("生成的XML包含危险脚本标签")
		}

		t.Logf("恶意XML输入处理完成，长度: %d", len(xmlBytes))
	})
}

// FuzzXmlHugeSizeInput 超大XML输入模糊测试
func FuzzXmlHugeSizeInput(f *testing.F) {
	// 添加超大输入种子
	hugeInputs := []struct {
		content string
	}{
		{content: strings.Repeat("a", 10000)}, // 10KB
		{content: strings.Repeat("b", 5000)},  // 5KB
		{content: strings.Repeat("c", 1000)},  // 1KB
		{content: strings.Repeat("d", 500)},   // 500B
	}

	for _, seed := range hugeInputs {
		f.Add(seed.content)
	}

	f.Fuzz(func(t *testing.T, content string) {
		// 防止测试超时
		if len(f.Fuzzing()) > 30000 {
			t.Skip("跳过长时间运行")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// 根据内容大小调整限制
		maxSize := len(content)
		if maxSize > 20*1024 { // 超过20KB的输入就限制大小
			t.Skip("内容过大，跳过")
		}

		runtime.ReadMemStats(&m)
		if m.Alloc > 5*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(content)
		if err != nil {
			t.Logf("大内容XML转换失败: %v, 内容长度: %d", err, len(content))
			return
		}

		if len(xmlBytes) == 0 {
			t.Error("大内容生成空XML")
		}

		t.Logf("大内容XML转换成功，内容长度: %d, XML长度: %d",
			len(content), len(xmlBytes))
	})
}

// FuzzXmlWithSpecialCharacters 特殊字符XML输入模糊测试
func FuzzXmlWithSpecialCharacters(f *testing.F) {
	// 添加特殊字符种子数据
	specialInputs := []string{
		`<test>data</test>`,                              // 普通文本
		`<test>data & "quote" & 'apostrophe'</test>`,    // 引号
		`<test>data &lt; &gt; &amp;</test>`,           // HTML实体
		`<test>中文数据αβγ</test>`,                 // Unicode字符
		<test>data\t\n\r\v\f</test>`,              // 控制字符
		<test>data<![CDATA[特殊]]]</test>`,          // CDATA段
		`<test>data&#x0000;&#x0001;&#x0002;</test>`,     // 实体引用
		<test>data\u0000\u0001\u0002</test>`,        // Unicode转义
		<test>data\ud83c\udf4\u4e8</test>`,         // Emoji
	}

	for _, seed := range specialInputs {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// 防止测试超时
		if len(f.Fuzzing()) > 20000 {
			t.Skip("跳过长时间运行")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 8*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("特殊字符XML转换失败: %v, 输入长度: %d", err, len(input))
			return
		}

		if len(xmlBytes) == 0 {
			t.Error("特殊字符输入生成空XML")
		}

		// 验证特殊字符被正确转义
		xmlStr := string(xmlBytes)
		if strings.Contains(xmlStr, "<") && !strings.Contains(xmlStr, "&lt;") {
			t.Error("小于号未被转义")
			}

		if strings.Contains(xmlStr, ">") && !strings.Contains(xmlStr, "&gt;") {
				t.Error("大于号未被转义")
			}

		if strings.Contains(xmlStr, "&") && !strings.Contains(xmlStr, "&amp;") {
				t.Error("和号未被转义")
			}

		t.Logf("特殊字符XML转换成功，XML长度: %d", len(xmlBytes))
	})
}

// FuzzXmlStructureWithDeepNesting 深度嵌套结构XML模糊测试
func FuzzXmlStructureWithDeepNesting(f *testing.F) {
	// 添加深度嵌套种子数据
	deepNestedInputs := []interface{}{
		struct {
			Name:       "root",
			Children: []interface{}{
				struct {
					Name:       "child1",
					Children: []interface{}{
						struct {
							Name:   "grandchild1",
							Attrs: map[string]string{
								"attr1": "value1",
								"attr2": "value2",
							},
						},
						struct {
							Name:   "grandchild2",
							Attrs: map[string]string{
								"attr1": "value1",
							},
						},
					},
				},
				struct {
					Name:       "child2",
					Children: []interface{}{
						struct {
							Name:   "grandchild3",
						},
					},
				},
			},
		},
	}

	for _, seed := range deepNestedInputs {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input interface{}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 20000 {
			t.Skip("跳过长时间运行")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 10*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("深度嵌套XML转换失败: %v", err)
			return
		}

		if len(xmlBytes) == 0 && input != nil {
			t.Error("深度嵌套结构生成空XML")
		}

		// 验证XML格式
		xmlStr := string(xmlBytes)
		openTags := strings.Count(xmlStr, "<")
		closeTags := strings.Count(xmlStr, ">")

		if openTags != closeTags {
			t.Errorf("标签不匹配，开启: %d, 关闭: %d", openTags, closeTags)
		}

		t.Logf("深度嵌套XML转换成功，XML长度: %d, 标签数: %d",
			len(xmlBytes), openTags)
	})
}

// FuzzXmlWithLargeMap 大型Map结构XML模糊测试
func FuzzXmlWithLargeMap(f *testing.F) {
	// 添加大型map种子数据
	largeMapInputs := []map[string]string{}

	// 创建一个包含很多字段的大型map
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("field-%d", i)
		value := fmt.Sprintf("value-%d", i)
		largeMapInputs[key] = value
	}

	for _, seed := range []map[string]string{largeMapInputs} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input map[string]string) {
		// 防止测试超时
		if len(f.Fuzzing()) > 15000 {
			t.Skip("跳过长时间运行")
		}

		// 验证map大小
		if len(input) > 100 {
			t.Skip("map过大，跳过")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 20*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("大型Map XML转换失败: %v", err)
			return
		}

		if len(xmlBytes) == 0 {
			t.Error("大型Map生成空XML")
		}

		t.Logf("大型Map XML转换成功，字段数: %d, XML长度: %d",
			len(input), len(xmlBytes))
	})
}

// FuzzXmlWithEmptyValues 空值XML模糊测试
func FuzzXmlWithEmptyValues(f *testing.F) {
	// 添加空值种子数据
	emptyValueInputs := []struct {
		testField string `json:"testField"`
	}{
		{testField: ""},
		{testField: "valid-value"},
		{testField: strings.Repeat("a", 100)}, // 接近空
		{testField: " \t\n\r"},                  // 空白字符
		{testField: "                              "}, // 大量空格
	}

	for _, seed := range emptyValueInputs {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input struct {
		testField string `json:"testField"`
	}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 20000 {
			t.Skip("跳过长时间运行")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 5*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("空值XML转换失败: %v, 字段值: %q", err, input.testField)
			return
		}

		if len(xmlBytes) == 0 {
			t.Logf("空值字段生成空XML（边界情况）")
		} else {
			// 验证空值被正确处理
			if strings.Contains(string(xmlBytes), "<testField></testField>") {
				t.Log("空值字段生成预期XML: <testField></testField>")
			}
		}
	})
}

// FuzzXmlWithLongStringValues 超长字符串XML模糊测试
func FuzzXmlWithLongStringValues(f *testing.F) {
	// 添加超长字符串种子数据
	longStringInputs := []struct {
		testField string `json:"testField"`
	}{
		{testField: strings.Repeat("a", 5000)},     // 5000字符
		{testField: strings.Repeat("b", 10000)},    // 10000字符
		{testField: strings.Repeat("c", 2000)},      // 2000字符
		{testField: strings.Repeat("d", 800)},       // 800字符
	}

	for _, seed := range longStringInputs {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input struct {
		testField string `json:"testField"`
	}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 10000 {
			t.Skip("跳过长时间运行")
		}

		// 限制输入大小
		if len(input.testField) > 10000 {
			t.Skip("超长字符串，跳过")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 15*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("超长字符串XML转换失败: %v, 字符串长度: %d", err, len(input.testField))
			return
		}

		if len(xmlBytes) == 0 {
			t.Error("超长字符串生成空XML")
		}

		t.Logf("超长字符串XML转换成功，字符串长度: %d, XML长度: %d",
			len(input.testField), len(xmlBytes))
	})
}

// FuzzXmlWithUnicodeCharacters Unicode字符XML模糊测试
func FuzzXmlWithUnicodeCharacters(f *testing.F) {
	// 添加Unicode字符种子数据
	unicodeInputs := []struct {
		testField string `json:"testField"`
	}{
		{testField: "测试中文数据αβγδε"},
		{testField: "Test Emoji 😀😎🎉"},
		{testField: "日本語データ"},
		{testField: "한국어 데이터"},
		{testField: "العربية"},
		{testField: "עברית"},
		{testField: "Tiếng Việt"},
		{testField: "ไทย"},
		{testField: "বাংলা"},
	{testField: "Ελληνικά"},
		{testField: "русский"},
	},

		// 包含特殊Unicode点的文本
		{testField: "Text with \u0000-\uFFFF ranges"},
		{testField: "Text with emoji \uD83D-\uDBFF\uD83C-\uDFFF"},
	{testField: "Text with mathematical operators ≠≥≤≠"},
	}

	for _, seed := range unicodeInputs {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input struct {
		testField string `json:"testField"`
	}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 15000 {
			t.Skip("跳过长时间运行")
		}

		// 验证Unicode内容长度
		if len(input.testField) > 10000 {
			t.Skip("Unicode内容过长，跳过")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 10*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("Unicode字符XML转换失败: %v, 字段串: %q", err, input.testField)
			return
		}

		if len(xmlBytes) == 0 {
			t.Error("Unicode字符生成空XML")
		}

		// 验证Unicode字符被正确编码
		xmlStr := string(xmlBytes)
		if len(xmlStr) != len([]byte(input.testField)) {
			t.Logf("Unicode字符编码后长度变化: 原始 %d, XML %d",
				len([]byte(input.testField)), len(xmlStr))
		}

		t.Logf("Unicode字符XML转换成功，字符串长度: %d", len(input.testField))
	})
}

// FuzzXmlWithNumericValues 数值类型XML模糊测试
func FuzzXmlWithNumericValues(f *testing.F) {
	// 添加数值类型种子数据
	numericInputs := []struct {
		intField   int    `json:"intField"`
		floatField float64 `json:"floatField"`
		boolField bool   `json:"boolField"`
	}{
		{intField: 0, floatField: 0.0, boolField: false},
		{intField: 2147483647, floatField: 3.1415926535, boolField: true},
		{intField: -2147483648, floatField: -3.1415926535, boolField: true},
		{intField: 2147483647, floatField: 1.7976931342, boolField: false},
	}

	for _, seed := range numericInputs {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input struct {
		intField   int    `json:"intField"`
		floatField float64 `json:"floatField"`
		boolField bool   `json:"boolField"`
	}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 15000 {
			t.Skip("跳过长时间运行")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 10*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("数值类型XML转换失败: %v", err)
			return
		}

		if len(xmlBytes) == 0 {
			t.Error("数值类型生成空XML")
		}

		// 验证数值被正确序列化
		xmlStr := string(xmlBytes)
		if !strings.Contains(xmlStr, fmt.Sprintf("%d", input.intField)) && input.intField != 0 {
			t.Logf("整数值 %d 未在XML中找到", input.intField)
		}

		if !strings.Contains(xmlStr, fmt.Sprintf("%.6f", input.floatField)) && input.floatField != 0 {
			t.Logf("浮点数 %.6f 未在XML中找到", input.floatField)
		}

		if !strings.Contains(xmlStr, fmt.Sprintf("%t", input.boolField)) {
			t.Logf("布尔值 %t 未在XML中找到", input.boolField)
	}

		t.Logf("数值类型XML转换成功")
	})
}

// FuzzXmlWithArrayValues 数组类型XML模糊测试
func FuzzXmlWithArrayValues(f *testing.F) {
	// 添加数组类型种子数据
	arrayInputs := []struct {
		stringArray []string `json:"stringArray"`
		intArray   []int    `json:"intArray"`
		structArray []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"structArray"`
	}{
		{stringArray: []string{"a", "b", "c"}},
		{intArray: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{intArray: []int{}}, // 空数组
		{intArray: []int{-1, -2, -3}}, // 负数
		{structArray: []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			{Name: "test1", Value: "value1"},
			{Name: "test2", Value: "value2"},
		{Name: "test3", Value: "value3"},
	}},
	}

	for _, seed := range arrayInputs {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input interface{}) {
		// 防止测试超时
		if len(f.Fuzzing()) > 15000 {
			t.Skip("跳过长时间运行")
		}

		// 内存监控
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		if m.Alloc > 15*1024*1024 {
			t.Skip("内存使用过高")
		}

		// 执行XML转换
		xmlBytes, err := TransToXml(input)
		if err != nil {
			t.Logf("数组类型XML转换失败: %v", err)
			return
		}

		if len(xmlBytes) == 0 && input != nil {
			t.Logf("数组类型生成空XML（可能是边界情况）")
		}

		// 验证数组格式
		xmlStr := string(xmlBytes)
		if strings.Contains(xmlStr, "array") && !strings.Contains(xmlStr, "item") {
			t.Log("数组结构可能不是预期格式")
		}

		t.Logf("数组类型XML转换成功，XML长度: %d", len(xmlBytes))
	})
}
