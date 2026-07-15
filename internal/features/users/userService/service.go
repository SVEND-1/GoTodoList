package userService

import (
	"TodoList/internal/core/domain"
	"context"
)

type UserService struct {
	userRepository UserRepository
}

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=service.go -destination=mocks/mock_user_service.go -package=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, patch domain.User) (domain.User, error)
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}
