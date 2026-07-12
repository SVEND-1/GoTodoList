package userApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type GetUserResponse UserDTOResponse

// GetUser 	godoc
// @Summary 	Получение пользователя
// @Description Получение пользователя по id
// @Tags 		Users
// @Produce 	json
// @Param		id path int true 					"ID пользователя"
// @Success 	200 {object} GetUserResponse 		"Успешно найден пользователь по Id"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "User not found"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [get]
func (c *UserController) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	userId, err := requests.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err, "User Id is required",
		)
		return
	}

	user, err := c.UserService.GetUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "User not found",
		)
		return
	}

	response := GetUserResponse(convertUserDTOFromDomain(user))
	responseHandler.JsonResponse(response, http.StatusOK)
}
