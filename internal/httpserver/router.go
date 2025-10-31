package httpserver

import (
	"github.com/AlexAnd012/BookFinder/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlexAnd012/BookFinder/internal/logging"
)

type BookHandlers interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
}

func NewRouter(log logging.Logger, books BookHandlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Use(AccessLog(log))

	r.Get("/health", handlers.Health(log)) //

	r.Route("/api", func(api chi.Router) {
		api.Get("/books", books.Search)
		api.Get("/books/{id}", books.Get)
		api.Post("/books", books.Create)
	})
	return r
}
