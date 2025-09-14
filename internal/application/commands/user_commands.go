package commands

import (
	"context"

	"go-protos/internal/domain"

	"github.com/google/uuid"
)

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	Username     string
	Email        string
	PasswordHash string
}

// CreateUserCommandHandler 创建用户命令处理器
type CreateUserCommandHandler struct {
	userRepo      domain.UserRepository
	userDomainSvc *domain.UserDomainService
}

// NewCreateUserCommandHandler 创建命令处理器
func NewCreateUserCommandHandler(
	userRepo domain.UserRepository,
	userDomainSvc *domain.UserDomainService,
) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		userRepo:      userRepo,
		userDomainSvc: userDomainSvc,
	}
}

// Handle 处理创建用户命令
func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
	// 验证用户唯一性
	if err := h.userDomainSvc.ValidateUserUniqueness(ctx, cmd.Username, cmd.Email); err != nil {
		return nil, err
	}

	// 生成ID
	var id string = uuid.New().String()

	// 创建用户实体
	user, err := domain.NewUser(id, cmd.Username, cmd.Email, cmd.PasswordHash)
	if err != nil {
		return nil, err
	}

	// 保存到仓储
	if err := h.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserEmailCommand 更新用户邮箱命令
type UpdateUserEmailCommand struct {
	UserID string
	Email  string
}

// UpdateUserEmailCommandHandler 更新用户邮箱命令处理器
type UpdateUserEmailCommandHandler struct {
	userRepo      domain.UserRepository
	userDomainSvc *domain.UserDomainService
}

// NewUpdateUserEmailCommandHandler 创建命令处理器
func NewUpdateUserEmailCommandHandler(
	userRepo domain.UserRepository,
	userDomainSvc *domain.UserDomainService,
) *UpdateUserEmailCommandHandler {
	return &UpdateUserEmailCommandHandler{
		userRepo:      userRepo,
		userDomainSvc: userDomainSvc,
	}
}

// Handle 处理更新邮箱命令
func (h *UpdateUserEmailCommandHandler) Handle(ctx context.Context, cmd UpdateUserEmailCommand) error {
	// 查找用户
	user, err := h.userRepo.FindById(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	// 验证邮箱唯一性
	if cmd.Email != "" {
		emailUnique, err := h.userDomainSvc.IsEmailUnique(ctx, cmd.Email)
		if err != nil {
			return err
		}
		if !emailUnique {
			return domain.ErrEmailExists
		}
	}

	// 更新邮箱
	if err := user.UpdateEmail(cmd.Email); err != nil {
		return err
	}

	// 保存到仓储
	return h.userRepo.Save(ctx, user)
}
