package userRepository

import (
	"TodoList/internal/core/repository/pool/postgres"
)

type UserRepository struct {
	pool postgres.Pool
}

func NewUserRepository(pool postgres.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}
