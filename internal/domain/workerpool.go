package domain

import "context"

type ErrGroupWorkerPool interface {
	Go(f func() error) error
	Wait() error
	Context() context.Context
	Stats() (active, size int)
	Release()
	Reboot()
}
