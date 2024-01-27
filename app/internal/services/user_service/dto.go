package user_service

import "time"

type UserDTO struct {
	Id              int       `json:"id"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedAt       time.Time `json:"createdAt"`
}

type EmailVerifyingDTO struct {
	EmailVerifyingLink string `json:"emailVerifyingLink"`
}
