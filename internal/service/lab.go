package service

import (
	"context"
	"sync"

	"github.com/jb843051627/foram-bench/internal/clock"
	"github.com/jb843051627/foram-bench/internal/ingest"
	"github.com/jb843051627/foram-bench/internal/metrics"
	"github.com/jb843051627/foram-bench/internal/store"
)

type Lab struct {
	store      *store.Store
	clock      clock.Clock
	queue      *ingest.Queue
	metrics    *metrics.Registry
	cacheMu    sync.RWMutex
	batchCache map[string]batchCacheEntry
	stateMu    sync.Mutex
	closed     chan struct{}
}

type batchCacheEntry struct {
	status   string
	revision int
}

func NewLab(repository *store.Store) *Lab {
	return NewLabWithClock(repository, clock.System{})
}

func NewLabWithClock(repository *store.Store, now clock.Clock) *Lab {
	return &Lab{
		store:      repository,
		clock:      now,
		queue:      ingest.New(32),
		metrics:    metrics.New(),
		batchCache: make(map[string]batchCacheEntry),
		closed:     make(chan struct{}),
	}
}

func (l *Lab) Close() {
	select {
	case <-l.closed:
		return
	default:
		close(l.closed)
	}
	l.queue.Close()
}

func (l *Lab) Metrics() map[string]int64 {
	return l.metrics.Snapshot()
}

func (l *Lab) cacheBatch(id, status string, revision int) {
	l.cacheMu.Lock()
	l.batchCache[id] = batchCacheEntry{status: status, revision: revision}
	l.cacheMu.Unlock()
}

func (l *Lab) cachedBatch(id string) (batchCacheEntry, bool) {
	l.cacheMu.RLock()
	entry, ok := l.batchCache[id]
	l.cacheMu.RUnlock()
	return entry, ok
}

func checkContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}
