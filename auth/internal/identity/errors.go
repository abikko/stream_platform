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

type FieldCode string

const (
	FieldCodeRequired      FieldCode = "required"
	FieldCodeTooLong       FieldCode = "too_long"
	FieldCodeInvalidFormat FieldCode = "invalid_format"
)

type ValidationError struct {
	Fields []FieldError
}

type FieldError struct {
	Field string
	Code  FieldCode
}

func (v ValidationError) Error() string {
	return "validation_error"
}
