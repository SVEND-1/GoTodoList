package userApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/response"
	"TodoList/internal/core/transport/http/utils"
	"net/http"
)

func (c *UserController) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	userId, err := utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id path value")
		return
	}
	if err := c.UserService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}
	responseHandler.NoContentResponse()
}
