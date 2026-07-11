package taskApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/transport/http/server"
	"context"
)

type TaskController struct {
	taskService TaskService
}

type TaskService interface {
	CreateTask(ctx context.Context, taskDomain domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
	PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error)
}

func NewTaskController(taskService TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

func (c *TaskController) Routers() []server.Route {
	return []server.Route{
		{
			Method:  "POST",
			Path:    "/tasks",
			Handler: c.CreateTask,
		},
		{
			Method:  "GET",
			Path:    "/tasks",
			Handler: c.GetTasks,
		},
		{
			Method:  "GET",
			Path:    "/tasks/{id}",
			Handler: c.GetTask,
		},
		{
			Method:  "DELETE",
			Path:    "/tasks/{id}",
			Handler: c.DeleteTask,
		},
		{
			Method:  "PATCH",
			Path:    "/tasks/{id}",
			Handler: c.PatchTask,
		},
	}
}
