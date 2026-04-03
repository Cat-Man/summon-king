package dungeon

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrDungeonRunNotFound  = errors.New("dungeon run not found")
	ErrCultivationNotFound = errors.New("cultivation not found")
)

type Repository interface {
	GetWorldMap(ctx context.Context) (WorldMap, error)
	TeleportToCity(ctx context.Context, cityID int64) (MapCity, error)
	EnterDungeon(ctx context.Context, playerID, dungeonID int64) (DungeonRun, error)
	RollDungeonDice(ctx context.Context, playerID int64) (DungeonRun, error)
	GetDungeonRun(ctx context.Context, playerID int64) (DungeonRun, error)
	StartCultivation(ctx context.Context, playerID int64) (CultivationStatus, error)
	ClaimCultivation(ctx context.Context, playerID int64) (CultivationStatus, error)
	GetCultivation(ctx context.Context, playerID int64) (CultivationStatus, error)
}

type MemoryRepository struct {
	mu          sync.Mutex
	worldMap    WorldMap
	runs        map[int64]DungeonRun
	cultivation map[int64]CultivationStatus
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		worldMap: WorldMap{
			Name: "玄境",
			Cities: []MapCity{
				{CityID: 1, Name: "晨曦城", Region: "东境", LocX: 110.5, LocY: 220.4},
				{CityID: 2, Name: "霞光堡", Region: "南境", LocX: 190.8, LocY: 180.1},
			},
		},
		runs:        make(map[int64]DungeonRun),
		cultivation: make(map[int64]CultivationStatus),
	}
}

func (r *MemoryRepository) GetWorldMap(_ context.Context) (WorldMap, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.worldMap, nil
}

func (r *MemoryRepository) TeleportToCity(_ context.Context, cityID int64) (MapCity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, city := range r.worldMap.Cities {
		if city.CityID == cityID {
			return city, nil
		}
	}
	return MapCity{}, errors.New("city not found")
}

func (r *MemoryRepository) EnterDungeon(_ context.Context, playerID, dungeonID int64) (DungeonRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	run := DungeonRun{
		PlayerID:     playerID,
		DungeonID:    dungeonID,
		RemainDice:   15,
		CurrentFloor: 1,
		Status:       "ongoing",
		StartedAt:    time.Now(),
	}
	r.runs[playerID] = run
	return run, nil
}

func (r *MemoryRepository) RollDungeonDice(_ context.Context, playerID int64) (DungeonRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	run, ok := r.runs[playerID]
	if !ok {
		return DungeonRun{}, ErrDungeonRunNotFound
	}
	if run.RemainDice <= 0 {
		run.Status = "exhausted"
		r.runs[playerID] = run
		return run, nil
	}
	run.RemainDice--
	run.CurrentFloor++
	if run.CurrentFloor%5 == 0 {
		run.Status = "boss"
	} else {
		run.Status = "ongoing"
	}
	r.runs[playerID] = run
	return run, nil
}

func (r *MemoryRepository) GetDungeonRun(_ context.Context, playerID int64) (DungeonRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	run, ok := r.runs[playerID]
	if !ok {
		return DungeonRun{}, ErrDungeonRunNotFound
	}
	return run, nil
}

func (r *MemoryRepository) StartCultivation(_ context.Context, playerID int64) (CultivationStatus, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	status := CultivationStatus{
		PlayerID:    playerID,
		SpiritPower: 0,
		State:       "cultivating",
		StartAt:     time.Now(),
		ClaimableAt: time.Now().Add(1 * time.Hour),
	}
	r.cultivation[playerID] = status
	return status, nil
}

func (r *MemoryRepository) ClaimCultivation(_ context.Context, playerID int64) (CultivationStatus, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	status, ok := r.cultivation[playerID]
	if !ok {
		return CultivationStatus{}, errors.New("cultivation not started")
	}
	status.SpiritPower += 10
	status.State = "idle"
	status.ClaimableAt = time.Time{}
	r.cultivation[playerID] = status
	return status, nil
}

func (r *MemoryRepository) GetCultivation(_ context.Context, playerID int64) (CultivationStatus, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	status, ok := r.cultivation[playerID]
	if !ok {
		return CultivationStatus{}, ErrCultivationNotFound
	}
	return status, nil
}
