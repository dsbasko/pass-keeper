package app

import (
	"context"
	gogs "github.com/dsbasko/go-gs"
	"github.com/dsbasko/pass-keeper/internal/server/config"
	"github.com/dsbasko/pass-keeper/internal/server/endpoint/grpc"
	"github.com/dsbasko/pass-keeper/pkg/errors"
	"github.com/dsbasko/pass-keeper/pkg/logger"
	"syscall"
)

type Options struct {
	Cfg    *config.Config
	Logger logger.Logger
}

func Run(opts Options) (err error) {
	defer errors.ErrorPtrWithOP(&err, "app.Run")

	gs, ctx, cancel := gogs.NewContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()

	err = grpc.Run(ctx, grpc.Options{
		Cfg:    opts.Cfg,
		Logger: opts.Logger,
		GS:     gs,
	})
	if err != nil {
		err = errors.ErrorWithOP(err, "grpc.Run")
	}

	gs.Wait()
	opts.Logger.Info("server is stopped")

	return nil
}

func MustRun(opts Options) {
	if err := Run(opts); err != nil {
		panic(err)
	}
}
