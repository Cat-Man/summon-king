package tower

import (
	"context"
	"sync"
)

type Repository interface {
	StartChallenge(ctx context.Context, playerID int64, tower string) (TowerResult, error)
	GetInfo(ctx context.Context, tower string) TowerInfo
}

type MemoryRepository struct {
	mu       sync.Mutex
	progress map[int64]int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		progress: make(map[int64]int),
	}
}

func (r *MemoryRepository) StartChallenge(_ context.Context, playerID int64, tower string) (TowerResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	floor := r.progress[playerID]
	floor++
	r.progress[playerID] = floor
	reward := "chest"
	if tower == "spirit" {
		reward = "spirit shard"
	}
	return TowerResult{
		PlayerID: playerID,
		Floor:    floor,
		Reward:   reward,
	}, nil
}

func (r *MemoryRepository) GetInfo(_ context.Context, _ string) TowerInfo {
	return TowerInfo{
		Level:    "Calm",
		MaxFloor: 10,
	}
}
