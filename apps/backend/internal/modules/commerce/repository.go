package commerce

import (
	"context"
	"sync"

	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

type Repository interface {
	GetState(ctx context.Context, playerID int64) (State, error)
	SaveState(ctx context.Context, playerID int64, state State) error
}

type State struct {
	SigninLastDate string `json:"signin_last_date"`
	SigninStreak   int    `json:"signin_streak"`
	VIPLevel       int    `json:"vip_level"`
	VIPDailyDate   string `json:"vip_daily_date"`
}

type MemoryRepository struct {
	mu     sync.Mutex
	states map[int64]State
}

type MySQLRepository struct {
	store mysqlstore.ModuleStateStore
}

const mysqlModuleName = "commerce"

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		states: make(map[int64]State),
	}
}

func NewMySQLRepository(store mysqlstore.ModuleStateStore) *MySQLRepository {
	return &MySQLRepository{store: store}
}

func (r *MemoryRepository) GetState(_ context.Context, playerID int64) (State, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.ensureStateLocked(playerID), nil
}

func (r *MemoryRepository) SaveState(_ context.Context, playerID int64, state State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.states[playerID] = normalizeState(state)
	return nil
}

func (r *MemoryRepository) ensureStateLocked(playerID int64) State {
	if state, exists := r.states[playerID]; exists {
		return normalizeState(state)
	}

	state := defaultState()
	r.states[playerID] = state
	return state
}

func (r *MySQLRepository) GetState(ctx context.Context, playerID int64) (State, error) {
	var state State
	ok, err := r.store.LoadModuleState(ctx, playerID, mysqlModuleName, &state)
	if err != nil {
		return State{}, err
	}
	if !ok {
		return defaultState(), nil
	}
	return normalizeState(state), nil
}

func (r *MySQLRepository) SaveState(ctx context.Context, playerID int64, state State) error {
	return r.store.SaveModuleState(ctx, playerID, mysqlModuleName, normalizeState(state))
}

func defaultState() State {
	return State{VIPLevel: 0}
}

func normalizeState(state State) State {
	if state.VIPLevel < 0 {
		state.VIPLevel = 0
	}
	if state.SigninStreak < 0 {
		state.SigninStreak = 0
	}
	return state
}
