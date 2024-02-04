package user_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/email-auth/app/pkg/e"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	GetUserByEmail(email string) (*UserDAO, error)
	NewUser(user *UserCreationData) error
	MarkEmailAsVerified(email string) error
	ResetEmailVerifyingCode(email string) error
	UpdateEmailVerifyingCode(email string, code int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetUserByEmail(email string) (user *UserDAO, err error) {
	defer func() { err = wrapDbErrIfNotNil(err, "error while getting user by email") }()

	row := u.db.QueryRow(
		"SELECT * FROM users WHERE email = $1",
		email,
	)

	err = row.Err()

	if err != nil {
		return nil, err
	}

	user = &UserDAO{}
	err = row.Scan(&user.Id, &user.Email, &user.IsEmailVerified, &user.CreatedAt, &user.EmailVerifyingCode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) NewUser(user *UserCreationData) (err error) {
	defer func() { err = wrapDbErrIfNotNil(err, "error while creating new user") }()

	_, err = u.db.Exec(
		"INSERT INTO users (email, email_verifying_code) VALUES ($1, $2)",
		user.Email,
		user.EmailVerifyingCode,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) MarkEmailAsVerified(email string) (err error) {
	defer func() { err = wrapDbErrIfNotNil(err, "error while marking email as verified") }()

	_, err = u.db.Exec(
		"UPDATE users SET is_email_verified=true WHERE email=$1",
		email,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) ResetEmailVerifyingCode(email string) (err error) {
	defer func() { err = wrapDbErrIfNotNil(err, "error while reseting email verifyng code") }()

	_, err = u.db.Exec(
		"UPDATE users SET email_verifying_code=-1 WHERE email=$1",
		email,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) UpdateEmailVerifyingCode(email string, code int) (err error) {
	defer func() { err = wrapDbErrIfNotNil(err, "error while updating email verifyng code") }()

	_, err = u.db.Exec(
		"UPDATE users SET email_verifying_code=$1 WHERE email=$2",
		code,
		email,
	)

	if err != nil {
		return err
	}

	return nil
}

func wrapDbErrIfNotNil(err error, description string) error {
	const usersTableName = "users"

	fullDescription := fmt.Sprintf("%s table: %s", usersTableName, description)

	return e.WrapIfNotNil(err, fullDescription)
}
