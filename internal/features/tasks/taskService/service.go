package taskService

import (
	"TodoList/internal/core/domain"
	"context"
)

type TaskService struct {
	taskRepository TaskRepository
	txManager      TxManager
}

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=service.go -destination=mocks/mock_task_service.go -package=mocks
type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, taskId int) (domain.Task, error)
	DeleteTask(ctx context.Context, taskId int) error
	PatchTask(ctx context.Context, taskId int, task domain.Task) (domain.Task, error)
}

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=service.go -destination=mocks/mock_task_service.go -package=mocks
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTaskService(taskRepository TaskRepository, manager TxManager) *TaskService {
	return &TaskService{taskRepository: taskRepository, txManager: manager}
}
