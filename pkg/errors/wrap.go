package errors

import (
	"fmt"
)

func ErrorPtrWithOP(err *error, op string) { //nolint:gocritic
	if *err != nil {
		*err = fmt.Errorf("%s -> %w", op, *err)
	}
}

func ErrorWithOP(err error, op string) error {
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func LogWithOP(err error, op string) string {
	if err == nil {
		return op
	}
	return fmt.Sprintf("%s: %v", op, err)
}
