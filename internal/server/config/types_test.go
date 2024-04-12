package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Timeout(t *testing.T) {
	cfg := &config{
		Controller: controller{
			GRPC: controllerGRPC{
				TimeoutMs: 1000,
			},
		},
	}

	assert.Equal(t, time.Duration(1000)*time.Millisecond, cfg.Controller.GRPC.Timeout())
}
