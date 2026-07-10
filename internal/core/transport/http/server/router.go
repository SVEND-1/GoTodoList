package server

import (
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
	apiVersion ApiVersion
}

func NewAPIVersionRouter(version ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: version,
	}
}

func (r *APIVersionRouter) RegisterRouters(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.ServeMux.Handle(pattern, route.Handler)
	}
}
