package httpapi

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var errEmailAlreadyInUse = errors.New("email already in use")
var errPasswordTooLong = errors.New("password too long")
var errPasswordIsEmpty = errors.New("password is empty")
var errEmailIsEmpty = errors.New("email is empty")
var errUnableToHash = errors.New("unable to hash password")

type RegisterService struct {
	Store BaseStore
}

func (service RegisterService) InsertUser(ctx context.Context, req RegisterRequest) (string, error) {
	if req.Email == "" {
		return "", errEmailIsEmpty
	}
	if req.Password == "" {
		return "", errPasswordIsEmpty
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", errPasswordTooLong
		}
		return "", fmt.Errorf("Unable to hash password %w", err)
	}

	id, err := service.Store.CreateUser(ctx, passwordHash, req.Email)

	if err != nil {
		if errors.Is(err, errEmailAlreadyInUse) {
			return "", errEmailAlreadyInUse
		}

		return "", fmt.Errorf("register user: %w", err)
	}

	return id, nil

}
