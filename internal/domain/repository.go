package domain

import "context"

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindById(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) error
}
