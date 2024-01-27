package user_repo

import "time"

type UserDAO struct {
	Id                 int
	Email              string
	IsEmailVerified    bool
	EmailVerifyingHash string
	CreatedAt          time.Time
}

type UserCreationData struct {
	Email              string
	EmailVerifyingHash string
}
