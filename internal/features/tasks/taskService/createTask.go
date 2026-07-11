package taskService

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (s *TaskService) CreateTask(ctx context.Context, taskDomain domain.Task) (domain.Task, error) {
	if err := taskDomain.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}
	task, err := s.taskRepository.CreateTask(ctx, taskDomain)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}
