package authService

import (
	"TodoList/internal/core/domain"
	"context"
)

type AuthService struct {
	authRepository AuthRepository
}

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=service.go -destination=mocks/mock_auth_service.go -package=mocks
type AuthRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (domain.User, error)
}

func NewAuthService(authRepository AuthRepository) *AuthService {
	return &AuthService{authRepository: authRepository}
}
