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
	GetRefreshVersion(ctx context.Context, playerID int64) (int, error)
	IncrementRefreshVersion(ctx context.Context, playerID int64) (int, error)
}

type MemoryRepository struct {
	mu              sync.Mutex
	records         map[int64]DailyRecord
	refreshVersions map[int64]int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		records:         make(map[int64]DailyRecord),
		refreshVersions: make(map[int64]int),
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

func (r *MemoryRepository) GetRefreshVersion(_ context.Context, playerID int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.refreshVersions[playerID], nil
}

func (r *MemoryRepository) IncrementRefreshVersion(_ context.Context, playerID int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.refreshVersions[playerID]++
	return r.refreshVersions[playerID], nil
}
