package domain

import (
	"context"
)

// UserDomainService 用户领域服务
type UserDomainService struct {
	userRepo UserRepository
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(userRepo UserRepository) *UserDomainService {
	return &UserDomainService{
		userRepo: userRepo,
	}
}

// IsUsernameUnique 检查用户名是否唯一
func (s *UserDomainService) IsUsernameUnique(ctx context.Context, username string) (bool, error) {
	existing, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	return existing == nil, nil
}

// IsEmailUnique 检查邮箱是否唯一
func (s *UserDomainService) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return true, nil
	}
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return existing == nil, nil
}

// ValidateUserUniqueness 验证用户唯一性
func (s *UserDomainService) ValidateUserUniqueness(ctx context.Context, username, email string) error {
	// 检查用户名唯一性
	usernameUnique, err := s.IsUsernameUnique(ctx, username)
	if err != nil {
		return err
	}
	if !usernameUnique {
		return ErrUsernameExists
	}

	// 检查邮箱唯一性
	if email != "" {
		emailUnique, err := s.IsEmailUnique(ctx, email)
		if err != nil {
			return err
		}
		if !emailUnique {
			return ErrEmailExists
		}
	}

	return nil
}
