package requests

import (
	core_errors "TodoList/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("no key=%s in path value,%w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	value, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("path values=%s by key=%s not valid integer: %w",
			value, key, core_errors.ErrInvalidArgument,
		)
	}

	return value, nil
}
