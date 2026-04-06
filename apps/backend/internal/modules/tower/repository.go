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
	SaveLastReward(ctx context.Context, playerID int64, tower string, reward string, delta TowerRewardDelta) error
}

type MemoryRepository struct {
	mu              sync.Mutex
	progress        map[string]map[int64]int
	lastReward      map[string]map[int64]string
	lastRewardDelta map[string]map[int64]TowerRewardDelta
}

type MySQLRepository struct {
	store mysqlstore.ModuleStateStore
}

type mysqlState struct {
	Progress        map[string]int              `json:"progress"`
	LastReward      map[string]string           `json:"last_reward,omitempty"`
	LastRewardDelta map[string]TowerRewardDelta `json:"last_reward_delta,omitempty"`
}

const mysqlModuleName = "tower"

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		progress: map[string]map[int64]int{
			"pagoda": {},
			"spirit": {},
		},
		lastReward: map[string]map[int64]string{
			"pagoda": {},
			"spirit": {},
		},
		lastRewardDelta: map[string]map[int64]TowerRewardDelta{
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

	status := buildTowerStatus(tower, floor, r.lastReward[tower][playerID], r.lastRewardDelta[tower][playerID])
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
	return buildTowerStatus(
		tower,
		progress[playerID],
		r.ensureTowerLastReward(tower)[playerID],
		r.ensureTowerLastRewardDelta(tower)[playerID],
	)
}

func (r *MemoryRepository) SaveLastReward(_ context.Context, playerID int64, tower string, reward string, delta TowerRewardDelta) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ensureTowerLastReward(tower)[playerID] = reward
	r.ensureTowerLastRewardDelta(tower)[playerID] = delta
	return nil
}

func (r *MemoryRepository) ensureTowerProgress(tower string) map[int64]int {
	if _, ok := r.progress[tower]; !ok {
		r.progress[tower] = make(map[int64]int)
	}
	return r.progress[tower]
}

func (r *MemoryRepository) ensureTowerLastReward(tower string) map[int64]string {
	if _, ok := r.lastReward[tower]; !ok {
		r.lastReward[tower] = make(map[int64]string)
	}
	return r.lastReward[tower]
}

func (r *MemoryRepository) ensureTowerLastRewardDelta(tower string) map[int64]TowerRewardDelta {
	if _, ok := r.lastRewardDelta[tower]; !ok {
		r.lastRewardDelta[tower] = make(map[int64]TowerRewardDelta)
	}
	return r.lastRewardDelta[tower]
}

func buildTowerStatus(tower string, currentFloor int, lastReward string, lastRewardDelta TowerRewardDelta) TowerStatus {
	remaining := 5 - currentFloor
	if remaining < 0 {
		remaining = 0
	}

	status := TowerStatus{
		Tower:               tower,
		CurrentFloor:        currentFloor,
		MaxFloor:            10,
		RemainingChallenges: remaining,
		LastReward:          lastReward,
		LastRewardDelta:     lastRewardDelta,
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

	status := buildTowerStatus(tower, floor, state.LastReward[tower], state.LastRewardDelta[tower])
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
		return buildTowerStatus(tower, 0, "", TowerRewardDelta{})
	}
	return buildTowerStatus(tower, state.Progress[tower], state.LastReward[tower], state.LastRewardDelta[tower])
}

func (r *MySQLRepository) SaveLastReward(ctx context.Context, playerID int64, tower string, reward string, delta TowerRewardDelta) error {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return err
	}
	state.LastReward[tower] = reward
	state.LastRewardDelta[tower] = delta
	return r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) loadState(ctx context.Context, playerID int64) (mysqlState, error) {
	state := mysqlState{
		Progress:        map[string]int{"pagoda": 0, "spirit": 0},
		LastReward:      map[string]string{"pagoda": "", "spirit": ""},
		LastRewardDelta: map[string]TowerRewardDelta{"pagoda": {}, "spirit": {}},
	}
	ok, err := r.store.LoadModuleState(ctx, playerID, mysqlModuleName, &state)
	if err != nil {
		return mysqlState{}, err
	}
	if !ok || state.Progress == nil {
		state.Progress = map[string]int{"pagoda": 0, "spirit": 0}
	}
	if state.LastReward == nil {
		state.LastReward = map[string]string{"pagoda": "", "spirit": ""}
	}
	if state.LastRewardDelta == nil {
		state.LastRewardDelta = map[string]TowerRewardDelta{"pagoda": {}, "spirit": {}}
	}
	if _, ok := state.Progress["pagoda"]; !ok {
		state.Progress["pagoda"] = 0
	}
	if _, ok := state.Progress["spirit"]; !ok {
		state.Progress["spirit"] = 0
	}
	if _, ok := state.LastReward["pagoda"]; !ok {
		state.LastReward["pagoda"] = ""
	}
	if _, ok := state.LastReward["spirit"]; !ok {
		state.LastReward["spirit"] = ""
	}
	if _, ok := state.LastRewardDelta["pagoda"]; !ok {
		state.LastRewardDelta["pagoda"] = TowerRewardDelta{}
	}
	if _, ok := state.LastRewardDelta["spirit"]; !ok {
		state.LastRewardDelta["spirit"] = TowerRewardDelta{}
	}
	return state, nil
}

func (r *MySQLRepository) saveState(ctx context.Context, playerID int64, state mysqlState) error {
	if state.Progress == nil {
		state.Progress = map[string]int{"pagoda": 0, "spirit": 0}
	}
	if state.LastReward == nil {
		state.LastReward = map[string]string{"pagoda": "", "spirit": ""}
	}
	if state.LastRewardDelta == nil {
		state.LastRewardDelta = map[string]TowerRewardDelta{"pagoda": {}, "spirit": {}}
	}
	return r.store.SaveModuleState(ctx, playerID, mysqlModuleName, state)
}
