package taskApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

// DeleteTask 	godoc
// @Summary 	Удаление задачи
// @Description Удалить задачу в системе
// @Tags 		Tasks
// @Produce 	json
// @Param 		id path int true 					"ID задачи"
// @Success 	204 								"Успешно удаленная задача по Id"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "Task not found"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [delete]
func (c *TaskController) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	taskId, err := requests.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskId patch value")
		return
	}

	if err := c.taskService.DeleteTask(ctx, taskId); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete task")
		return
	}
	responseHandler.NoContentResponse()
}
