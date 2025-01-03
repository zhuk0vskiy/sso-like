package auth

import (
	"context"

	dtoService "sso-like/internal/service/dto"
	dtoStorage "sso-like/internal/storage/dto"
	appStorage "sso-like/internal/storage/sqlite/app"
	userStorage "sso-like/internal/storage/sqlite/user"
	"sso-like/pkg/jwt"
	"sso-like/pkg/logger"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log         *logger.Logger
	userStorage userStorage.UserInterface
	appStorage  appStorage.AppInterface

	tokenTTL time.Duration
}

func NewAuthService(log *logger.Logger, userStrg userStorage.UserInterface, appStrg appStorage.AppInterface, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:         log,
		userStorage: userStrg,
		appStorage:  appStrg,
		tokenTTL:    tokenTTL,
	}
}

func (a *AuthService) SignUp(ctx context.Context, request *dtoService.SignUpRequest) (int64, error) {


	passHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		a.log.Errorf("failed to generate password hash: %w", err)
		return 0, err
	}

	id, err := a.userStorage.Insert(ctx, &dtoStorage.InsertUserRequest{
		PassHash: passHash,
	})
	if err != nil {
		a.log.Errorf("failed to save user: %w", err)

		return 0, err
	}

	return id, nil
}

func (a *AuthService) LogIn(ctx context.Context, request *dtoService.LogInRequest) (string, error) {

	user, err := a.userStorage.Get(ctx, &dtoStorage.GetUserRequest{
		Email: request.Email,
	})
	if err != nil {
		a.log.Errorf("failed to get user: %w", err)

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(request.Password)); err != nil {
		a.log.Infof("invalid credentials: %w", err)

		return "", err
	}

	app, err := a.appStorage.Get(ctx, &dtoStorage.GetAppRequest{
		Id: request.AppId,
	})
	if err != nil {
		a.log.Errorf("failed to get app: %w", err)
		return "", err
	}

	a.log.Infof("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Errorf("failed to generate token", err)

		return "", err
	}

	return token, nil
}
