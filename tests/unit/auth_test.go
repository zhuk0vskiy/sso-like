package unit

import (
	"context"
	"testing"

	"sso-like/internal/storage/postgres/user/mocks"
	"sso-like/internal/service/auth"
	dtoService "sso-like/internal/service/dto"
	loggerMocks "sso-like/pkg/logger/mocks"	
	dtoStorage "sso-like/internal/storage/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUpSuccess_01(t *testing.T) {
	mockUserStorage := new(mocks.UserInterface)
	mockLogger := new(loggerMocks.Interface)

	authService := auth.NewAuthService(mockLogger, mockUserStorage)

	ctx := context.Background()
	req := &dtoService.SignUpRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)

	mockUserStorage.On("Insert", ctx, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			r := args.Get(1).(*dtoStorage.InsertUserRequest)
			assert.Equal(t, req.Email, r.Email)
			assert.NotEmpty(t, r.Password)
			assert.NotEmpty(t, r.TotpSecret)
		})

	secret, err := authService.SignUp(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
	// assert.Equal(t, "sso", issuer)
}

func TestSignUpFailed_01(t *testing.T) {
	mockUserStorage := new(mocks.UserInterface)
	mockLogger := new(loggerMocks.Interface)

	authService := auth.NewAuthService(mockLogger, mockUserStorage)

	ctx := context.Background()
	req := &dtoService.SignUpRequest{
		Email:    "test@example.com",
		Password: "",
	}

	mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)

	mockUserStorage.On("Insert", ctx, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			r := args.Get(1).(*dtoStorage.InsertUserRequest)
			assert.Equal(t, req.Email, r.Email)
			assert.NotEmpty(t, r.Password)
			assert.NotEmpty(t, r.TotpSecret)
		})

	secret, err := authService.SignUp(ctx, req)

	assert.Error(t, err)
	assert.Empty(t, secret)
	// assert.Equal(t, "sso", issuer)
}

