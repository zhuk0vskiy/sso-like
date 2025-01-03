package auth

import (
	"context"
	"sso-like/internal/service/dto"
)

type AuthInterface interface {
	SignUp(ctx context.Context, request *dto.SignUpRequest) (int64, error)
	LogIn(ctx context.Context, request *dto.LogInRequest) (string, error)
}
