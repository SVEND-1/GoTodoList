package userApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"  example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79994445533"`
}

type CreateUserResponse UserDTOResponse

// CreateUser 	godoc
// @Summary 	Создание пользователя
// @Description Создать нового пользователя в системе
// @Tags 		Users
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateUserRequest true "Тело запроса"
// @Success 	201 {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/users [post]
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
