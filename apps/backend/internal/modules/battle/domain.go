package battle

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrBattleNotFound = errors.New("battle not found")

type BattleInput struct {
	BattleNo string `json:"battle_no"`
	Left     []Unit `json:"left"`
	Right    []Unit `json:"right"`
}

type Unit struct {
	Name     string         `json:"name"`
	HP       int            `json:"hp"`
	MaxHP    int            `json:"max_hp"`
	Attack   int            `json:"attack"`
	Defense  int            `json:"defense"`
	Speed    int            `json:"speed"`
	Skills   []Skill        `json:"skills,omitempty"`
	Statuses []StatusEffect `json:"statuses,omitempty"`
}

type Skill struct {
	Name        string  `json:"name"`
	Enabled     bool    `json:"enabled"`
	TriggerRate float64 `json:"trigger_rate"`
	Multiplier  float64 `json:"multiplier"`
}

type BattleRound struct {
	Number  int            `json:"number"`
	Actions []BattleAction `json:"actions"`
}

type BattleAction struct {
	Actor          string `json:"actor"`
	Target         string `json:"target"`
	Damage         int    `json:"damage"`
	Skill          string `json:"skill,omitempty"`
	TargetHP       int    `json:"target_hp"`
	TargetDefeated bool   `json:"target_defeated"`
}

type BattleResult struct {
	BattleNo  string        `json:"battle_no"`
	Winner    string        `json:"winner"`
	CreatedAt time.Time     `json:"created_at"`
	Left      []Unit        `json:"left"`
	Right     []Unit        `json:"right"`
	Rounds    []BattleRound `json:"rounds"`
}

type Replay struct {
	BattleNo string        `json:"battle_no"`
	Winner   string        `json:"winner"`
	Rounds   []BattleRound `json:"rounds"`
}

type Store interface {
	Save(ctx context.Context, result BattleResult) error
	GetDetail(ctx context.Context, battleNo string) (BattleResult, error)
	GetReplay(ctx context.Context, battleNo string) (Replay, error)
}

type MemoryStore struct {
	mu      sync.RWMutex
	records map[string]BattleResult
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{records: make(map[string]BattleResult)}
}

func (s *MemoryStore) Save(_ context.Context, result BattleResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[result.BattleNo] = result
	return nil
}

func (s *MemoryStore) GetDetail(_ context.Context, battleNo string) (BattleResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.records[battleNo]
	if !ok {
		return BattleResult{}, ErrBattleNotFound
	}

	return result, nil
}

func (s *MemoryStore) GetReplay(ctx context.Context, battleNo string) (Replay, error) {
	result, err := s.GetDetail(ctx, battleNo)
	if err != nil {
		return Replay{}, err
	}

	return Replay{
		BattleNo: result.BattleNo,
		Winner:   result.Winner,
		Rounds:   result.Rounds,
	}, nil
}
