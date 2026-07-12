package userApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

// DeleteUser 	godoc
// @Summary 	Удаление пользователя
// @Description Удалить пользователя в системе
// @Tags 		Users
// @Produce 	json
// @Param		id path int true 					"ID пользователя"
// @Success 	204 								"Успешно удаленный пользователь по Id"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "User not found"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [delete]
func (c *UserController) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	userId, err := requests.GetIntPathValue(r, "id")
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
