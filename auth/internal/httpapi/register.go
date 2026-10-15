package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var errEmailAlreadyInUse = errors.New("email already in use")

const maxByteLength = 72
const emailAlreadyInDB = "23505"

type RegisterResponse struct {
	ID string `json:"id"`
}

type RegisterHandler struct {
	pool *pgxpool.Pool
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegisterHandler(pool *pgxpool.Pool) *RegisterHandler {
	return &RegisterHandler{pool: pool}
}

func (h *RegisterHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Unable to parse", http.StatusBadRequest)
		return
	}

	if len(req.Password) > maxByteLength {
		http.Error(w, "Password is too long", http.StatusBadRequest)
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Unable to hash", http.StatusInternalServerError)
		log.Printf("password hash: %v", err)
		return
	}

	id, err := insertUser(r.Context(), passwordHash, req.Email, h.pool)

	if err != nil {
		if errors.Is(err, errEmailAlreadyInUse) {
			http.Error(w, "email already in use", http.StatusConflict)
			return
		}
		http.Error(w, "unable to save", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(RegisterResponse{ID: id}); err != nil {
		log.Printf("encode RegisterResponse: %v", err)
	}
}

func insertUser(ctx context.Context, passwordHash []byte, email string, pool *pgxpool.Pool) (id string, insertError error) {
	query := "INSERT INTO users(email, password_hash) VALUES($1,$2) RETURNING id"
	if err := pool.QueryRow(ctx, query, email, passwordHash).Scan(&id); err != nil {
		if dbError, found := errors.AsType[*pgconn.PgError](err); found && dbError.Code == emailAlreadyInDB {
			return "", errEmailAlreadyInUse
		}

		return "", fmt.Errorf("insert new user: %w", err)
	}

	return id, nil
}
