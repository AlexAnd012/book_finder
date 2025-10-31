package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct{ Pool *pgxpool.Pool }

// NewPostgres Создаём пул по DSN
func NewPostgres(dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Close() { p.Pool.Close() }
