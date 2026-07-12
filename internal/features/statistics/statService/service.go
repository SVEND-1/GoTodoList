package statService

import (
	"TodoList/internal/core/domain"
	"context"
	"time"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userId *int, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(r StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: r,
	}
}
