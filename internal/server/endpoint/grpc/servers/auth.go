package servers

import (
	"context"
	"errors"

	"github.com/dsbasko/pass-keeper/internal/server/provider/postgre"
	"github.com/dsbasko/pass-keeper/internal/server/service/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiV1 "github.com/dsbasko/pass-keeper/api/v1"
	"github.com/dsbasko/pass-keeper/internal/model"
	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
	logWrapper "github.com/dsbasko/pass-keeper/pkg/log-wrapper"
	"github.com/dsbasko/pass-keeper/pkg/logger"
)

//go:generate ../../../../../bin/mockgen -destination=../mocks/auth-mutator.go -package=grpc_mock github.com/dsbasko/pass-keeper/internal/server/endpoint/grpc/servers AuthMutator
type AuthMutator interface {
	CreateUser(ctx context.Context, email, password string) (model.User, error)
}

type AuthServer struct {
	apiV1.UnimplementedAuthServer
	log     logger.Logger
	mutator AuthMutator
}

type AuthOptions struct {
	Logger  logger.Logger
	Mutator AuthMutator
}

func NewAuthServer(opts AuthOptions) (_ AuthServer, err error) {
	defer errWrapper.PtrWithOP(&err, "auth.NewAuthServer")

	// Валидация аргументов
	switch {
	case opts.Logger == nil:
		return AuthServer{}, ErrMissingLogger
	case opts.Mutator == nil:
		return AuthServer{}, ErrMissingMutator
	}

	return AuthServer{
		log:     opts.Logger,
		mutator: opts.Mutator,
	}, nil
}

func (s AuthServer) Login(_ context.Context, _ *apiV1.LoginRequest) (_ *apiV1.LoginResponse, err error) {
	s.log.ErrorF(logWrapper.WithOP(err, "not implemented login method"))
	return nil, status.Errorf(codes.Unimplemented, "method not implemented")
}

func (s AuthServer) Register(ctx context.Context, dto *apiV1.RegisterRequest) (_ *apiV1.RegisterResponse, err error) {
	createdUser, err := s.mutator.CreateUser(ctx, dto.Email, dto.Password)
	if err != nil {
		s.log.Error(err.Error())

		// Валидация
		switch {
		case errors.Is(err, postgre.ErrEmailExists):
			return nil, status.Errorf(codes.AlreadyExists, postgre.ErrEmailExists.Error())
		case errors.Is(err, auth.ErrValidationEmail):
			return nil, status.Errorf(codes.InvalidArgument, auth.ErrValidationEmail.Error())
		case errors.Is(err, auth.ErrValidationPassMinLen):
			return nil, status.Errorf(codes.InvalidArgument, auth.ErrValidationPassMinLen.Error())
		case errors.Is(err, auth.ErrValidationPassMaxLen):
			return nil, status.Errorf(codes.InvalidArgument, auth.ErrValidationPassMaxLen.Error())
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %v", dto.Email)
	}

	return &apiV1.RegisterResponse{
		UserId: createdUser.ID,
	}, nil
}
