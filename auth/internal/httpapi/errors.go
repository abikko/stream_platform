package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeEmailAlreadyInUse ErrorCode = "email_already_in_use"
	ErrorCodeValidation        ErrorCode = "validation_failed"
	ErrorInternalError         ErrorCode = "internal_error"
	ErrorPasswordTooLong       ErrorCode = "password_too_long"
)

type ErrorResponse struct {
	Code    ErrorCode            `json:"code"`
	Message string               `json:"message"`
	Fields  []ErrorFieldResponse `json:"fields,omitempty"`
}

type ErrorFieldResponse struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteInternalServerError(w http.ResponseWriter, logErr error) {
	ErrorResponse{Code: ErrorInternalError, Message: "Internal Server Error"}.WriteError(w, http.StatusInternalServerError)

	log.Print(logErr)
}

func (er ErrorResponse) WriteError(w http.ResponseWriter, statusCode int) {
	setHeaders(w)
	w.WriteHeader(statusCode)

	writeJSON(w, er)
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func writeJSON(w http.ResponseWriter, response any) {
	encodeErr := json.NewEncoder(w).Encode(response)

	if encodeErr != nil {
		log.Printf("unable to encode: %v", encodeErr)
	}
}
