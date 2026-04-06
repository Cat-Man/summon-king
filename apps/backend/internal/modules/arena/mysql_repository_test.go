package arena

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMySQLRepository_PersistsArenaRecordAndRefreshVersion(t *testing.T) {
	store := &fakeModuleStateStore{}
	repo := NewMySQLRepository(store)

	record, err := repo.GetDailyRecord(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected default record success, got %v", err)
	}
	if record.CurrentStreak != 0 || record.LastWin {
		t.Fatalf("expected default arena record, got %+v", record)
	}

	if err := repo.SetCurrentStreak(context.Background(), 1001, 4); err != nil {
		t.Fatalf("expected set streak success, got %v", err)
	}
	if err := repo.RecordBattleResult(context.Background(), 1001, false); err != nil {
		t.Fatalf("expected record result success, got %v", err)
	}

	version, err := repo.IncrementRefreshVersion(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected refresh version success, got %v", err)
	}
	if version != 1 {
		t.Fatalf("expected refresh version 1, got %d", version)
	}

	record, err = repo.GetDailyRecord(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected persisted record success, got %v", err)
	}
	if record.CurrentStreak != 0 || record.LastWin {
		t.Fatalf("expected reset record after lose, got %+v", record)
	}

	version, err = repo.GetRefreshVersion(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected get refresh version success, got %v", err)
	}
	if version != 1 {
		t.Fatalf("expected persisted refresh version 1, got %d", version)
	}
}

type arenaStatePayload struct {
	Record         DailyRecord `json:"record"`
	RefreshVersion int         `json:"refresh_version"`
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
