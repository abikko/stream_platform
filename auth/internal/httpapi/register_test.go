package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockUserStore struct {
}

func (store MockUserStore) CreateUser(ctx context.Context, passwordHash []byte, email string) (string, error) {
	if email == "test123.kz" {
		return "", errEmailAlreadyInUse
	}

	if email == "" {
		return "", errEmailIsEmpty
	}

	if string(passwordHash) == "" {
		return "", errPasswordIsEmpty
	}

	if len(passwordHash) > 72 {
		return "", errPasswordTooLong
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
			Email: "test123.kz", Password: "passik", statusCode: http.StatusConflict, name: "Email already in use"},
		{
			Email: "", Password: "passik", statusCode: http.StatusBadRequest, name: "Email is empty",
		},
		{
			Email: "", Password: "", statusCode: http.StatusBadRequest, name: "Email and password is Empty",
		}, {
			Email: "emailisnotempty@gmail.kz", Password: strings.Repeat("1231233331", 10), statusCode: http.StatusBadRequest, name: "Password is too long",
		}, {
			Email: "toktok@gmail.kz", Password: "", statusCode: http.StatusBadRequest, name: "password is empty",
		}, {
			Email: "toktok@gmail.kz", Password: "passik", statusCode: http.StatusCreated, name: "User created successfully",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			testJson, err := json.Marshal(test)
			if err != nil {
				t.Errorf("Error to deserialize: %v", err)
			}
			request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(testJson))
			NewRegisterHandler(RegisterService{
				Store: MockUserStore{},
			}).ServeHTTP(recorder, request)

			if recorder.Code != test.statusCode {
				t.Errorf("%v Want statusCode: %v got: %v", test.name, test.statusCode, recorder.Code)
			}
		})
	}
}
