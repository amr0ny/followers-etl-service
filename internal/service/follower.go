package service

import (
	"context"
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"github.com/google/uuid"
	"time"
)

type FollowerService struct {
	rep domain.FollowerRepository
}

func NewFollowerService(rep domain.FollowerRepository) *FollowerService {
	return &FollowerService{rep: rep}
}

func (service *FollowerService) UpsertStagingFollowers(ctx context.Context, followers []domain.Follower, loadId uuid.UUID) error {
	err := service.rep.BulkUpsertStagingFollowers(ctx, followers, loadId)
	if err != nil {
		return err
	}
	return nil
}

func (service *FollowerService) RemoveOutdatedFollowers(ctx context.Context, timestamp time.Time) error {
	return service.rep.RemoveFollowersByTimestamp(ctx, timestamp)
}

func (service *FollowerService) CleanStagingTable(ctx context.Context, loadId uuid.UUID) error {
	return service.rep.DiscardStagingChangesByLoadId(ctx, loadId)
}

func (service *FollowerService) MoveFollowersToProductionTable(ctx context.Context, loadId uuid.UUID) error {
	return service.rep.MoveFollowersToProductionTable(ctx, loadId)
}
