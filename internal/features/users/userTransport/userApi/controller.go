package userApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/transport/http/server"
	"context"
	"net/http"
)

type UserController struct {
	UserService UserService
}

type UserService interface {
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.User, error)
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (c *UserController) Routers() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: c.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: c.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: c.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: c.PatchUser,
		},
	}
}
