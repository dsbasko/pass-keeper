package servers

import (
	"context"
	apiV1 "github.com/dsbasko/pass-keeper/api/v1"
	"github.com/dsbasko/pass-keeper/internal/model"
	"github.com/dsbasko/pass-keeper/pkg/errors"
	"github.com/dsbasko/pass-keeper/pkg/logger"
)

type AuthMutator interface {
	CreateUser(ctx context.Context, email, passwordHash string) (model.User, error)
}

//go:generate ../../../../../bin/mockgen -destination=../mocks/auth-mutator.go -package=grpc_mock github.com/dsbasko/pass-keeper/internal/server/endpoint/grpc/servers AuthMutator

type AuthServer struct {
	apiV1.UnimplementedAuthServer
	log     logger.Logger
	mutator AuthMutator
}

type AuthOptions struct {
	Logger  logger.Logger
	Mutator AuthMutator
}

func NewAuthServer(opts AuthOptions) (server AuthServer, err error) {
	defer errors.ErrorPtrWithOP(&err, "auth.NewAuthServer")

	switch {
	case opts.Logger == nil:
		err = ErrMissingLogger
		return
	case opts.Mutator == nil:
		err = ErrMissingMutator
		return
	}

	return AuthServer{
		log:     opts.Logger,
		mutator: opts.Mutator,
	}, nil
}

func (s AuthServer) Login(ctx context.Context, dto *apiV1.LoginRequest) (*apiV1.LoginResponse, error) {
	s.log.Error("not implemented login method")
	panic("not implemented")
}
