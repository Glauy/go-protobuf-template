package grpc

import (
	"fmt"
	"net"

	"go-protos/internal/application"
	"go-protos/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server gRPC服务器
type Server struct {
	grpcServer  *grpc.Server
	userService *UserGrpcService
}

// NewServer 创建gRPC服务器
func NewServer(appService *application.UserAppService) *Server {
	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// 创建用户服务
	userService := NewUserGrpcService(appService)

	// 注册服务
	userpb.RegisterUserServiceServer(grpcServer, userService)

	// 启用反射（方便调试）
	reflection.Register(grpcServer)

	return &Server{
		grpcServer:  grpcServer,
		userService: userService,
	}
}

// Start 启动gRPC服务器
func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	fmt.Printf("gRPC server starting on %s\n", addr)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

// Stop 优雅停止gRPC服务器
func (s *Server) Stop() {
	fmt.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}
