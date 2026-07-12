package statRepository

import "TodoList/internal/core/repository/pool/postgres"

type StatisticsRepository struct {
	pool postgres.Pool
}

func NewStatisticsRepository(pool postgres.Pool) *StatisticsRepository {
	return &StatisticsRepository{pool}
}
