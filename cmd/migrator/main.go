package main

import (
	"errors"
	"fmt"
	"sso-like/config"


	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// var storagePath, migrationsPath, migrationsTable string

	cfg, err := config.New()
	if err != nil {
		fmt.Println("config error: %w", err)
		return
	}
	// Валидация параметров


	// Создаем объект мигратора, передав креды нашей БД
	m, err := migrate.New(
		"file://"+cfg.DB.Sqlite.MigrationPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", cfg.DB.Sqlite.StoragePath, cfg.DB.Sqlite.MigrationTable),
	)
	if err != nil {
		panic(err)
	}

	// Выполняем миграции до последней версии
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}
}
