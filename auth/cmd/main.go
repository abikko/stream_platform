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

	log.Println("Сервер запущен: http://localhost:8080")

	timeoutDuration := time.Duration(time.Duration.Seconds(15))

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  timeoutDuration,
		WriteTimeout: timeoutDuration,
		IdleTimeout:  -1,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("listen and: %v", err)
	}
}
