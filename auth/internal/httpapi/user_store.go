package httpapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseStore interface {
	CreateUser(ctx context.Context, passwordHash []byte, email string) (id string, err error)
}

type UserStore struct {
	Pool *pgxpool.Pool
}

const emailAlreadyInDB = "23505"

func (userStore *UserStore) CreateUser(ctx context.Context, passwordHash []byte, email string) (id string, err error) {
	query := "INSERT INTO users(email, password_hash) VALUES($1,$2) RETURNING id"
	if err := userStore.Pool.QueryRow(ctx, query, email, passwordHash).Scan(&id); err != nil {
		if dbError, found := errors.AsType[*pgconn.PgError](err); found && dbError.Code == emailAlreadyInDB {
			return "", errEmailAlreadyInUse
		}

		return "", fmt.Errorf("insert new user: %w", err)
	}

	return id, nil
}
