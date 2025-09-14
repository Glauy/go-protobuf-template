package grpc

import (
	"context"
	"fmt"
	"time"

	"go-protos/internal/application"
	"go-protos/internal/domain"
	"go-protos/proto/userpb"
)

// UserGrpcService 用户gRPC服务实现
type UserGrpcService struct {
	userpb.UnimplementedUserServiceServer
	appService *application.UserAppService
}

// NewUserGrpcService 创建用户gRPC服务
func NewUserGrpcService(appService *application.UserAppService) *UserGrpcService {
	return &UserGrpcService{
		appService: appService,
	}
}

// GetUserByUsername 根据用户名获取用户
func (s *UserGrpcService) GetUserByUsername(ctx context.Context, req *userpb.GetUserByUsernameRequest) (*userpb.GetUserByUsernameResponse, error) {
	fmt.Printf("gRPC GetUserByUsername called with username: %s\n", req.Username)

	user, err := s.appService.GetUserByUsername(ctx, req.Username)
	if err != nil {
		fmt.Printf("gRPC GetUserByUsername failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("gRPC GetUserByUsername success for user: %s\n", req.Username)
	return &userpb.GetUserByUsernameResponse{
		User: toProtoUser(user),
	}, nil
}

// GetUserById 根据ID获取用户
func (s *UserGrpcService) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserByIdResponse, error) {
	fmt.Printf("gRPC GetUserById called with id: %s\n", req.Id)

	user, err := s.appService.GetUserById(ctx, req.Id)
	if err != nil {
		fmt.Printf("gRPC GetUserById failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("gRPC GetUserById success for user: %s\n", req.Id)
	return &userpb.GetUserByIdResponse{
		User: toProtoUser(user),
	}, nil
}

// CreateUser 创建用户
func (s *UserGrpcService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	fmt.Printf("gRPC CreateUser called with username: %s, email: %s\n", req.Username, req.Email)

	user, err := s.appService.CreateUser(ctx, req.Username, req.Email, req.PasswordHash)
	if err != nil {
		fmt.Printf("gRPC CreateUser failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("gRPC CreateUser success for user: %s\n", req.Username)
	return &userpb.CreateUserResponse{
		User: toProtoUser(user),
	}, nil
}

// UpdateUserEmail 更新用户邮箱
func (s *UserGrpcService) UpdateUserEmail(ctx context.Context, req *userpb.UpdateUserEmailRequest) (*userpb.UpdateUserEmailResponse, error) {
	fmt.Printf("gRPC UpdateUserEmail called with userID: %s, email: %s\n", req.UserId, req.Email)

	err := s.appService.UpdateUserEmail(ctx, req.UserId, req.Email)
	if err != nil {
		fmt.Printf("gRPC UpdateUserEmail failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("gRPC UpdateUserEmail success for user: %s\n", req.UserId)
	return &userpb.UpdateUserEmailResponse{
		Success: true,
	}, nil
}

// toProtoUser 将领域用户转换为protobuf用户
func toProtoUser(u *domain.User) *userpb.User {
	if u == nil {
		return nil
	}
	return &userpb.User{
		Id:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
	}
}
