package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/abikko/stream_platform/auth/internal/identity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	Pool *pgxpool.Pool
}

const (
	emailAlreadyInDB = "23505"
	emailNotFound    = "20000"
)

func (userStore UserStore) CreateUser(ctx context.Context, passwordHash []byte, email string) (id string, err error) {
	query := "INSERT INTO users(email, password_hash) VALUES($1,$2) RETURNING id"
	if err := userStore.Pool.QueryRow(ctx, query, email, passwordHash).Scan(&id); err != nil {
		if dbError, found := errors.AsType[*pgconn.PgError](err); found && dbError.Code == emailAlreadyInDB {
			return "", identity.ErrEmailAlreadyInUse
		}

		return "", fmt.Errorf("insert new user: %w", err)
	}

	return id, nil
}

func (userStore UserStore) FindUser(ctx context.Context, email string) (identity.User, error) {
	query := "SELECT id,email,password_hash,created_at FROM users WHERE email=$1"

	queryRow, err := userStore.Pool.Query(ctx, query, email)
	if err != nil {
		return identity.User{}, fmt.Errorf("load find user query: %w", err)
	}
	p, err := pgx.CollectOneRow(queryRow, pgx.RowToStructByName[UserRow])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return identity.User{}, identity.ErrUserDoesNotExist
	}

	return p.toUser(), nil
}
