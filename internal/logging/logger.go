package logging

import (
	"log/slog"
	"os"
)

// Logger не привязываемся к конкретной реализации логгера
type Logger interface {
	With(kv ...any) Logger
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Error(msg string, kv ...any)
}

type SLogger struct{ l *slog.Logger }

// New Создаём JSON-хендлер, пишем в stdout. Level управляет минимальным уровнем логирования
func New(level slog.Leveler) *SLogger {
	return &SLogger{
		l: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})),
	}
}

// With Возвращаем новый логгер с постоянными полями.
func (s *SLogger) With(kv ...any) Logger {
	return &SLogger{l: s.l.With(kv...)}
}
func (s *SLogger) Debug(msg string, kv ...any) {
	s.l.Debug(msg, kv...)
}
func (s *SLogger) Info(msg string, kv ...any) {
	s.l.Info(msg, kv...)
}
func (s *SLogger) Error(msg string, kv ...any) {
	s.l.Error(msg, kv...)
}
