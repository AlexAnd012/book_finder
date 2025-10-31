package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexAnd012/BookFinder/internal/config"
	"github.com/AlexAnd012/BookFinder/internal/handlers"
	"github.com/AlexAnd012/BookFinder/internal/httpserver"
	"github.com/AlexAnd012/BookFinder/internal/logging"
	"github.com/AlexAnd012/BookFinder/internal/repo"
)

func main() {
	// 1) Конфиг
	cfg := config.Load()

	// 2) Логгер (JSON). Поставь LevelDebug в dev, если нужно подробнее.
	log := logging.New(slog.LevelInfo).With("service", "book-finder", "env", "dev")

	// 3) DSN: берем DB_DSN, иначе склеиваем из HOST/USER/PASS + cfg.Dbname
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		host := envOr("DB_HOST", "localhost") // если API в Docker – обычно "db"
		user := envOr("DB_USER", "postgres")
		pass := envOr("DB_PASSWORD", "12345")
		dsn = "postgres://" + user + ":" + pass + "@" + host + ":5432/" + cfg.Dbname + "?sslmode=disable"
	}

	// 4) Подключение к БД
	pg, err := repo.NewPostgres(dsn)
	if err != nil {
		log.Error("db_connect_failed", "err", err.Error())
		os.Exit(1)
	}
	defer pg.Close()

	// 5) Хендлеры домена
	bookRepo := repo.NewBookRepo(pg)
	books := handlers.NewBookHTTP(bookRepo, log)

	// 6) Роутер принимает только готовые хендлеры
	r := httpserver.NewRouter(log, books)

	// 7) HTTP-сервер
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8) Старт + graceful shutdown
	errCh := make(chan error, 1)
	go func() {
		log.Info("http_server_start", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// ждём сигнал или фатальную ошибку сервера
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Info("shutdown_signal", "signal", sig.String())
	case err := <-errCh:
		log.Error("http_server_error", "err", err.Error())
	}

	// плавная остановка
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("http_server_shutdown_error", "err", err.Error())
	} else {
		log.Info("http_server_stopped", "timeout_sec", 10)
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
