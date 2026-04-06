package arena

import (
	"context"
	"errors"
	"sync"

	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
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

type MySQLRepository struct {
	store mysqlstore.ModuleStateStore
}

type mysqlState struct {
	Record         DailyRecord `json:"record"`
	RefreshVersion int         `json:"refresh_version"`
}

const mysqlModuleName = "arena"

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		records:         make(map[int64]DailyRecord),
		refreshVersions: make(map[int64]int),
	}
}

func NewMySQLRepository(store mysqlstore.ModuleStateStore) *MySQLRepository {
	return &MySQLRepository{store: store}
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

func (r *MySQLRepository) SetCurrentStreak(ctx context.Context, playerID int64, streak int) error {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return err
	}
	state.Record.PlayerID = playerID
	state.Record.CurrentStreak = streak
	state.Record.LastWin = streak > 0
	return r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) RecordBattleResult(ctx context.Context, playerID int64, won bool) error {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return err
	}
	state.Record.PlayerID = playerID
	state.Record.LastWin = won
	if won {
		state.Record.CurrentStreak++
	} else {
		state.Record.CurrentStreak = 0
	}
	return r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) GetDailyRecord(ctx context.Context, playerID int64) (DailyRecord, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return DailyRecord{}, err
	}
	state.Record.PlayerID = playerID
	return state.Record, nil
}

func (r *MySQLRepository) GetRefreshVersion(ctx context.Context, playerID int64) (int, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return 0, err
	}
	return state.RefreshVersion, nil
}

func (r *MySQLRepository) IncrementRefreshVersion(ctx context.Context, playerID int64) (int, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return 0, err
	}
	state.RefreshVersion++
	return state.RefreshVersion, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) loadState(ctx context.Context, playerID int64) (mysqlState, error) {
	var state mysqlState
	ok, err := r.store.LoadModuleState(ctx, playerID, mysqlModuleName, &state)
	if err != nil {
		return mysqlState{}, err
	}
	if !ok {
		state.Record = DailyRecord{
			PlayerID:      playerID,
			CurrentStreak: 0,
			LastWin:       false,
		}
	}
	if state.Record.PlayerID == 0 {
		state.Record.PlayerID = playerID
	}
	return state, nil
}

func (r *MySQLRepository) saveState(ctx context.Context, playerID int64, state mysqlState) error {
	if state.Record.PlayerID == 0 {
		state.Record.PlayerID = playerID
	}
	return r.store.SaveModuleState(ctx, playerID, mysqlModuleName, state)
}
