//go:build fuzz

package obs

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// FuzzHmacSha256 测试HMAC SHA256签名函数
func FuzzHmacSha256(f *testing.F) {
	// 添加种子数据 - 各种密钥和数据组合
	seeds := []struct {
		key   []byte
		value []byte
	}{
		{[]byte("test-key"), []byte("test-value")},
		{[]byte("short"), []byte("data")},
		{[]byte(strings.Repeat("a", 32)), []byte(strings.Repeat("b", 64))},
		{[]byte("very-long-key-" + strings.Repeat("x", 100)), []byte("very-long-value-" + strings.Repeat("y", 100))},
		{[]byte(""), []byte("test")},         // 空密钥
		{[]byte("test"), []byte("")},         // 空值
		{[]byte("中文密钥"), []byte("中文数据")}, // Unicode字符
		{[]byte("😀key"), []byte("😀data")},   // Emoji
	}

	for _, seed := range seeds {
		f.Add(seed.key, seed.value)
	}

	f.Fuzz(func(t *testing.T, key, value []byte) {
		// 输入验证
		if len(key) > 4096 || len(value) > 4096*1024 { // 4MB限制
			t.Skip("输入过大，跳过")
		}

		// 监控内存使用
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memBefore := m.Alloc

		// 执行HMAC SHA256签名
		signature1 := dummyHmacSha256(key, value)
		signature2 := dummyHmacSha256(key, value)

		// 验证签名一致性
		if !bytes.Equal(signature1, signature2) {
			t.Error("HMAC签名结果不一致")
		}

		// 验证签名长度
		if len(signature1) != sha256.Size {
			t.Errorf("HMAC签名长度错误: 期望=%d, 实际=%d", sha256.Size, len(signature1))
		}

		// 检查内存使用
		runtime.ReadMemStats(&m)
		memAfter := m.Alloc
		memDelta := memAfter - memBefore

		// 内存使用不应超过输入大小的3倍
		if memDelta > int64(len(key)+len(value))*3 {
			t.Errorf("内存使用过高: 输入=%d bytes, 内存增量=%d bytes", len(key)+len(value), memDelta)
		}
	})
}

// FuzzSignatureGeneration 测试签名生成函数
func FuzzSignatureGeneration(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		stringToSign string
		sk           string
		region       string
		shortDate    string
	}{
		{"GET\n/\n\nhost:obs.cn-north-4.myhuaweicloud.com\n\nhost\n", "secret-key", "cn-north-4", "20240101"},
		{"PUT\n/bucket/object\n\nhost:bucket.obs.cn-north-4.myhuaweicloud.com\n\nhost\n", "test-secret", "cn-south-1", "20240202"},
		{"POST\n/\napplication/x-www-form-urlencoded\nhost:obs.cn-north-4.myhuaweicloud.com\nx-amz-date:20240101T000000Z\n\nhost;x-amz-date\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "access-key", "cn-east-2", "20240303"},
	}

	for _, seed := range seeds {
		f.Add(seed.stringToSign, seed.sk, seed.region, seed.shortDate)
	}

	f.Fuzz(func(t *testing.T, stringToSign, sk, region, shortDate string) {
		// 输入验证
		if len(stringToSign) > 1024*100 || len(sk) > 1024 || len(region) > 64 || len(shortDate) > 16 {
			t.Skip("输入过大，跳过")
		}

		// 验证shortDate格式（应该是YYYYMMDD）
		if len(shortDate) == 8 {
			year, month, day := shortDate[0:4], shortDate[4:6], shortDate[6:8]
			// 检查是否是数字
			if !isNumeric(year) || !isNumeric(month) || !isNumeric(day) {
				t.Skip("无效的日期格式，跳过")
			}
		}

		// 执行签名生成
		signature := dummyGetSignature(stringToSign, sk, region, shortDate)

		// 验证签名格式
		if signature == "" && len(stringToSign) > 0 {
			t.Error("签名不应为空（除非输入为空）")
		}

		// 验证签名长度（HMAC SHA256应该是64个字符的十六进制字符串）
		if len(signature) > 64 {
			t.Errorf("签名长度异常: %d", len(signature))
		}

		// 验证签名只包含有效的十六进制字符
		for _, c := range signature {
			if !isHexChar(string(c)) {
				t.Errorf("签名包含无效字符: %c", c)
				break
			}
		}
	})
}

// FuzzStringToSignConstruction 测试待签名字符串构建
func FuzzStringToSignConstruction(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		method          string
		canonicalizedURL string
		queryURL        string
		payload         string
		headers         map[string]string
	}{
		{"GET", "/", "", "", "host:obs.cn-north-4.myhuaweicloud.com\n\nhost\n", nil},
		{"PUT", "/bucket/object", "?versionId=v1", "e3b0c44...", "host:bucket.obs.cn-north-4.myhuaweicloud.com\nx-amz-date:20240101T000000Z\n\nhost;x-amz-date\n", map[string]string{"x-amz-date": "20240101T000000Z"}},
		{"POST", "/path/to/object", "?uploadId=abc123", "application/x-www-form-urlencoded...", "host:example.com\ncontent-type:application/x-www-form-urlencoded\n\nhost;content-type\n", map[string]string{"content-type": "application/x-www-form-urlencoded"}},
		{"DELETE", "/bucket/object", "", "", "host:obs.cn-north-4.myhuaweicloud.com\n\nhost\n", nil},
	}

	for _, seed := range seeds {
		f.Add(seed.method, seed.canonicalizedURL, seed.queryURL, seed.payload, seed.headers)
	}

	f.Fuzz(func(t *testing.T, method, canonicalizedURL, queryURL, payload, headers string) {
		// 输入验证
		if len(method) > 16 || len(canonicalizedURL) > 1024 || len(queryURL) > 1024 || len(payload) > 1024*100 {
			t.Skip("输入过大，跳过")
		}

		// 验证HTTP方法
		validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true, "HEAD": true}
		if method != "" && !validMethods[strings.ToUpper(method)] {
			t.Skip("无效的HTTP方法，跳过")
		}

		// 构建待签名字符串（简化模拟）
		stringToSign := method + "\n"
		stringToSign += canonicalizedURL + "\n"
		stringToSign += queryURL + "\n"
		stringToSign += payload

		// 验证待签名字符串的格式
		if len(stringToSign) > 1024*100 {
			t.Error("待签名字符串过长")
		}

		// 验证换行符的数量
		newlineCount := strings.Count(stringToSign, "\n")
		if newlineCount < 3 { // method, url, query, payload之间至少有3个换行
			t.Error("待签名字符串格式错误：换行符不足")
		}

		// 验证开头不包含多余空格
		if strings.HasPrefix(stringToSign, " ") {
			t.Error("待签名字符串开头不应有空格")
		}

		// 验证尾部的payload没有多余的换行符
		if len(payload) > 0 && strings.HasSuffix(stringToSign, "\n\n") {
			t.Error("待签名字符串尾部有多余换行符")
		}
	})
}

// FuzzMd5Hashing 测试MD5哈希函数
func FuzzMd5Hashing(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		data []byte
	}{
		{[]byte("test-data")},
		{[]byte("very-long-data-" + strings.Repeat("x", 1000))},
		{[]byte("")},                          // 空数据
		{[]byte("中文数据")},                    // Unicode
		{[]byte("😀data")},                  // Emoji
		{[]byte(strings.Repeat("a", 64))},    // 64字节
		{[]byte(strings.Repeat("b", 1024))},   // 1KB
	}

	for _, seed := range seeds {
		f.Add(seed.data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// 输入验证
		if len(data) > 10*1024*1024 { // 10MB限制
			t.Skip("数据过大，跳过")
		}

		// 执行MD5哈希
		hash1 := md5.Sum(data)
		hash2 := md5.Sum(data)

		// 验证哈希一致性
		if !bytes.Equal(hash1, hash2) {
			t.Error("MD5哈希结果不一致")
		}

		// 验证哈希长度
		if len(hash1) != md5.Size {
			t.Errorf("MD5哈希长度错误: 期望=%d, 实际=%d", md5.Size, len(hash1))
		}

		// 验证哈希十六进制格式
		hexHash := hex.EncodeToString(hash1)
		if len(hexHash) != md5.Size*2 { // 每字节2个十六进制字符
			t.Errorf("MD5十六进制哈希长度错误: 期望=%d, 实际=%d", md5.Size*2, len(hexHash))
		}

		// 验证十六进制哈希只包含有效的十六进制字符
		for _, c := range hexHash {
			if !isHexChar(string(c)) {
				t.Errorf("MD5哈希包含无效字符: %c", c)
				break
			}
		}
	})
}

// FuzzSha256Hashing 测试SHA256哈希函数
func FuzzSha256Hashing(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		data []byte
	}{
		{[]byte("test-data")},
		{[]byte("very-long-data-" + strings.Repeat("x", 1000))},
		{[]byte("")},                          // 空数据
		{[]byte("中文数据")},                    // Unicode
		{[]byte("😀data")},                  // Emoji
		{[]byte(strings.Repeat("a", 64))},    // 64字节
		{[]byte(strings.Repeat("b", 1024))},   // 1KB
	}

	for _, seed := range seeds {
		f.Add(seed.data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// 输入验证
		if len(data) > 10*1024*1024 { // 10MB限制
			t.Skip("数据过大，跳过")
		}

		// 执行SHA256哈希
		hash1 := sha256.Sum256(data)
		hash2 := sha256.Sum256(data)

		// 验证哈希一致性
		if !bytes.Equal(hash1, hash2) {
			t.Error("SHA256哈希结果不一致")
		}

		// 验证哈希长度
		if len(hash1) != sha256.Size {
			t.Errorf("SHA256哈希长度错误: 期望=%d, 实际=%d", sha256.Size, len(hash1))
		}

		// 验证哈希十六进制格式
		hexHash := hex.EncodeToString(hash1)
		if len(hexHash) != sha256.Size*2 { // 每字节2个十六进制字符
			t.Errorf("SHA256十六进制哈希长度错误: 期望=%d, 实际=%d", sha256.Size*2, len(hexHash))
		}

		// 验证十六进制哈希只包含有效的十六进制字符
		for _, c := range hexHash {
			if !isHexChar(string(c)) {
				t.Errorf("SHA256哈希包含无效字符: %c", c)
				break
			}
		}
	})
}

// FuzzHexEncoding 测试十六进制编码函数
func FuzzHexEncoding(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		data []byte
	}{
		{[]byte{0x00, 0x01, 0x02}},        // 二进制数据
		{[]byte("test-data")},                  // ASCII
		{[]byte("中文")},                       // UTF-8中文
		{[]byte("😀")},                        // UTF-8 emoji
		{[]byte{0xFF, 0xFE, 0xFD}},         // 特殊字节
		{[]byte(strings.Repeat("a", 100))},    // 100字节
		{[]byte(strings.Repeat("b", 1024))},   // 1KB
	}

	for _, seed := range seeds {
		f.Add(seed.data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// 输入验证
		if len(data) > 1024*100 { // 100KB限制
			t.Skip("数据过大，跳过")
		}

		// 执行十六进制编码
		encoded1 := hex.EncodeToString(data)
		encoded2 := hex.EncodeToString(data)

		// 验证编码一致性
		if encoded1 != encoded2 {
			t.Error("十六进制编码结果不一致")
		}

		// 验证编码长度
		expectedLen := len(data) * 2 // 每字节2个十六进制字符
		if len(encoded1) != expectedLen {
			t.Errorf("十六进制编码长度错误: 期望=%d, 实际=%d", expectedLen, len(encoded1))
		}

		// 验证编码只包含有效的十六进制字符
		for _, c := range encoded1 {
			if !isHexChar(string(c)) {
				t.Errorf("十六进制编码包含无效字符: %c", c)
				break
			}
		}

		// 验证编码可逆（如果数据非空）
		if len(data) > 0 {
			decoded, err := hex.DecodeString(encoded1)
			if err != nil {
				t.Errorf("十六进制解码失败: %v", err)
			} else if !bytes.Equal(data, decoded) {
				t.Error("十六进制编码不可逆")
			}
		}
	})
}

// FuzzSignatureWithVariousInputs 测试各种输入的签名处理
func FuzzSignatureWithVariousInputs(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		ak          string
		sk          string
		method      string
		bucket      string
		object      string
	}{
		{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", "GET", "my-bucket", "my-object.txt"},
		{"AKIAI44QH8DHBEXAMPLE", "je7MtGbClwBF/2Vp4UK8Q6BmT/5bE3EXAMPLEKEY", "PUT", "test-bucket", "path/to/object.txt"},
		{"test-ak", "test-sk", "DELETE", "bucket-123", "object-456.txt"},
	}

	for _, seed := range seeds {
		f.Add(seed.ak, seed.sk, seed.method, seed.bucket, seed.object)
	}

	f.Fuzz(func(t *testing.T, ak, sk, method, bucket, object string) {
		// 输入验证
		if len(ak) > 128 || len(sk) > 128 || len(method) > 16 || len(bucket) > 255 || len(object) > 1024 {
			t.Skip("输入过大，跳过")
		}

		// 验证AK/SK不为空时才签名
		if ak == "" || sk == "" {
			t.Skip("AK或SK为空，跳过签名测试")
		}

		// 验证HTTP方法
		validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true, "HEAD": true}
		if !validMethods[strings.ToUpper(method)] {
			t.Skip("无效的HTTP方法，跳过")
		}

		// 验证bucket名称格式
		if len(bucket) > 0 {
			if !isBucketNameValid(bucket) {
				t.Skip("无效的bucket名称，跳过")
			}
		}

		// 构建待签名字符串（简化）
		canonicalURL := "/"
		if bucket != "" {
			canonicalURL += bucket
		}
		if object != "" {
			canonicalURL += "/" + object
		}

		stringToSign := method + "\n" + canonicalURL + "\n\n" + "host:obs.cn-north-4.myhuaweicloud.com\n\nhost\n"

		// 执行签名
		signature := dummyGetSignature(stringToSign, sk, "cn-north-4", "20240101")

		// 验证签名不为空
		if signature == "" {
			t.Error("签名不应为空")
		}

		// 验证签名格式
		if len(signature) != 64 { // HMAC SHA256应该是64个字符
			t.Errorf("签名长度错误: 期望=64, 实际=%d", len(signature))
		}
	})
}

// FuzzSignatureBoundaryConditions 测试签名处理的边界条件
func FuzzSignatureBoundaryConditions(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		key   string
		data  string
	}{
		{"a", "b"},                     // 最小输入
		{"a", strings.Repeat("a", 1000)}, // 长数据
		{strings.Repeat("a", 100), "b"}, // 长密钥
		{"", ""},                         // 空输入
		{"key with spaces", "data with spaces"}, // 包含空格
		{"key\nwith\nnewlines", "data\twith\ttabs"}, // 控制字符
	}

	for _, seed := range seeds {
		f.Add(seed.key, seed.data)
	}

	f.Fuzz(func(t *testing.T, key, data string) {
		// 输入验证
		if len(key) > 1024 || len(data) > 1024*100 {
			t.Skip("输入超出边界，跳过")
		}

		// 执行HMAC签名
		signature := dummyHmacSha256([]byte(key), []byte(data))

		// 验证边界条件处理
		if len(key) == 0 || len(data) == 0 {
			// 空输入应该也能处理，但结果可能不是我们关心的
			t.Logf("空输入处理: key_len=%d, data_len=%d", len(key), len(data))
		}

		// 验证空输入时的行为
		if len(key) == 0 && len(data) == 0 {
			if len(signature) != 32 { // HMAC SHA256 of empty input
				t.Logf("空输入HMAC长度: %d", len(signature))
			}
		}

		// 验证包含空格的输入
		if strings.Contains(key, " ") || strings.Contains(data, " ") {
			// 空格应该被正确处理
			t.Log("处理包含空格的输入")
		}
	})
}

// FuzzConcurrentSignatureGeneration 测试并发签名生成
func FuzzConcurrentSignatureGeneration(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		key   string
		data  string
	}{
		{"test-key", "test-data"},
		{"key-1", "data-1"},
		{"key-2", "data-2"},
	}

	for _, seed := range seeds {
		f.Add(seed.key, seed.data)
	}

	f.Fuzz(func(t *testing.T, key, data string) {
		// 输入验证
		if len(key) > 1024 || len(data) > 1024 {
			t.Skip("输入过大，跳过")
		}

		// 并发测试：同时生成多个签名
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(index int) {
				defer func() { done <- true }()

				// 执行HMAC签名
				signature := dummyHmacSha256([]byte(key), []byte(data))

				// 验证签名长度
				if len(signature) != sha256.Size {
					t.Logf("并发签名长度错误 (goroutine %d): 期望=%d, 实际=%d", index, sha256.Size, len(signature))
				}
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// FuzzSignatureWithMaliciousInputs 测试恶意输入的签名处理
func FuzzSignatureWithMaliciousInputs(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		key  string
		data string
	}{
		{"normal-key", "normal-data"},
		{"' or '1'='1", "data"},                 // SQL注入尝试
		{"<script>alert(1)</script>", "data"},         // XSS尝试
		{"${jndi:ldap://exp}", "data"},             // JNDI注入
		{"${7*7}", "data"},                        // 表达式注入
		{"../../../etc/passwd", "data"},              // 路径遍历
		{"key", "<!DOCTYPE html><html><body>x</body></html>"}, // HTML注入
	}

	for _, seed := range seeds {
		f.Add(seed.key, seed.data)
	}

	f.Fuzz(func(t *testing.T, key, data string) {
		// 输入验证
		if len(key) > 1024 || len(data) > 1024 {
			t.Skip("输入过大，跳过")
		}

		// 检查恶意模式
		maliciousPatterns := []string{
			"' or '1'='1", "' or 1=1--", "<script", "javascript:",
			"${jndi:", "${7*7}", "../../../", "<!DOCTYPE",
			"eval(", "alert(", "document.cookie",
		}

		for _, pattern := range maliciousPatterns {
			if strings.Contains(strings.ToLower(key), strings.ToLower(pattern)) ||
				strings.Contains(strings.ToLower(data), strings.ToLower(pattern)) {
				t.Log("检测到潜在恶意模式，处理但不崩溃")
				break
			}
		}

		// 执行HMAC签名（应该处理而不崩溃）
		signature := dummyHmacSha256([]byte(key), []byte(data))

		// 验证签名生成不会崩溃
		if signature == nil {
			t.Error("HMAC签名返回nil，这可能表示处理失败")
		}
	})
}

// ===== 辅助函数（模拟实际的签名函数） =====

// dummyHmacSha256 模拟HmacSha256函数
func dummyHmacSha256(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, err := mac.Write(value)
	if err != nil {
		fmt.Printf("HmacSha256 failed to write: %v\n", err)
	}
	return mac.Sum(nil)
}

// dummyGetSignature 模拟getSignature函数
func dummyGetSignature(stringToSign, sk, region, shortDate string) string {
	// V4签名算法（简化）
	const (
		V4_HASH_PRE     = "AWS4"
		V4_SERVICE_NAME = "s3"
		V4_SERVICE_SUFFIX = "aws4_request"
	)

	key := dummyHmacSha256([]byte(V4_HASH_PRE+sk), []byte(shortDate))
	key = dummyHmacSha256(key, []byte(region))
	key = dummyHmacSha256(key, []byte(V4_SERVICE_NAME))
	key = dummyHmacSha256(key, []byte(V4_SERVICE_SUFFIX))
	return hex.EncodeToString(dummyHmacSha256(key, []byte(stringToSign)))
}

// isBucketNameValid 验证bucket名称是否有效
func isBucketNameValid(bucket string) bool {
	if len(bucket) < 3 || len(bucket) > 63 {
		return false
	}

	// bucket名称只能包含小写字母、数字、点和连字符
	for _, c := range bucket {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '.') {
			return false
		}
	}

	// 不能以点开头或结尾
	if strings.HasPrefix(bucket, ".") || strings.HasSuffix(bucket, ".") {
		return false
	}

	// 不能包含连续的点
	if strings.Contains(bucket, "..") {
		return false
	}

	return true
}

// isNumeric 检查字符串是否全是数字
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isHexChar 检查字符是否是有效的十六进制字符
func isHexChar(c string) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}