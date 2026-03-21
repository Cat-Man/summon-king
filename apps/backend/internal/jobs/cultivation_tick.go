package jobs

import (
	"context"
	"time"
)

type CultivationTickJob struct {
	store *MemoryStateStore
	now   func() time.Time
}

func NewCultivationTickJob(store *MemoryStateStore, now func() time.Time) *CultivationTickJob {
	if now == nil {
		now = time.Now
	}
	return &CultivationTickJob{store: store, now: now}
}

func (j *CultivationTickJob) Run(_ context.Context) error {
	j.store.mu.Lock()
	defer j.store.mu.Unlock()
	current := j.now()
	for playerID, task := range j.store.Cultivations {
		if !task.Claimable && !current.Before(task.ReadyAt) {
			task.Claimable = true
			j.store.Cultivations[playerID] = task
		}
	}
	return nil
}
