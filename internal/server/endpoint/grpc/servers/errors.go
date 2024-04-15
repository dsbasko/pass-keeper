package servers

import "errors"

var (
	ErrMissingLogger  = errors.New("missing logger")
	ErrMissingMutator = errors.New("missing mutator")
)
