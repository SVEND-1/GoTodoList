package userService

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (s *UserService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user from repository: %w", err)
	}
	return user, nil
}
