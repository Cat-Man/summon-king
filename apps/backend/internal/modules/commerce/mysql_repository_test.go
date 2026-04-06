package commerce

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMySQLRepository_PersistsSigninAndVIPState(t *testing.T) {
	store := &fakeModuleStateStore{}
	repo := NewMySQLRepository(store)

	state, err := repo.GetState(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected default commerce state success, got %v", err)
	}
	if state.VIPLevel != 0 || state.SigninStreak != 0 {
		t.Fatalf("expected default commerce state, got %+v", state)
	}

	if err := repo.SaveState(context.Background(), 1001, State{
		SigninLastDate: "2026-04-08",
		SigninStreak:   3,
		VIPLevel:       1,
		VIPDailyDate:   "2026-04-08",
	}); err != nil {
		t.Fatalf("expected save commerce state success, got %v", err)
	}

	loaded, err := repo.GetState(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected load commerce state success, got %v", err)
	}
	if loaded.SigninLastDate != "2026-04-08" || loaded.SigninStreak != 3 || loaded.VIPLevel != 1 || loaded.VIPDailyDate != "2026-04-08" {
		t.Fatalf("expected persisted commerce state, got %+v", loaded)
	}
}

type fakeModuleStateStore struct {
	states map[string]map[int64]json.RawMessage
}

func (s *fakeModuleStateStore) LoadModuleState(_ context.Context, playerID int64, module string, target any) (bool, error) {
	if s.states == nil {
		return false, nil
	}
	players, ok := s.states[module]
	if !ok {
		return false, nil
	}
	raw, ok := players[playerID]
	if !ok {
		return false, nil
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return false, err
	}
	return true, nil
}

func (s *fakeModuleStateStore) SaveModuleState(_ context.Context, playerID int64, module string, value any) error {
	if s.states == nil {
		s.states = make(map[string]map[int64]json.RawMessage)
	}
	if _, ok := s.states[module]; !ok {
		s.states[module] = make(map[int64]json.RawMessage)
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.states[module][playerID] = raw
	return nil
}
