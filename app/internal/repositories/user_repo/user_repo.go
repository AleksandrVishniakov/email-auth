package user_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/pkg/e"
	"log"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	GetUserByEmail(email string) (*UserDAO, error)
	NewUser(user *UserCreationData) error
	MarkEmailAsVerified(email string) error
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
	log.Println("db:", err, errors.Is(err, sql.ErrNoRows))

	if err != nil {
		return nil, err
	}

	user = &UserDAO{}
	err = row.Scan(&user.Id, &user.Email, &user.IsEmailVerified, &user.EmailVerifyingHash, &user.CreatedAt)
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
		"INSERT INTO users (email, email_verifying_hash) VALUES ($1, $2)",
		user.Email,
		user.EmailVerifyingHash,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) MarkEmailAsVerified(email string) (err error) {
	defer func() { err = wrapDbErrIfNotNil(err, "error while marking email as verified") }()

	_, err = u.db.Exec(
		"UPDATE users SET is_email_verified=true, email_verifying_hash='' WHERE email=$1",
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
