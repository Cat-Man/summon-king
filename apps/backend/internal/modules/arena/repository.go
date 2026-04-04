package arena

import (
	"context"
	"errors"
	"sync"
)

var ErrArenaRecordNotFound = errors.New("arena record not found")

type Repository interface {
	SetCurrentStreak(ctx context.Context, playerID int64, streak int) error
	RecordBattleResult(ctx context.Context, playerID int64, won bool) error
	GetDailyRecord(ctx context.Context, playerID int64) (DailyRecord, error)
}

type MemoryRepository struct {
	mu      sync.Mutex
	records map[int64]DailyRecord
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		records: make(map[int64]DailyRecord),
	}
}

func (r *MemoryRepository) SetCurrentStreak(_ context.Context, playerID int64, streak int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	record := r.records[playerID]
	record.PlayerID = playerID
	record.CurrentStreak = streak
	record.LastWin = streak > 0
	r.records[playerID] = record
	return nil
}

func (r *MemoryRepository) RecordBattleResult(_ context.Context, playerID int64, won bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	record := r.records[playerID]
	record.PlayerID = playerID
	record.LastWin = won
	if won {
		record.CurrentStreak++
	} else {
		record.CurrentStreak = 0
	}
	r.records[playerID] = record
	return nil
}

func (r *MemoryRepository) GetDailyRecord(_ context.Context, playerID int64) (DailyRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.records[playerID]
	if !ok {
		return DailyRecord{
			PlayerID:      playerID,
			CurrentStreak: 0,
			LastWin:       false,
		}, nil
	}
	return record, nil
}
