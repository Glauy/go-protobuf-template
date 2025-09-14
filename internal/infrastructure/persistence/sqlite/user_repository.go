package sqlite

import (
	"context"

	"go-protos/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// 构造函数
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// 保存用户
func (r *UserRepository) Save(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// 根据ID查找
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// 根据用户名查找
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// 根据邮箱查找
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// 更新状态
func (r *UserRepository) UpdateStatus(ctx context.Context, id string, status int8) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("status", status).Error
}
