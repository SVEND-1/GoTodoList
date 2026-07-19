package authRepository

import "TodoList/internal/core/repository/pool/postgres"

type AuthRepository struct {
	pool postgres.Pool
}

func NewAuthRepository(pool postgres.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}
