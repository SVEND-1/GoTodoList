package taskRepository

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/core/repository/pool/postgres"
	"context"
	"errors"
	"fmt"
)

func (r *TaskRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (title, description,completed,created_at,completed_at,author_user_id)  
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id,version,title, description,completed,created_at,completed_at,author_user_id;
`
	row := r.pool.QueryRow(ctx, query,
		task.Title, task.Description, task.Completed, task.CreatedAt, task.CompletedAt, task.UserId,
	)
	var entity TaskEntity

	err := row.Scan(
		&entity.Id,
		&entity.Version,
		&entity.Title,
		&entity.Description,
		&entity.Completed,
		&entity.CreatedAt,
		&entity.CompletedAt,
		&entity.UserId)
	if err != nil {
		if errors.Is(err, postgres.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf("%v task violates foreign key: %w", err, core_errors.ErrNotFound)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := convertDomainFromEntity(entity)
	return taskDomain, nil
}
