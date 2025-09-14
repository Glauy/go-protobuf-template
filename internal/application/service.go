package application

import (
	"context"

	"go-protos/internal/application/commands"
	"go-protos/internal/application/queries"
	"go-protos/internal/domain"
)

// UserAppService 用户应用服务
type UserAppService struct {
	// 命令处理器
	createUserHandler      *commands.CreateUserCommandHandler
	updateUserEmailHandler *commands.UpdateUserEmailCommandHandler

	// 查询处理器
	getUserByIdHandler       *queries.GetUserByIdQueryHandler
	getUserByUsernameHandler *queries.GetUserByUsernameQueryHandler
	getUserByEmailHandler    *queries.GetUserByEmailQueryHandler
}

// NewUserAppService 创建用户应用服务
func NewUserAppService(
	userRepo domain.UserRepository,
	userDomainSvc *domain.UserDomainService,
) *UserAppService {
	return &UserAppService{
		createUserHandler:        commands.NewCreateUserCommandHandler(userRepo, userDomainSvc),
		updateUserEmailHandler:   commands.NewUpdateUserEmailCommandHandler(userRepo, userDomainSvc),
		getUserByIdHandler:       queries.NewGetUserByIdQueryHandler(userRepo),
		getUserByUsernameHandler: queries.NewGetUserByUsernameQueryHandler(userRepo),
		getUserByEmailHandler:    queries.NewGetUserByEmailQueryHandler(userRepo),
	}
}

// 命令方法
func (s *UserAppService) CreateUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error) {
	cmd := commands.CreateUserCommand{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
	return s.createUserHandler.Handle(ctx, cmd)
}

func (s *UserAppService) UpdateUserEmail(ctx context.Context, userID, email string) error {
	cmd := commands.UpdateUserEmailCommand{
		UserID: userID,
		Email:  email,
	}
	return s.updateUserEmailHandler.Handle(ctx, cmd)
}

// 查询方法
func (s *UserAppService) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	query := queries.GetUserByIdQuery{UserID: id}
	return s.getUserByIdHandler.Handle(ctx, query)
}

func (s *UserAppService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := queries.GetUserByUsernameQuery{Username: username}
	return s.getUserByUsernameHandler.Handle(ctx, query)
}

func (s *UserAppService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := queries.GetUserByEmailQuery{Email: email}
	return s.getUserByEmailHandler.Handle(ctx, query)
}
