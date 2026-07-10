package userApi

import "TodoList/internal/core/domain"

type UserDTOResponse struct {
	Id          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func convertUserDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		Id:          user.Id,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func convertUserDTOsFromDomains(user []domain.User) []UserDTOResponse {
	userDTOs := make([]UserDTOResponse, len(user))
	for i, userDTO := range user {
		userDTOs[i] = convertUserDTOFromDomain(userDTO)
	}
	return userDTOs
}
