package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abikko/stream_platform/auth/internal/identity"
)

type MockUserStore struct{}

func (store MockUserStore) FindUser(_ context.Context, _ string) (identity.User, error) {
	return identity.User{}, nil
}

func (store MockUserStore) CreateUser(_ context.Context, passwordHash []byte, email string) (string, error) {
	if email == "test123.kz" {
		return "", identity.ErrEmailAlreadyInUse
	}

	if email == "" {
		return "", identity.ErrEmailIsEmpty
	}

	if string(passwordHash) == "" {
		return "", identity.ErrPasswordIsEmpty
	}

	if len(passwordHash) > 72 {
		return "", identity.ErrPasswordTooLong
	}

	return "123", nil
}

func Test_RegisterEmail(t *testing.T) {
	tests := []struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		statusCode int
		name       string
	}{
		{
			Email: "test123.kz", Password: "passik", statusCode: http.StatusConflict, name: "Email already in use",
		},
		{
			Email: "", Password: "passik", statusCode: http.StatusBadRequest, name: "Email is empty",
		},
		{
			Email: "", Password: "", statusCode: http.StatusBadRequest, name: "Email and password is Empty",
		},
		{
			Email: "emailisnotempty@gmail.kz", Password: strings.Repeat("1231233331", 10), statusCode: http.StatusBadRequest, name: "Password is too long",
		},
		{
			Email: "toktok@gmail.kz", Password: "", statusCode: http.StatusBadRequest, name: "password is empty",
		},
		{
			Email: "toktok@gmail.kz", Password: "passik", statusCode: http.StatusCreated, name: "User created successfully",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			testJSON, err := json.Marshal(test)
			if err != nil {
				t.Errorf("Error to deserialize: %v", err)
			}
			request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/register", bytes.NewReader(testJSON))
			NewRegisterHandler(identity.RegisterService{
				Store: MockUserStore{},
			}).ServeHTTP(recorder, request)

			if recorder.Code != test.statusCode {
				t.Errorf("Want statusCode: %v got: %v", test.statusCode, recorder.Code)
			}

			var responseJSON struct {
				ID string `json:"id"`
			}
			if err := json.NewDecoder(recorder.Body).Decode(&responseJSON); err != nil {
				http.Error(recorder, "unable to parse", http.StatusInternalServerError)
			}

			if responseJSON.ID != "123" && test.statusCode == http.StatusCreated {
				t.Errorf("expected id 123, got %v", responseJSON.ID)
			}
		})
	}
}
