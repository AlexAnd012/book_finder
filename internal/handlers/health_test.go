package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexAnd012/BookFinder/internal/handlers"
	"github.com/AlexAnd012/BookFinder/internal/logging"
)

type nopLogger struct{}

func (nopLogger) With(...any) logging.Logger {
	return nopLogger{}
}
func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}

func TestHealth_OK(t *testing.T) {
	h := handlers.Health(nopLogger{}) // если у тебя фабрика; если объект — создай NewHealthHTTP(nopLogger{}).Handle
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want content-type application/json, got %q", ct)
	}
	if body := w.Body.String(); body != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}
