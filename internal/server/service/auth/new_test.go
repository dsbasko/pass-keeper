package auth

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	authServiceMock "github.com/dsbasko/pass-keeper/internal/server/service/auth/mocks"
)

type ServiceSuite struct {
	suite.Suite

	mutator *authServiceMock.MockMutator
	service *Service
	fake    SuiteFake
}

type SuiteFake struct {
	Email         string
	Password      string
	PasswordShort string
	PasswordLong  string
	Error         error
}

func (s *ServiceSuite) SetupSuite() {
	t := s.T()
	ctrl := gomock.NewController(t)

	s.mutator = authServiceMock.NewMockMutator(ctrl)
	s.service = MustNew(Options{
		Mutator: s.mutator,
	})
	s.fake = SuiteFake{
		Email:         gofakeit.Email(),
		Password:      gofakeit.Password(true, true, true, true, false, 16),
		PasswordShort: gofakeit.Password(true, true, true, true, false, 1),
		PasswordLong:  gofakeit.Password(true, true, true, true, false, 108),
		Error:         errors.New("some error"),
	}
}

func (s *ServiceSuite) Test_New() {
	t := s.T()

	t.Run("HappyPath", func(t *testing.T) {
		mockService := Service{
			mutator: s.mutator,
		}

		assert.NotNil(t, s.service)
		assert.Equal(t, mockService, *s.service)
	})

	t.Run("MissingMutator", func(t *testing.T) {
		_, err := New(Options{})
		assert.Error(t, err)
	})
}

func Test_Auth_Service(t *testing.T) {
	suite.Run(t, &ServiceSuite{Suite: suite.Suite{}})
}

func (s *ServiceSuite) Test_MustNew() {
	t := s.T()
	assert.Panics(t, func() { MustNew(Options{}) })
}
