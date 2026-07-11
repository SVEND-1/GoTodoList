package taskApi

import (
	"TodoList/internal/core/domain"
	"time"
)

type TaskResponseDTO struct {
	Id          int        `json:"id"`
	Version     int        `json:"version"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
	UserId      int        `json:"userId"`
}

func convertTaskDtoFromDomain(task domain.Task) TaskResponseDTO {
	return TaskResponseDTO{
		Id:          task.Id,
		Version:     task.Version,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		CompletedAt: task.CompletedAt,
		UserId:      task.UserId,
	}
}

func converterTaskDTOsFromDomains(tasks []domain.Task) []TaskResponseDTO {
	dtos := make([]TaskResponseDTO, len(tasks))
	for i, task := range tasks {
		dtos[i] = convertTaskDtoFromDomain(task)
	}
	return dtos
}
