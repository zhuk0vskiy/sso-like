package app

import (
	"database/sql"
	"time"

	grpc "sso-like/internal/grpc"
	"sso-like/internal/service/auth"
	appStorage "sso-like/internal/storage/sqlite/app"
	userStorage "sso-like/internal/storage/sqlite/user"
	"sso-like/pkg/logger"
)

type App struct {
	GrpcApp *grpc.GrpcApp
}

func NewApp(
	log *logger.Logger,
	dbConnector *sql.DB,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	
	appStrg := appStorage.NewAppStorage(dbConnector)
	userStrg := userStorage.NewUserStorage(dbConnector)

	authService := auth.NewAuthService(log, userStrg, appStrg, tokenTTL)

	grpcApp := grpc.NewGrpcApp(log, authService, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}
}
