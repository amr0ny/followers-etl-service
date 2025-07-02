package app

import (
	"context"
	"fmt"
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"github.com/amr0ny/followers-etl-service/internal/service"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"time"
)

type App struct {
	Logger            domain.Logger
	FollowerService   *service.FollowerService
	HelperService     *service.HelperService
	FSManagerService  *service.FSManagerService
	WorkerPoolService *service.ErrGroupWorkerPoolService
}

func (a *App) discardChanges(ctx context.Context, err error, loadId uuid.UUID, timestamp *time.Time) error {
	var innerErr error
	if timestamp != nil {
		innerErr = a.HelperService.SetLastTableActivityTimestamp(ctx, "followers", *timestamp)
	}
	innerErr = a.FollowerService.CleanStagingTable(ctx, loadId)
	if innerErr != nil {
		return fmt.Errorf("an error occured:%v, failed to discard changes: %v", err, innerErr)
	}
	return fmt.Errorf("an error occured, discarding changes: %v", err)
}

func (a *App) processBatches(ctx context.Context, batches <-chan [][]string) (uuid.UUID, error) {
	loadId, err := uuid.NewV7()
	if err != nil {
		return loadId, err
	}
	for batch := range batches {
		batch := batch
		err := a.WorkerPoolService.Submit(func() error {
			var followers []domain.Follower
			for _, record := range batch {
				record := record
				follower, err := domain.NewFollower(record[0], record[1])
				if err != nil {
					return fmt.Errorf("failed to created follower: %v", err)
				}
				followers = append(followers, follower)
			}
			return a.FollowerService.UpsertStagingFollowers(ctx, followers, loadId)
		})
		if err != nil {
			return loadId, fmt.Errorf("failed to submit a task to a worker: %v", err)
		}
	}
	return loadId, a.WorkerPoolService.Wait()
}

func (a *App) StartCron(ctx context.Context, schedule string) error {
	c := cron.New(cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))

	_, err := c.AddFunc(schedule, func() {
		if err := a.Run(ctx); err != nil {
			a.Logger.Error("cron task failed", "error", err)
		}
	})
	if err != nil {
		return err
	}

	c.Start()
	a.Logger.Info("Cron started")
	go func() {
		<-ctx.Done()
		c.Stop()
	}()

	return nil
}

func (a *App) Run(ctx context.Context) error {
	a.Logger.Info("Starting syncronization procedure...")
	a.WorkerPoolService.Reboot()
	batches, err := a.FSManagerService.ReadCSV(ctx)
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}
	loadId, err := a.processBatches(ctx, batches)
	if err != nil {
		return a.discardChanges(ctx, err, loadId, nil)
	}
	a.Logger.Info("Staging table filled with values")
	timestamp, err := a.HelperService.GetLastTableActivityTimestamp(ctx, "followers")
	if err != nil {
		return a.discardChanges(ctx, err, loadId, nil)
	}
	a.Logger.Info("Last table syncronization: %v", timestamp.String())
	err = a.HelperService.SetLastTableActivityTimestamp(ctx, "followers", time.Now())
	if err != nil {
		return a.discardChanges(ctx, err, loadId, nil)
	}
	err = a.FollowerService.MoveFollowersToProductionTable(ctx, loadId)
	if err != nil {
		return a.discardChanges(ctx, err, loadId, &timestamp)
	}
	a.Logger.Info("Changes written to main table")
	timestamp, err = a.HelperService.GetLastTableActivityTimestamp(ctx, "followers")
	if err != nil {
		return a.discardChanges(ctx, err, loadId, &timestamp)
	}
	err = a.FollowerService.RemoveOutdatedFollowers(ctx, timestamp)
	if err != nil {
		return a.discardChanges(ctx, err, loadId, &timestamp)
	}
	a.Logger.Info("Outdated data removed")
	err = a.FollowerService.CleanStagingTable(ctx, loadId)
	if err != nil {
		return fmt.Errorf("failed to clean up staging table: %v", err)
	}
	a.Logger.Info("Synchronization completed")
	return nil
}
