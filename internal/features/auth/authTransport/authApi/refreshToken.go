package authApi

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/cookies"
	"TodoList/internal/core/transport/http/response"
	"net/http"
	"time"
)

// RefreshToken 	godoc
// @Summary 		Refresh для генерации jwt
// @Description 	Refresh для генерации jwt access токена
// @Tags 			Auth
// @Accept 			json
// @Produce 		json
// @Success 		200  "Успешное обновление токена"
// @Failure 		400 {object} response.ErrorResponse "Bad request"
// @Failure 		500 {object} response.ErrorResponse "Internal server error"
// @Router 			/auth/refresh [post]
func (c *AuthController) Refresh(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	refreshTokenStr, err := cookies.GetTokenFromCookie(r, cookies.RefreshTokenCookie)
	if err != nil {
		responseHandler.ErrorResponse(err, "missing refresh token cookie")
		return
	}

	claims, err := c.jwtProvider.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid refresh token")
		return
	}

	newAccessToken, err := c.jwtProvider.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to generate access token")
		return
	}

	newRefreshToken, err := c.jwtProvider.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to generate refresh token")
		return
	}

	cookies.SetAuthCookies(
		rw,
		newAccessToken,
		newRefreshToken,
		15*time.Minute,
		30*24*time.Hour,
		cookies.Options{
			Secure:   true,
			Domain:   "localhost",
			SameSite: http.SameSiteLaxMode,
		},
	)

	responseHandler.NoContentResponse()
}
