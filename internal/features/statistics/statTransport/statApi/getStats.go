package statApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"fmt"
	"net/http"
	"time"
)

type GetStatResponse struct {
	TaskCreated              int      `json:"task_created"`
	TaskCompleted            int      `json:"task_completed"`
	TaskCompletedRate        *float64 `json:"task_completed_rate"`
	TaskAverageCompletedTime *string  `json:"task_average_completed_time"`
}

// GetStats 	godoc
// @Summary 	Получение статистики
// @Description Получение статистики по задачам с фильтрацией по пользователю и периоду
// @Description ### Формат дат:
// @Description Даты передаются в формате **"2006-01-02 15:04:05"** (например: 2026-07-12 00:00:00)
// @Tags 		Statistics
// @Produce 	json
// @Param 		userId query int false 				"ID пользователя"
// @Param 		from query string false 			"Начало периода (2006-01-02 15:04:05)"
// @Param 		to query string false 				"Конец периода (2006-01-02 15:04:05)"
// @Success 	200 {object} GetStatResponse 		"Успешное получение статистики"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/statistics [get]
func (c *StatisticsController) GetStats(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	userId, from, to, err := c.getUserIdFromToQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "GET userId/From/To QueryParam fail")
		return
	}

	statDomain, err := c.statService.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "GET Statistics fail")
		return
	}

	response := convertResponseFromDomain(statDomain)
	responseHandler.JsonResponse(response, http.StatusOK)
}

func convertResponseFromDomain(statistics domain.Statistics) GetStatResponse {
	var avgTime *string
	if statistics.TaskAverageCompletedTime != nil {
		duration := statistics.TaskAverageCompletedTime.String()
		avgTime = &duration
	}
	return GetStatResponse{
		TaskCreated:              statistics.TaskCreated,
		TaskCompleted:            statistics.TaskCompleted,
		TaskCompletedRate:        statistics.TaskCompletedRate,
		TaskAverageCompletedTime: avgTime,
	}
}

func (c *StatisticsController) getUserIdFromToQueryParam(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIdQueryParamKey = "userId"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	userId, err := requests.GetIntQueryParams(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get userId query param: %w", err)
	}

	from, err := requests.GetTimeQueryParams(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get from query param: %w", err)
	}

	to, err := requests.GetTimeQueryParams(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get to query param: %w", err)
	}

	return userId, from, to, nil
}
