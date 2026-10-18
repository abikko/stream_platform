package httpapi

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"time"
)

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := rand.Text()
		w.Header().Set("X-Request-ID", requestID)

		recorder := statusRecorder{ResponseWriter: w}
		start := time.Now()
		h.ServeHTTP(&recorder, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "X-Request-ID", requestID, "status", recorder.status, "duration", time.Since(start))
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
