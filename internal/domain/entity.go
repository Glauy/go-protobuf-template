package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// User 用户聚合根
type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(36);not null" json:"id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // 密码哈希不序列化
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// NewUser 创建新用户（工厂方法）
func NewUser(id, username, email, passwordHash string) (*User, error) {
	user := &User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 验证用户数据
	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// BeforeCreate GORM钩子：创建前验证
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.validate()
}

// BeforeUpdate GORM钩子：更新前验证
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 更新时只验证非空字段
	if u.Username != "" && len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if u.Email != "" && !u.isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if u.PasswordHash != "" && strings.TrimSpace(u.PasswordHash) == "" {
		return errors.New("password hash cannot be empty")
	}
	return nil
}

// 验证用户数据
func (u *User) validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return ErrInvalidUsername
	}

	if len(u.Username) < 3 {
		return ErrInvalidUsername
	}

	if u.Email != "" && !u.isValidEmail(u.Email) {
		return ErrInvalidEmail
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		return ErrInvalidPassword
	}

	return nil
}

// 验证邮箱格式
func (u *User) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(email string) error {
	if email != "" && !u.isValidEmail(email) {
		return ErrInvalidEmail
	}
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// UpdatePassword 更新密码
func (u *User) UpdatePassword(passwordHash string) error {
	if strings.TrimSpace(passwordHash) == "" {
		return ErrInvalidPassword
	}
	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()
	return nil
}

// IsValid 检查用户是否有效
func (u *User) IsValid() bool {
	return u.validate() == nil
}

// TableName 指定表名
func (User) TableName() string {
	return "users_test"
}
