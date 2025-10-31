package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AlexAnd012/BookFinder/internal/data"
	"github.com/AlexAnd012/BookFinder/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type BookStore interface {
	Create(ctx context.Context, b data.Book) (int64, error)
	Get(ctx context.Context, id int64) (data.BookWithMeta, error)
	Search(ctx context.Context, q *string, genre *string, yearFrom, yearTo *int, limit, offset int32) ([]data.BookWithMeta, error)
}

type BookHTTP struct {
	repo BookStore // ← было: *repo.BookRepo
	log  logging.Logger
}

func NewBookHTTP(r BookStore, l logging.Logger) *BookHTTP { // ← было: *repo.BookRepo
	return &BookHTTP{repo: r, log: l}
}

func (h *BookHTTP) Create(w http.ResponseWriter, r *http.Request) {
	reqLog := h.log.With("req_id", middleware.GetReqID(r.Context()), "route", "POST /v1/books")

	var book data.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		reqLog.Error("bad_json", "err", err.Error())
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if book.Title == "" {
		reqLog.Info("validation_failed", "field", "title")
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	id, err := h.repo.Create(r.Context(), book)
	if err != nil {
		reqLog.Error("db_create_failed", "err", err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	reqLog.Info("book_created", "book_id", id, "title", book.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

// GET /v1/books/{id}
func (h *BookHTTP) Get(w http.ResponseWriter, r *http.Request) {
	reqLog := h.log.With("req_id", middleware.GetReqID(r.Context()), "route", "GET /v1/books/{id}")

	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		reqLog.Info("invalid_id", "raw", raw)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	out, err := h.repo.Get(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		reqLog.Info("book_not_found", "book_id", id)
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}
	if err != nil {
		reqLog.Error("db_get_failed", "err", err.Error(), "book_id", id)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// GET /v1/books?q=...&limit=&offset=
func (h *BookHTTP) Search(w http.ResponseWriter, r *http.Request) {
	reqLog := h.log.With("req_id", middleware.GetReqID(r.Context()), "route", "GET /v1/books")

	q := r.URL.Query().Get("q")
	var qPtr *string
	if q != "" {
		qPtr = &q
	}

	items, err := h.repo.Search(r.Context(), qPtr, nil, nil, nil, 20, 0)
	if err != nil {
		reqLog.Error("db_search_failed", "err", err.Error(), "q", q)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}
