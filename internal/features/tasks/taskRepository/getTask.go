package taskRepository

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/core/repository/pool/postgres"
	"context"
	"errors"
	"fmt"
)

func (r *TaskRepository) GetTask(ctx context.Context, taskId int) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id,version,title,description,completed,created_at,completed_at,author_user_id FROM todoapp.tasks
	WHERE id = $1;
`
	row := r.pool.QueryRow(ctx, query, taskId)

	var result domain.Task
	err := row.Scan(
		&result.Id,
		&result.Version,
		&result.Title,
		&result.Description,
		&result.Completed,
		&result.CreatedAt,
		&result.CompletedAt,
		&result.UserId)

	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("get with id=%d: %w", taskId, core_errors.ErrNotFound)
		}
		return result, fmt.Errorf("get task from repository: %w", err)
	}

	taskDomain := domain.Task{
		Id:          result.Id,
		Version:     result.Version,
		Title:       result.Title,
		Description: result.Description,
		Completed:   result.Completed,
		CreatedAt:   result.CreatedAt,
		CompletedAt: result.CompletedAt,
		UserId:      result.UserId,
	}
	return taskDomain, nil
}
