package grpc

import "errors"

var (
	ErrMissingContext     = errors.New("missing context")
	ErrMissingConfig      = errors.New("missing config")
	ErrMissingLogger      = errors.New("missing logger")
	ErrMissingAuthMutator = errors.New("missing auth mutator")
	ErrMissingGS          = errors.New("missing graceful shutdowner")
)
