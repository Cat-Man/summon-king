package dungeon

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

func TestEnterDungeon_GivesInitialDice(t *testing.T) {
	svc := newTestDungeonService(t)
	ctx := context.Background()
	playerID := int64(1002)
	dungeonID := int64(1)

	run, err := svc.EnterDungeon(ctx, playerID, dungeonID)
	if err != nil {
		t.Fatalf("expected enter dungeon to succeed, got %v", err)
	}
	if run.RemainDice != 15 {
		t.Fatalf("expected remain dice 15, got %d", run.RemainDice)
	}
}

func TestEnterDungeon_ReturnsWalletSnapshot(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), growthRepo)
	playerID := int64(1003)

	run, err := svc.EnterDungeon(ctx, playerID, 1)
	if err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	if run.WalletSnapshot.PlayerID != playerID {
		t.Fatalf("expected wallet snapshot player id %d, got %d", playerID, run.WalletSnapshot.PlayerID)
	}
	if run.WalletSnapshot.SpiritPower != 100 {
		t.Fatalf("expected wallet snapshot spirit 100, got %d", run.WalletSnapshot.SpiritPower)
	}
}

func TestClaimCultivation_UpdatesSpiritWallet(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), growthRepo)
	playerID := int64(1009)

	if _, err := svc.StartCultivation(ctx, playerID); err != nil {
		t.Fatalf("expected cultivation start success, got %v", err)
	}
	if _, err := svc.ClaimCultivation(ctx, playerID); err != nil {
		t.Fatalf("expected cultivation claim success, got %v", err)
	}

	wallet, err := growthRepo.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet load success, got %v", err)
	}
	if wallet.SpiritPower != 110 {
		t.Fatalf("expected spirit power 110 after claim, got %d", wallet.SpiritPower)
	}
}

func TestRollDice_UpdatesSpiritWallet(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), growthRepo)
	playerID := int64(1010)

	if _, err := svc.EnterDungeon(ctx, playerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}

	if _, err := svc.RollDice(ctx, playerID); err != nil {
		t.Fatalf("expected roll dice success, got %v", err)
	}

	wallet, err := growthRepo.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet load success, got %v", err)
	}
	if wallet.SpiritPower != 105 {
		t.Fatalf("expected spirit power 105 after roll reward, got %d", wallet.SpiritPower)
	}
}

func TestRollDice_BossFloorAddsSoulReward(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), growthRepo)
	playerID := int64(1011)

	if _, err := svc.EnterDungeon(ctx, playerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}

	for i := 0; i < 4; i++ {
		if _, err := svc.RollDice(ctx, playerID); err != nil {
			t.Fatalf("expected roll dice success, got %v", err)
		}
	}

	wallet, err := growthRepo.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet load success, got %v", err)
	}
	if wallet.SoulPieces != 1 {
		t.Fatalf("expected soul pieces 1 after boss reward, got %d", wallet.SoulPieces)
	}
}

func TestRollDice_ReturnsRewardAndWalletSnapshot(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), growthRepo)
	playerID := int64(1012)

	if _, err := svc.EnterDungeon(ctx, playerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}

	run, err := svc.RollDice(ctx, playerID)
	if err != nil {
		t.Fatalf("expected roll dice success, got %v", err)
	}

	if run.LastReward.SpiritPower != 5 {
		t.Fatalf("expected spirit reward 5, got %d", run.LastReward.SpiritPower)
	}
	if run.LastReward.SoulPieces != 0 {
		t.Fatalf("expected soul reward 0, got %d", run.LastReward.SoulPieces)
	}
	if run.WalletSnapshot.SpiritPower != 105 {
		t.Fatalf("expected wallet snapshot spirit 105, got %d", run.WalletSnapshot.SpiritPower)
	}
}

func newTestDungeonService(t *testing.T) *Service {
	t.Helper()
	repo := NewMemoryRepository()
	return NewService(repo)
}
