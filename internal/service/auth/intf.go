package auth

import (
	"context"
	"sso-like/internal/service/dto"
)

type AuthInterface interface {
	SignUp(ctx context.Context, request *dto.SignUpRequest) (string, error)
	LogIn(ctx context.Context, request *dto.LogInRequest) (string, error)
}
