package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type FollowerRepository interface {
	BulkUpsertStagingFollowers(ctx context.Context, followers []Follower, loadId uuid.UUID) error
	RemoveFollowersByTimestamp(ctx context.Context, timestamp time.Time) error
	DiscardStagingChangesByLoadId(ctx context.Context, loadId uuid.UUID) error
	MoveFollowersToProductionTable(ctx context.Context, loadId uuid.UUID) error
}

type Follower struct {
	Email      string
	FullName   string
	ModifiedAt time.Time
}

func NewFollower(email, fullName string) (Follower, error) {
	return Follower{
		Email:      email,
		FullName:   fullName,
		ModifiedAt: time.Now(),
	}, nil
}
