package userService

import (
	"context"
	"fmt"
)

func (s *UserService) DeleteUser(ctx context.Context, userId int) error {
	if err := s.userRepository.DeleteUser(ctx, userId); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
