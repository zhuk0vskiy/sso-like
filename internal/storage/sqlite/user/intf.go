package user

import (
	"context"
	"sso-like/internal/storage/dto"
	"sso-like/internal/model"
)

type UserInterface interface {
	Insert(ctx context.Context, request *dto.InsertUserRequest) (int64, error)
	Get(ctx context.Context, request *dto.GetUserRequest) (model.User, error)
}
