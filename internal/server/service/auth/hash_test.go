package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HashPasswordAndCompare(t *testing.T) {
	var (
		pass1 = []byte("pass1")
		pass2 = []byte("pass2")
	)

	hash, err := hashPassword(pass1)
	assert.NoError(t, err)

	isCompare, err := compareHashAndPassword(hash, pass1)
	assert.NoError(t, err)
	assert.True(t, isCompare)

	isCompare, err = compareHashAndPassword(hash, pass2)
	assert.Error(t, err)
	assert.False(t, isCompare)
}
