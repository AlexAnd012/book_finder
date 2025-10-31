package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexAnd012/BookFinder/internal/httpserver"
	"github.com/AlexAnd012/BookFinder/internal/logging"
)

type tlog struct{}

func (tlog) With(...any) logging.Logger { return tlog{} }
func (tlog) Debug(string, ...any)       {}
func (tlog) Info(string, ...any)        {}
func (tlog) Error(string, ...any)       {}

type fakeBooks struct{}

func (fakeBooks) Create(http.ResponseWriter, *http.Request) {}
func (fakeBooks) Get(http.ResponseWriter, *http.Request)    {}
func (fakeBooks) Search(http.ResponseWriter, *http.Request) {}

func TestRouter_Health(t *testing.T) {
	r := httpserver.NewRouter(tlog{}, fakeBooks{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
}
