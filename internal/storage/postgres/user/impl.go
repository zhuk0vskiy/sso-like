package user

import (
	"context"
	"fmt"
	"sso-like/internal/model"
	"sso-like/internal/storage/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) UserInterface {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Insert(ctx context.Context, request *dto.InsertUserRequest) (err error) {
	query := `insert into users(email, password, totp_secret) values ($1, $2, $3)`

	_, err = s.db.Exec(
		ctx,
		query,
		request.Email,
		request.Password,
		request.TotpSecret,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return err
}

func (s *UserStorage) Get(ctx context.Context, request *dto.GetUserRequest) (user *model.User, err error) {

	query := `select id, email, password, totp_secret from users where email = $1`

	user = new(model.User)

	err = s.db.QueryRow(
		ctx,
		query,
		request.Email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.TotpSecret,
	)

	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, err
}
