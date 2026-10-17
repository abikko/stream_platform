package postgres

import (
	"stream_platform/auth/internal/identity"
	"time"
)

type UserRow struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func (userRow UserRow) toUser() identity.User {
	return identity.User{
		ID:           userRow.ID,
		Email:        userRow.Email,
		PasswordHash: userRow.PasswordHash,
		CreatedAt:    userRow.CreatedAt,
	}
}
