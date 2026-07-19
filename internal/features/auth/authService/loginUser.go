package authService

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (s *AuthService) LoginUser(ctx context.Context, phone string) (domain.User, error) {
	user, err := s.authRepository.GetUserByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by phone %w", err)
	}

	return user, nil
}
