package userRepository

import "TodoList/internal/core/repository/conn"

type UserRepository struct {
	pool conn.Pool
}

func NewUserRepository(pool conn.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}
