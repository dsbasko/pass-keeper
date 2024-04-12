package loggermock

import (
	"github.com/dsbasko/pass-keeper/pkg/logger"
)

type Log struct{}

func NewMock() (logger.Logger, error) { return &Log{}, nil }

func (l *Log) Debug(string)          {}
func (l *Log) DebugF(string, ...any) {}
func (l *Log) DebugW(string, ...any) {}

func (l *Log) Info(string)          {}
func (l *Log) InfoF(string, ...any) {}
func (l *Log) InfoW(string, ...any) {}

func (l *Log) Warn(string)          {}
func (l *Log) WarnF(string, ...any) {}
func (l *Log) WarnW(string, ...any) {}

func (l *Log) Error(string)          {}
func (l *Log) ErrorF(string, ...any) {}
func (l *Log) ErrorW(string, ...any) {}
