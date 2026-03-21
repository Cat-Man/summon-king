package jobs

import (
	"context"
	"time"
)

type ManorTickJob struct {
	store *MemoryStateStore
	now   func() time.Time
}

func NewManorTickJob(store *MemoryStateStore, now func() time.Time) *ManorTickJob {
	if now == nil {
		now = time.Now
	}
	return &ManorTickJob{store: store, now: now}
}

func (j *ManorTickJob) Run(_ context.Context) error {
	j.store.mu.Lock()
	defer j.store.mu.Unlock()
	current := j.now()
	for plotID, plot := range j.store.ManorPlots {
		if plot.Status == "growing" && !current.Before(plot.MatureAt) {
			plot.Status = "ready"
			j.store.ManorPlots[plotID] = plot
		}
	}
	return nil
}
