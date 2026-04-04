package arena

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

func TestArena_FailResetsCurrentStreak(t *testing.T) {
	ctx := context.Background()
	svc := newTestArenaService(t)
	playerID := int64(1001)

	if err := svc.SetCurrentStreak(ctx, playerID, 10); err != nil {
		t.Fatalf("expected SetCurrentStreak success, got %v", err)
	}

	if _, err := svc.RecordBattleResult(ctx, playerID, false); err != nil {
		t.Fatalf("expected RecordBattleResult success, got %v", err)
	}

	rec, err := svc.GetDailyRecord(ctx, playerID)
	if err != nil {
		t.Fatalf("expected GetDailyRecord success, got %v", err)
	}
	if rec.CurrentStreak != 0 {
		t.Fatalf("expected streak reset to 0, got %d", rec.CurrentStreak)
	}
}

func newTestArenaService(t *testing.T) *Service {
	t.Helper()
	growthRepo := growth.NewMemoryRepository()
	return NewService(NewMemoryRepository(), asset.NewService(growthRepo))
}

func TestArena_NewPlayerStartsWithEmptyRecord(t *testing.T) {
	ctx := context.Background()
	svc := newTestArenaService(t)
	playerID := int64(2002)

	rec, err := svc.GetDailyRecord(ctx, playerID)
	if err != nil {
		t.Fatalf("expected empty record success, got %v", err)
	}
	if rec.PlayerID != playerID {
		t.Fatalf("expected player id %d, got %d", playerID, rec.PlayerID)
	}
	if rec.CurrentStreak != 0 {
		t.Fatalf("expected current streak 0, got %d", rec.CurrentStreak)
	}
}

func TestArena_WinBattleAppliesRewardAndReturnsSnapshot(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
	playerID := int64(3003)

	before, err := growthRepo.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet fetch success, got %v", err)
	}

	result, err := svc.RecordBattleResult(ctx, playerID, true)
	if err != nil {
		t.Fatalf("expected record battle result success, got %v", err)
	}

	if result.RewardDelta.SpiritPower == 0 {
		t.Fatal("expected spirit power reward on win")
	}
	if result.WalletSnapshot.SpiritPower <= before.SpiritPower {
		t.Fatalf("expected spirit power to increase, before=%d after=%d", before.SpiritPower, result.WalletSnapshot.SpiritPower)
	}
	if result.Record.CurrentStreak != 1 {
		t.Fatalf("expected current streak 1, got %d", result.Record.CurrentStreak)
	}
}
