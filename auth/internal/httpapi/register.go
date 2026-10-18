package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
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
		if validationError, found := errors.AsType[identity.ValidationError](err); found {
			var fieldErrorResponse []ErrorFieldResponse
			for _, fieldError := range validationError.Fields {
				fieldErrorResponse = append(fieldErrorResponse, ErrorFieldResponse{
					Field:   fieldError.Field,
					Code:    string(fieldError.Code),
					Message: "error validate",
				})
			}
			ErrorResponse{
				Code:    ErrorCodeValidation,
				Message: "Validation error",
				Fields:  fieldErrorResponse,
			}.WriteError(w, http.StatusBadRequest)
			return
		}
		switch {
		case errors.Is(err, identity.ErrEmailAlreadyInUse):
			{
				ErrorResponse{Code: ErrorCodeEmailAlreadyInUse, Message: "email already in use"}.WriteError(w, http.StatusConflict)
				return
			}
		case errors.Is(err, identity.ErrPasswordTooLong):
			{
				ErrorResponse{Code: ErrorPasswordTooLong, Message: "password too long"}.WriteError(w, http.StatusBadRequest)
				return
			}
		case errors.Is(err, identity.ErrEmailIsEmpty):
			{
				ErrorResponse{Code: "invalid_data", Message: "invalid credentials"}.WriteError(w, http.StatusBadRequest)
				return
			}
		case errors.Is(err, identity.ErrPasswordIsEmpty):
			{
				ErrorResponse{Code: "invalid_data", Message: "invalid credentials"}.WriteError(w, http.StatusBadRequest)
				return
			}
		default:
			{
				WriteInternalServerError(w, fmt.Errorf("create user: %w", err))
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
