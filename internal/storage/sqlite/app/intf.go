package app

import (
	"context"
	"sso-like/internal/model"
	"sso-like/internal/storage/dto"
)

type AppInterface interface {
	Get(ctx context.Context, request *dto.GetAppRequest) (model.App, error)
}
