package main

import (
	"path"

	"github.com/dsbasko/pass-keeper/internal/server/app"
	"github.com/dsbasko/pass-keeper/internal/server/config"
	"github.com/dsbasko/pass-keeper/pkg/logger"
)

func main() {
	config.MustInit(path.Join("configs", "server.env"))
	cfg := config.Get()

	log := logger.MustNew(cfg.Env, cfg.AppName)

	app.MustRun(app.Options{
		Cfg:    cfg,
		Logger: log,
	})
}
