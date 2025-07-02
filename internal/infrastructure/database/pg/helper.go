package pg

import (
	"context"
	"fmt"
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type HelperRepositoryPG struct {
	*BaseRepositoryPG
}

func NewHelperRepositoryPG(pool *pgxpool.Pool) domain.HelperRepository {
	return &HelperRepositoryPG{&BaseRepositoryPG{pool: pool}}
}

func (rep *HelperRepositoryPG) RetrieveLastTableActivity(ctx context.Context, tableName string) (time.Time, error) {
	conn, err := rep.pool.Acquire(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to acquire connection: %v", err)
	}
	defer conn.Release()
	query := "SELECT last_activity FROM etl_info_schema.table_status WHERE table_name = $1"
	row := conn.QueryRow(ctx, query, tableName)
	var lastActivity time.Time
	if err := row.Scan(&lastActivity); err != nil {
		if err == pgx.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("failed to parse resulting row: %v", err)
	}

	return lastActivity, nil
}

func (rep *HelperRepositoryPG) UpdateLastTableActivity(ctx context.Context, tableName string, timestamp time.Time) error {
	conn, err := rep.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %v", err)
	}

	defer conn.Release()
	query := `INSERT INTO etl_info_schema.table_status (table_name, last_activity)
				VALUES ($1, $2)
				ON CONFLICT (table_name)
				DO UPDATE
					SET last_activity = EXCLUDED.last_activity`
	_, err = conn.Exec(ctx, query, tableName, timestamp)
	if err != nil {
		return fmt.Errorf("failed to perform a query: %v", err)
	}
	return nil
}
