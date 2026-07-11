package taskRepository

import "TodoList/internal/core/repository/pool/postgres"

type TaskRepository struct {
	pool postgres.Pool
}

func NewTaskRepository(pool postgres.Pool) *TaskRepository {
	return &TaskRepository{pool}
}
