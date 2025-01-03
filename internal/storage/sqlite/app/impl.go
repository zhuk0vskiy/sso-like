package app

import (
	"context"
	"database/sql"
	"fmt"
	"sso-like/internal/model"
	"sso-like/internal/storage/dto"
)

type AppStorage struct {
	db *sql.DB
}

func NewAppStorage(db *sql.DB) AppInterface {
	return &AppStorage{
		db: db,
	}
}

func (s *AppStorage) Get(ctx context.Context, request *dto.GetAppRequest) (model.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, request.Id)

	var app model.App
	err = row.Scan(&app.Id, &app.Name, &app.Secret)
	if err != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	return model.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		// }

		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
