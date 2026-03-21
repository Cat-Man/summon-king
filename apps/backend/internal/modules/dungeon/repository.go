package dungeon

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrDungeonRunNotFound    = errors.New("dungeon run not found")
	ErrNoDiceRemaining       = errors.New("no dice remaining")
	ErrCultivationNotFound   = errors.New("cultivation record not found")
	ErrCultivationUnfinished = errors.New("cultivation not finished")
	ErrCultivationClaimed    = errors.New("cultivation already claimed")
)

type Repository interface {
	GetWorldMap(ctx context.Context, playerID int64) (WorldMap, error)
	Teleport(ctx context.Context, playerID int64, cityID string) (City, error)
	CreateDungeonRun(ctx context.Context, playerID, dungeonID int64, now time.Time) (DungeonRun, error)
	GetDungeonRun(ctx context.Context, runID string) (DungeonRun, error)
	SaveDungeonRun(ctx context.Context, run DungeonRun) error
	StartCultivation(ctx context.Context, playerID int64, mapID string, hours int, now time.Time) (CultivationRecord, error)
	GetCultivation(ctx context.Context, playerID int64, recordID string) (CultivationRecord, error)
	SaveCultivation(ctx context.Context, record CultivationRecord) error
}

type MemoryRepository struct {
	mu                sync.Mutex
	nextRunID         int64
	nextCultivationID int64
	runs              map[string]DungeonRun
	cultivations      map[string]CultivationRecord
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextRunID:         1,
		nextCultivationID: 1,
		runs:              make(map[string]DungeonRun),
		cultivations:      make(map[string]CultivationRecord),
	}
}

func (r *MemoryRepository) GetWorldMap(_ context.Context, playerID int64) (WorldMap, error) {
	return WorldMap{
		PlayerID: playerID,
		Region:   "东灵大陆",
		Cities: []City{
			{CityID: "qingyun", CityName: "青云城", Unlocked: true},
			{CityID: "fire-mountain", CityName: "火焰山", Unlocked: true},
			{CityID: "frozen-bay", CityName: "寒霜港", Unlocked: false},
		},
	}, nil
}

func (r *MemoryRepository) Teleport(ctx context.Context, playerID int64, cityID string) (City, error) {
	world, _ := r.GetWorldMap(ctx, playerID)
	for _, city := range world.Cities {
		if city.CityID == cityID {
			return city, nil
		}
	}
	return City{}, ErrDungeonRunNotFound
}

func (r *MemoryRepository) CreateDungeonRun(_ context.Context, playerID, dungeonID int64, now time.Time) (DungeonRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	run := DungeonRun{
		RunID:        fmt.Sprintf("run-%d", r.nextRunID),
		PlayerID:     playerID,
		DungeonID:    dungeonID,
		CurrentFloor: 1,
		RemainDice:   15,
		CreatedAt:    now,
	}
	r.nextRunID++
	r.runs[run.RunID] = run
	return run, nil
}

func (r *MemoryRepository) GetDungeonRun(_ context.Context, runID string) (DungeonRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if run, ok := r.runs[runID]; ok {
		return run, nil
	}
	return DungeonRun{}, ErrDungeonRunNotFound
}

func (r *MemoryRepository) SaveDungeonRun(_ context.Context, run DungeonRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runs[run.RunID] = run
	return nil
}

func (r *MemoryRepository) StartCultivation(_ context.Context, playerID int64, mapID string, hours int, now time.Time) (CultivationRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record := CultivationRecord{
		RecordID:     fmt.Sprintf("cult-%d", r.nextCultivationID),
		PlayerID:     playerID,
		MapID:        mapID,
		Hours:        hours,
		Status:       "running",
		RewardCoins:  int64(hours * 120),
		RewardPetExp: int64(hours * 80),
		StartedAt:    now,
		FinishedAt:   now.Add(time.Duration(hours) * time.Hour),
	}
	r.nextCultivationID++
	r.cultivations[record.RecordID] = record
	return record, nil
}

func (r *MemoryRepository) GetCultivation(_ context.Context, playerID int64, recordID string) (CultivationRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.cultivations[recordID]
	if !ok || record.PlayerID != playerID {
		return CultivationRecord{}, ErrCultivationNotFound
	}
	return record, nil
}

func (r *MemoryRepository) SaveCultivation(_ context.Context, record CultivationRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cultivations[record.RecordID] = record
	return nil
}
