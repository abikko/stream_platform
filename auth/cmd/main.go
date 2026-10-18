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
	pool, bootstrapErr := bootstrapDatabase()
	if bootstrapErr != nil {
		log.Fatalf("bootstrap: %v", bootstrapErr)
	}

	store := postgres.UserStore{Pool: pool}

	router := httpapi.NewRouter(
		identity.RegisterService{Store: store},
		identity.LoginService{Store: store},
	)

	timeoutDuration := 15 * time.Second

	server := http.Server{
		Addr:              ":8080",
		Handler:           router,
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
