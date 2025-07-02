package workerpool

import (
	"context"
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"github.com/panjf2000/ants/v2"
	"sync"
)

type AntsErrGroupCapacity int

type AntsErrGroup struct {
	pool    *ants.Pool
	wg      sync.WaitGroup
	mu      sync.Mutex
	err     error
	ctx     context.Context
	cancel  context.CancelFunc
	errOnce sync.Once
	errChan chan error
}

func NewAntsErrGroup(capacity AntsErrGroupCapacity, opts ...ants.Option) (domain.ErrGroupWorkerPool, error) {
	p, err := ants.NewPool(int(capacity), opts...)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &AntsErrGroup{
		pool:    p,
		ctx:     ctx,
		cancel:  cancel,
		errChan: make(chan error, int(capacity)),
	}, nil
}

func (g *AntsErrGroup) Go(f func() error) error {
	if g.ctx.Err() != nil {
		return g.ctx.Err()
	}

	g.wg.Add(1)
	err := g.pool.Submit(func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.mu.Lock()
				g.err = err
				g.mu.Unlock()
				g.cancel()
			})
		}
	})
	if err != nil {
		g.wg.Done()
		return err
	}
	return nil
}

func (g *AntsErrGroup) Wait() error {
	g.wg.Wait()
	g.Release()
	return g.err
}

func (g *AntsErrGroup) Context() context.Context {
	return g.ctx
}

func (p *AntsErrGroup) Release() {
	p.pool.Release()
}

func (p *AntsErrGroup) Stats() (active, size int) {
	return p.pool.Running(), p.pool.Cap()
}

func (p *AntsErrGroup) Reboot() {
	p.pool.Reboot()
}
