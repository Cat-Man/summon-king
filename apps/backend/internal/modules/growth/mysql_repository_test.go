package growth

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMySQLRepository_PersistsWalletProgress(t *testing.T) {
	store := &fakeModuleStateStore{}
	repo := NewMySQLRepository(store)

	wallet, err := repo.GetWallet(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected wallet load success, got %v", err)
	}
	if wallet.SpiritPower != 100 || wallet.SpiritFreeWash != 3 {
		t.Fatalf("expected default wallet, got %+v", wallet)
	}

	free, err := repo.BeginFreeWash(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected free wash success, got %v", err)
	}
	if !free {
		t.Fatal("expected first wash to consume free count")
	}

	if err := repo.CommitWash(context.Background(), 1001, -5); err != nil {
		t.Fatalf("expected commit wash success, got %v", err)
	}
	if err := repo.UpdateSpiritPower(context.Background(), 1001, 12); err != nil {
		t.Fatalf("expected spirit power update success, got %v", err)
	}

	upgradedBone, err := repo.UpgradeBoneLevel(context.Background(), 1001, 1)
	if err != nil {
		t.Fatalf("expected bone upgrade success, got %v", err)
	}
	if upgradedBone.BoneLevel != 2 {
		t.Fatalf("expected bone level 2, got %d", upgradedBone.BoneLevel)
	}

	upgradedSoul, err := repo.UpgradeSoulPower(context.Background(), 1001, 2)
	if err != nil {
		t.Fatalf("expected soul upgrade success, got %v", err)
	}
	if upgradedSoul.SoulPieces != 2 {
		t.Fatalf("expected soul pieces 2, got %d", upgradedSoul.SoulPieces)
	}

	finalWallet, err := repo.GetWallet(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected final wallet success, got %v", err)
	}
	if finalWallet.SpiritPower != 107 || finalWallet.SpiritBonusPower != 112 || finalWallet.SpiritFreeWash != 2 || finalWallet.BoneLevel != 2 || finalWallet.SoulPieces != 2 {
		t.Fatalf("expected persisted wallet progress, got %+v", finalWallet)
	}
}

func TestMySQLRepository_PersistsManorPlots(t *testing.T) {
	store := &fakeModuleStateStore{}
	repo := NewMySQLRepository(store)

	plots, err := repo.GetManorPlots(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected get manor plots success, got %v", err)
	}
	if len(plots) != 2 {
		t.Fatalf("expected 2 plots, got %d", len(plots))
	}

	planted, err := repo.PlantManor(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected plant manor success, got %v", err)
	}
	if planted[0].State != "成长中" || planted[1].State != "成长中" {
		t.Fatalf("expected planted plots to be 成长中, got %+v", planted)
	}

	harvested, err := repo.HarvestManor(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected harvest manor success, got %v", err)
	}
	if harvested[0].State != "冷却中" || harvested[1].State != "冷却中" {
		t.Fatalf("expected harvested plots to be 冷却中, got %+v", harvested)
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
