package http

import (
	"context"
	"fmt"
	"net/http"

	"go-protos/config"
	"go-protos/internal/application"
)

// Server HTTP服务器
type Server struct {
	appService *application.UserAppService
	config     *config.AppConfig
	server     *http.Server
}

// NewServer 创建HTTP服务器
func NewServer(appService *application.UserAppService, cfg *config.AppConfig) *Server {
	return &Server{
		appService: appService,
		config:     cfg,
	}
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	readTimeout, err := s.config.GetReadTimeout()
	if err != nil {
		return fmt.Errorf("failed to parse read timeout: %w", err)
	}

	writeTimeout, err := s.config.GetWriteTimeout()
	if err != nil {
		return fmt.Errorf("failed to parse write timeout: %w", err)
	}

	s.server = &http.Server{
		Addr:         s.config.GetAddress(),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      s.setupRoutes(),
	}

	fmt.Printf("HTTP server starting on %s\n", s.config.GetAddress())
	return s.server.ListenAndServe()
}

// Stop 停止HTTP服务器
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// 健康检查
	mux.HandleFunc("/health", s.healthCheck)

	// API路由
	mux.HandleFunc("/api/users/", s.handleUsers)

	return mux
}

// healthCheck 健康检查
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleUsers 处理用户相关请求
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	// 实现HTTP API
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not implemented"))
}
