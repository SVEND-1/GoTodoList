package taskApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type CreateTaskRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=100" example:"Стать senior"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Стать senior разработчиком за 3 часа"`
	UserID      int     `json:"user_id" validate:"required"`
}

type CreateTaskResponse TaskResponseDTO

// CreateTask 	godoc
// @Summary 	Создание задачи
// @Description Создать новую задачу в системе
// @Tags 		Tasks
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateTaskRequest true "Тело запроса"
// @Success 	201 {object} CreateTaskResponse 	"Успешно созданная задача"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "Author not found"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/tasks [post]
func (c *TaskController) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	var req CreateTaskRequest
	if err := requests.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode the request")
		return
	}

	domain := domain.NewTaskUninitialized(
		req.Title, req.Description, req.UserID,
	)
	taskDomain, err := c.taskService.CreateTask(ctx, domain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}
	response := CreateTaskResponse(convertTaskDtoFromDomain(taskDomain))
	responseHandler.JsonResponse(response, http.StatusCreated)
}
