package main

import (
	"context"
	"path"
	"syscall"

	goGS "github.com/dsbasko/go-gs"

	"github.com/dsbasko/pass-keeper/internal/client/app"
	"github.com/dsbasko/pass-keeper/internal/client/config"
)

func main() {
	_, ctx, cancel := goGS.NewContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	config.MustInit(path.Join("configs", "client.env"))
	cfg := config.Get()
	_ = cfg // TODO Don't forget to remove нахрен

	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}
