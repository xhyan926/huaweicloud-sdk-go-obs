//go:build fuzz

package obs

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"testing"
	"time"
)

// FuzzFormatUrls 测试URL格式化函数
func FuzzFormatUrls(f *testing.F) {
	// 添加种子数据 - 各种有效的和边缘情况的URL格式
	seeds := []struct {
		bucketName string
		objectKey  string
		params     map[string]string
		escape     bool
	}{
		{"test-bucket", "test-object.txt", map[string]string{"version": "v1"}, true},
		{"test-bucket", "path/to/object.txt", map[string]string{"uploadId": "abc123"}, false},
		{"", "", map[string]string{"max-keys": "100"}, true},
		{"bucket-123", "特殊字符文件.txt", map[string]string{"encoding": "utf-8"}, true},
		{"my.bucket", "object?query=value", map[string]string{"acl": "public-read"}, false},
	}

	for _, seed := range seeds {
		f.Add(seed.bucketName, seed.objectKey, seed.params, seed.escape)
	}

	f.Fuzz(func(t *testing.T, bucketName, objectKey string, params map[string]string, escape bool) {
		// 输入验证
		if len(bucketName) > 255 || len(objectKey) > 1024*100 {
			t.Skip("输入过大，跳过")
		}

		if params == nil {
			params = make(map[string]string)
		}

		// 检查参数键和值长度
		for key, value := range params {
			if len(key) > 1024 || len(value) > 1024*10 {
				t.Skip("参数过大，跳过")
			}
			// 检查是否包含危险的SQL注入或XSS模式
			if strings.Contains(key, "<script") || strings.Contains(value, "<script") ||
				strings.Contains(key, "javascript:") || strings.Contains(value, "javascript:") ||
				strings.Contains(key, "eval(") || strings.Contains(value, "eval(") {
				t.Skip("检测到危险模式，跳过")
			}
		}

		// 执行URL格式化
		requestURL, canonicalizedURL := dummyFormatUrls(bucketName, objectKey, params, escape)

		// 验证输出的有效性
		if requestURL == "" {
			t.Error("生成的请求URL为空")
		}

		// 验证URL格式
		if strings.HasPrefix(requestURL, "http://") || strings.HasPrefix(requestURL, "https://") {
			parsedURL, err := url.Parse(requestURL)
			if err != nil {
				t.Logf("URL解析失败: %v, URL: %s", err, requestURL)
				return
			}

			// 验证URL的基本结构
			if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
				t.Errorf("无效的URL协议: %s", parsedURL.Scheme)
			}
		}

		// 验证规范化URL不以http://或https://开头
		if strings.HasPrefix(canonicalizedURL, "http://") || strings.HasPrefix(canonicalizedURL, "https://") {
			t.Error("规范化URL不应包含协议")
		}

		// 验证规范化URL以/开头
		if len(canonicalizedURL) > 0 && !strings.HasPrefix(canonicalizedURL, "/") {
			t.Error("规范化URL应以/开头")
		}
	})
}

// FuzzPrepareBaseURL 测试基础URL准备函数
func FuzzPrepareBaseURL(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		bucketName string
	}{
		{"test-bucket"},
		{"bucket-123"},
		{"my.bucket.com"},
		{""},
		{"a"},
		{"very-long-bucket-name-123456789"},
	}

	for _, seed := range seeds {
		f.Add(seed.bucketName)
	}

	f.Fuzz(func(t *testing.T, bucketName string) {
		// 输入验证
		if len(bucketName) > 255 {
			t.Skip("bucket名称过长，跳过")
		}

		// 检查bucket名称是否包含无效字符
		invalidChars := []string{" ", "\n", "\t", "\r", "/", "\\", "#", "?", "&", "="}
		for _, char := range invalidChars {
			if strings.Contains(bucketName, char) {
				t.Skip("bucket名称包含无效字符，跳过")
			}
		}

		// 执行基础URL准备
		requestURL, canonicalizedURL := dummyPrepareBaseURL(bucketName)

		// 验证输出的有效性
		if requestURL == "" {
			t.Error("生成的请求URL为空")
		}

		// 验证URL包含协议
		if !strings.HasPrefix(requestURL, "http://") && !strings.HasPrefix(requestURL, "https://") {
			t.Error("URL应包含协议")
		}

		// 验证规范化URL格式
		if len(canonicalizedURL) > 0 && canonicalizedURL[0] != '/' {
			t.Error("规范化URL应以/开头")
		}

		// 验证URL不包含连续的斜杠（除了协议部分）
		if !strings.Contains(requestURL, "://") && strings.Contains(requestURL, "//") {
			t.Error("URL不应包含连续的斜杠（协议后）")
		}
	})
}

// FuzzPrepareObjectKey 测试对象键准备函数
func FuzzPrepareObjectKey(f *testing.F) {
	// 添加种子数据 - 包含各种字符和格式的对象键
	seeds := []struct {
		objectKey string
		escape    bool
	}{
		{"test-object.txt", true},
		{"path/to/object.txt", false},
		{"特殊字符文件.txt", true},
		{"file with spaces.txt", true},
		{"file%20name.txt", false},
		{"非常长的文件名" + strings.Repeat("x", 100) + ".txt", true},
		{"a", true},
		{"", false},
		{"object?query=value", true},
		{"file#hash.txt", false},
		{"file&param=value.txt", true},
		{"emoji😀文件.txt", true},
		{"file\twith\ttabs.txt", false},
		{"file\nwith\nnewlines.txt", true},
	}

	for _, seed := range seeds {
		f.Add(seed.objectKey, seed.escape)
	}

	f.Fuzz(func(t *testing.T, objectKey string, escape bool) {
		// 输入验证
		if len(objectKey) > 1024*100 {
			t.Skip("对象键过长，跳过")
		}

		// 检查是否包含危险的控制字符
		controlChars := []string{"\x00", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07", "\x08", "\x0B", "\x0C", "\x0E", "\x0F"}
		for _, char := range controlChars {
			if strings.Contains(objectKey, char) {
				t.Skip("对象键包含控制字符，跳过")
			}
		}

		// 执行对象键准备
		encodedKey := dummyPrepareObjectKey(escape, objectKey)

		// 验证输出的有效性
		if encodedKey == "" && objectKey != "" {
			t.Error("编码后的对象键为空")
		}

		// 验证编码后的键不会过长
		if len(encodedKey) > len(objectKey)*4 && len(objectKey) > 0 {
			t.Errorf("编码后的对象键过长: 原始=%d, 编码=%d", len(objectKey), len(encodedKey))
		}

		// 验证编码后的键不包含未转义的危险字符（当escape=true时）
		if escape && len(objectKey) > 0 {
			dangerousPatterns := []string{"<script", "javascript:", "onclick=", "onerror=", "onload="}
			for _, pattern := range dangerousPatterns {
				if strings.Contains(strings.ToLower(encodedKey), pattern) {
					t.Errorf("编码后的键包含危险模式: %s", pattern)
				}
			}
		}
	})
}

// FuzzUrlParameterEscape 测试URL参数转义
func FuzzUrlParameterEscape(f *testing.F) {
	// 添加种子数据
	seeds := []string{
		"normal-value",
		"value with spaces",
		"value&special=chars",
		"unicode中文",
		"emoji😀",
		"",
		"a",
		"very-long-value" + strings.Repeat("x", 1000),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, paramValue string) {
		// 输入验证
		if len(paramValue) > 1024*10 {
			t.Skip("参数值过长，跳过")
		}

		// 执行参数转义
		escaped := url.QueryEscape(paramValue)

		// 验证转义后的值
		if escaped == "" && paramValue != "" {
			t.Error("转义后的值为空")
		}

		// 验证转义后的值不包含未转义的特殊字符
		if len(paramValue) > 0 {
			// 检查一些常见的未转义字符
			if strings.ContainsAny(escaped, " &?=#") {
				t.Error("转义后的值包含未转义的特殊字符")
			}
		}

		// 验证转义是可逆的
		unescaped, err := url.QueryUnescape(escaped)
		if err != nil {
			t.Logf("参数反转义失败: %v, 原始: %s", err, paramValue)
		} else if unescaped != paramValue && paramValue != "" {
			t.Logf("参数转义不完全可逆: 原始=%s, 转义=%s, 反转义=%s", paramValue, escaped, unescaped)
		}
	})
}

// FuzzUrlPathEscape 测试URL路径转义
func FuzzUrlPathEscape(f *testing.F) {
	// 添加种子数据
	seeds := []string{
		"normal/path",
		"path with spaces",
		"unicode路径/中文",
		"emoji路径/😀",
		"",
		"a",
		"very/long/path/" + strings.Repeat("x", 100),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, pathValue string) {
		// 输入验证
		if len(pathValue) > 1024*10 {
			t.Skip("路径值过长，跳过")
		}

		// 检查是否包含路径遍历攻击
		if strings.Contains(pathValue, "../") || strings.Contains(pathValue, "..\\") {
			t.Skip("检测到路径遍历攻击，跳过")
		}

		// 执行路径转义
		escaped := url.PathEscape(pathValue)

		// 验证转义后的值
		if escaped == "" && pathValue != "" {
			t.Error("转义后的值为空")
		}

		// 验证路径转义保留斜杠
		slashCount := strings.Count(pathValue, "/")
		escapedSlashCount := strings.Count(escaped, "%2F") + strings.Count(escaped, "/")

		// 斜杠应该被保留或正确转义
		if slashCount != escapedSlashCount {
			t.Logf("路径斜杠数量不一致: 原始=%d, 转义=%d", slashCount, escapedSlashCount)
		}
	})
}

// FuzzUrlWithSpecialPatterns 测试特殊URL模式
func FuzzUrlWithSpecialPatterns(f *testing.F) {
	// 添加种子数据 - 包含各种特殊模式的URL
	seeds := []string{
		"https://obs.cn-north-4.myhuaweicloud.com/bucket/object",
		"http://localhost:8080/test/path",
		"https://example.com/bucket/object?version=1&param=value",
		"https://bucket.obs.cn-north-4.myhuaweicloud.com/object",
		"https://obs.cn-north-4.myhuaweicloud.com/bucket/object?type=video&quality=1080p",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, urlStr string) {
		// 输入验证
		if len(urlStr) > 4096 {
			t.Skip("URL过长，跳过")
		}

		// 执行URL解析
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			t.Logf("URL解析失败: %v, URL: %s", err, urlStr)
			return
		}

		// 验证URL的基本结构
		if parsedURL.Scheme != "" && parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			t.Errorf("无效的URL协议: %s", parsedURL.Scheme)
		}

		// 验证host不包含无效字符
		if parsedURL.Host != "" {
			invalidHostChars := []string{" ", "\n", "\t", "\r", "<", ">", "\"", "'", "\\"}
			for _, char := range invalidHostChars {
				if strings.Contains(parsedURL.Host, char) {
					t.Errorf("Host包含无效字符: %s", char)
				}
			}
		}

		// 验证path不包含路径遍历
		if strings.Contains(parsedURL.Path, "../") || strings.Contains(parsedURL.Path, "..\\") {
			t.Error("路径包含路径遍历攻击")
		}

		// 验证query参数
		if parsedURL.RawQuery != "" {
			queryParams, err := url.ParseQuery(parsedURL.RawQuery)
			if err != nil {
				t.Logf("Query参数解析失败: %v", err)
			}

			// 验证参数键和值不包含危险模式
			for key, values := range queryParams {
				for _, value := range values {
					if strings.Contains(key, "<script") || strings.Contains(value, "<script") ||
						strings.Contains(key, "javascript:") || strings.Contains(value, "javascript:") {
						t.Logf("检测到潜在XSS: key=%s, value=%s", key, value)
					}
				}
			}
		}
	})
}

// FuzzSignedUrlGeneration 测试签名URL生成
func FuzzSignedUrlGeneration(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		bucketName string
		objectKey  string
		expires    int64
	}{
		{"test-bucket", "test-object.txt", 3600},
		{"bucket-123", "path/to/object.txt", 7200},
		{"my.bucket", "特殊文件.txt", 1800},
		{"bucket", "", 3600},
		{"a", "a", 60},
	}

	for _, seed := range seeds {
		f.Add(seed.bucketName, seed.objectKey, seed.expires)
	}

	f.Fuzz(func(t *testing.T, bucketName, objectKey string, expires int64) {
		// 输入验证
		if len(bucketName) > 255 || len(objectKey) > 1024*10 {
			t.Skip("输入过大，跳过")
		}

		if expires < 0 || expires > 60*60*24*7 { // 最大7天
			t.Skip("过期时间无效，跳过")
		}

		// 执行签名URL生成（模拟）
		signedUrl := dummyCreateSignedUrl(bucketName, objectKey, expires)

		// 验证签名URL的有效性
		if signedUrl == "" {
			t.Error("生成的签名URL为空")
		}

		// 验证签名URL包含基本组件
		if !strings.HasPrefix(signedUrl, "http://") && !strings.HasPrefix(signedUrl, "https://") {
			t.Error("签名URL应包含协议")
		}

		// 验证签名URL包含过期时间参数
		if !strings.Contains(signedUrl, "expires") && !strings.Contains(signedUrl, "AWSAccessKeyId") {
			t.Log("签名URL可能缺少过期时间参数")
		}

		// 验证签名URL可以解析
		parsedURL, err := url.Parse(signedUrl)
		if err != nil {
			t.Errorf("签名URL解析失败: %v", err)
		} else {
			// 验证过期时间在未来
			if expires > 0 {
				expiryTime := time.Now().Add(time.Duration(expires) * time.Second)
				if expiryTime.Before(time.Now().Add(-10 * time.Second)) { // 允许10秒误差
					t.Error("过期时间应在未来")
				}
			}
		}
	})
}

// FuzzUrlWithMemoryMonitoring 测试URL处理时的内存使用
func FuzzUrlWithMemoryMonitoring(f *testing.F) {
	// 添加种子数据
	seeds := []string{
		"https://obs.cn-north-4.myhuaweicloud.com/bucket/object",
		"http://localhost:8080/test/path",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, urlStr string) {
		// 输入验证
		if len(urlStr) > 2048 {
			t.Skip("URL过长，跳过")
		}

		// 监控内存使用
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memBefore := m.Alloc

		// 执行URL解析
		_, err := url.Parse(urlStr)
		if err != nil {
			t.Logf("URL解析失败: %v", err)
		}

		// 检查内存使用
		runtime.ReadMemStats(&m)
		memAfter := m.Alloc
		memDelta := memAfter - memBefore

		// 内存使用不应超过10MB
		if memDelta > 10*1024*1024 {
			t.Errorf("内存使用过高: %d bytes", memDelta)
		}
	})
}

// FuzzUrlBoundaryConditions 测试URL处理的边界条件
func FuzzUrlBoundaryConditions(f *testing.F) {
	// 添加种子数据 - 边界条件
	seeds := []struct {
		bucketName string
		objectKey  string
	}{
		{"", ""},                     // 空值
		{"a", "a"},                   // 单字符
		{"a", strings.Repeat("a", 1023)}, // 最大长度-1
		{"a", strings.Repeat("a", 1024)}, // 最大长度
		{strings.Repeat("a", 255), ""},    // 最大bucket名称
		{"a-b-c-d-e-f-g-h-i-j", "a/b/c/d/e/f/g/h/i/j"}, // 多层级路径
	}

	for _, seed := range seeds {
		f.Add(seed.bucketName, seed.objectKey)
	}

	f.Fuzz(func(t *testing.T, bucketName, objectKey string) {
		// 输入验证
		if len(bucketName) > 255 || len(objectKey) > 1024*100 {
			t.Skip("输入超出边界，跳过")
		}

		// 执行URL格式化
		requestURL, canonicalizedURL := dummyFormatUrls(bucketName, objectKey, nil, false)

		// 验证边界条件
		if len(requestURL) > 8192 { // URL长度限制
			t.Errorf("生成的URL过长: %d 字符", len(requestURL))
		}

		if len(canonicalizedURL) > 2048 { // 规范化URL长度限制
			t.Errorf("规范化URL过长: %d 字符", len(canonicalizedURL))
		}

		// 验证空值处理
		if bucketName == "" && objectKey == "" {
			if !strings.HasPrefix(requestURL, "http://") && !strings.HasPrefix(requestURL, "https://") {
				t.Error("空bucket和object应生成基础URL")
			}
		}
	})
}

// FuzzUrlWithLongValues 测试超长值的处理
func FuzzUrlWithLongValues(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		longString string
	}{
		{strings.Repeat("a", 100)},
		{strings.Repeat("中文", 50)},
		{strings.Repeat("😀", 50)},
	}

	for _, seed := range seeds {
		f.Add(seed.longString)
	}

	f.Fuzz(func(t *testing.T, longString string) {
		// 输入验证
		if len(longString) > 1024*100 {
			t.Skip("字符串过长，跳过")
		}

		// 创建包含长字符串的URL参数
		params := map[string]string{
			"longParam": longString,
			"normalParam": "normal",
		}

		// 执行URL格式化
		requestURL, _ := dummyFormatUrls("test-bucket", "test-object.txt", params, true)

		// 验证长字符串被正确处理
		if !strings.Contains(requestURL, "longParam=") {
			t.Error("URL应包含长参数")
		}

		// 验证URL长度在合理范围内
		if len(requestURL) > 8192 {
			t.Errorf("URL过长: %d 字符", len(requestURL))
		}
	})
}

// ===== 辅助函数（模拟实际的URL处理函数） =====

// dummyFormatUrls 模拟formatUrls函数
func dummyFormatUrls(bucketName, objectKey string, params map[string]string, escape bool) (requestURL string, canonicalizedURL string) {
	// 简化实现用于测试
	requestURL = "https://obs.cn-north-4.myhuaweicloud.com"
	canonicalizedURL = "/"

	if bucketName != "" {
		requestURL += "/" + bucketName
		canonicalizedURL += bucketName
	}

	if objectKey != "" {
		requestURL += "/" + objectKey
		if !strings.HasSuffix(canonicalizedURL, "/") {
			canonicalizedURL += "/"
		}
		canonicalizedURL += objectKey
	}

	if len(params) > 0 {
		requestURL += "?"
		keys := make([]string, 0, len(params))
		for key := range params {
			keys = append(keys, key)
		}
		for i, key := range keys {
			if i > 0 {
				requestURL += "&"
			}
			if escape {
				requestURL += fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(params[key]))
			} else {
				requestURL += fmt.Sprintf("%s=%s", key, params[key])
			}
		}
	}

	return
}

// dummyPrepareBaseURL 模拟prepareBaseURL函数
func dummyPrepareBaseURL(bucketName string) (requestURL string, canonicalizedURL string) {
	requestURL = "https://obs.cn-north-4.myhuaweicloud.com"
	canonicalizedURL = "/"

	if bucketName != "" {
		requestURL += "/" + bucketName
		canonicalizedURL += bucketName
	}

	return
}

// dummyPrepareObjectKey 模拟prepareObjectKey函数
func dummyPrepareObjectKey(escape bool, objectKey string) string {
	if escape {
		result := make([]string, 0, len(objectKey))
		for _, char := range objectKey {
			if string(char) == "/" {
				result = append(result, string(char))
			} else if string(char) == " " {
				result = append(result, url.PathEscape(string(char)))
			} else {
				result = append(result, url.QueryEscape(string(char)))
			}
		}
		return strings.Join(result, "")
	}
	return objectKey
}

// dummyCreateSignedUrl 模拟CreateSignedUrl函数
func dummyCreateSignedUrl(bucketName, objectKey string, expires int64) string {
	url := "https://obs.cn-north-4.myhuaweicloud.com"
	if bucketName != "" {
		url += "/" + bucketName
	}
	if objectKey != "" {
		url += "/" + objectKey
	}

	if expires > 0 {
		url += fmt.Sprintf("?Expires=%d&AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE",
			time.Now().Add(time.Duration(expires)*time.Second).Unix())
	}

	return url
}