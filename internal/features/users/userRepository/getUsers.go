package userRepository

import (
	"TodoList/internal/core/domain"
	"context"
	"fmt"
)

func (r *UserRepository) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id,version,full_name,phone_number FROM todoapp.users 
	ORDER BY id ASC
	LIMIT $1 
	OFFSET $2;
`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var users []UserEntity
	for rows.Next() {
		var user UserEntity
		err := rows.Scan(&user.ID, &user.Version, &user.FullName, &user.PhoneNum)
		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	userDomains := convertUserDomainsFromEntities(users)

	return userDomains, nil
}
