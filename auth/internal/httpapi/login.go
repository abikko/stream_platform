package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/abikko/stream_platform/auth/internal/identity"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginHandler struct {
	service identity.LoginService
}

func NewLoginHandler(service identity.LoginService) *LoginHandler {
	return &LoginHandler{service: service}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decodeError := decoder.Decode(&req)
	if decodeError != nil {
		http.Error(w, "unable to parse", http.StatusBadRequest)
		return
	}

	err := h.service.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, identity.ErrInvalidCredentials):
			{
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			}
		case errors.Is(err, identity.ErrEmailIsEmpty):
			{
				http.Error(w, "Provided email is empty", http.StatusBadRequest)
			}
		case errors.Is(err, identity.ErrPasswordIsEmpty):
			{
				http.Error(w, "Provided password is empty", http.StatusBadRequest)
			}
		default:
			{
				log.Printf("login: %v", err)
				http.Error(w, "unable to login", http.StatusInternalServerError)
			}
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
