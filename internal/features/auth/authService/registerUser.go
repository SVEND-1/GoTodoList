package authService

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (s *AuthService) RegisterUser(ctx context.Context, user domain.User) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("invalid user: %w", err)
	}

	user, err := s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
