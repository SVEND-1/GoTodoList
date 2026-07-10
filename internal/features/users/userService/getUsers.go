package userService

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"context"
	"fmt"
)

func (s UserService) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be negative: %w", core_errors.ErrInvalidArgument)
	}

	users, err := s.userRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users from repository: %w", err)
	}
	return users, nil
}
