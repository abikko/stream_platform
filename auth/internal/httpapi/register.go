package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"stream_platform/auth/internal/identity"
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

func NewRegisterHandler(service identity.RegisterService) *RegisterHandler {
	return &RegisterHandler{service: service}
}

func (h *RegisterHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}

	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Unable to parse", http.StatusBadRequest)
		return
	}

	id, err := h.service.InsertUser(r.Context(), req.Email, req.Password)

	if err != nil {
		switch err {
		case identity.ErrEmailAlreadyInUse:
			{
				http.Error(w, "email already in use", http.StatusConflict)
				return
			}
		case identity.ErrPasswordTooLong:
			{
				http.Error(w, "password too long", http.StatusBadRequest)
				return
			}
		case identity.ErrEmailIsEmpty:
			{
				http.Error(w, "email is empty", http.StatusBadRequest)
				return
			}
		case identity.ErrPasswordIsEmpty:
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
