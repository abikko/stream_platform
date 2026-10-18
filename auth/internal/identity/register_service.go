package identity

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	Store BaseStore
}

func (service RegisterService) InsertUser(ctx context.Context, email string, password string) (string, error) {
	if email == "" {
		return "", ErrEmailIsEmpty
	}
	if password == "" {
		return "", ErrPasswordIsEmpty
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", ErrPasswordTooLong
		}
		return "", fmt.Errorf("unable to hash password %w", err)
	}

	id, err := service.Store.CreateUser(ctx, passwordHash, email)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyInUse) {
			return "", ErrEmailAlreadyInUse
		}

		return "", fmt.Errorf("register user: %w", err)
	}

	return id, nil
}
