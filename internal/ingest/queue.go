package ingest

import (
	"context"
	"sync"
)

type Job struct {
	ID   string
	Ctx  context.Context
	Run  func(context.Context) error
	Done chan error
}
type Queue struct {
	jobs chan Job
	stop chan struct{}
	once sync.Once
	wg   sync.WaitGroup
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
			err := job.Run(ctx)
			select {
			case job.Done <- err:
			case <-q.stop:
			default:
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
	select {
	case <-q.stop:
		return context.Canceled
	default:
	}
	select {
	case q.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-q.stop:
		return context.Canceled
	}
}
func (q *Queue) Close() { q.once.Do(func() { close(q.stop); q.wg.Wait() }) }
