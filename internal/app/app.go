package app

import (

	"time"

	grpc "sso-like/internal/grpc"
	"sso-like/internal/service/auth"

	// appStorage "sso-like/internal/storage/postgres/app"z
	userStorage "sso-like/internal/storage/postgres/user"
	"sso-like/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	GrpcApp *grpc.GrpcApp
}

func NewApp(
	log logger.Interface,
	dbConnector *pgxpool.Pool,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	
	// appStrg := appStorage.NewAppStorage(dbConnector)
	userStrg := userStorage.NewUserStorage(dbConnector)

	authService := auth.NewAuthService(log, userStrg)

	grpcApp := grpc.NewGrpcApp(log, authService, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}
}
