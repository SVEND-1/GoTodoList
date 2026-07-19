package authApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/transport/http/cookies"
	"TodoList/internal/core/transport/http/server"
	"TodoList/internal/features/auth/authService/jwt"
	"context"
	"fmt"
	"net/http"
	"time"
)

type AuthController struct {
	authService AuthService
	jwtProvider jwt.JwtProvider // СДЕЛАТЬ ИНТЕРФЕЙС
}

type AuthService interface {
	RegisterUser(ctx context.Context, user domain.User) (domain.User, error)
	LoginUser(ctx context.Context, phone string) (domain.User, error)
}

func NewAuthController(authService AuthService, jwtProvider jwt.JwtProvider) *AuthController {
	return &AuthController{
		authService: authService,
		jwtProvider: jwtProvider,
	}
}

func (c *AuthController) Routers() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: c.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: c.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: c.Logout,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: c.Refresh,
		},
	}
}

func (c *AuthController) issueAuthCookies(rw http.ResponseWriter, userId int, role domain.Role) error {
	accessToken, err := c.jwtProvider.GenerateAccessToken(userId, role)
	if err != nil {
		return fmt.Errorf("failed to generate access token %w", err)
	}

	refreshToken, err := c.jwtProvider.GenerateRefreshToken(userId, role)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token, %w", err)
	}

	cookies.SetAuthCookies(
		rw,
		accessToken,
		refreshToken,
		15*time.Minute,
		30*24*time.Hour,
		cookies.Options{
			Secure:   true,
			Domain:   "localhost",
			SameSite: http.SameSiteLaxMode,
		},
	)
	return nil
}
