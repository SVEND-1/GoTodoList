package requests

import (
	core_errors "TodoList/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

func GetTimeQueryParams(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}
	timeParam, err := time.Parse("2006-01-02 15:04:05", param)

	if err != nil {
		return nil, fmt.Errorf(
			"parameter %s by key:%s is not a valid time: %v: %w",
			param, key,
			err, core_errors.ErrInvalidArgument,
		)
	}

	return &timeParam, nil
}
