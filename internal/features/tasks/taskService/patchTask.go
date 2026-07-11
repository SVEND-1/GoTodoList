package taskService

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (s *TaskService) PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error) {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return task, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return task, fmt.Errorf("apply task patch: %w", err)
	}

	taskDomain, err := s.taskRepository.PatchTask(ctx, id, task)
	if err != nil {
		return taskDomain, fmt.Errorf("patch task: %w", err)
	}
	return taskDomain, nil
}
