package utils

import (
	core_errors "TodoList/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}
	intParam, err := strconv.Atoi(param)

	if err != nil {
		return nil, fmt.Errorf(
			"parameter %s by key:%s is not a valid integer: %v: %w",
			param, key,
			err, core_errors.ErrInvalidArgument,
		)
	}

	return &intParam, nil
}
