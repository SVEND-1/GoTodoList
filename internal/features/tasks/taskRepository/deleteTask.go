package taskRepository

import (
	core_errors "TodoList/internal/core/errors"
	"context"
	"fmt"
)

func (r *TaskRepository) DeleteTask(ctx context.Context, taskId int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.tasks
	WHERE id=$1
`
	cmdTag, err := r.pool.Exec(ctx, query, taskId)

	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task with id=%d: %w", taskId, core_errors.ErrNotFound)
	}
	return nil
}
