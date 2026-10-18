package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Fields  []ErrorFieldResponse `json:"fields"`
}

type ErrorFieldResponse struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (er ErrorResponse) WriteError(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(er)
	if err != nil {
		log.Printf("unable to encode: %v", err)
		return
	}
}
