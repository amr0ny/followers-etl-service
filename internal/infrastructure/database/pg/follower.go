package pg

import (
	"context"
	"fmt"
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type FollowerRepositoryPG struct {
	*BaseRepositoryPG
}

func NewFollowerRepositoryPG(pool *pgxpool.Pool) domain.FollowerRepository {
	return &FollowerRepositoryPG{&BaseRepositoryPG{pool: pool}}
}
func (r *FollowerRepositoryPG) BulkUpsertStagingFollowers(ctx context.Context, followers []domain.Follower, loadId uuid.UUID) error {
	if len(followers) == 0 {
		return nil
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	batch := &pgx.Batch{}

	query := `INSERT INTO public.staging_followers (email, full_name, load_id) 
				VALUES ($1, $2, $3);`
	for _, follower := range followers {
		batch.Queue(query, follower.Email, follower.FullName, loadId)
	}
	br := tx.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return fmt.Errorf("batch item %d failed: %w", i, err)
		}
	}
	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close batch results: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}
	return nil
}

func (r *FollowerRepositoryPG) RemoveFollowersByTimestamp(ctx context.Context, timestamp time.Time) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	query := "DELETE FROM public.followers WHERE modified_at < $1"
	_, err = conn.Exec(ctx, query, timestamp)
	if err != nil {
		return fmt.Errorf("failed to perform a query: %v", err)
	}
	return nil
}

func (r *FollowerRepositoryPG) DiscardStagingChangesByLoadId(ctx context.Context, loadId uuid.UUID) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	query := "DELETE FROM public.staging_followers WHERE load_id = $1"
	_, err = conn.Exec(ctx, query, loadId)
	if err != nil {
		return fmt.Errorf("failed to perform a query: %v", err)
	}
	return nil
}

func (r *FollowerRepositoryPG) MoveFollowersToProductionTable(ctx context.Context, loadId uuid.UUID) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	query := `
		WITH batch AS (
			SELECT * FROM staging_followers WHERE load_id = $1
		)
		INSERT INTO followers (email, full_name)
		SELECT email, full_name FROM batch
		ON CONFLICT (email) DO UPDATE SET full_name=EXCLUDED.full_name, modified_at=DEFAULT
	`
	_, err = conn.Exec(ctx, query, loadId)
	if err != nil {
		return fmt.Errorf("failed to perform a query: %v", err)
	}
	return nil
}
