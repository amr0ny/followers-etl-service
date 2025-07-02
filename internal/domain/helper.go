package domain

import (
	"context"
	"time"
)

type HelperRepository interface {
	RetrieveLastTableActivity(ctx context.Context, tableName string) (time.Time, error)
	UpdateLastTableActivity(ctx context.Context, tableName string, timestamp time.Time) error
}
