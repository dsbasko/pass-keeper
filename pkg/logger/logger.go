package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

var (
	ErrEnv = errors.New("invalid run mode")
)

type Logger interface {
	Debug(format string)
	DebugF(format string, args ...any)
	DebugW(format string, args ...any)

	Info(format string)
	InfoF(format string, args ...any)
	InfoW(format string, args ...any)

	Warn(format string)
	WarnF(format string, args ...any)
	WarnW(format string, args ...any)

	Error(format string)
	ErrorF(format string, args ...any)
	ErrorW(format string, args ...any)
}

type Log struct {
	log *slog.Logger
}

func New(env, serviceName string) (Logger, error) {
	switch env {
	case "dev":
		return &Log{
			log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})),
		}, nil
	case "prod":
		return &Log{
			log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})).With("service", serviceName),
		}, nil
	default:
		return nil, ErrEnv
	}
}

func Must(log Logger, err error) Logger {
	if err != nil {
		panic(err)
	}
	return log
}

func (l *Log) Debug(format string) {
	l.log.Debug(format)
}

func (l *Log) DebugF(format string, args ...any) {
	l.log.Debug(fmt.Sprintf(format, args...))
}

func (l *Log) DebugW(format string, args ...any) {
	l.log.Debug(format, args...)
}

func (l *Log) Info(format string) {
	l.log.Info(format)
}

func (l *Log) InfoF(format string, args ...any) {
	l.log.Info(fmt.Sprintf(format, args...))
}

func (l *Log) InfoW(format string, args ...any) {
	l.log.Info(format, args...)
}

func (l *Log) Warn(format string) {
	l.log.Warn(format)
}

func (l *Log) WarnF(format string, args ...any) {
	l.log.Warn(fmt.Sprintf(format, args...))
}

func (l *Log) WarnW(format string, args ...any) {
	l.log.Warn(format, args...)
}

func (l *Log) Error(format string) {
	l.log.Error(format)
}

func (l *Log) ErrorF(format string, args ...any) {
	l.log.Error(fmt.Sprintf(format, args...))
}

func (l *Log) ErrorW(format string, args ...any) {
	l.log.Error(format, args...)
}
