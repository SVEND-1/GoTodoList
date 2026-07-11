package taskService

import (
	"TodoList/internal/core/domain"
	"context"
)

type TaskService struct {
	taskRepository TaskRepository
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, taskId int) (domain.Task, error)
	DeleteTask(ctx context.Context, taskId int) error
	PatchTask(ctx context.Context, taskId int, task domain.Task) (domain.Task, error)
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}
