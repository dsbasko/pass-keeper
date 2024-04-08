package main

import (
	"path"

	"github.com/dsbasko/pass-keeper/internal/client/config"
)

func main() {
	config.MustInit(path.Join("configs", "client.env"))
	cfg := config.Get()

	_ = cfg // TODO Don't forget to remove нахрен
}
