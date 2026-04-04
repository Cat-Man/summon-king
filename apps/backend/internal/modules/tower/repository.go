package tower

import (
	"context"
	"sync"
)

type Repository interface {
	StartChallenge(ctx context.Context, playerID int64, tower string) (TowerResult, error)
	GetStatus(ctx context.Context, playerID int64, tower string) TowerStatus
}

type MemoryRepository struct {
	mu       sync.Mutex
	progress map[string]map[int64]int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		progress: map[string]map[int64]int{
			"pagoda": {},
			"spirit": {},
		},
	}
}

func (r *MemoryRepository) StartChallenge(_ context.Context, playerID int64, tower string) (TowerResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	progress := r.ensureTowerProgress(tower)
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
