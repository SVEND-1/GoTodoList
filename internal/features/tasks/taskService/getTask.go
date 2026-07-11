package taskService

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (s *TaskService) GetTask(ctx context.Context, id int) (domain.Task, error) {
	task, err := s.taskRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}
	return task, nil
}
