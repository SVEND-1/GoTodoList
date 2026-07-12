package taskApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type TaskResponse TaskResponseDTO

// GetTask 	godoc
// @Summary 	Получение задачи
// @Description Получение задачи по id
// @Tags 		Tasks
// @Produce 	json
// @Param 		id path int true 					"ID задачи"
// @Success 	200 {object} TaskResponse 			"Успешно найденная задача по Id"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "Task not found"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [get]
func (c *TaskController) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)
	taskId, err := requests.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskId patch value")
		return
	}

	taskDomain, err := c.taskService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}
	response := TaskResponse(taskDomain)
	responseHandler.JsonResponse(response, http.StatusOK)
}
