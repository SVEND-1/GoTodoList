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
	Title       core_types.Nullable[string] `json:"title"`
	Description core_types.Nullable[string] `json:"description"`
	Completed   core_types.Nullable[bool]   `json:"completed"`
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
