package taskRepository

import (
	"TodoList/internal/core/domain"
	"time"
)

type TaskEntity struct {
	Id          int
	Version     int
	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
	UserId      int
}

func convertDomainsFromEntities(entities []TaskEntity) []domain.Task {
	result := make([]domain.Task, len(entities))
	for i, entity := range entities {
		result[i] = convertDomainFromEntity(entity)
	}
	return result
}

func convertDomainFromEntity(entity TaskEntity) domain.Task {
	return domain.NewTask(
		entity.Id,
		entity.Version,
		entity.Title,
		entity.Description,
		entity.Completed,
		entity.CreatedAt,
		entity.CompletedAt,
		entity.UserId,
	)
}
