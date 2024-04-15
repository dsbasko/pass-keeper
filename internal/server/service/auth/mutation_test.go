package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dsbasko/pass-keeper/internal/model"
)

func (s *ServiceSuite) Test_CreateUser() {
	t := s.T()

	tests := []struct {
		name     string
		initCfg  func()
		ctx      context.Context
		email    string
		password string
		wantResp model.User
		wantErr  error
	}{
		{
			name: "HappyPath",
			initCfg: func() {
				s.mutator.EXPECT().
					CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(model.User{Email: s.fake.Email}, nil)
			},
			ctx:      context.Background(),
			email:    s.fake.Email,
			password: s.fake.Password,
			wantResp: model.User{
				Email: s.fake.Email,
			},
			wantErr: nil,
		},
		{
			name:     "Missing Context",
			initCfg:  func() {},
			ctx:      nil,
			email:    "",
			password: "",
			wantResp: model.User{},
			wantErr:  ErrMissingContext,
		},
		{
			name:     "Invalid Email",
			initCfg:  func() {},
			ctx:      context.Background(),
			email:    "42",
			password: "",
			wantResp: model.User{},
			wantErr:  ErrValidationEmail,
		},
		{
			name:     "Short Password",
			initCfg:  func() {},
			ctx:      context.Background(),
			email:    s.fake.Email,
			password: s.fake.PasswordShort,
			wantResp: model.User{},
			wantErr:  ErrValidationPassMinLen,
		},
		{
			name:     "Long Password",
			initCfg:  func() {},
			ctx:      context.Background(),
			email:    s.fake.Email,
			password: s.fake.PasswordLong,
			wantResp: model.User{},
			wantErr:  ErrValidationPassMaxLen,
		},
		{
			name: "Mutator Error",
			initCfg: func() {
				s.mutator.EXPECT().
					CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(model.User{}, s.fake.Error)
			},
			ctx:      context.Background(),
			email:    s.fake.Email,
			password: s.fake.Password,
			wantResp: model.User{},
			wantErr:  s.fake.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initCfg()
			resp, err := s.service.CreateUser(tt.ctx, tt.email, tt.password)
			if tt.wantErr != nil || err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
			assert.Equal(t, tt.wantResp, resp)
		})
	}
}
