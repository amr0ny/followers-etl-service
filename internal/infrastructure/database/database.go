package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPoolOptions struct {
	DSN      string
	MinConns int32
	MaxConns int32
}

func NewPgPool(ctx context.Context, opts PgPoolOptions) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config: %v", err)
	}

	pgxCfg.MinConns = opts.MinConns
	pgxCfg.MaxConns = opts.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection pool: %v", err)
	}

	return pool, nil
}
