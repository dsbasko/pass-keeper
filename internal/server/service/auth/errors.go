package auth

import (
	"errors"
	"fmt"
)

var (
	ErrMissingContext = errors.New("missing context")
	ErrMissingMutator = errors.New("missing mutator")
)

var (
	ErrValidationEmail      = errors.New("invalid email")
	ErrValidationPassMinLen = fmt.Errorf("password must be at least %v characters long", ValidationPassMinLen)
	ErrValidationPassMaxLen = fmt.Errorf("password must be at most %v characters long", ValidationPassMaxLen)
)
