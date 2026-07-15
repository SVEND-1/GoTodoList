package statService

import (
	"TodoList/internal/core/domain"
	"context"
	"time"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=service.go -destination=mocks/mock_stat_service.go -package=mocks
type StatisticsRepository interface {
	GetTasks(ctx context.Context, userId *int, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(r StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: r,
	}
}
