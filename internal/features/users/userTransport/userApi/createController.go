package userApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDTOResponse

func (c *UserController) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	var request CreateUserRequest
	if err := requests.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDto(request)
	userDomain, err := c.UserService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(convertUserDTOFromDomain(userDomain))
	responseHandler.JsonResponse(response, http.StatusCreated)
}

func domainFromDto(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
