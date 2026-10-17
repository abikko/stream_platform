package identity

import (
	"context"
	"time"
)

type BaseStore interface {
	CreateUser(ctx context.Context, passwordHash []byte, email string) (id string, err error)

	FindUser(ctx context.Context, email string) (User, error)
}

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
