package httpserver

import (
	"net/http"
	"time"

	"github.com/AlexAnd012/BookFinder/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

// WriteHeader .чтобы знать какой статус отправили
func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func AccessLog(log logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Сохраняем старт времени
			start := time.Now()

			//Создаём наш statusWriter и передаём его дальше
			sw := &statusWriter{ResponseWriter: w}
			next.ServeHTTP(sw, r)
			log.Info("http_access",
				"req_id", middleware.GetReqID(r.Context()),
				"method", r.Method,
				"path", r.URL.RequestURI(),
				"status", sw.status,
				"bytes", sw.bytes,
				"dur_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
