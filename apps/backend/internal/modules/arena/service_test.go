package arena

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
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
	if result.BattleResult.BattleType != "arena" {
		t.Fatalf("expected arena battle summary, got %s", result.BattleResult.BattleType)
	}
	if result.BattleResult.Result != "success" {
		t.Fatalf("expected success battle result, got %s", result.BattleResult.Result)
	}
	if result.BattleResult.BattleNo == "" {
		t.Fatal("expected arena battle number")
	}
	if result.BattleResult.WinnerSide != "attacker" {
		t.Fatalf("expected arena winner side attacker, got %s", result.BattleResult.WinnerSide)
	}
	if result.BattleResult.AttackerPower <= 0 {
		t.Fatalf("expected attacker power > 0, got %d", result.BattleResult.AttackerPower)
	}
}

func TestArena_LoseBattleReturnsFailedBattleSummary(t *testing.T) {
	ctx := context.Background()
	svc := newTestArenaService(t)
	playerID := int64(3004)

	result, err := svc.RecordBattleResult(ctx, playerID, false)
	if err != nil {
		t.Fatalf("expected record battle result success, got %v", err)
	}
	if result.BattleResult.Result != "fail" {
		t.Fatalf("expected fail battle result, got %s", result.BattleResult.Result)
	}
	if result.BattleResult.BattleNo == "" {
		t.Fatal("expected arena battle number")
	}
	if result.BattleResult.WinnerSide != "defender" {
		t.Fatalf("expected arena winner side defender, got %s", result.BattleResult.WinnerSide)
	}
}

func TestArena_UsesInjectedBattleTeamReader(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository())
	if _, err := petSvc.SetMainPet(ctx, 3101, 31012); err != nil {
		t.Fatalf("expected set main pet success, got %v", err)
	}

	svc := NewService(
		NewMemoryRepository(),
		asset.NewService(growthRepo),
		WithBattleTeamReader(petSvc),
	)

	result, err := svc.RecordBattleResult(ctx, 3101, true)
	if err != nil {
		t.Fatalf("expected record battle success, got %v", err)
	}
	if result.BattleResult.AttackerPower != 156 {
		t.Fatalf("expected attacker power 156, got %d", result.BattleResult.AttackerPower)
	}
}

func TestArena_UsesSavedTeamPower(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository())
	if _, err := petSvc.SaveTeam(ctx, 3201, []int64{32011, 32012}); err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}

	svc := NewService(
		NewMemoryRepository(),
		asset.NewService(growthRepo),
		WithBattleTeamReader(petSvc),
	)

	result, err := svc.RecordBattleResult(ctx, 3201, true)
	if err != nil {
		t.Fatalf("expected record battle success, got %v", err)
	}
	if result.BattleResult.AttackerPower != 276 {
		t.Fatalf("expected attacker power 276, got %d", result.BattleResult.AttackerPower)
	}
}

func TestArena_UsesGrowthBoostedTeamPower(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	if _, err := growthRepo.UpgradeBoneLevel(ctx, 3301, 1); err != nil {
		t.Fatalf("expected bone upgrade success, got %v", err)
	}
	petSvc := pet.NewService(pet.NewMemoryRepository(), pet.WithGrowthReader(growthRepo))

	svc := NewService(
		NewMemoryRepository(),
		asset.NewService(growthRepo),
		WithBattleTeamReader(petSvc),
	)

	result, err := svc.RecordBattleResult(ctx, 3301, true)
	if err != nil {
		t.Fatalf("expected record battle success, got %v", err)
	}
	if result.BattleResult.AttackerPower != 144 {
		t.Fatalf("expected attacker power 144 with bone bonus, got %d", result.BattleResult.AttackerPower)
	}
}

func TestArena_BattleSummaryUsesPreRewardTeamPower(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository(), pet.WithGrowthReader(growthRepo))

	svc := NewService(
		NewMemoryRepository(),
		asset.NewService(growthRepo),
		WithBattleTeamReader(petSvc),
	)

	result, err := svc.RecordBattleResult(ctx, 3401, true)
	if err != nil {
		t.Fatalf("expected record battle success, got %v", err)
	}
	if result.BattleResult.AttackerPower != 120 {
		t.Fatalf("expected attacker power 120 before arena reward write-back, got %d", result.BattleResult.AttackerPower)
	}
}
