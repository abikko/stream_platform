package main

import (
	"log"
	"net/http"

	"stream_platform/auth/internal/httpapi"
	"stream_platform/auth/internal/identity"
	"stream_platform/auth/internal/postgres"
)

func main() {
	pool, bootstrapErr := bootstrap()
	if bootstrapErr != nil {
		log.Fatalf("bootstrap: %v", bootstrapErr)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", httpapi.NewHealthHandler().ServeHTTP)

	mux.HandleFunc("/register", httpapi.NewRegisterHandler(identity.RegisterService{Store: &postgres.UserStore{Pool: pool}}).ServeHTTP)
	mux.HandleFunc("/login", httpapi.NewLoginHandler(identity.LoginService{Store: &postgres.UserStore{Pool: pool}}).ServeHTTP)

	log.Println("Сервер запущен: http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("listen and: %v", err)
	}
}
