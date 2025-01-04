package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	

)

// Конструктор Storage
func NewSqliteDb(storagePath string) (*sql.DB, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
