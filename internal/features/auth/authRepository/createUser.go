package authRepository

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/features/users/userRepository"
	"context"
	"fmt"
)

func (r *AuthRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout()) //TODO пароли хешировать надо когда они будут
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name,phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;
`
	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userEntity userRepository.UserEntity
	err := row.Scan(&userEntity.ID, &userEntity.Version, &userEntity.FullName, &userEntity.PhoneNum)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan user error: %w", err)
	}

	userDomain := domain.NewUser(userEntity.ID, userEntity.Version, userEntity.FullName, userEntity.PhoneNum)
	return userDomain, nil
}
