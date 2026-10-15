package main

import (
	"log"
	"net/http"

	"stream_platform/auth/internal/httpapi"
)

func main() {
	pool, bootstrapErr := bootstrap()
	if bootstrapErr != nil {
		log.Fatalf("bootstrap: %v", bootstrapErr)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", httpapi.NewHealthHandler().ServeHTTP)

	mux.HandleFunc("/register", httpapi.NewRegisterHandler(pool).ServeHTTP)

	log.Println("Сервер запущен: http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
