package postgre

import "errors"

var (
	ErrMissingContext = errors.New("missing context")
	ErrMissingLogger  = errors.New("missing logger")
	ErrMissingCfg     = errors.New("missing config")
	ErrMissingGS      = errors.New("missing graceful shutdowner")
	ErrEmailExists    = errors.New("email already exists")
)
