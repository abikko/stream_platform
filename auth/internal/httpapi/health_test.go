package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
	}{
		{"Check GET method", http.MethodGet, 200},
		{"Check POST method", http.MethodPost, 405},
		{"Check DELETE method", http.MethodDelete, 405},
		{"Check PATCH method", http.MethodPatch, 405},
		{"Check PUT method", http.MethodPut, 405},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			NewHealthHandler().ServeHTTP(
				recorder, httptest.NewRequest(test.method, "/health", nil),
			)

			if recorder.Code != test.expectedCode {
				t.Errorf("expected %v got %v", test.expectedCode, recorder.Code)
			}
		})
	}
}
