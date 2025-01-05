package unit

import (
	"context"
	"time"

	// "sso-like/internal/model"
	"sso-like/internal/service/auth"
	dtoService "sso-like/internal/service/dto"
	dtoStorage "sso-like/internal/storage/dto"
	"sso-like/internal/storage/postgres/user/mocks"

	// "sso-like/pkg/crypt"
	"sso-like/internal/model"
	"sso-like/pkg/crypt"
	loggerMock "sso-like/pkg/logger/mocks"

	// "github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type AuthSuite struct {
	suite.Suite
}

func (suite *AuthSuite) TestSignUpSuccess_01(t provider.T) {
	t.Title("[SignUp] correct email and correct password")
	t.Tags("auth", "signup", "failed")
	t.Parallel()

	t.WithNewStep("successfully signup new user", func(sCtx provider.StepCtx) {
		mockUserStorage := new(mocks.UserInterface)
		mockLogger := new(loggerMock.Interface)

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
		// userStorage.AssertExpectations(sCtx.T())
	})
}

// Empty email provided in request
func (suite *AuthSuite) TestSignUpFailed_01(t provider.T) {
	t.Title("[SignUp] empty email provided")
	t.Tags("auth", "signup", "validation")
	t.Parallel()

	t.WithNewStep("signup fails with empty email", func(sCtx provider.StepCtx) {
		mockUserStorage := new(mocks.UserInterface)
		mockLogger := new(loggerMock.Interface)

		authService := auth.NewAuthService(mockLogger, mockUserStorage)

		ctx := context.Background()
		req := &dtoService.SignUpRequest{
			Email:    "",
			Password: "password123",
		}

		mockLogger.On("Errorf", "email is empty").Return()

		mockUserStorage.On("Insert", ctx, mock.Anything).Times(0)

		secret, err := authService.SignUp(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "email is empty", err.Error())
		assert.Empty(t, secret)
	})
}

func (suite *AuthSuite) TestSignUpFailed_02(t provider.T) {
	t.Title("[SignUp] empty password provided")
	t.Tags("auth", "signup", "validation")
	t.Parallel()

	t.WithNewStep("signup fails with empty email", func(sCtx provider.StepCtx) {
		mockUserStorage := new(mocks.UserInterface)
		mockLogger := new(loggerMock.Interface)

		authService := auth.NewAuthService(mockLogger, mockUserStorage)

		ctx := context.Background()
		req := &dtoService.SignUpRequest{
			Email:    "test@test.com",
			Password: "",
		}

		mockLogger.On("Errorf", "password is empty").Return()

		mockUserStorage.On("Insert", ctx, mock.Anything).Times(0)

		secret, err := authService.SignUp(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "password is empty", err.Error())
		assert.Empty(t, secret)
	})
}

// Valid email, password and TOTP token returns a Paseto token
func (suite *AuthSuite) TestLoginSuccess_01(t provider.T) {
	t.Title("[Login] valid credentials and TOTP token")
	t.Tags("auth", "login", "success")
	t.Parallel()

	t.WithNewStep("successfully login user", func(sCtx provider.StepCtx) {
		mockUserStorage := new(mocks.UserInterface)
		mockLogger := new(loggerMock.Interface)

		authService := auth.NewAuthService(mockLogger, mockUserStorage)

		ctx := context.Background()
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		// encryptedSecret, _ := crypt.AesEncrypt(crypt.KEY, []byte("TOTP_SECRET"))
		token, _ := totp.GenerateCode(string("TOTP_SECRET"), time.Now())

		req := &dtoService.LogInRequest{
			Email:    "test@example.com",
			Password: "password123",
			Token:    token,
		}

		totpToken, _ := totp.Generate(totp.GenerateOpts{
			Issuer: "sso",
			AccountName: req.Email,
		})

		totpEncypt, _ := crypt.AesEncrypt(crypt.KEY, []byte(totpToken.Secret()))
		
		mockUser := &model.User{
			Email:      req.Email,
			Password:   hashedPass,
			TotpSecret: totpEncypt,
		}
		req.Token, _ = totp.GenerateCode(totpToken.Secret(), time.Now())

		mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
		mockUserStorage.On("Get", ctx, mock.Anything).Return(mockUser, nil)

		token, err := authService.LogIn(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func (suite *AuthSuite) TestLoginFailed_01(t provider.T) {
	t.Title("[Login] invalid credentials")
	t.Tags("auth", "login", "failed")
	t.Parallel()

	t.WithNewStep("successfully login user", func(sCtx provider.StepCtx) {
		mockUserStorage := new(mocks.UserInterface)
		mockLogger := new(loggerMock.Interface)

		authService := auth.NewAuthService(mockLogger, mockUserStorage)

		ctx := context.Background()
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		// encryptedSecret, _ := crypt.AesEncrypt(crypt.KEY, []byte("TOTP_SECRET"))
		token, _ := totp.GenerateCode(string("TOTP_SECRET"), time.Now())

		req := &dtoService.LogInRequest{
			Email:    "test@example.com",
			Password: "password123",
			Token:    token,
		}

		totpToken, _ := totp.Generate(totp.GenerateOpts{
			Issuer: "sso",
			AccountName: req.Email,
		})

		totpEncypt, _ := crypt.AesEncrypt(crypt.KEY, []byte(totpToken.Secret()))
		
		mockUser := &model.User{
			Email:      "test@test.com",
			Password:   hashedPass,
			TotpSecret: totpEncypt,
		}
		req.Token, _ = totp.GenerateCode(totpToken.Secret(), time.Now())

		mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
		mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
		mockUserStorage.On("Get", ctx, mock.Anything).Return(mockUser, nil)

		token, err := authService.LogIn(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func (suite *AuthSuite) TestLoginFailed_02(t provider.T) {
	t.Title("[Login] invalid credentials")
	t.Tags("auth", "login", "failed")
	t.Parallel()

	t.WithNewStep("successfully login user", func(sCtx provider.StepCtx) {
		mockUserStorage := new(mocks.UserInterface)
		mockLogger := new(loggerMock.Interface)

		authService := auth.NewAuthService(mockLogger, mockUserStorage)

		ctx := context.Background()
		// hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		// encryptedSecret, _ := crypt.AesEncrypt(crypt.KEY, []byte("TOTP_SECRET"))
		token, _ := totp.GenerateCode(string("TOTP_SECRET"), time.Now())

		req := &dtoService.LogInRequest{
			Email:    "test@example.com",
			Password: "password123",
			Token:    token,
		}

		totpToken, _ := totp.Generate(totp.GenerateOpts{
			Issuer: "sso",
			AccountName: req.Email,
		})

		totpEncypt, _ := crypt.AesEncrypt(crypt.KEY, []byte(totpToken.Secret()))
		
		mockUser := &model.User{
			Email:      "test@example.com",
			Password:   []byte("a"),
			TotpSecret: totpEncypt,
		}
		req.Token, _ = totp.GenerateCode(totpToken.Secret(), time.Now())

		mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
		mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
		mockUserStorage.On("Get", ctx, mock.Anything).Return(mockUser, nil)

		token, err := authService.LogIn(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func (suite *AuthSuite) TestLoginFailed_03(t provider.T) {
	t.Title("[Login] invalid token")
	t.Tags("auth", "login", "failed")
	t.Parallel()

	t.WithNewStep("successfully login user", func(sCtx provider.StepCtx) {
		mockUserStorage := new(mocks.UserInterface)
		mockLogger := new(loggerMock.Interface)

		authService := auth.NewAuthService(mockLogger, mockUserStorage)

		ctx := context.Background()
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		// encryptedSecret, _ := crypt.AesEncrypt(crypt.KEY, []byte("TOTP_SECRET"))
		token, _ := totp.GenerateCode(string("TOTP_SECRET"), time.Date(2025, 1, 1, 0,0,0,0, time.UTC))

		req := &dtoService.LogInRequest{
			Email:    "test@example.com",
			Password: "password123",
			Token:    token,
		}

		totpToken, _ := totp.Generate(totp.GenerateOpts{
			Issuer: "sso",
			AccountName: req.Email,
		})

		totpEncypt := []byte("tt")
		
		mockUser := &model.User{
			Email:      "test@example.com",
			Password:   hashedPass,
			TotpSecret: totpEncypt,
		}
		req.Token, _ = totp.GenerateCode(totpToken.Secret(), time.Now())

		mockLogger.On("Infof", mock.Anything, mock.Anything).Times(0)
		mockLogger.On("Errorf", mock.Anything, mock.Anything).Times(0)
		mockUserStorage.On("Get", ctx, mock.Anything).Return(mockUser, nil)

		token, err := authService.LogIn(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
