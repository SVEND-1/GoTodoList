package statRepository

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
	"strings"
	"time"
)

func (r *StatisticsRepository) GetTasks(ctx context.Context, userId *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query strings.Builder

	query.WriteString(`
	SELECT id,version,title,description,completed,created_at,completed_at,author_user_id FROM todoapp.tasks
`)
	args := []any{}
	conditions := []string{}

	if userId != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, userId)
	}
	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, from)
	}
	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}
	query.WriteString(" ORDER BY id ASC;")

	rows, err := r.pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from repository: %w", err)
	}
	defer rows.Close()

	var tasks []TaskEntity
	for rows.Next() {
		var task TaskEntity
		err := rows.Scan(
			&task.Id, &task.Version,
			&task.Title, &task.Description,
			&task.Completed, &task.CreatedAt,
			&task.CompletedAt, &task.UserId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	tasksDomain := convertDomainsFromEntities(tasks)
	return tasksDomain, nil
}
