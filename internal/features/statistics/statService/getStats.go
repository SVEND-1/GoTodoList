package statService

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"context"
	"fmt"
	"time"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, userId *int, from *time.Time, to *time.Time) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("to must be before from: %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userId, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get Task from repository: %w", err)
	}

	statistics := calcStatistics(tasks)
	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{}
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completedDuration := task.CompletedDuration()
		if completedDuration != nil {
			totalCompletedDuration += *completedDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletedTime *time.Duration
	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)
		tasksAverageCompletedTime = &avg
	}

	return domain.Statistics{
		tasksCreated, tasksCompleted,
		&tasksCompletedRate, tasksAverageCompletedTime,
	}
}
