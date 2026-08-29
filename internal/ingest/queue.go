package ingest

import (
	"context"
	"sync"

	"github.com/jb843051627/foram-bench/internal/model"
)

type Job struct {
	ID   string
	Ctx  context.Context
	Run  func(context.Context) error
	Done chan error
}
type Queue struct {
	jobs   chan Job
	stop   chan struct{}
	once   sync.Once
	mu     sync.RWMutex
	closed bool
	wg     sync.WaitGroup
}

func New(size int) *Queue {
	q := &Queue{jobs: make(chan Job, size), stop: make(chan struct{})}
	q.wg.Add(1)
	go q.loop()
	return q
}
func (q *Queue) loop() {
	defer q.wg.Done()
	for {
		select {
		case job := <-q.jobs:
			ctx := jobContext(job)
			var err error
			if job.Run == nil {
				err = model.ErrInvalidInput
			} else {
				err = job.Run(ctx)
			}
			if job.Done == nil {
				continue
			}
			select {
			case job.Done <- err:
			case <-q.stop:
			}
		case <-q.stop:
			return
		}
	}
}

func jobContext(job Job) context.Context {
	if job.Ctx != nil {
		return job.Ctx
	}
	return context.Background()
}
func (q *Queue) Submit(ctx context.Context, job Job) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	q.mu.RLock()
	closed := q.closed
	q.mu.RUnlock()
	if closed {
		return model.ErrQueueClosed
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.jobs <- job:
		return nil
	case <-q.stop:
		return model.ErrQueueClosed
	}
}
func (q *Queue) Close() {
	q.once.Do(func() {
		q.mu.Lock()
		q.closed = true
		close(q.stop)
		q.mu.Unlock()
		q.wg.Wait()
	})
}
