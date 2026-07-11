package userRepository

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/core/repository/pool/postgres"
	"context"
	"errors"
	"fmt"
)

func (r *UserRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET full_name=$1, phone_number=$2,version=version+1
	WHERE id=$3 AND version=$4
	RETURNING id,full_name,phone_number,version;
`
	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, id, user.Version)
	var userEntity UserEntity
	err := row.Scan(
		&userEntity.ID,
		&userEntity.FullName,
		&userEntity.PhoneNum,
		&userEntity.Version,
	)

	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id=%d concurrently accessed: %w", id, core_errors.ErrConflict)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.User{
		userEntity.ID,
		userEntity.Version,
		userEntity.FullName,
		userEntity.PhoneNum,
	}

	return userDomain, nil
}
