package logwrapper

import (
	"fmt"
)

func WithOP(err error, op string) string {
	if err == nil {
		return op
	}
	return fmt.Sprintf("%s: %v", op, err)
}
