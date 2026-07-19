package authApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/cookies"
	"TodoList/internal/core/transport/http/response"
	"net/http"
)

// Logout 			godoc
// @Summary 		Выход с аккаунта
// @Description 	Выход с аккаунта и удаление jwt токена из кук
// @Tags 			Auth
// @Accept 			json
// @Produce 		json
// @Success 		200  "Успешый выход с аккаунта"
// @Failure 		400 {object} response.ErrorResponse "Bad request"
// @Failure 		500 {object} response.ErrorResponse "Internal server error"
// @Router 			/auth/logout [post]
func (c *AuthController) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	cookies.ClearAuthCookies(rw, cookies.Options{
		Secure:   true,
		Domain:   "localhost",
		SameSite: http.SameSiteLaxMode,
	})

	responseHandler.NoContentResponse()
}
