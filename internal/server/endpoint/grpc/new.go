package grpc

import (
	"context"
	"net"

	goGS "github.com/dsbasko/go-gs"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/dsbasko/pass-keeper/api/v1"
	"github.com/dsbasko/pass-keeper/internal/server/config"
	"github.com/dsbasko/pass-keeper/internal/server/endpoint/grpc/servers"
	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
	"github.com/dsbasko/pass-keeper/pkg/logger"
)

type Options struct {
	Logger      logger.Logger
	Cfg         *config.Config
	GS          goGS.GracefulShutdowner
	AuthMutator servers.AuthMutator
}

func Run(ctx context.Context, opts Options) (err error) {
	defer errWrapper.PtrWithOP(&err, "grpc.Run")

	// Валидация аргументов
	switch {
	case ctx == nil:
		return ErrMissingContext
	case opts.Cfg == nil:
		return ErrMissingConfig
	case opts.Logger == nil:
		return ErrMissingLogger
	case opts.GS == nil:
		return ErrMissingGS
	case opts.AuthMutator == nil:
		return ErrMissingAuthMutator
	}

	// Резервация порта
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: opts.Cfg.Endpoint.GRPC.Port,
	})
	if err != nil {
		return errWrapper.WithOP(err, "net.ListenTCP")
	}

	// Создание сервера с интерцепторами
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
		),
	)

	// Регистрация сервера авторизации
	authServer, err := servers.NewAuthServer(servers.AuthOptions{
		Logger:  opts.Logger,
		Mutator: opts.AuthMutator,
	})
	apiV1.RegisterAuthServer(srv, authServer)

	// Подключение рефлексии
	if opts.Cfg.Env == "dev" {
		reflection.Register(srv)
	}

	// Запуск сервера
	go runServer(ctx, listen, srv)
	opts.Logger.InfoF("grpc server is running on port %d", opts.Cfg.Endpoint.GRPC.Port)

	// GS отключение
	opts.GS.Subscribe()
	go gsStop(ctx, opts.GS, srv)

	return nil
}

func MustRun(ctx context.Context, opts Options) {
	if err := Run(ctx, opts); err != nil {
		panic(err)
	}
}

func runServer(ctx context.Context, listen *net.TCPListener, srv *grpc.Server) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := srv.Serve(listen); err != nil {
		cancel()
	}
}

func gsStop(ctx context.Context, gs goGS.GracefulShutdowner, srv *grpc.Server) {
	defer gs.UnsubscribeFn(func() {
		srv.GracefulStop()
	})
	<-ctx.Done()
}
