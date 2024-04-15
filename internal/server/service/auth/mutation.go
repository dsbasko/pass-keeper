package auth

import (
	"context"
	"strconv"

	goValidator "github.com/asaskevich/govalidator"

	"github.com/dsbasko/pass-keeper/internal/model"
	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
)

//go:generate ../../../../bin/mockgen -destination=./mocks/auth-mutator.go -package=auth_service_mock github.com/dsbasko/pass-keeper/internal/server/service/auth Mutator
type Mutator interface {
	CreateUser(ctx context.Context, email, passwordHash string) (model.User, error)
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (_ model.User, err error) {
	defer errWrapper.PtrWithOP(&err, "auth.Service.CreateUser")

	// Валидация аргументов
	switch {
	case ctx == nil:
		return model.User{}, ErrMissingContext
	case !goValidator.IsEmail(email):
		return model.User{}, ErrValidationEmail
	case !goValidator.MinStringLength(password, strconv.Itoa(ValidationPassMinLen)):
		return model.User{}, ErrValidationPassMinLen
	case !goValidator.MaxStringLength(password, strconv.Itoa(ValidationPassMaxLen)):
		return model.User{}, ErrValidationPassMaxLen
	}

	passwordHash, err := hashPassword([]byte(password))
	if err != nil {
		return model.User{}, err
	}

	createdUser, err := s.mutator.CreateUser(ctx, email, string(passwordHash))
	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}
