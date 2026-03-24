package jobs

import (
	"context"
	"time"
)

type AllianceWarPhaseJob struct {
	store *MemoryStateStore
	now   func() time.Time
}

func NewAllianceWarPhaseJob(store *MemoryStateStore, now func() time.Time) *AllianceWarPhaseJob {
	if now == nil {
		now = time.Now
	}
	return &AllianceWarPhaseJob{store: store, now: now}
}

func (j *AllianceWarPhaseJob) Run(_ context.Context) error {
	j.store.mu.Lock()
	defer j.store.mu.Unlock()
	current := j.now()
	for roundID, round := range j.store.AllianceWars {
		switch {
		case round.Phase == "signup" && !current.Before(round.LockAt):
			round.Phase = "locked"
		case round.Phase == "locked" && !current.Before(round.BattleAt):
			round.Phase = "battle"
		case round.Phase == "battle" && !current.Before(round.ResultAt):
			round.Phase = "result"
		}
		j.store.AllianceWars[roundID] = round
	}
	return nil
}
