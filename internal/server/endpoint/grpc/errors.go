package grpc

import "errors"

var (
	ErrMissingConfig = errors.New("missing config")
	ErrMissingLogger = errors.New("missing logger")
	ErrMissingGS     = errors.New("missing graceful shutdowner")
)
