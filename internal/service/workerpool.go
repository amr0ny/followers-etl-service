package service

import "github.com/amr0ny/followers-etl-service/internal/domain"

type ErrGroupWorkerPoolService struct {
	workerPool domain.ErrGroupWorkerPool
}

func NewWorkerPoolService(workerPool domain.ErrGroupWorkerPool) *ErrGroupWorkerPoolService {
	return &ErrGroupWorkerPoolService{workerPool: workerPool}
}

// Изменяем Submit, чтобы принимать func() error
func (service *ErrGroupWorkerPoolService) Submit(task func() error) error {
	return service.workerPool.Go(task)
}

func (service *ErrGroupWorkerPoolService) Wait() error {
	return service.workerPool.Wait()
}

func (service *ErrGroupWorkerPoolService) Stats() (active, size int) {
	return service.workerPool.Stats()
}

func (service *ErrGroupWorkerPoolService) Release() {
	service.workerPool.Release()
}

func (service *ErrGroupWorkerPoolService) Reboot() {
	service.workerPool.Reboot()
}
