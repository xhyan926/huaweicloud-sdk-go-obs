package mock_server

import (
	"fmt"
	"net/http"
)

// MockHandler Mock HTTP处理器
type MockHandler struct {
	server *Server
}

// NewMockHandler 创建新的Mock处理器
func NewMockHandler(server *Server) *MockHandler {
	return &MockHandler{
		server: server,
	}
}

// RegisterOBSRoutes 注册OBS相关的路由
func (h *MockHandler) RegisterOBSRoutes() {
	// 服务相关的API
	h.server.RegisterFuncHandler("/v1/", h.serviceHandler)
	h.server.RegisterFuncHandler("/", h.serviceHandler)

	// 桶相关的API
	h.server.RegisterFuncHandler("/v1/", h.bucketHandler)
	h.server.RegisterFuncHandler("/", h.bucketHandler)

	// 对象相关的API
	h.server.RegisterFuncHandler("/v1/", h.objectHandler)
	h.server.RegisterFuncHandler("/", h.objectHandler)
}

// serviceHandler 处理服务相关请求
func (h *MockHandler) serviceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodHead:
		h.handleServiceHead(w, r)
	default:
		h.handleServiceGeneric(w, r)
	}
}

// bucketHandler 处理桶相关请求
func (h *MockHandler) bucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := extractBucketName(r.URL.Path)
	if bucketName == "" {
		h.sendError(w, http.StatusBadRequest, "Invalid bucket name")
		return
	}

	switch r.Method {
	case http.MethodHead:
		h.handleBucketHead(w, r, bucketName)
	case http.MethodGet:
		h.handleBucketGet(w, r, bucketName)
	case http.MethodPut:
		h.handleBucketPut(w, r, bucketName)
	default:
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// objectHandler 处理对象相关请求
func (h *MockHandler) objectHandler(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := extractBucketAndObject(r.URL.Path)
	if bucketName == "" || objectKey == "" {
		h.sendError(w, http.StatusBadRequest, "Invalid bucket or object name")
		return
	}

	switch r.Method {
	case http.MethodHead:
		h.handleObjectHead(w, r, bucketName, objectKey)
	case http.MethodGet:
		h.handleObjectGet(w, r, bucketName, objectKey)
	case http.MethodPut:
		h.handleObjectPut(w, r, bucketName, objectKey)
	case http.MethodDelete:
		h.handleObjectDelete(w, r, bucketName, objectKey)
	default:
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleServiceHead 处理服务检查HEAD请求
func (h *MockHandler) handleServiceHead(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"x-obs-request-id": {"mock-request-id"},
			"x-obs-id-2":      {"mock-id-2"},
		},
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
}

// handleServiceGeneric 处理通用服务请求
func (h *MockHandler) handleServiceGeneric(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"Content-Type":      {"application/xml"},
			"x-obs-request-id": {"mock-request-id"},
		},
		Body: `<?xml version="1.0" encoding="UTF-8"?>
<CreateBucketConfiguration>
  <StorageClass>Standard</StorageClass>
</CreateBucketConfiguration>`,
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(resp.Body))
}

// handleBucketHead 处理桶检查HEAD请求
func (h *MockHandler) handleBucketHead(w http.ResponseWriter, r *http.Request, bucketName string) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"x-obs-request-id":       {"mock-request-id"},
			"x-obs-id-2":            {"mock-id-2"},
			"Content-Length":        {"0"},
			"x-obs-storage-class":   {"Standard"},
			"x-obs-creation-date":   {"2023-01-01T00:00:00.000Z"},
		},
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
}

// handleBucketGet 处理桶获取请求（列出对象）
func (h *MockHandler) handleBucketGet(w http.ResponseWriter, r *http.Request, bucketName string) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"Content-Type":   {"application/xml"},
			"x-obs-request-id": {"mock-request-id"},
		},
		Body: `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
  <Name>test-bucket</Name>
  <Prefix></Prefix>
  <Marker></Marker>
  <MaxKeys>1000</MaxKeys>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>test-object.txt</Key>
    <LastModified>2023-01-01T00:00:00.000Z</LastModified>
    <ETag>"mock-etag"</ETag>
    <Size>1024</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>test-owner-id</ID>
      <DisplayName>test-owner</DisplayName>
    </Owner>
  </Contents>
</ListBucketResult>`,
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(resp.Body))
}

// handleBucketPut 处理桶创建请求
func (h *MockHandler) handleBucketPut(w http.ResponseWriter, r *http.Request, bucketName string) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"x-obs-request-id": {"mock-request-id"},
			"x-obs-id-2":      {"mock-id-2"},
			"Location":        {fmt.Sprintf("/%s", bucketName)},
		},
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
}

// handleObjectHead 处理对象检查HEAD请求
func (h *MockHandler) handleObjectHead(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"x-obs-request-id":     {"mock-request-id"},
			"Content-Length":      {"1024"},
			"ETag":               {"\"mock-etag\""},
			"Last-Modified":      {"Wed, 01 Jan 2023 00:00:00 GMT"},
			"Content-Type":       {"text/plain"},
			"x-obs-storage-class": {"Standard"},
		},
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
}

// handleObjectGet 处理对象获取请求
func (h *MockHandler) handleObjectGet(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"x-obs-request-id":     {"mock-request-id"},
			"Content-Length":      {"1024"},
			"ETag":               {"\"mock-etag\""},
			"Last-Modified":      {"Wed, 01 Jan 2023 00:00:00 GMT"},
			"Content-Type":       {"text/plain"},
			"x-obs-storage-class": {"Standard"},
		},
		Body: "This is a mock object content for testing purposes.",
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(resp.Body))
}

// handleObjectPut 处理对象上传请求
func (h *MockHandler) handleObjectPut(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"x-obs-request-id": {"mock-request-id"},
			"ETag":            {"\"mock-etag-put\""},
			"Content-Length":  {"0"},
		},
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
}

// handleObjectDelete 处理对象删除请求
func (h *MockHandler) handleObjectDelete(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	resp := &ResponseRecord{
		StatusCode: http.StatusNoContent,
		Headers: map[string][]string{
			"x-obs-request-id": {"mock-request-id"},
		},
	}

	h.server.RecordRequest(r, resp)

	h.setHeaders(w, resp)
	w.WriteHeader(resp.StatusCode)
}

// sendError 发送错误响应
func (h *MockHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	errorBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>%s</Code>
  <Message>%s</Message>
  <Resource></Resource>
  <RequestId>mock-error-request-id</RequestId>
</Error>`, http.StatusText(statusCode), message)

	resp := &ResponseRecord{
		StatusCode: statusCode,
		Headers: map[string][]string{
			"Content-Type": {"application/xml"},
		},
		Body: errorBody,
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(errorBody))
}

// setHeaders 设置响应头
func (h *MockHandler) setHeaders(w http.ResponseWriter, resp *ResponseRecord) {
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}

// resp declared but not used - 删除未使用的变量
// h.handleObjectDelete 中没有使用resp变量，可以直接删除

// extractBucketName 从URL路径中提取桶名
func extractBucketName(path string) string {
	// 简单实现，假设路径格式为 /bucket 或 /bucket/object
	if len(path) > 1 {
		// 找到第一个/后面的部分，直到下一个/
		for i := 1; i < len(path); i++ {
			if path[i] == '/' {
				return path[1:i]
			}
		}
		return path[1:]
	}
	return ""
}

// extractBucketAndObject 从URL路径中提取桶名和对象键
func extractBucketAndObject(path string) (string, string) {
	// 移除开头的/
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// 找到第一个/
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			return path[:i], path[i+1:]
		}
	}

	// 如果没有/，整个路径是桶名，对象键为空
	return path, ""
}