package jobs

import (
	"context"
	"sync"
	"time"
)

type DailyState struct {
	PlayerID           int64 `json:"player_id"`
	SigninStatus       int   `json:"signin_status"`
	VIPChestClaimed    bool  `json:"vip_chest_claimed"`
	ArenaRewardClaimed bool  `json:"arena_reward_claimed"`
}

type CultivationTask struct {
	PlayerID  int64     `json:"player_id"`
	ReadyAt   time.Time `json:"ready_at"`
	Claimable bool      `json:"claimable"`
}

type ManorPlot struct {
	PlotID   string    `json:"plot_id"`
	MatureAt time.Time `json:"mature_at"`
	Status   string    `json:"status"`
}

type AllianceWarRound struct {
	RoundID  string    `json:"round_id"`
	Phase    string    `json:"phase"`
	LockAt   time.Time `json:"lock_at"`
	BattleAt time.Time `json:"battle_at"`
	ResultAt time.Time `json:"result_at"`
}

type MemoryStateStore struct {
	mu           sync.Mutex
	DailyStates  map[int64]DailyState
	Cultivations map[int64]CultivationTask
	ManorPlots   map[string]ManorPlot
	AllianceWars map[string]AllianceWarRound
}

func NewMemoryStateStore() *MemoryStateStore {
	return &MemoryStateStore{
		DailyStates:  map[int64]DailyState{},
		Cultivations: map[int64]CultivationTask{},
		ManorPlots:   map[string]ManorPlot{},
		AllianceWars: map[string]AllianceWarRound{},
	}
}

func newTestDeps(_ interface{ Helper() }) *MemoryStateStore {
	return NewMemoryStateStore()
}

type DailyResetJob struct{ store *MemoryStateStore }

func NewDailyResetJob(store *MemoryStateStore) *DailyResetJob { return &DailyResetJob{store: store} }

func (j *DailyResetJob) Run(_ context.Context) error {
	j.store.mu.Lock()
	defer j.store.mu.Unlock()
	for playerID, state := range j.store.DailyStates {
		state.SigninStatus = 0
		state.VIPChestClaimed = false
		state.ArenaRewardClaimed = false
		j.store.DailyStates[playerID] = state
	}
	return nil
}
