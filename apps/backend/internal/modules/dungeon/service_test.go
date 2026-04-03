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

func newTestDungeonService(t *testing.T) *Service {
	t.Helper()
	repo := NewMemoryRepository()
	return NewService(repo)
}
