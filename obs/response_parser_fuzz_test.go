//go:build fuzz

package obs

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"testing"
	"time"
)

// FuzzParseResponseToBaseModel 测试响应解析函数
func FuzzParseResponseToBaseModel(f *testing.F) {
	// 添加种子数据 - 各种响应类型
	seeds := []struct {
		statusCode int
		body       []byte
		headers    map[string]string
		xmlResult  bool
	}{
		{200, []byte(`<ListBucketResult xmlns="http://obs.cn-north-4.myhuaweicloud.com/doc/2015-06-30/"><Name>test-bucket</Name></ListBucketResult>`), map[string]string{"Content-Type": "application/xml"}, true},
		{200, []byte(`{"bucket":"test-bucket","prefix":"test/"}`), map[string]string{"Content-Type": "application/json"}, false},
		{404, []byte(`<Error><Code>NoSuchBucket</Code><Message>The specified bucket does not exist</Message></Error>`), map[string]string{"Content-Type": "application/xml"}, true},
		{200, []byte(`<CreateBucketOutput><Location>test-bucket.obs.cn-north-4.myhuaweicloud.com</Location></CreateBucketOutput>`), map[string]string{"Content-Type": "application/xml"}, true},
		{500, []byte(`<Error><Code>InternalError</Code><Message>We encountered an internal error. Please try again.</Message></Error>`), map[string]string{"Content-Type": "application/xml"}, true},
	}

	for _, seed := range seeds {
		f.Add(seed.statusCode, seed.body, seed.headers, seed.xmlResult)
	}

	f.Fuzz(func(t *testing.T, statusCode int, body []byte, headers map[string]string, xmlResult bool) {
		// 输入验证
		if statusCode < 100 || statusCode > 599 {
			t.Skip("无效的状态码，跳过")
		}

		if len(body) > 10*1024*1024 { // 10MB限制
			t.Skip("响应体过大，跳过")
		}

		if headers == nil {
			headers = make(map[string]string)
		}

		// 检查header长度
		for key, value := range headers {
			if len(key) > 1024 || len(value) > 1024*10 {
				t.Skip("header过大，跳过")
			}
		}

		// 检查危险内容
		bodyStr := string(body)
		dangerousPatterns := []string{
			"<!ENTITY", "<!DOCTYPE", "<xsl:", "<script:",
			"<iframe:", "<svg:", "<style:", "<link:", "<meta:",
			"javascript:", "eval(", "onclick=", "onerror=", "onload=",
		}
		for _, pattern := range dangerousPatterns {
			if strings.Contains(strings.ToLower(bodyStr), strings.ToLower(pattern)) {
				t.Skip("检测到危险模式，跳过")
			}
		}

		// 模拟HTTP响应
		resp := &http.Response{
			StatusCode: statusCode,
			Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}

		// 设置headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		// 创建BaseModel实例（使用简单的ObsError）
		baseModel := &ObsError{}

		// 执行响应解析（模拟）
		err := dummyParseResponseToBaseModel(resp, baseModel, xmlResult, false)

		// 验证解析结果
		if err != nil && statusCode >= 200 && statusCode < 300 {
			t.Logf("成功状态码解析失败: %v", err)
		}

		// 验证状态码被正确设置
		if baseModel.StatusCode() != statusCode {
			t.Errorf("状态码未正确设置: 期望=%d, 实际=%d", statusCode, baseModel.StatusCode())
		}
	})
}

// FuzzParseResponseToObsError 测试错误响应解析
func FuzzParseResponseToObsError(f *testing.F) {
	// 添加种子数据 - 各种错误响应
	seeds := []struct {
		statusCode int
		body       []byte
		headers    map[string]string
	}{
		{404, []byte(`<Error><Code>NoSuchBucket</Code><Message>The specified bucket does not exist</Message><Resource>/test-bucket</Resource><RequestId>TX000001</RequestId></Error>`), map[string]string{"Content-Type": "application/xml"}},
		{403, []byte(`<Error><Code>AccessDenied</Code><Message>Access Denied</Message><RequestId>TX000002</RequestId></Error>`), map[string]string{"Content-Type": "application/xml"}},
		{500, []byte(`<Error><Code>InternalError</Code><Message>Internal Server Error</Message><HostId>obs.cn-north-4.myhuaweicloud.com</HostId></Error>`), map[string]string{"Content-Type": "application/xml"}},
		{400, []byte(`<Error><Code>InvalidBucketName</Code><Message>Bucket name is invalid</Message></Error>`), map[string]string{"Content-Type": "application/xml"}},
		{403, []byte(`{"code":"AccessDenied","message":"Access Denied"}`), map[string]string{"Content-Type": "application/json"}},
	}

	for _, seed := range seeds {
		f.Add(seed.statusCode, seed.body, seed.headers)
	}

	f.Fuzz(func(t *testing.T, statusCode int, body []byte, headers map[string]string) {
		// 输入验证
		if statusCode < 400 || statusCode > 599 {
			t.Skip("无效的错误状态码，跳过")
		}

		if len(body) > 1024*100 { // 100KB限制
			t.Skip("错误响应体过大，跳过")
		}

		if headers == nil {
			headers = make(map[string]string)
		}

		// 检查header长度
		for key, value := range headers {
			if len(key) > 1024 || len(value) > 1024 {
				t.Skip("header过大，跳过")
			}
		}

		// 模拟HTTP响应
		resp := &http.Response{
			StatusCode: statusCode,
			Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}

		// 设置headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		// 执行错误响应解析（模拟）
		err := dummyParseResponseToObsError(resp, false)

		// 验证解析结果
		if err == nil {
			obsError, ok := err.(ObsError)
			if !ok {
				t.Error("应返回ObsError类型")
			} else {
				// 验证错误的基本结构
				if obsError.Status == "" {
					t.Error("错误状态未设置")
				}

				// 验证错误代码和消息
				if obsError.Code == "" && obsError.Message == "" {
					t.Log("错误代码和消息都为空")
				}

				// 验证状态码匹配
				if obsError.StatusCode() != statusCode {
					t.Errorf("状态码不匹配: 期望=%d, 实际=%d", statusCode, obsError.StatusCode())
				}
			}
		} else {
			t.Logf("错误解析失败: %v", err)
		}
	})
}

// FuzzXMLResponseParsing 测试XML响应解析
func FuzzXMLResponseParsing(f *testing.F) {
	// 添加种子数据 - 各种XML响应
	seeds := []struct {
		xmlContent string
	}{
		{`<ListBucketResult><Name>test-bucket</Name></ListBucketResult>`},
		{`<GetBucketLocationOutput><Location>cn-north-4</Location></GetBucketLocationOutput>`},
		{`<Object><Key>test-object.txt</Key><Size>1024</Size><ETag>"abc123"</ETag></Object>`},
		{`<CompleteMultipartUploadResult><Location>test-bucket.obs.cn-north-4.myhuaweicloud.com/test-object.txt</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"xyz789"</ETag></CompleteMultipartUploadResult>`},
		{`<ListMultipartUploadsResult><Bucket>test-bucket</Bucket></ListMultipartUploadsResult>`},
		{`<Error><Code>TestError</Code><Message>Test message</Message></Error>`},
	}

	for _, seed := range seeds {
		f.Add(seed.xmlContent)
	}

	f.Fuzz(func(t *testing.T, xmlContent string) {
		// 输入验证
		if len(xmlContent) > 1024*100 { // 100KB限制
			t.Skip("XML内容过长，跳过")
		}

		// 检查危险模式
		dangerousPatterns := []string{
			"<!ENTITY", "<!DOCTYPE", "<xsl:", "<script:",
			"<iframe:", "<svg:", "<style:", "<link:", "<meta:",
			"<!ATTLIST", "<!ELEMENT", "<!NOTATION",
		}
		for _, pattern := range dangerousPatterns {
			if strings.Contains(strings.ToLower(xmlContent), strings.ToLower(pattern)) {
				t.Skip("检测到危险XML模式，跳过")
			}
		}

		// 执行XML解析
		var result map[string]interface{}
		err := xml.Unmarshal([]byte(xmlContent), &result)

		// 验证解析结果
		if err != nil {
			t.Logf("XML解析失败: %v, XML: %s", err, xmlContent[:min(len(xmlContent), 100)])
			return
		}

		// 验证解析后的结果不为空
		if len(result) == 0 {
			t.Log("XML解析结果为空")
		}
	})
}

// FuzzJSONResponseParsing 测试JSON响应解析
func FuzzJSONResponseParsing(f *testing.F) {
	// 添加种子数据 - 各种JSON响应
	seeds := []struct {
		jsonContent string
	}{
		{`{"bucket":"test-bucket","prefix":"test/"}`},
		{`{"code":"AccessDenied","message":"Access Denied"}`},
		{`{"version":"2015-06-30","isTruncated":false,"marker":"","nextMarker":"","contents":[]}`},
		{`{"bucket":"my-bucket","key":"test-object.txt","etag":"\"abc123\""}`},
		{`{"location":"cn-north-4"}`},
	}

	for _, seed := range seeds {
		f.Add(seed.jsonContent)
	}

	f.Fuzz(func(t *testing.T, jsonContent string) {
		// 输入验证
		if len(jsonContent) > 1024*100 { // 100KB限制
			t.Skip("JSON内容过长，跳过")
		}

		// 检查危险模式
		dangerousPatterns := []string{
			"<script", "javascript:", "eval(", "document.cookie",
			"innerHTML", "document.write", "alert(",
		}
		for _, pattern := range dangerousPatterns {
			if strings.Contains(strings.ToLower(jsonContent), strings.ToLower(pattern)) {
				t.Skip("检测到危险JSON模式，跳过")
			}
		}

		// 执行JSON解析
		var result map[string]interface{}
		err := json.Unmarshal([]byte(jsonContent), &result)

		// 验证解析结果
		if err != nil {
			t.Logf("JSON解析失败: %v, JSON: %s", err, jsonContent[:min(len(jsonContent), 100)])
			return
		}

		// 验证解析后的结果不为空
		if len(result) == 0 && len(jsonContent) > 2 { // 至少有 {}
			t.Log("JSON解析结果为空")
		}
	})
}

// FuzzResponseHeaders 测试响应头处理
func FuzzResponseHeaders(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		headers map[string][]string
	}{
		{map[string][]string{
			"Content-Type":       []string{"application/xml"},
			"Content-Length":     []string{"1024"},
			"ETag":             []string{"\"abc123\""},
			"Last-Modified":      []string{"Wed, 01 Jan 2024 00:00:00 GMT"},
			"Request-Id":        []string{"TX000001"},
			"x-obs-request-id": []string{"TX000001"},
		}},
		{map[string][]string{
			"Content-Type":   []string{"application/json"},
			"Cache-Control":  []string{"no-cache"},
			"Content-Type":   []string{"text/plain"},
		}},
		{map[string][]string{}}, // 空headers
	}

	for _, seed := range seeds {
		f.Add(seed.headers)
	}

	f.Fuzz(func(t *testing.T, headers map[string][]string) {
		// 输入验证
		if headers == nil {
			headers = make(map[string][]string)
		}

		// 检查header数量和长度
		totalLength := 0
		for key, values := range headers {
			if len(key) > 1024 {
				t.Skip("header key过长，跳过")
			}
			for _, value := range values {
				if len(value) > 1024*10 {
					t.Skip("header value过长，跳过")
				}
				totalLength += len(key) + len(value)
			}
		}

		// 总header长度限制
		if totalLength > 64*1024 { // 64KB
			t.Skip("headers总长度过大，跳过")
		}

		// 执行header处理（模拟）
		cleanedHeaders := dummyCleanHeaderPrefix(headers, false)

		// 验证处理后的headers
		if cleanedHeaders == nil {
			t.Error("处理后的headers不应为nil")
		}

		// 验证原始headers和清理后的headers数量一致
		if len(cleanedHeaders) > len(headers)*2 { // 允许一定的扩展
			t.Errorf("header数量异常: 原始=%d, 清理后=%d", len(headers), len(cleanedHeaders))
		}
	})
}

// FuzzErrorResponseWithVariousStatusCodes 测试不同状态码的错误响应
func FuzzErrorResponseWithVariousStatusCodes(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		statusCode int
		errorCode  string
		errorMsg   string
	}{
		{400, "InvalidBucketName", "Bucket name is invalid"},
		{403, "AccessDenied", "Access Denied"},
		{404, "NoSuchBucket", "The specified bucket does not exist"},
		{409, "BucketAlreadyExists", "Bucket already exists"},
		{500, "InternalError", "Internal Server Error"},
		{503, "ServiceUnavailable", "Service Unavailable"},
	}

	for _, seed := range seeds {
		f.Add(seed.statusCode, seed.errorCode, seed.errorMsg)
	}

	f.Fuzz(func(t *testing.T, statusCode int, errorCode, errorMsg string) {
		// 输入验证
		if statusCode < 100 || statusCode > 599 {
			t.Skip("无效的状态码，跳过")
		}

		if len(errorCode) > 256 || len(errorMsg) > 1024 {
			t.Skip("错误信息过长，跳过")
		}

		// 创建错误响应XML
		errorXML := fmt.Sprintf(`<Error><Code>%s</Code><Message>%s</Message></Error>`,
			errorCode, errorMsg)

		// 模拟HTTP响应
		resp := &http.Response{
			StatusCode: statusCode,
			Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(errorXML))),
			Header:     make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/xml")

		// 执行错误响应解析
		err := dummyParseResponseToObsError(resp, false)

		// 验证解析结果
		if err != nil {
			obsError, ok := err.(ObsError)
			if !ok {
				t.Error("应返回ObsError类型")
			} else {
				// 验证状态码
				if obsError.StatusCode() != statusCode {
					t.Errorf("状态码不匹配: 期望=%d, 实际=%d", statusCode, obsError.StatusCode())
				}

				// 验证状态字符串
				if !strings.HasPrefix(obsError.Status, fmt.Sprintf("%d", statusCode)) {
					t.Errorf("状态字符串错误: %s", obsError.Status)
				}
			}
		}
	})
}

// FuzzLargeResponsePayload 测试大响应载荷处理
func FuzzLargeResponsePayload(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		payloadSize int
	}{
		{1024},          // 1KB
		{10 * 1024},     // 10KB
		{100 * 1024},    // 100KB
		{1024 * 1024},   // 1MB
		{10 * 1024 * 1024}, // 10MB
	}

	for _, seed := range seeds {
		f.Add(seed.payloadSize)
	}

	f.Fuzz(func(t *testing.T, payloadSize int) {
		// 输入验证
		if payloadSize <= 0 || payloadSize > 50*1024*1024 { // 最大50MB
			t.Skip("无效的载荷大小，跳过")
		}

		// 生成模拟响应体
		responseBody := fmt.Sprintf(`<Result><Size>%d</Size><Data>%s</Data></Result>`,
			payloadSize, strings.Repeat("x", min(payloadSize, 1000)))

		// 监控内存使用
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memBefore := m.Alloc

		// 模拟HTTP响应
		resp := &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(responseBody))),
			Header:     make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/xml")
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(responseBody)))

		// 创建BaseModel实例
		baseModel := &ObsError{}

		// 执行响应解析（模拟）
		err := dummyParseResponseToBaseModel(resp, baseModel, true, false)

		// 检查内存使用
		runtime.ReadMemStats(&m)
		memAfter := m.Alloc
		memDelta := memAfter - memBefore

		// 内存使用不应超过载荷大小的2倍
		if memDelta > int64(payloadSize)*2 {
			t.Errorf("内存使用过高: 载荷=%d bytes, 内存增量=%d bytes", payloadSize, memDelta)
		}

		// 验证解析结果
		if err != nil && payloadSize < 1024*1024 { // 小载荷应该能解析
			t.Logf("大载荷解析失败: %v", err)
		}
	})
}

// FuzzMalformedResponse 测试格式错误的响应处理
func FuzzMalformedResponse(f *testing.F) {
	// 添加种子数据 - 各种格式错误
	seeds := []string{
		`<Error><Code>Test</Code><Message>Test</Error>`, // 缺少闭合标签
		`{invalid json}`,                          // 无效JSON
		`<Error><Code>Test</Code></Error>`,         // 缺少Message
		`<Error>Test</Error>`,                       // 缺少Code和Message
		``,                                         // 空响应
		`<Error><Code>Test</Code><Message>Test</Message>`, // 正常XML
		`{"code":"Test"}`,                           // 缺少message的JSON
		`<NotAnError><Code>Test</Code></NotAnError>`, // 不同的根元素
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, responseBody string) {
		// 输入验证
		if len(responseBody) > 1024*100 { // 100KB限制
			t.Skip("响应体过长，跳过")
		}

		// 模拟HTTP响应
		resp := &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(responseBody))),
			Header:     make(http.Header),
		}

		// 根据内容类型设置Content-Type
		if strings.HasPrefix(responseBody, "{") || strings.HasPrefix(responseBody, "[") {
			resp.Header.Set("Content-Type", "application/json")
		} else {
			resp.Header.Set("Content-Type", "application/xml")
		}

		// 创建BaseModel实例
		baseModel := &ObsError{}

		// 执行响应解析（模拟）
		err := dummyParseResponseToBaseModel(resp, baseModel, true, false)

		// 验证对格式错误的处理
		if err != nil {
			t.Logf("格式错误响应处理: %v", err)
			// 这应该是预期的行为
		} else if len(responseBody) > 0 {
			// 如果没有错误，验证基本结构
			if baseModel.StatusCode() != 200 {
				t.Log("状态码未正确设置")
			}
		}
	})
}

// FuzzResponseWithQueryParameters 测试包含查询参数的响应
func FuzzResponseWithQueryParameters(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		baseURL  string
		params    map[string]string
	}{
		{"https://obs.cn-north-4.myhuaweicloud.com/bucket", map[string]string{"version": "v1", "acl": "public-read"}},
		{"https://example.com/object", map[string]string{"uploadId": "abc123", "partNumber": "1"}},
		{"https://test.com/path", map[string]string{"response-content-type": "image/jpeg", "response-cache-control": "max-age=3600"}},
	}

	for _, seed := range seeds {
		f.Add(seed.baseURL, seed.params)
	}

	f.Fuzz(func(t *testing.T, baseURL string, params map[string]string) {
		// 输入验证
		if len(baseURL) > 2048 {
			t.Skip("URL过长，跳过")
		}

		if params == nil {
			params = make(map[string]string)
		}

		// 检查参数长度
		for key, value := range params {
			if len(key) > 256 || len(value) > 1024 {
				t.Skip("参数过长，跳过")
			}
		}

		// 构建完整URL
		fullURL := baseURL
		if len(params) > 0 {
			fullURL += "?"
			values := url.Values{}
			for key, value := range params {
				values.Add(key, value)
			}
			fullURL += values.Encode()
		}

		// 验证URL有效性
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			t.Logf("URL解析失败: %v", err)
			return
		}

		// 验证查询参数
		if parsedURL.RawQuery != "" {
			queryParams, err := url.ParseQuery(parsedURL.RawQuery)
			if err != nil {
				t.Errorf("查询参数解析失败: %v", err)
			} else {
				// 验证参数数量一致
				if len(queryParams) != len(params) {
					t.Logf("参数数量不一致: 期望=%d, 实际=%d", len(params), len(queryParams))
				}
			}
		}
	})
}

// FuzzConcurrentResponseParsing 测试并发响应解析
func FuzzConcurrentResponseParsing(f *testing.F) {
	// 添加种子数据
	seeds := []struct {
		xmlContent string
	}{
		{`<ListBucketResult><Name>test-bucket</Name></ListBucketResult>`},
		{`<Error><Code>TestError</Code><Message>Test message</Message></Error>`},
	}

	for _, seed := range seeds {
		f.Add(seed.xmlContent)
	}

	f.Fuzz(func(t *testing.T, xmlContent string) {
		// 输入验证
		if len(xmlContent) > 10*1024 { // 10KB限制
			t.Skip("XML内容过长，跳过")
		}

		// 并发测试：同时解析多个响应
		done := make(chan bool, 5)

		for i := 0; i < 5; i++ {
			go func(index int) {
				defer func() { done <- true }()

				// 模拟HTTP响应
				resp := &http.Response{
					StatusCode: 200,
					Status:     "200 OK",
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(xmlContent))),
					Header:     make(http.Header),
				}
				resp.Header.Set("Content-Type", "application/xml")

				// 创建BaseModel实例
				baseModel := &ObsError{}

				// 执行响应解析
				err := dummyParseResponseToBaseModel(resp, baseModel, true, false)
				if err != nil {
					t.Logf("并发解析失败 (goroutine %d): %v", index, err)
				}
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < 5; i++ {
			<-done
		}
	})
}

// ===== 辅助函数（模拟实际的响应解析函数） =====

// dummyParseResponseToBaseModel 模拟ParseResponseToBaseModel函数
func dummyParseResponseToBaseModel(resp *http.Response, baseModel IBaseModel, xmlResult bool, isObs bool) error {
	defer func() {
		errMsg := resp.Body.Close()
		if errMsg != nil {
			fmt.Printf("Failed to close response body: %v\n", errMsg)
		}
	}()

	var body []byte
	var err error
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		return err
	}

	if xmlResult {
		err = xml.Unmarshal(body, baseModel)
	} else {
		err = json.Unmarshal(body, baseModel)
	}

	if err != nil {
		fmt.Printf("Unmarshal error: %v, body: %s\n", err, body)
		return err
	}

	// 设置状态码
	if bm, ok := baseModel.(interface{ setStatusCode(int) }); ok {
		bm.setStatusCode(resp.StatusCode)
	}

	return nil
}

// dummyParseResponseToObsError 模拟ParseResponseToObsError函数
func dummyParseResponseToObsError(resp *http.Response, isObs bool) error {
	isJson := false
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		isJson = strings.Contains(contentType, "json")
	}

	obsError := ObsError{}
	respError := dummyParseResponseToBaseModel(resp, &obsError, !isJson, isObs)
	if respError != nil {
		fmt.Printf("Parse response to BaseModel with error: %v\n", respError)
	}

	obsError.Status = resp.Status

	// 处理header中的错误信息
	if errMsg := resp.Header.Get("x-obs-error-message"); errMsg != "" {
		obsError.Message = errMsg
	}
	if errCode := resp.Header.Get("x-obs-error-code"); errCode != "" {
		obsError.Code = errCode
	}
	if indicator := resp.Header.Get("x-obs-error-indicator"); indicator != "" {
		obsError.Indicator = indicator
	}

	return obsError
}

// dummyCleanHeaderPrefix 模拟cleanHeaderPrefix函数
func dummyCleanHeaderPrefix(headers map[string][]string, isObs bool) map[string][]string {
	cleaned := make(map[string][]string)
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if isObs {
			if strings.HasPrefix(lowerKey, "x-obs-") {
				cleaned[key] = values
			}
		} else {
			if strings.HasPrefix(lowerKey, "x-amz-") {
				cleaned[key] = values
			}
		}
		// 保留标准header
		if lowerKey == "content-type" || lowerKey == "content-length" ||
			lowerKey == "etag" || lowerKey == "last-modified" {
			cleaned[key] = values
		}
	}
	return cleaned
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}