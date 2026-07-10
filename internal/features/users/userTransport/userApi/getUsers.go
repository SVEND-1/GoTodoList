package userApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/response"
	"TodoList/internal/core/transport/http/utils"
	"fmt"
	"net/http"
)

type GetUsersResponse []UserDTOResponse

func (c *UserController) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit and offset query params")
		return
	}

	userDomains, err := c.UserService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(convertUserDTOsFromDomains(userDomains))
	responseHandler.JsonResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := utils.GetIntQueryParams(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get limit query param: %w", err)
	}

	offset, err := utils.GetIntQueryParams(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get offset query param: %w", err)
	}

	return limit, offset, nil
}
