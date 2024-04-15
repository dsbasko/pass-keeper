package errwrapper

import (
	"fmt"
)

func PtrWithOP(err *error, op string) { //nolint:gocritic
	if *err != nil {
		*err = fmt.Errorf("%s -> %w", op, *err)
	}
}

func WithOP(err error, op string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
