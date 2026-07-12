package taskApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	core_types "TodoList/internal/core/transport/http/types"
	"fmt"
	"net/http"
)

type TaskPatchRequest struct {
	Title       core_types.Nullable[string] `json:"title" swaggertype:"string" example:"Стать java/go senior"`
	Description core_types.Nullable[string] `json:"description" swaggertype:"string" example:"Стать java/go senior,но уже за 20 минут"`
	Completed   core_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"true"`
}

func (r *TaskPatchRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("Title can`t be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("Title length must be between 1 and 100")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("Description length must be between 1 and 1000")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("Completed can`t be NULL")
		}
	}
	return nil
}

type TaskPatchResponse TaskResponseDTO

// PatchTask 	godoc
// @Summary 	Изменение задачи
// @Description Изменение информации о задаче
// @Description ### Логика обновления полей:
// @Description 1.**Поле не переданно** значение в бд не меняется
// @Description 2.**Передан null** удаление поля из бд (для title/completed недопустимо)
// @Description 3.**Передано значение** обновление в бд
// @Description 3.**title** и **completed** не могут быть null
// @Tags 		Tasks
// @Accept 		json
// @Produce 	json
// @Param 		id path int true 					"ID задачи"
// @Param 		request body TaskPatchRequest true 	"Тело запроса"
// @Success 	200 {object} TaskPatchResponse 		"Успешно измененная задача"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "Task not found"
// @Failure 	409 {object} response.ErrorResponse "Conflict"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [patch]
func (c *TaskController) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	taskId, err := requests.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskId patch value")
		return
	}

	var taskPatchRequest TaskPatchRequest
	if err := requests.DecodeAndValidateRequest(r, &taskPatchRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	taskPatch := convertPatchFromRequest(taskPatchRequest)
	taskDomain, err := c.taskService.PatchTask(ctx, taskId, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := TaskPatchResponse(convertTaskDtoFromDomain(taskDomain))
	responseHandler.JsonResponse(response, http.StatusOK)
}

func convertPatchFromRequest(request TaskPatchRequest) domain.TaskPatch {
	return domain.TaskPatch{
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	}
}
