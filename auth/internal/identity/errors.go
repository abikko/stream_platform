package identity

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDoesNotExist   = errors.New("user does not exist")

	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrPasswordTooLong   = errors.New("password too long")
	ErrPasswordIsEmpty   = errors.New("password is empty")
	ErrEmailIsEmpty      = errors.New("email is empty")
	ErrUnableToHash      = errors.New("unable to hash password")
)
