package userApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type GetUserResponse UserDTOResponse

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
