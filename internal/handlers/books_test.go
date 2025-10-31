package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexAnd012/BookFinder/internal/data"
	"github.com/AlexAnd012/BookFinder/internal/handlers"
	"github.com/AlexAnd012/BookFinder/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type fakeBookRepo struct {
	CreateFn func(context.Context, data.Book) (int64, error)
	GetFn    func(context.Context, int64) (data.BookWithMeta, error)
	SearchFn func(context.Context, *string, *string, *int, *int, int32, int32) ([]data.BookWithMeta, error)
}

func (f *fakeBookRepo) Create(ctx context.Context, b data.Book) (int64, error) {
	if f.CreateFn == nil {
		return 0, nil
	}
	return f.CreateFn(ctx, b)
}
func (f *fakeBookRepo) Get(ctx context.Context, id int64) (data.BookWithMeta, error) {
	if f.GetFn == nil {
		return data.BookWithMeta{}, nil
	}
	return f.GetFn(ctx, id)
}
func (f *fakeBookRepo) Search(ctx context.Context, q *string, g *string, yf, yt *int, limit, offset int32) ([]data.BookWithMeta, error) {
	if f.SearchFn == nil {
		return nil, nil
	}
	return f.SearchFn(ctx, q, g, yf, yt, limit, offset)
}

// простейший JSON-логгер
type testLogger struct{}

func (testLogger) With(...any) logging.Logger { return testLogger{} }
func (testLogger) Debug(string, ...any)       {}
func (testLogger) Info(string, ...any)        {}
func (testLogger) Error(string, ...any)       {}

func TestBooks_Create_201(t *testing.T) {
	repo := &fakeBookRepo{
		CreateFn: func(_ context.Context, b data.Book) (int64, error) {
			if b.Title != "New Book" {
				t.Fatalf("unexpected title: %s", b.Title)
			}
			return 42, nil
		},
	}
	h := handlers.NewBookHTTP(repo, testLogger{})

	body := `{"title":"New Book"}`
	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if _, ok := resp["id"]; !ok {
		t.Fatalf("response must contain id")
	}
}

func TestBooks_Create_400_BadJSON(t *testing.T) {
	h := handlers.NewBookHTTP(&fakeBookRepo{}, testLogger{})
	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBufferString("{bad}"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestBooks_Get_200(t *testing.T) {
	repo := &fakeBookRepo{
		GetFn: func(_ context.Context, id int64) (data.BookWithMeta, error) {
			return data.BookWithMeta{
				Book:    data.Book{ID: id, Title: "X"},
				Authors: []string{"A"},
				Genres:  []string{"G"},
			}, nil
		},
	}
	h := handlers.NewBookHTTP(repo, testLogger{})

	r := chi.NewRouter()
	r.Get("/api/books/{id}", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/books/7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
}

func TestBooks_Get_404(t *testing.T) {
	repo := &fakeBookRepo{
		GetFn: func(_ context.Context, _ int64) (data.BookWithMeta, error) {
			return data.BookWithMeta{}, pgx.ErrNoRows
		},
	}
	h := handlers.NewBookHTTP(repo, testLogger{})

	r := chi.NewRouter()
	r.Get("/api/books/{id}", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/books/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestBooks_Get_400_InvalidID(t *testing.T) {
	h := handlers.NewBookHTTP(&fakeBookRepo{}, testLogger{})

	r := chi.NewRouter()
	r.Get("/api/books/{id}", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/books/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestBooks_Search_200(t *testing.T) {
	repo := &fakeBookRepo{
		SearchFn: func(_ context.Context, q *string, _ *string, _ *int, _ *int, limit, offset int32) ([]data.BookWithMeta, error) {
			if q == nil || *q != "PS" {
				return nil, errors.New("q mismatch")
			}
			return []data.BookWithMeta{{Book: data.Book{ID: 1, Title: "PS Book"}}}, nil
		},
	}
	h := handlers.NewBookHTTP(repo, testLogger{})

	req := httptest.NewRequest(http.MethodGet, "/api/books?q=PS&limit=5&offset=0", nil)
	w := httptest.NewRecorder()

	// Можно вызывать напрямую:
	h.Search(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
}
