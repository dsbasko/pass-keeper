package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	err := Init(path.Join("stubs", "stub.env"))
	assert.Nil(t, err)
}

func Test_MustInit(t *testing.T) {
	assert.NotPanics(t, func() {
		MustInit(path.Join("stubs", "stub.env"))
	})
}

func Test_Get(t *testing.T) {
	MustInit(path.Join("stubs", "stub.env"))
	cfg := Get()
	assert.Equal(t, "test__client", cfg.AppName)
}
