package user

import (
	"context"
	"sso-like/internal/model"
	"sso-like/internal/storage/dto"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UserInterface
type UserInterface interface {
	Insert(ctx context.Context, request *dto.InsertUserRequest) (err error)
	Get(ctx context.Context, request *dto.GetUserRequest) (user *model.User, err error)
}
