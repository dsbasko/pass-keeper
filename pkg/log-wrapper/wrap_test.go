package logwrapper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithOP(t *testing.T) {
	t.Run("Should Return", func(t *testing.T) {
		err := fmt.Errorf("some error")
		errString := WithOP(err, "wrap")
		assert.Equal(t, "wrap: some error", errString)
	})

	t.Run("Should Not Return", func(t *testing.T) {
		var err error
		errString := WithOP(err, "wrap")
		assert.Equal(t, "wrap", errString)
	})
}
