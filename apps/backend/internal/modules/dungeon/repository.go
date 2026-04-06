package dungeon

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

var (
	ErrDungeonRunNotFound  = errors.New("dungeon run not found")
	ErrCultivationNotFound = errors.New("cultivation not found")
	ErrCultivationNotReady = errors.New("cultivation not ready")
	ErrDungeonLocked       = errors.New("dungeon locked")
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

type MySQLRepository struct {
	store mysqlstore.ModuleStateStore
}

type mysqlState struct {
	Run         *DungeonRun        `json:"run"`
	Cultivation *CultivationStatus `json:"cultivation"`
}

const mysqlModuleName = "dungeon"

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		worldMap:    defaultWorldMap(),
		runs:        make(map[int64]DungeonRun),
		cultivation: make(map[int64]CultivationStatus),
	}
}

func NewMySQLRepository(store mysqlstore.ModuleStateStore) *MySQLRepository {
	return &MySQLRepository{store: store}
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
	now := currentTime()
	status := CultivationStatus{
		PlayerID:    playerID,
		SpiritPower: 0,
		State:       "cultivating",
		StartAt:     now,
		ClaimableAt: now.Add(1 * time.Hour),
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

func (r *MySQLRepository) GetWorldMap(_ context.Context) (WorldMap, error) {
	return defaultWorldMap(), nil
}

func (r *MySQLRepository) TeleportToCity(_ context.Context, cityID int64) (MapCity, error) {
	for _, city := range defaultWorldMap().Cities {
		if city.CityID == cityID {
			return city, nil
		}
	}
	return MapCity{}, errors.New("city not found")
}

func (r *MySQLRepository) EnterDungeon(ctx context.Context, playerID, dungeonID int64) (DungeonRun, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return DungeonRun{}, err
	}
	run := DungeonRun{
		PlayerID:     playerID,
		DungeonID:    dungeonID,
		RemainDice:   15,
		CurrentFloor: 1,
		Status:       "ongoing",
		StartedAt:    time.Now(),
	}
	state.Run = &run
	return run, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) RollDungeonDice(ctx context.Context, playerID int64) (DungeonRun, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return DungeonRun{}, err
	}
	if state.Run == nil {
		return DungeonRun{}, ErrDungeonRunNotFound
	}

	run := *state.Run
	if run.RemainDice <= 0 {
		run.Status = "exhausted"
		state.Run = &run
		return run, r.saveState(ctx, playerID, state)
	}
	run.RemainDice--
	run.CurrentFloor++
	if run.CurrentFloor%5 == 0 {
		run.Status = "boss"
	} else {
		run.Status = "ongoing"
	}
	state.Run = &run
	return run, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) GetDungeonRun(ctx context.Context, playerID int64) (DungeonRun, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return DungeonRun{}, err
	}
	if state.Run == nil {
		return DungeonRun{}, ErrDungeonRunNotFound
	}
	return *state.Run, nil
}

func (r *MySQLRepository) StartCultivation(ctx context.Context, playerID int64) (CultivationStatus, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, err
	}
	now := currentTime()
	status := CultivationStatus{
		PlayerID:    playerID,
		SpiritPower: 0,
		State:       "cultivating",
		StartAt:     now,
		ClaimableAt: now.Add(1 * time.Hour),
	}
	state.Cultivation = &status
	return status, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) ClaimCultivation(ctx context.Context, playerID int64) (CultivationStatus, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, err
	}
	if state.Cultivation == nil {
		return CultivationStatus{}, errors.New("cultivation not started")
	}

	status := *state.Cultivation
	status.SpiritPower += 10
	status.State = "idle"
	status.ClaimableAt = time.Time{}
	state.Cultivation = &status
	return status, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) GetCultivation(ctx context.Context, playerID int64) (CultivationStatus, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, err
	}
	if state.Cultivation == nil {
		return CultivationStatus{}, ErrCultivationNotFound
	}
	return *state.Cultivation, nil
}

func (r *MySQLRepository) loadState(ctx context.Context, playerID int64) (mysqlState, error) {
	var state mysqlState
	_, err := r.store.LoadModuleState(ctx, playerID, mysqlModuleName, &state)
	if err != nil {
		return mysqlState{}, err
	}
	return state, nil
}

func (r *MySQLRepository) saveState(ctx context.Context, playerID int64, state mysqlState) error {
	return r.store.SaveModuleState(ctx, playerID, mysqlModuleName, state)
}

func defaultWorldMap() WorldMap {
	return WorldMap{
		Name: "玄境",
		Cities: []MapCity{
			{
				CityID: 1,
				Name:   "晨曦城",
				Region: "东境",
				LocX:   110.5,
				LocY:   220.4,
				Dungeons: []MapDungeon{
					{DungeonID: 1, DungeonName: "妖窟试炼", UnlockSpiritPower: 0},
					{
						DungeonID:         2,
						DungeonName:       "寒渊裂隙",
						UnlockSpiritPower: 120,
						UnlockBoneLevel:   2,
						UnlockSoulPieces:  3,
					},
				},
			},
			{
				CityID: 2,
				Name:   "霞光堡",
				Region: "南境",
				LocX:   190.8,
				LocY:   180.1,
				Dungeons: []MapDungeon{
					{
						DungeonID:         2,
						DungeonName:       "寒渊裂隙",
						UnlockSpiritPower: 120,
						UnlockBoneLevel:   2,
						UnlockSoulPieces:  3,
					},
					{DungeonID: 1, DungeonName: "妖窟试炼", UnlockSpiritPower: 0},
				},
			},
		},
	}
}
