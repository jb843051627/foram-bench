package ingest

import (
	"context"
	"time"
)

func Retry(ctx context.Context, attempts int, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = ctx.Err(); err != nil {
			return err
		}
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return err
}

func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	<-timer.C
	return nil
}
