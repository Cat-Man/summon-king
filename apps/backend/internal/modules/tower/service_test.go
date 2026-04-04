package tower

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

func TestService_TracksProgressPerTower(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository(), asset.NewService(growth.NewMemoryRepository()))
	playerID := int64(1001)

	pagodaResult, err := svc.StartChallenge(ctx, playerID, "pagoda")
	if err != nil {
		t.Fatalf("expected pagoda challenge success, got %v", err)
	}
	spiritResult, err := svc.StartChallenge(ctx, playerID, "spirit")
	if err != nil {
		t.Fatalf("expected spirit challenge success, got %v", err)
	}

	if pagodaResult.Floor != 1 {
		t.Fatalf("expected pagoda first floor 1, got %d", pagodaResult.Floor)
	}
	if spiritResult.Floor != 1 {
		t.Fatalf("expected spirit first floor 1, got %d", spiritResult.Floor)
	}
}

func TestService_StartChallengeAppliesTowerRewardsToWallet(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
	playerID := int64(1002)

	before, err := growthRepo.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet fetch, got %v", err)
	}

	result, err := svc.StartChallenge(ctx, playerID, "pagoda")
	if err != nil {
		t.Fatalf("expected challenge success, got %v", err)
	}

	if result.RewardDelta.BoneLevel != 1 {
		t.Fatalf("expected bone delta 1, got %d", result.RewardDelta.BoneLevel)
	}

	if result.WalletSnapshot.BoneLevel != before.BoneLevel+1 {
		t.Fatalf("expected bone level %d, got %d", before.BoneLevel+1, result.WalletSnapshot.BoneLevel)
	}
	if result.BattleResult.BattleType != "tower" {
		t.Fatalf("expected tower battle summary, got %s", result.BattleResult.BattleType)
	}
	if result.BattleResult.Result != "success" {
		t.Fatalf("expected success battle result, got %s", result.BattleResult.Result)
	}
	if result.BattleResult.BattleNo == "" {
		t.Fatal("expected tower battle number")
	}
	if result.BattleResult.WinnerSide != "attacker" {
		t.Fatalf("expected tower winner side attacker, got %s", result.BattleResult.WinnerSide)
	}
	if result.BattleResult.AttackerPower <= 0 {
		t.Fatalf("expected attacker power > 0, got %d", result.BattleResult.AttackerPower)
	}
}

func TestService_StartChallengeRejectsWhenNoChallengesRemain(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
	playerID := int64(1003)

	for i := 0; i < 5; i++ {
		if _, err := svc.StartChallenge(ctx, playerID, "pagoda"); err != nil {
			t.Fatalf("expected challenge %d success, got %v", i+1, err)
		}
	}

	before, err := growthRepo.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet fetch success, got %v", err)
	}

	if _, err := svc.StartChallenge(ctx, playerID, "pagoda"); err == nil {
		t.Fatal("expected sixth challenge to fail when remaining challenges are exhausted")
	}

	after, err := growthRepo.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet fetch success, got %v", err)
	}
	if after.BoneLevel != before.BoneLevel {
		t.Fatalf("expected wallet to stay at bone level %d, got %d", before.BoneLevel, after.BoneLevel)
	}
}
