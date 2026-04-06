package pet

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMySQLRepository_PersistsTeamChanges(t *testing.T) {
	store := &fakeModuleStateStore{}
	repo := NewMySQLRepository(store)

	team, err := repo.GetBattleTeam(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected default team success, got %v", err)
	}
	if len(team.Pets) != 1 || team.Pets[0].PetID != 10011 {
		t.Fatalf("expected default active pet 10011, got %+v", team.Pets)
	}

	if err := repo.SaveTeam(context.Background(), 1001, []int64{10011, 10012}); err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}
	team, err = repo.GetBattleTeam(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected persisted team success, got %v", err)
	}
	if len(team.Pets) != 2 || team.Pets[1].PetID != 10012 {
		t.Fatalf("expected persisted dual team, got %+v", team.Pets)
	}

	if err := repo.SetMainPet(context.Background(), 1001, 10012); err != nil {
		t.Fatalf("expected set main pet success, got %v", err)
	}
	team, err = repo.GetBattleTeam(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected reordered team success, got %v", err)
	}
	if team.Pets[0].PetID != 10012 {
		t.Fatalf("expected main pet 10012, got %+v", team.Pets)
	}
}

func TestMySQLRepository_GrantActiveTeamExperience(t *testing.T) {
	store := &fakeModuleStateStore{}
	repo := NewMySQLRepository(store)

	if err := repo.SaveTeam(context.Background(), 1001, []int64{10011, 10012}); err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}

	team, err := repo.GrantActiveTeamExperience(context.Background(), 1001, 120)
	if err != nil {
		t.Fatalf("expected grant exp success, got %v", err)
	}
	if team.Pets[0].Level != 2 || team.Pets[0].Exp != 20 || team.Pets[0].Power != 144 {
		t.Fatalf("expected leveled main pet, got %+v", team.Pets[0])
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
