package mock_server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Server Mock服务器
type Server struct {
	httpServer   *http.Server
	mux          *http.ServeMux
	requests     []RequestRecord
	mu           sync.RWMutex
	shutdownChan chan struct{}
}

// RequestRecord 请求记录
type RequestRecord struct {
	Method     string
	URL        string
	Headers    map[string][]string
	Body       string
	Timestamp  time.Time
	Response   *ResponseRecord
}

// ResponseRecord 响应记录
type ResponseRecord struct {
	StatusCode int
	Headers    map[string][]string
	Body       string
}

// NewServer 创建新的Mock服务器
func NewServer() *Server {
	mux := http.NewServeMux()
	server := &Server{
		mux:          mux,
		requests:     make([]RequestRecord, 0),
		shutdownChan: make(chan struct{}),
	}

	// 注册默认路由
	server.RegisterFuncHandler("/", server.defaultHandler)
	server.RegisterFuncHandler("/health", server.healthHandler)

	return server
}

// RegisterHandler 注册路由处理器
func (s *Server) RegisterHandler(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// RegisterFuncHandler 注册函数处理器
func (s *Server) RegisterFuncHandler(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	s.mux.Handle(pattern, http.HandlerFunc(handlerFunc))
}

// Start 启动Mock服务器
func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Mock server error: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止Mock服务器
func (s *Server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}
	return nil
}

// GetRequests 获取所有请求记录
func (s *Server) GetRequests() []RequestRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回副本
	requests := make([]RequestRecord, len(s.requests))
	copy(requests, s.requests)
	return requests
}

// ClearRequests 清除请求记录
func (s *Server) ClearRequests() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = make([]RequestRecord, 0)
}

// GetRequestsByMethod 按方法获取请求记录
func (s *Server) GetRequestsByMethod(method string) []RequestRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var requests []RequestRecord
	for _, req := range s.requests {
		if req.Method == method {
			requests = append(requests, req)
		}
	}
	return requests
}

// GetLastRequest 获取最后一个请求
func (s *Server) GetLastRequest() *RequestRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.requests) == 0 {
		return nil
	}
	return &s.requests[len(s.requests)-1]
}

// RecordRequest 记录请求
func (s *Server) RecordRequest(req *http.Request, resp *ResponseRecord) {
	record := RequestRecord{
		Method:    req.Method,
		URL:       req.URL.String(),
		Headers:   make(map[string][]string),
		Body:      "", // 注意：这里应该读取请求体，但需要处理EOF
		Timestamp: time.Now(),
		Response:  resp,
	}

	// 复制请求头
	for key, values := range req.Header {
		record.Headers[key] = make([]string, len(values))
		copy(record.Headers[key], values)
	}

	s.mu.Lock()
	s.requests = append(s.requests, record)
	s.mu.Unlock()
}

// defaultHandler 默认处理器
func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers:    make(map[string][]string),
		Body:       "Mock Server Response",
	}

	// 记录请求
	s.RecordRequest(r, resp)

	// 返回响应
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(resp.Body))
}

// healthHandler 健康检查处理器
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseRecord{
		StatusCode: http.StatusOK,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: `{"status": "ok", "timestamp": "2023-01-01T00:00:00Z"}`,
	}

	s.RecordRequest(r, resp)

	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(resp.Body))
}

// GetStats 获取服务器统计信息
func (s *Server) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalRequests := len(s.requests)
	methodStats := make(map[string]int)

	for _, req := range s.requests {
		methodStats[req.Method]++
	}

	return map[string]interface{}{
		"total_requests": totalRequests,
		"method_stats":   methodStats,
		"uptime":        time.Since(s.requests[0].Timestamp).String(),
	}
}

// Reset 重置服务器状态
func (s *Server) Reset() {
	s.ClearRequests()
	// 这里可以添加其他需要重置的状态
}

// IsRunning 检查服务器是否正在运行
func (s *Server) IsRunning() bool {
	return s.httpServer != nil && s.httpServer != nil
}