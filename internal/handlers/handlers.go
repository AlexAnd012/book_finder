package handlers

import (
	"github.com/AlexAnd012/BookFinder/internal/logging"
	"github.com/AlexAnd012/BookFinder/internal/repo"
)

type Handlers struct {
	Books *BookHTTP
}

func New(pg *repo.Postgres, log logging.Logger) *Handlers {
	return &Handlers{
		Books: NewBookHTTP(repo.NewBookRepo(pg), log),
	}
}
