package taskService

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"context"
	"fmt"
)

func (s *TaskService) GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be negative: %w", core_errors.ErrInvalidArgument)
	}

	tasks, err := s.taskRepository.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from repository: %w", err)
	}
	return tasks, nil
}
