package user_service

import (
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/email-auth/app/internal/repositories/user_repo"
	"github.com/AleksandrVishniakov/email-auth/app/internal/services/email_service"
	"math/rand"
	"time"
)

type UserService interface {
	GetUserByEmail(email string) (*UserDTO, error)
	AuthUser(email string) (bool, error)
	VerifyEmail(email string, code int) (bool, error)
}

type userService struct {
	userRepository user_repo.UserRepository
	emailService   email_service.EmailService
}

func NewUserService(repo user_repo.UserRepository, emailService email_service.EmailService) UserService {
	return &userService{
		userRepository: repo,
		emailService:   emailService,
	}
}

func (u *userService) GetUserByEmail(email string) (*UserDTO, error) {
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return mapUserDTOFromDAO(user), nil
}

func (u *userService) AuthUser(email string) (bool, error) {
	var isUserExists bool
	user, err := u.userRepository.GetUserByEmail(email)

	if err == nil && user != nil {
		isUserExists = true
	}

	if errors.Is(err, user_repo.ErrUserNotFound) {
		err = nil
		isUserExists = false
	}

	if err != nil {
		return false, err
	}

	var isUserAuthenticated = isUserExists && user.IsEmailVerified

	code := generateCode()

	if isUserAuthenticated || isUserExists {
		err = u.userRepository.UpdateEmailVerifyingCode(email, code)
		if err != nil {
			return isUserAuthenticated, err
		}
	} else {
		err = u.userRepository.NewUser(&user_repo.UserCreationData{
			Email:              email,
			EmailVerifyingCode: code,
		})

		if err != nil {
			return isUserAuthenticated, err
		}
	}

	err = u.emailService.Write(&email_service.EmailContent{
		To:      email,
		Subject: createEmailSubject(code),
		Body:    createEmailMessage(code),
	})

	return isUserAuthenticated, err
}

func (u *userService) VerifyEmail(email string, code int) (bool, error) {
	user, err := u.userRepository.GetUserByEmail(email)

	if err != nil {
		return false, err
	}

	if user.EmailVerifyingCode != code {
		return false, nil
	}

	err = u.userRepository.MarkEmailAsVerified(email)
	if err != nil {
		return false, err
	}

	err = u.userRepository.ResetEmailVerifyingCode(email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func generateCode() int {
	const minN = 100_000
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(900_000) + minN
}

func createEmailSubject(code int) string {
	codeStr := fmt.Sprintf("[%d]", code)
	return fmt.Sprintf("Auth %s", codeStr)
}

func createEmailMessage(code int) string {
	codeStr := fmt.Sprintf("[%d]", code)
	return fmt.Sprintf("Thank you for registering!\n We are glad to see you in our application. To complete, enter this code into registration form:\n%s", codeStr)
}
