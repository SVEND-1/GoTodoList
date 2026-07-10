package userRepository

import "TodoList/internal/core/domain"

type UserEntity struct {
	ID       int
	Version  int
	FullName string
	PhoneNum *string
}

func convertUserDomainsFromEntities(users []UserEntity) []domain.User {
	usersDomains := make([]domain.User, len(users))
	for i, user := range users {
		usersDomains[i] = domain.NewUser(
			user.ID, user.Version,
			user.FullName, user.PhoneNum,
		)
	}
	return usersDomains
}
