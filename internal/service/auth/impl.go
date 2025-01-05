package auth

import (
	"context"
	"fmt"

	dtoService "sso-like/internal/service/dto"
	dtoStorage "sso-like/internal/storage/dto"

	userStorage "sso-like/internal/storage/postgres/user"
	"sso-like/pkg/crypt"
	"sso-like/pkg/token/paseto"
	"sso-like/pkg/logger"
	"time"

	"github.com/pquerna/otp/totp"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	logger         logger.Interface
	userStorage userStorage.UserInterface

}

func NewAuthService(logger logger.Interface, userStrg userStorage.UserInterface) *AuthService {
	return &AuthService{
		logger:         logger,
		userStorage: userStrg,
	}
}

func (a *AuthService) SignUp(ctx context.Context, request *dtoService.SignUpRequest) (string, error) {

	if request.Email == "" {
		a.logger.Errorf("email is empty")
		return "", fmt.Errorf("email is empty")
	}

	if request.Password == "" {
		a.logger.Errorf("password is empty")
		return "", fmt.Errorf("password is empty")
	}

	totp, err := totp.Generate(totp.GenerateOpts{
		Issuer: "sso",
		AccountName: request.Email,
	})
	if err != nil {
		a.logger.Errorf("failed to generate topt: %w", err)
		return "", err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		a.logger.Errorf("failed to generate password hash: %w", err)
		return "", err
	}

	totpSecretEncypt, err := crypt.AesEncrypt(crypt.KEY, []byte(totp.Secret()))
	if err != nil {
		a.logger.Errorf("failed to encrypt totp secret: %w", err)
		return "", err
	}

	err = a.userStorage.Insert(ctx, &dtoStorage.InsertUserRequest{
		Email:      request.Email,
		Password:   passHash,
		TotpSecret: totpSecretEncypt,
	})
	if err != nil {
		a.logger.Errorf("failed to save user: %w", err)

		return "", err
	}
	fmt.Println(totp.URL())
	return totp.Secret(), nil
}

func (a *AuthService) LogIn(ctx context.Context, request *dtoService.LogInRequest) (string, error) {

	user, err := a.userStorage.Get(ctx, &dtoStorage.GetUserRequest{
		Email: request.Email,
	})
	if err != nil {
		a.logger.Errorf("failed to get user: %w", err)

		return "", err
	}

	if request.Email != user.Email {
		a.logger.Infof("invalid email credentials")

		return "", fmt.Errorf("invalid email credentials")
	}
	

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password))
	if err != nil {
		a.logger.Infof("invalid password credentials: %w", err)

		return "", err
	}

	totpSecret, err := crypt.AesDecrypt(crypt.KEY, user.TotpSecret)
	if err != nil {
		a.logger.Errorf("failed to decrypt totp secret: %w", err)
		return "", err
	}

		// app, err := a.appStorage.Get(ctx, &dtoStorage.GetAppRequest{
		// 	Id: request.AppId,
		// })
		// if err != nil {
		// 	a.log.Errorf("failed to get app: %w", err)
		// 	return "", err
		// }

		// a.log.Infof("user logged in successfully")

	if !totp.Validate(request.Token, string(totpSecret)) {
		a.logger.Infof("invalid token")
		return "", fmt.Errorf("invalid token")
	}

	pasetoStruct, err := paseto.NewPaseto(paseto.KEY)
	if err != nil {
		a.logger.Errorf("failed to generate paseto struct", err)
		return "", err
	}
	token, err := pasetoStruct.CreateToken(request.Email, 1 * time.Hour)

	// token, err := pas.NewToken(user, 1 * time.Hour)
	if err != nil {
		a.logger.Errorf("failed to generate paseto token", err)
		return "", err
	}

	return token, nil
}
