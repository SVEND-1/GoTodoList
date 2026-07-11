package server

import (
	"TodoList/internal/core/transport/http/middleware"
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion  ApiVersion
	middlewares []middleware.Middleware
}

func NewAPIVersionRouter(
	version ApiVersion,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: version,
	}
}

func NewAPIVersionRouterWithMiddlewares(
	version ApiVersion,
	middlewares []middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:    http.NewServeMux(),
		apiVersion:  version,
		middlewares: middlewares,
	}
}

func (r *APIVersionRouter) RegisterRouters(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.ServeMux.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return middleware.ChainMiddleware(
		r,
		r.middlewares...,
	)
}
