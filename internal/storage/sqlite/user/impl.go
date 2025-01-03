package user

import (
	"context"
	"database/sql"
	"fmt"
	"sso-like/internal/model"
	"sso-like/internal/storage/dto"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) UserInterface {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Insert(ctx context.Context, request *dto.InsertUserRequest) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем запрос, передав параметры
	res, err := stmt.ExecContext(ctx, request.Email, request.PassHash)
	if err != nil {
		// var sqliteErr sqlite3.Error

		// Небольшое кунг-фу для выявления ошибки ErrConstraintUnique
		// (см. подробности ниже)
		// if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		// 	return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		// }

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *UserStorage) Get(ctx context.Context, request *dto.GetUserRequest) (model.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, request.Email)

	var user model.User
	err = row.Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		// }

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
