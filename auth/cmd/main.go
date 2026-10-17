package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abikko/stream_platform/auth/internal/httpapi"
	"github.com/abikko/stream_platform/auth/internal/identity"
	"github.com/abikko/stream_platform/auth/internal/postgres"
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

	timeoutDuration := 15 * time.Second

	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 8 * time.Second,
		ReadTimeout:       timeoutDuration,
		WriteTimeout:      timeoutDuration,
		IdleTimeout:       time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
	log.Println("Сервер запущен: http://localhost:8080")
}
