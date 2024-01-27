package user_service

import (
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/repositories/user_repo"
)

func mapUserDTOFromDAO(dao *user_repo.UserDAO) *UserDTO {
	return &UserDTO{
		Id:              dao.Id,
		Email:           dao.Email,
		IsEmailVerified: dao.IsEmailVerified,
		CreatedAt:       dao.CreatedAt,
	}
}
