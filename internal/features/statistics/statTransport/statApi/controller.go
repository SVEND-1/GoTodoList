package statApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/transport/http/server"
	"context"
	"time"
)

type StatisticsController struct {
	statService StatisticsService
}

type StatisticsService interface {
	GetStatistics(ctx context.Context, userId *int, from *time.Time, to *time.Time) (domain.Statistics, error)
}

func NewStatisticsController(s StatisticsService) *StatisticsController {
	return &StatisticsController{s}
}

func (c *StatisticsController) Routes() []server.Route {
	return []server.Route{
		{
			Method:  "GET",
			Path:    "/statistics",
			Handler: c.GetStats,
		},
	}
}
