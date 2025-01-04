package storage

import (
	"context"
	"fmt"

	"sso-like/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDbConn(ctx context.Context, cfg *config.PostgresConfig) (pool *pgxpool.Pool, err error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		cfg.Driver,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("connectoing to db: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return pool, nil
}
