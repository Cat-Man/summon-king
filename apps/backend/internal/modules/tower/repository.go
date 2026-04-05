package tower

import (
	"context"
	"errors"
	"sync"

	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

var ErrNoChallengesRemaining = errors.New("tower challenges exhausted")

type Repository interface {
	StartChallenge(ctx context.Context, playerID int64, tower string) (TowerResult, error)
	GetStatus(ctx context.Context, playerID int64, tower string) TowerStatus
}

type MemoryRepository struct {
	mu       sync.Mutex
	progress map[string]map[int64]int
}

type MySQLRepository struct {
	store mysqlstore.ModuleStateStore
}

type mysqlState struct {
	Progress map[string]int `json:"progress"`
}

const mysqlModuleName = "tower"

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		progress: map[string]map[int64]int{
			"pagoda": {},
			"spirit": {},
		},
	}
}

func NewMySQLRepository(store mysqlstore.ModuleStateStore) *MySQLRepository {
	return &MySQLRepository{store: store}
}

func (r *MemoryRepository) StartChallenge(_ context.Context, playerID int64, tower string) (TowerResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	progress := r.ensureTowerProgress(tower)
	if progress[playerID] >= 5 {
		return TowerResult{}, ErrNoChallengesRemaining
	}
	floor := progress[playerID] + 1
	progress[playerID] = floor

	status := buildTowerStatus(tower, floor)
	return TowerResult{
		PlayerID:            playerID,
		Tower:               status.Tower,
		Floor:               floor,
		Reward:              status.RewardPreview,
		RemainingChallenges: status.RemainingChallenges,
	}, nil
}

func (r *MemoryRepository) GetStatus(_ context.Context, playerID int64, tower string) TowerStatus {
	r.mu.Lock()
	defer r.mu.Unlock()

	progress := r.ensureTowerProgress(tower)
	return buildTowerStatus(tower, progress[playerID])
}

func (r *MemoryRepository) ensureTowerProgress(tower string) map[int64]int {
	if _, ok := r.progress[tower]; !ok {
		r.progress[tower] = make(map[int64]int)
	}
	return r.progress[tower]
}

func buildTowerStatus(tower string, currentFloor int) TowerStatus {
	remaining := 5 - currentFloor
	if remaining < 0 {
		remaining = 0
	}

	status := TowerStatus{
		Tower:               tower,
		CurrentFloor:        currentFloor,
		MaxFloor:            10,
		RemainingChallenges: remaining,
	}

	if tower == "spirit" {
		status.Label = "战灵塔"
		status.RewardPreview = "灵魂碎片"
		status.MaxFloor = 12
		return status
	}

	status.Tower = "pagoda"
	status.Label = "通天塔"
	status.RewardPreview = "战骨锻造"
	return status
}

func (r *MySQLRepository) StartChallenge(ctx context.Context, playerID int64, tower string) (TowerResult, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return TowerResult{}, err
	}

	if state.Progress[tower] >= 5 {
		return TowerResult{}, ErrNoChallengesRemaining
	}
	floor := state.Progress[tower] + 1
	state.Progress[tower] = floor

	status := buildTowerStatus(tower, floor)
	return TowerResult{
		PlayerID:            playerID,
		Tower:               status.Tower,
		Floor:               floor,
		Reward:              status.RewardPreview,
		RemainingChallenges: status.RemainingChallenges,
	}, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) GetStatus(ctx context.Context, playerID int64, tower string) TowerStatus {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return buildTowerStatus(tower, 0)
	}
	return buildTowerStatus(tower, state.Progress[tower])
}

func (r *MySQLRepository) loadState(ctx context.Context, playerID int64) (mysqlState, error) {
	state := mysqlState{Progress: map[string]int{"pagoda": 0, "spirit": 0}}
	ok, err := r.store.LoadModuleState(ctx, playerID, mysqlModuleName, &state)
	if err != nil {
		return mysqlState{}, err
	}
	if !ok || state.Progress == nil {
		state.Progress = map[string]int{"pagoda": 0, "spirit": 0}
	}
	if _, ok := state.Progress["pagoda"]; !ok {
		state.Progress["pagoda"] = 0
	}
	if _, ok := state.Progress["spirit"]; !ok {
		state.Progress["spirit"] = 0
	}
	return state, nil
}

func (r *MySQLRepository) saveState(ctx context.Context, playerID int64, state mysqlState) error {
	if state.Progress == nil {
		state.Progress = map[string]int{"pagoda": 0, "spirit": 0}
	}
	return r.store.SaveModuleState(ctx, playerID, mysqlModuleName, state)
}
