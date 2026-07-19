package authApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type RegisterUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"  example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79994445533"`
}

type RegisterUserResponse UserDTOResponse

// RegisterUser 	godoc
// @Summary 		Регистрация пользователя
// @Description 	Регистрация пользователя в систему
// @Tags 			Auth
// @Accept 			json
// @Produce 		json
// @Param 			request body RegisterUserRequest true "Тело запроса"
// @Success 		201 {object} RegisterUserResponse "Успешно регистрации пользователя пользователь"
// @Failure 		400 {object} response.ErrorResponse "Bad request"
// @Failure 		500 {object} response.ErrorResponse "Internal server error"
// @Router 			/auth/register [post]
func (c *AuthController) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	var request RegisterUserRequest
	if err := requests.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDto(request)
	userDomain, err := c.authService.RegisterUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	err = c.issueAuthCookies(rw, userDomain.Id, domain.RoleUser)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create jwt")
	}

	response := RegisterUserResponse(convertUserDTOFromDomain(userDomain))
	responseHandler.JsonResponse(response, http.StatusCreated)
}

func domainFromDto(dto RegisterUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
