package grpc

import (
	"context"
	goGS "github.com/dsbasko/go-gs"
	apiV1 "github.com/dsbasko/pass-keeper/api/v1"
	"github.com/dsbasko/pass-keeper/internal/server/config"
	grpc_mock "github.com/dsbasko/pass-keeper/internal/server/endpoint/grpc/mocks"
	"github.com/dsbasko/pass-keeper/internal/server/endpoint/grpc/servers"
	"github.com/dsbasko/pass-keeper/pkg/errors"
	"github.com/dsbasko/pass-keeper/pkg/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Options struct {
	Logger logger.Logger
	Cfg    *config.Config
	GS     goGS.GracefulShutdowner
}

func Run(ctx context.Context, opts Options) (err error) {
	defer errors.ErrorPtrWithOP(&err, "grpc.Run")

	switch {
	case opts.Cfg == nil:
		err = ErrMissingConfig
		return
	case opts.Logger == nil:
		err = ErrMissingLogger
		return
	case opts.GS == nil:
		err = ErrMissingGS
		return
	}

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: opts.Cfg.Endpoint.GRPC.Port,
	})
	if err != nil {
		err = errors.ErrorWithOP(err, "net.ListenTCP")
		return
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
		),
	)

	authServer, err := servers.NewAuthServer(servers.AuthOptions{
		Logger:  opts.Logger,
		Mutator: grpc_mock.NewMockAuthMutator(&gomock.Controller{}), // TODO поменять
	})
	apiV1.RegisterAuthServer(srv, authServer)

	if opts.Cfg.Env == "dev" {
		reflection.Register(srv)
	}

	go runServer(ctx, listen, srv)
	opts.Logger.InfoF("grpc server is running on port %d", opts.Cfg.Endpoint.GRPC.Port)

	opts.GS.Subscribe()
	go gsStop(ctx, opts.GS, srv)

	return nil
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
