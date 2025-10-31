package handlers

import (
	"github.com/AlexAnd012/BookFinder/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Health(log logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("health_check", "req_id", middleware.GetReqID(r.Context()))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}
}
