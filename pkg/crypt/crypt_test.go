package crypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var (
		key      = []byte("someKey")
		wrongKey = []byte("wrongKey")
		data     = []byte("someData")
		enc      []byte
	)

	t.Run("Encrypt", func(t *testing.T) {
		encryptData, err := Encrypt(key, data)
		assert.NoError(t, err)
		assert.NotEmpty(t, encryptData)
		enc = encryptData
	})

	t.Run("Decrypt", func(t *testing.T) {
		decryptData, err := Decrypt(key, enc)
		assert.NoError(t, err)
		assert.Equal(t, data, decryptData)
	})

	t.Run("Decrypt aesGCM.Open Error", func(t *testing.T) {
		_, err := Decrypt(wrongKey, enc)
		fmt.Println(err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "aesGCM.Open")
	})
}
