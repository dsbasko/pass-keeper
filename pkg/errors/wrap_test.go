package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrorPtrWithOP(t *testing.T) {
	t.Run("Should Return", func(t *testing.T) {
		err := fmt.Errorf("some error")
		ErrorPtrWithOP(&err, "wrap")
		assert.Equal(t, "wrap -> some error", err.Error())
	})

	t.Run("Should Not Return", func(t *testing.T) {
		var err error
		ErrorPtrWithOP(&err, "wrap")
		assert.Nil(t, err)
	})
}

func Test_ErrorWithOP(t *testing.T) {
	t.Run("Should Return", func(t *testing.T) {
		err := ErrorWithOP(fmt.Errorf("some error"), "wrap")
		assert.Equal(t, "wrap: some error", err.Error())
	})

	t.Run("Should Not Return", func(t *testing.T) {
		err := ErrorWithOP(nil, "wrap")
		assert.Nil(t, err)
	})
}

func Test_LogWithOP(t *testing.T) {
	t.Run("Should Return", func(t *testing.T) {
		err := fmt.Errorf("some error")
		errString := LogWithOP(err, "wrap")
		assert.Equal(t, "wrap: some error", errString)
	})

	t.Run("Should Not Return", func(t *testing.T) {
		var err error
		errString := LogWithOP(err, "wrap")
		assert.Equal(t, "wrap", errString)
	})
}
