package httpapi

import (
	"net/http"

	"github.com/abikko/stream_platform/auth/internal/identity"
)

func NewRouter(registerService identity.RegisterService, loginService identity.LoginService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", NewHealthHandler().ServeHTTP)

	mux.HandleFunc("POST /register", NewRegisterHandler(registerService).ServeHTTP)
	mux.HandleFunc("POST /login", NewLoginHandler(loginService).ServeHTTP)

	return LogMiddleware(mux)
}
