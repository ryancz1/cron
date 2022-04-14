package cron

import (
	"context"
	"sync"
	"time"
)

var (
	timers []Timer
	wg     sync.WaitGroup
)

func Every(interval time.Duration, fn func()) {
	timers = append(timers, newIntervalTimer(interval, fn))
}

func Minutely(fn func()) {
	Every(time.Minute, fn)
}

func Hourly(fn func()) {
	Every(time.Hour, fn)
}

func Daily(hour, min int, fn func()) {
	timers = append(timers, newDailyTimer(hour, min, fn))
}

func Start(ctx context.Context) {
	for _, t := range timers {
		wg.Add(1)
		go func(t Timer) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case <-t.C():
					t.Do()
					t.Next()
				}
			}
		}(t)
	}
}

func Run(ctx context.Context) error {
	Start(ctx)
	wg.Wait()
	return ctx.Err()
}
