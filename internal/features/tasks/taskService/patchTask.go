package taskService

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (s *TaskService) PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error) {
	var result domain.Task

	err := s.txManager.WithinTx(ctx, func(ctx context.Context) error {
		task, err := s.GetTask(ctx, id)
		if err != nil {
			return fmt.Errorf("get task: %w", err)
		}

		if err := task.ApplyPatch(patch); err != nil {
			return fmt.Errorf("apply task patch: %w", err)
		}

		result, err = s.taskRepository.PatchTask(ctx, id, task)
		if err != nil {
			return fmt.Errorf("patch task: %w", err)
		}
		return nil
	})

	return result, err
}
