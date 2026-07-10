package userRepository

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *UserRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id,version,full_name,phone_number FROM todoapp.users
	WHERE id = $1;
`
	row := r.pool.QueryRow(ctx, query, id)
	var user UserEntity
	err := row.Scan(
		&user.ID,
		&user.Version,
		&user.FullName,
		&user.PhoneNum,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id=%d: %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		user.ID,
		user.Version,
		user.FullName,
		user.PhoneNum,
	)
	return userDomain, nil
}
