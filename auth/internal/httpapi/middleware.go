package httpapi

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"time"
)

const xRequestIDKey = "request_id"

func CORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-Id")
		if isPreflight := r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != ""; isPreflight {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Max-Age", "240")

			w.WriteHeader(http.StatusNoContent)

			return
		}

		h.ServeHTTP(w, r)
	})
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := rand.Text()
		w.Header().Set("X-Request-Id", requestID)

		recorder := statusRecorder{ResponseWriter: w}
		start := time.Now()
		h.ServeHTTP(&recorder, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, xRequestIDKey, requestID, "status", recorder.status, "duration", time.Since(start))
	})
}

type statusRecorder struct {
	http.ResponseWriter // встраивание: все методы w доступны «как свои»
	status              int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
