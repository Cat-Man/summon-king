package dungeon

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
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
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
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

func TestGetWorldMap_IncludesPrimaryDungeonBinding(t *testing.T) {
	svc := newTestDungeonService(t)

	world, err := svc.GetWorldMap(context.Background())
	if err != nil {
		t.Fatalf("expected get world map success, got %v", err)
	}
	if len(world.Cities) != 2 {
		t.Fatalf("expected 2 cities, got %d", len(world.Cities))
	}
	if len(world.Cities[0].Dungeons) != 2 {
		t.Fatalf("expected first city 2 dungeons, got %d", len(world.Cities[0].Dungeons))
	}
	if world.Cities[0].Dungeons[0].DungeonID != 1 {
		t.Fatalf("expected first city first dungeon id 1, got %d", world.Cities[0].Dungeons[0].DungeonID)
	}
	if world.Cities[0].Dungeons[1].DungeonID != 2 {
		t.Fatalf("expected first city second dungeon id 2, got %d", world.Cities[0].Dungeons[1].DungeonID)
	}
	if world.Cities[0].Dungeons[1].UnlockSpiritPower != 120 {
		t.Fatalf("expected first city second dungeon unlock spirit 120, got %d", world.Cities[0].Dungeons[1].UnlockSpiritPower)
	}
	if len(world.Cities[1].Dungeons) != 2 {
		t.Fatalf("expected second city 2 dungeons, got %d", len(world.Cities[1].Dungeons))
	}
	if world.Cities[1].Dungeons[0].DungeonName != "寒渊裂隙" {
		t.Fatalf("expected second city first dungeon name 寒渊裂隙, got %s", world.Cities[1].Dungeons[0].DungeonName)
	}
}

func TestGetWorldMap_IncludesCompositeUnlockRequirements(t *testing.T) {
	svc := newTestDungeonService(t)

	world, err := svc.GetWorldMap(context.Background())
	if err != nil {
		t.Fatalf("expected get world map success, got %v", err)
	}

	if len(world.Cities) == 0 {
		t.Fatal("expected at least one city")
	}

	found := false
	for _, city := range world.Cities {
		for _, dungeon := range city.Dungeons {
			if dungeon.DungeonID == 2 {
				found = true
				if dungeon.UnlockBoneLevel != 2 {
					t.Fatalf("expected bone level 2, got %d", dungeon.UnlockBoneLevel)
				}
				if dungeon.UnlockSoulPieces != 3 {
					t.Fatalf("expected soul pieces 3, got %d", dungeon.UnlockSoulPieces)
				}
			}
		}
	}
	if !found {
		t.Fatal("expected to find dungeon id 2")
	}
}

func TestEnterDungeon_RejectsWhenCompositeRequirementsNotMet(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
	playerID := int64(1016)

	if err := growthRepo.UpdateSpiritPower(ctx, playerID, 40); err != nil {
		t.Fatalf("expected spirit update success, got %v", err)
	}

	_, err := svc.EnterDungeon(ctx, playerID, 2)
	if err == nil {
		t.Fatal("expected an error for locked dungeon")
	}
	if err != ErrDungeonLocked {
		t.Fatalf("expected ErrDungeonLocked, got %v", err)
	}
}

func TestEnterDungeon_AllowsWhenCompositeRequirementsAreMet(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
	playerID := int64(1017)

	if err := growthRepo.UpdateSpiritPower(ctx, playerID, 140); err != nil {
		t.Fatalf("expected spirit update success, got %v", err)
	}
	if _, err := growthRepo.UpgradeBoneLevel(ctx, playerID, 1); err != nil {
		t.Fatalf("expected bone upgrade success, got %v", err)
	}
	if _, err := growthRepo.UpgradeSoulPower(ctx, playerID, 3); err != nil {
		t.Fatalf("expected soul upgrade success, got %v", err)
	}

	run, err := svc.EnterDungeon(ctx, playerID, 2)
	if err != nil {
		t.Fatalf("expected dungeon enter success, got %v", err)
	}
	if run.DungeonID != 2 {
		t.Fatalf("expected dungeon 2, got %d", run.DungeonID)
	}
}

func TestEnterDungeon_RejectsLockedDungeonWithoutEnoughSpiritPower(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))

	_, err := svc.EnterDungeon(ctx, 1015, 2)
	if err == nil {
		t.Fatal("expected locked dungeon error, got nil")
	}
	if err != ErrDungeonLocked {
		t.Fatalf("expected ErrDungeonLocked, got %v", err)
	}
}

func TestClaimCultivation_UpdatesSpiritWallet(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
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
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
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
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
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
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
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

func TestResolveRollReward_OngoingStatus(t *testing.T) {
	reward := resolveRollReward(DungeonRun{Status: "ongoing"})

	if reward.Label != "怪物掉落" {
		t.Fatalf("expected reward label 怪物掉落, got %s", reward.Label)
	}
	if reward.SpiritPower != 5 {
		t.Fatalf("expected spirit reward 5, got %d", reward.SpiritPower)
	}
	if reward.SoulPieces != 0 {
		t.Fatalf("expected soul reward 0, got %d", reward.SoulPieces)
	}
}

func TestResolveRollReward_BossStatus(t *testing.T) {
	reward := resolveRollReward(DungeonRun{Status: "boss"})

	if reward.Label != "Boss掉落" {
		t.Fatalf("expected reward label Boss掉落, got %s", reward.Label)
	}
	if reward.SpiritPower != 8 {
		t.Fatalf("expected spirit reward 8, got %d", reward.SpiritPower)
	}
	if reward.SoulPieces != 1 {
		t.Fatalf("expected soul reward 1, got %d", reward.SoulPieces)
	}
}

func TestResolveRollReward_DungeonTwoOngoingStatus(t *testing.T) {
	reward := resolveRollReward(DungeonRun{DungeonID: 2, CurrentFloor: 2, Status: "ongoing"})

	if reward.Label != "寒渊裂隙掉落" {
		t.Fatalf("expected reward label 寒渊裂隙掉落, got %s", reward.Label)
	}
	if reward.SpiritPower != 7 {
		t.Fatalf("expected spirit reward 7, got %d", reward.SpiritPower)
	}
	if reward.SoulPieces != 0 {
		t.Fatalf("expected soul reward 0, got %d", reward.SoulPieces)
	}
}

func TestResolveRollReward_DungeonOneDeepBossStatus(t *testing.T) {
	reward := resolveRollReward(DungeonRun{DungeonID: 1, CurrentFloor: 10, Status: "boss"})

	if reward.Label != "深层妖窟Boss掉落" {
		t.Fatalf("expected reward label 深层妖窟Boss掉落, got %s", reward.Label)
	}
	if reward.SpiritPower != 10 {
		t.Fatalf("expected spirit reward 10, got %d", reward.SpiritPower)
	}
	if reward.SoulPieces != 2 {
		t.Fatalf("expected soul reward 2, got %d", reward.SoulPieces)
	}
}

func TestRollDice_ExhaustedRunDoesNotGrantReward(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
	playerID := int64(1013)

	if _, err := svc.EnterDungeon(ctx, playerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}

	for i := 0; i < 15; i++ {
		if _, err := svc.RollDice(ctx, playerID); err != nil {
			t.Fatalf("expected roll dice success, got %v", err)
		}
	}

	run, err := svc.RollDice(ctx, playerID)
	if err != nil {
		t.Fatalf("expected exhausted roll success, got %v", err)
	}

	if run.Status != "exhausted" {
		t.Fatalf("expected exhausted status, got %s", run.Status)
	}
	if run.LastReward.Label != "无掉落" {
		t.Fatalf("expected no-drop label, got %s", run.LastReward.Label)
	}
	if run.LastReward.SpiritPower != 0 {
		t.Fatalf("expected spirit reward 0, got %d", run.LastReward.SpiritPower)
	}
	if run.WalletSnapshot.SpiritPower != 193 {
		t.Fatalf("expected wallet spirit 193 without extra reward, got %d", run.WalletSnapshot.SpiritPower)
	}
}

func TestRollDice_DungeonTwoUsesConfiguredReward(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(NewMemoryRepository(), asset.NewService(growthRepo))
	playerID := int64(1014)

	if err := growthRepo.UpdateSpiritPower(ctx, playerID, 20); err != nil {
		t.Fatalf("expected wallet spirit update success, got %v", err)
	}
	if _, err := growthRepo.UpgradeBoneLevel(ctx, playerID, 1); err != nil {
		t.Fatalf("expected bone upgrade success before entering, got %v", err)
	}
	if _, err := growthRepo.UpgradeSoulPower(ctx, playerID, 3); err != nil {
		t.Fatalf("expected soul upgrade success before entering, got %v", err)
	}
	if _, err := svc.EnterDungeon(ctx, playerID, 2); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}

	run, err := svc.RollDice(ctx, playerID)
	if err != nil {
		t.Fatalf("expected roll dice success, got %v", err)
	}

	if run.LastReward.Label != "寒渊裂隙掉落" {
		t.Fatalf("expected reward label 寒渊裂隙掉落, got %s", run.LastReward.Label)
	}
	if run.WalletSnapshot.SpiritPower != 127 {
		t.Fatalf("expected wallet spirit 127, got %d", run.WalletSnapshot.SpiritPower)
	}
}

func newTestDungeonService(t *testing.T) *Service {
	t.Helper()
	repo := NewMemoryRepository()
	return NewService(repo, asset.NewService(growth.NewMemoryRepository()))
}
