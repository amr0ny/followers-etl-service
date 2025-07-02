package service

import (
	"context"
	"github.com/amr0ny/followers-etl-service/internal/domain"
)

type FSManagerService struct {
	fsManager domain.FSManager
}

func NewFSManagerService(fsManager domain.FSManager) *FSManagerService {
	return &FSManagerService{fsManager: fsManager}
}

func (service *FSManagerService) ReadCSV(ctx context.Context) (<-chan [][]string, error) {
	return service.fsManager.ReadFileBatched(ctx)
}
