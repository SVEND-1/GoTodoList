package taskRepository

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/core/repository/pool/postgres"
	"context"
	"errors"
	"fmt"
)

func (r *TaskRepository) PatchTask(ctx context.Context, taskId int, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET title=$1,description=$2,completed=$3,completed_at=$4,version = version + 1
	WHERE id=$5 AND version = $6
	RETURNING id,version,title, description,completed,created_at,completed_at,author_user_id;
`
	row := r.pool.QueryRow(
		ctx, query,
		task.Title, task.Description, task.Completed, task.CompletedAt,
		taskId, task.Version,
	)

	var taskEntity TaskEntity
	err := row.Scan(
		&taskEntity.Id,
		&taskEntity.Version,
		&taskEntity.Title,
		&taskEntity.Description,
		&taskEntity.Completed,
		&taskEntity.CreatedAt,
		&taskEntity.CompletedAt,
		&taskEntity.UserId)

	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id=%d concurrently accessed: %w", taskId, core_errors.ErrConflict)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", core_errors.ErrConflict)
	}

	taskDomain := convertDomainFromEntity(taskEntity)
	return taskDomain, nil
}
