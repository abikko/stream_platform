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
	var fieldErrors []FieldError
	if email == "" {
		fieldErrors = append(fieldErrors, FieldError{Field: "email", Code: FieldCodeRequired})
	}
	if password == "" {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Code: FieldCodeRequired})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			fieldErrors = append(fieldErrors, FieldError{Field: "password", Code: FieldCodeTooLong})
			return "", ValidationError{Fields: fieldErrors}
		}
		return "", fmt.Errorf("unable to hash password %w", err)
	}

	if len(fieldErrors) > 0 {
		return "", ValidationError{Fields: fieldErrors}
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
