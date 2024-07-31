package cron

import (
	"log/slog"
	"sync"
	"time"
)

const dailyInterval = 24 * time.Hour

// NewDailyJob constructs job which will  be triggered on the
// daily basis at the given point of time.
// NOTE: it expects time in UTC timezone. It's -3 hours compared to Kyiv Time:
// 12:00 UTC = 15:00 by Kyiv
func NewDailyJob(hour, min int, log *slog.Logger) *Job {
	return &Job{
		interval: dailyInterval,
		ticker:   time.NewTicker(calculateFirstTick(hour, min)),
		log:      log,
	}
}

// NewJob constructs job which will  be triggered after
// provided interval.
func NewJob(interval time.Duration, log *slog.Logger) *Job {
	return &Job{
		interval: interval,
		ticker:   time.NewTicker(interval),
		log:      log,
	}
}

// Job reporesents Cron Job which could be runned in some
// interval or on daily basis. It's a thin wrapper around stdlib's
// time.Ticker.
type Job struct {
	interval time.Duration
	ticker   *time.Ticker
	log      *slog.Logger
}

//go:generate mockgen -destination=./mocks/mock_doer.go -package=mocks . Doer
type Doer interface {
	Do() error
}

// Do method calls provided fn with the given interval.
// Does not stop on error, only logs it and then goes on.
// Reset is needed for the daily job, to change it to 24 hours and it's done only once
// after first run.
func (j *Job) Do(fn Doer) {
	var once sync.Once
	go func() {
		for range j.ticker.C {
			once.Do(func() {
				j.ticker.Reset(j.interval)
			})
			if err := fn.Do(); err != nil {
				j.log.Error("Failed to do cron task", slog.Any("err", err))
			}
		}
	}()
}

func calculateFirstTick(hour, min int) time.Duration {
	now := time.Now()
	tickAt := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.UTC)
	if now.After(tickAt) {
		tickAt = tickAt.Add(dailyInterval)
	}

	return tickAt.Sub(now)
}
