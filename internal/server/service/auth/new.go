package auth

import (
	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
)

const (
	ValidationPassMinLen = 8
	ValidationPassMaxLen = 32
)

type Service struct {
	mutator Mutator
}

type Options struct {
	Mutator Mutator
}

func New(opts Options) (_ *Service, err error) {
	defer errWrapper.PtrWithOP(&err, "auth.New")

	// Валидация аргументов
	switch { //nolint:gocritic
	case opts.Mutator == nil:
		return nil, ErrMissingMutator
	}

	return &Service{
		mutator: opts.Mutator,
	}, nil
}

func MustNew(opts Options) *Service {
	service, err := New(opts)
	if err != nil {
		panic(err)
	}
	return service
}
