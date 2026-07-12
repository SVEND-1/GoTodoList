package domain

import "time"

type Statistics struct {
	TaskCreated              int
	TaskCompleted            int
	TaskCompletedRate        *float64
	TaskAverageCompletedTime *time.Duration
}

func NewStatistics(
	taskCreated int,
	taskCompleted int,
	taskCompletedRate *float64,
	taskAverageCompletedTime *time.Duration,
) Statistics {
	return Statistics{
		taskCreated, taskCompleted,
		taskCompletedRate, taskAverageCompletedTime,
	}
}
