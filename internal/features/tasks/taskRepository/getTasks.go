package taskRepository

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (r *TaskRepository) GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := r.pool.ExecutorFromContext(ctx)

	query := `
	SELECT id,version,title,description,completed,created_at,completed_at,author_user_id FROM todoapp.tasks
	%s
	ORDER BY id ASC
	LIMIT $1 OFFSET $2;
`

	args := []interface{}{limit, offset}
	if userId != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id = $3")
		args = append(args, userId)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var tasks []TaskEntity
	for rows.Next() {
		var task TaskEntity

		err := rows.Scan(
			&task.Id,
			&task.Version,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.UserId,
		)

		if err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := convertDomainsFromEntities(tasks)
	return taskDomains, nil
}
