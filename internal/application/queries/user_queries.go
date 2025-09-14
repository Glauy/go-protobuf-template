package queries

import (
	"context"

	"go-protos/internal/domain"
)

// GetUserByIdQuery 根据ID查询用户
type GetUserByIdQuery struct {
	UserID string
}

// GetUserByIdQueryHandler 查询处理器
type GetUserByIdQueryHandler struct {
	userRepo domain.UserRepository
}

// NewGetUserByIdQueryHandler 创建查询处理器
func NewGetUserByIdQueryHandler(userRepo domain.UserRepository) *GetUserByIdQueryHandler {
	return &GetUserByIdQueryHandler{
		userRepo: userRepo,
	}
}

// Handle 处理查询
func (h *GetUserByIdQueryHandler) Handle(ctx context.Context, query GetUserByIdQuery) (*domain.User, error) {
	return h.userRepo.FindById(ctx, query.UserID)
}

// GetUserByUsernameQuery 根据用户名查询用户
type GetUserByUsernameQuery struct {
	Username string
}

// GetUserByUsernameQueryHandler 查询处理器
type GetUserByUsernameQueryHandler struct {
	userRepo domain.UserRepository
}

// NewGetUserByUsernameQueryHandler 创建查询处理器
func NewGetUserByUsernameQueryHandler(userRepo domain.UserRepository) *GetUserByUsernameQueryHandler {
	return &GetUserByUsernameQueryHandler{
		userRepo: userRepo,
	}
}

// Handle 处理查询
func (h *GetUserByUsernameQueryHandler) Handle(ctx context.Context, query GetUserByUsernameQuery) (*domain.User, error) {
	return h.userRepo.FindByUsername(ctx, query.Username)
}

// GetUserByEmailQuery 根据邮箱查询用户
type GetUserByEmailQuery struct {
	Email string
}

// GetUserByEmailQueryHandler 查询处理器
type GetUserByEmailQueryHandler struct {
	userRepo domain.UserRepository
}

// NewGetUserByEmailQueryHandler 创建查询处理器
func NewGetUserByEmailQueryHandler(userRepo domain.UserRepository) *GetUserByEmailQueryHandler {
	return &GetUserByEmailQueryHandler{
		userRepo: userRepo,
	}
}

// Handle 处理查询
func (h *GetUserByEmailQueryHandler) Handle(ctx context.Context, query GetUserByEmailQuery) (*domain.User, error) {
	return h.userRepo.FindByEmail(ctx, query.Email)
}
