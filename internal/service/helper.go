package service

import (
	"context"
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"time"
)

type HelperService struct {
	rep domain.HelperRepository
}

func NewHelperService(rep domain.HelperRepository) *HelperService {
	return &HelperService{rep: rep}
}

func (service *HelperService) GetLastTableActivityTimestamp(ctx context.Context, tableName string) (time.Time, error) {
	return service.rep.RetrieveLastTableActivity(ctx, tableName)
}

func (service *HelperService) SetLastTableActivityTimestamp(ctx context.Context, tableName string, timestamp time.Time) error {
	return service.rep.UpdateLastTableActivity(ctx, tableName, timestamp)
}
