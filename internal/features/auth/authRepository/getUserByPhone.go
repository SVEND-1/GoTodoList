package authRepository

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/core/repository/pool/postgres"
	"context"
	"errors"
	"fmt"
)

func (r *AuthRepository) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id,version,full_name,phone_number FROM todoapp.users
    WHERE phone_number = $1
`

	row := r.pool.QueryRow(ctx, query, phone)
	var user UserEntity
	err := row.Scan(
		&user.ID,
		&user.Version,
		&user.FullName,
		&user.PhoneNum,
	)

	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with phone=%s: %w", phone, core_errors.ErrNotFound)
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
