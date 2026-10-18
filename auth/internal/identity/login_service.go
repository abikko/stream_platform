package identity

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	Store BaseStore
}

// Скоро будет возвращать JWT, так что пока оставил логику если пароль правильный то error будет nil
func (service LoginService) LoginUser(ctx context.Context, email string, password string) error {
	if email == "" {
		return ErrEmailIsEmpty
	}

	if password == "" {
		return ErrPasswordIsEmpty
	}

	userDto, err := service.Store.FindUser(ctx, email)

	cryptError := bcrypt.CompareHashAndPassword([]byte(userDto.PasswordHash), []byte(password))

	if err != nil {
		if errors.Is(err, ErrUserDoesNotExist) {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("find user: %w", err)
	}

	if cryptError != nil {
		return ErrInvalidCredentials
	}

	return nil
}
