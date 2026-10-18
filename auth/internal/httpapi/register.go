package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/abikko/stream_platform/auth/internal/identity"
)

type RegisterResponse struct {
	ID string `json:"id"`
}

type RegisterHandler struct {
	service identity.RegisterService
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegisterHandler(service identity.RegisterService) (h *RegisterHandler) {
	return &RegisterHandler{service: service}
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

	id, err := h.service.InsertUser(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, identity.ErrEmailAlreadyInUse):
			{
				http.Error(w, "email already in use", http.StatusConflict)
				return
			}
		case errors.Is(err, identity.ErrPasswordTooLong):
			{
				http.Error(w, "password too long", http.StatusBadRequest)
				return
			}
		case errors.Is(err, identity.ErrEmailIsEmpty):
			{
				http.Error(w, "email is empty", http.StatusBadRequest)
				return
			}
		case errors.Is(err, identity.ErrPasswordIsEmpty):
			{
				http.Error(w, "password is empty", http.StatusBadRequest)
				return
			}
		default:
			{
				http.Error(w, "unable to save", http.StatusInternalServerError)
				log.Print(err)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(RegisterResponse{ID: id}); err != nil {
		log.Printf("encode RegisterResponse: %v", err)
	}
}
