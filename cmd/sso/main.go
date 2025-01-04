package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sso-like/config"
	"sso-like/internal/app"
	"sso-like/internal/grpc"
	"sso-like/internal/storage"
	"sso-like/pkg/logger"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		fmt.Println("config error: %w", err)
		return
	}

	fmt.Println("trying to create logger")
	loggerFile, err := os.OpenFile(
		cfg.Logger.File,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func(loggerFile *os.File) {
		err := loggerFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(loggerFile)

	log := logger.New(cfg.Logger.Level, loggerFile)

	fmt.Println("trying to connect db")
	ctx := context.Background()
	dbConnector, err := storage.NewDbConn(ctx, &cfg.DB.Postgres)
	if err != nil {

		fmt.Printf("cannot connect to db", err)
		return
	}

	a := app.NewApp(log, dbConnector, cfg.SSO.GRPC.Port, cfg.TokenTTL)


	grpc.Run(a.GrpcApp)
}
