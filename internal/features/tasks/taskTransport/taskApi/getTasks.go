package taskApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type TasksResponse []TaskResponseDTO

// GetTasks 	godoc
// @Summary 	Получение задач
// @Description Получение задач с фильтрацией по пользователю и паггинацией
// @Tags 		Tasks
// @Produce 	json
// @Param 		userId query int false 			"ID пользователя"
// @Param 		limit query int false 				"Размер страницы с задачами"
// @Param 		offset query int false 			"Смещение страницы с задачами"
// @Success 	200 {object} TasksResponse 			"Успешное получение списка задач"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/tasks [get]
func (c *TaskController) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)
	userId, limit, offset, err := getUserIdLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId/limit/offset query params")
		return
	}
	domain, err := c.taskService.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := TasksResponse(converterTaskDTOsFromDomains(domain))
	responseHandler.JsonResponse(response, http.StatusOK)
}

func getUserIdLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIdQueryParamKey = "userId"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	userId, err := requests.GetIntQueryParams(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get userId query param: %w", err)
	}

	limit, err := requests.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit query param: %w", err)
	}

	offset, err := requests.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get offset query param: %w", err)
	}

	return userId, limit, offset, nil
}
