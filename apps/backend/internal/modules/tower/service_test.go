package tower

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
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

func TestService_UsesInjectedBattleTeamReader(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository())
	if _, err := petSvc.SetMainPet(ctx, 4101, 41012); err != nil {
		t.Fatalf("expected set main pet success, got %v", err)
	}

	svc := NewService(
		NewMemoryRepository(),
		asset.NewService(growthRepo),
		WithBattleTeamReader(petSvc),
	)

	result, err := svc.StartChallenge(ctx, 4101, "pagoda")
	if err != nil {
		t.Fatalf("expected challenge success, got %v", err)
	}
	if result.BattleResult.AttackerPower != 156 {
		t.Fatalf("expected attacker power 156, got %d", result.BattleResult.AttackerPower)
	}
}

func TestService_PagodaRewardRaisesNextChallengeBattlePower(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository(), pet.WithGrowthReader(growthRepo))

	svc := NewService(
		NewMemoryRepository(),
		asset.NewService(growthRepo),
		WithBattleTeamReader(petSvc),
	)

	first, err := svc.StartChallenge(ctx, 4102, "pagoda")
	if err != nil {
		t.Fatalf("expected first challenge success, got %v", err)
	}
	second, err := svc.StartChallenge(ctx, 4102, "pagoda")
	if err != nil {
		t.Fatalf("expected second challenge success, got %v", err)
	}
	if first.BattleResult.AttackerPower != 120 {
		t.Fatalf("expected first challenge attacker power 120, got %d", first.BattleResult.AttackerPower)
	}
	if second.BattleResult.AttackerPower != 144 {
		t.Fatalf("expected second challenge attacker power 144 after pagoda reward, got %d", second.BattleResult.AttackerPower)
	}
}
