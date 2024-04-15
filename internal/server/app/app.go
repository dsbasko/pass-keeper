package app

import (
	"context"
	"syscall"

	goGS "github.com/dsbasko/go-gs"

	"github.com/dsbasko/pass-keeper/internal/server/config"
	"github.com/dsbasko/pass-keeper/internal/server/endpoint/grpc"
	"github.com/dsbasko/pass-keeper/internal/server/provider/postgre"
	"github.com/dsbasko/pass-keeper/internal/server/service/auth"
	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
	"github.com/dsbasko/pass-keeper/pkg/logger"
)

type Options struct {
	Cfg    *config.Config
	Logger logger.Logger
}

func Run(opts Options) (err error) {
	defer errWrapper.PtrWithOP(&err, "app.Run")

	// Graceful shutdown
	gs, ctx, cancel := goGS.NewContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()

	// Подключение к БД
	psql := postgre.MustNew(ctx, postgre.Options{
		Logger: opts.Logger,
		Cfg:    opts.Cfg,
		GS:     gs,
	})

	// Сервис авторизации
	authService := auth.MustNew(auth.Options{
		Mutator: psql,
	})

	// Запуск gRPC сервера
	grpc.MustRun(ctx, grpc.Options{
		Cfg:         opts.Cfg,
		Logger:      opts.Logger,
		GS:          gs,
		AuthMutator: authService,
	})

	// Ожидание завершения работы GS
	gs.Wait()
	opts.Logger.Info("server is stopped")

	return nil
}

func MustRun(opts Options) {
	if err := Run(opts); err != nil {
		panic(err)
	}
}
