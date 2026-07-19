package authApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

type LoginRequest struct {
	Phone string `json:"phone_number" validate:"required,min=10,max=15,startswith=+" example:"+79994443322"`
}

type LoginResponse UserDTOResponse

// Login godoc
// @Summary     Вход в аккаунт
// @Description Вход в аккаунт и создание jwt
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body LoginRequest true "Тело запроса"  // ← Здесь указывается структура
// @Success     200 {object} LoginResponse "Успешный вход в аккаунт"
// @Failure     400 {object} response.ErrorResponse "Bad request"
// @Failure     404 {object} response.ErrorResponse "User not found"
// @Failure     500 {object} response.ErrorResponse "Internal server error"
// @Router      /auth/login [post]
func (c *AuthController) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	var request LoginRequest
	if err := requests.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode request:",
		)
		return
	}

	userDomain, err := c.authService.LoginUser(ctx, request.Phone)
	if err != nil {
		responseHandler.ErrorResponse(err, "User not found")
		return
	}

	err = c.issueAuthCookies(rw, userDomain.Id, domain.RoleAdmin)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create jwt")
	}

	response := LoginResponse(userDomain)
	responseHandler.JsonResponse(response, http.StatusOK)
}
