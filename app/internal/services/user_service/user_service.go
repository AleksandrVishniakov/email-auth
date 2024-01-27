package user_service

import (
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/repositories/user_repo"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/services/email_service"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/pkg/str"
	"log"
)

var (
	ErrEmailValidation = errors.New("email validation failed")
	ErrEmailIsExists   = errors.New("email is already exists")
)

type UserService interface {
	GetUserByEmail(email string) (*UserDTO, error)
	NewUser(email string, hostname string) error
	VerifyEmail(email, verifyingHash string) error
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

func (u *userService) NewUser(email string, hostname string) error {
	user, err := u.userRepository.GetUserByEmail(email)
	log.Println("email:", email, "     err:", err)

	if err == nil && user != nil && user.IsEmailVerified {
		return ErrEmailIsExists
	}

	if errors.Is(err, user_repo.ErrUserNotFound) {
		log.Println("not found")
		err = nil
	}

	if err != nil {
		log.Println("fail")
		return err
	}

	if user != nil && !user.IsEmailVerified {
		link := createVerifyingLink(hostname, email, user.EmailVerifyingHash)

		err = u.emailService.Write(&email_service.EmailContent{
			To:      email,
			Subject: "Completion of registration",
			Body:    createEmailMessage(link),
		})

		return err
	}

	hash := str.Generate(64)
	userData := &user_repo.UserCreationData{
		Email:              email,
		EmailVerifyingHash: hash,
	}

	err = u.userRepository.NewUser(userData)
	if err != nil {
		return err
	}

	link := createVerifyingLink(hostname, email, hash)
	err = u.emailService.Write(&email_service.EmailContent{
		To:      email,
		Subject: "Completion of registration",
		Body:    createEmailMessage(link),
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *userService) VerifyEmail(email, verifyingHash string) error {
	user, err := u.userRepository.GetUserByEmail(email)

	if user.IsEmailVerified {
		return nil
	}

	log.Println(*user)
	log.Println(verifyingHash)
	if err != nil {
		return err
	}

	if user.EmailVerifyingHash != verifyingHash {
		return ErrEmailValidation
	}

	return u.userRepository.MarkEmailAsVerified(email)
}

func createEmailMessage(emailVerifyingLink string) string {
	return fmt.Sprintf("Thank you for registering!\n We are glad to see you in our application. To complete the registration, follow the link:\n%s", emailVerifyingLink)
}

func createVerifyingLink(hostname, email, hash string) string {
	return fmt.Sprintf("http://%s/verify/%s?h=%s", hostname, email, hash)
}
