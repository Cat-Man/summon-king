package home

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/ranking"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
)

func TestOverview_ReturnsCoreSections(t *testing.T) {
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "测试玩家")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	assetSvc := asset.NewService(growthRepo)
	dungeonSvc := dungeon.NewService(dungeonRepo, assetSvc)
	if _, err := dungeonSvc.EnterDungeon(ctx, guest.PlayerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	if _, err := dungeonSvc.StartCultivation(ctx, guest.PlayerID); err != nil {
		t.Fatalf("expected cultivation start success, got %v", err)
	}

	towerSvc := tower.NewService(tower.NewMemoryRepository(), assetSvc)
	if _, err := towerSvc.StartChallenge(ctx, guest.PlayerID, "pagoda"); err != nil {
		t.Fatalf("expected pagoda start, got %v", err)
	}
	petSvc := pet.NewService(pet.NewMemoryRepository())
	arenaSvc := arena.NewService(arena.NewMemoryRepository(), assetSvc)
	rankingSvc := ranking.NewService(accountRepo, growthRepo, dungeonSvc, arenaSvc)
	svc := NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc, arenaSvc, rankingSvc)

	overview, err := svc.GetOverview(ctx, guest.PlayerID, guest.Token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if overview.PlayerID == 0 {
		t.Fatal("expected player id")
	}
	if overview.Nickname != "测试玩家" {
		t.Fatalf("expected nickname 测试玩家, got %s", overview.Nickname)
	}
	if overview.Wallet.PlayerID != guest.PlayerID {
		t.Fatalf("expected wallet player id %d, got %d", guest.PlayerID, overview.Wallet.PlayerID)
	}
	if overview.Modules.MapLabel == "" {
		t.Fatal("expected map label")
	}
	if overview.Modules.Dungeon.CurrentFloor != 1 {
		t.Fatalf("expected current floor 1, got %d", overview.Modules.Dungeon.CurrentFloor)
	}
	if overview.Modules.Cultivation.State != "cultivating" {
		t.Fatalf("expected cultivation state cultivating, got %s", overview.Modules.Cultivation.State)
	}
	if overview.Modules.Pet.TotalPower <= 0 {
		t.Fatalf("expected pet total power > 0, got %d", overview.Modules.Pet.TotalPower)
	}
	if overview.Modules.Pet.ActiveCount != 1 {
		t.Fatalf("expected active pet count 1, got %d", overview.Modules.Pet.ActiveCount)
	}
	if overview.Modules.Pet.StarterPetName != "初始灵狐" {
		t.Fatalf("expected starter pet 初始灵狐, got %s", overview.Modules.Pet.StarterPetName)
	}
	if overview.Modules.Tower.Pagoda.CurrentFloor != 1 {
		t.Fatalf("expected pagoda floor 1, got %d", overview.Modules.Tower.Pagoda.CurrentFloor)
	}
	if overview.NextAction.Route != "/dungeon" {
		t.Fatalf("expected next action dungeon, got %s", overview.NextAction.Route)
	}
}

func TestOverview_SkipsExhaustedDungeonWhenChoosingNextAction(t *testing.T) {
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "疲劳修士")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	assetSvc := asset.NewService(growthRepo)
	dungeonSvc := dungeon.NewService(dungeonRepo, assetSvc)
	if _, err := dungeonSvc.EnterDungeon(ctx, guest.PlayerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	for i := 0; i < 16; i++ {
		if _, err := dungeonSvc.RollDice(ctx, guest.PlayerID); err != nil {
			t.Fatalf("expected roll success, got %v", err)
		}
	}

	towerSvc := tower.NewService(tower.NewMemoryRepository(), assetSvc)
	petSvc := pet.NewService(pet.NewMemoryRepository())
	arenaSvc := arena.NewService(arena.NewMemoryRepository(), assetSvc)
	rankingSvc := ranking.NewService(accountRepo, growthRepo, dungeonSvc, arenaSvc)
	svc := NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc, arenaSvc, rankingSvc)

	overview, err := svc.GetOverview(ctx, guest.PlayerID, guest.Token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if overview.Modules.Dungeon.Status != "exhausted" {
		t.Fatalf("expected exhausted dungeon, got %s", overview.Modules.Dungeon.Status)
	}
	if overview.Modules.Pet.TotalPower <= 0 {
		t.Fatalf("expected pet total power > 0, got %d", overview.Modules.Pet.TotalPower)
	}
	if overview.Modules.Pet.StarterPetName != "初始灵狐" {
		t.Fatalf("expected starter pet 初始灵狐, got %s", overview.Modules.Pet.StarterPetName)
	}
	if overview.NextAction.Route != "/tower/pagoda" {
		t.Fatalf("expected next action /tower/pagoda, got %s", overview.NextAction.Route)
	}
}

func TestOverview_ReflectsGrowthBoostedPetPower(t *testing.T) {
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "战骨修士")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	if _, err := growthRepo.UpgradeBoneLevel(ctx, guest.PlayerID, 1); err != nil {
		t.Fatalf("expected bone upgrade success, got %v", err)
	}
	assetSvc := asset.NewService(growthRepo)
	dungeonSvc := dungeon.NewService(dungeonRepo, assetSvc)
	towerSvc := tower.NewService(tower.NewMemoryRepository(), assetSvc)
	petSvc := pet.NewService(pet.NewMemoryRepository(), pet.WithGrowthReader(growthRepo))
	arenaSvc := arena.NewService(arena.NewMemoryRepository(), assetSvc)
	rankingSvc := ranking.NewService(accountRepo, growthRepo, dungeonSvc, arenaSvc)
	svc := NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc, arenaSvc, rankingSvc)

	overview, err := svc.GetOverview(ctx, guest.PlayerID, guest.Token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if overview.Modules.Pet.TotalPower != 144 {
		t.Fatalf("expected pet total power 144 with bone bonus, got %d", overview.Modules.Pet.TotalPower)
	}
}

func TestOverview_MarksCultivationClaimableAtExactClaimableAt(t *testing.T) {
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "整点修士")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	base := time.Date(2100, time.April, 6, 10, 0, 0, 0, time.UTC)
	current := base
	restoreDungeonNow := dungeon.SetNowForTesting(func() time.Time { return current })
	defer restoreDungeonNow()
	restoreHomeNow := SetNowForTesting(func() time.Time { return current })
	defer restoreHomeNow()

	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	assetSvc := asset.NewService(growthRepo)
	dungeonSvc := dungeon.NewService(dungeonRepo, assetSvc)
	status, err := dungeonSvc.StartCultivation(ctx, guest.PlayerID)
	if err != nil {
		t.Fatalf("expected cultivation start success, got %v", err)
	}
	current = status.ClaimableAt

	towerSvc := tower.NewService(tower.NewMemoryRepository(), assetSvc)
	petSvc := pet.NewService(pet.NewMemoryRepository())
	arenaSvc := arena.NewService(arena.NewMemoryRepository(), assetSvc)
	rankingSvc := ranking.NewService(accountRepo, growthRepo, dungeonSvc, arenaSvc)
	svc := NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc, arenaSvc, rankingSvc)

	overview, err := svc.GetOverview(ctx, guest.PlayerID, guest.Token)
	if err != nil {
		t.Fatalf("expected overview success, got %v", err)
	}
	if !overview.Modules.Cultivation.Claimable {
		t.Fatal("expected cultivation to be claimable at exact claimable_at")
	}
}

func TestOverview_RecommendsArenaAndExposesArenaRankingModulesAfterTowersExhausted(t *testing.T) {
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "斗法修士")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	assetSvc := asset.NewService(growthRepo)
	dungeonSvc := dungeon.NewService(dungeonRepo, assetSvc)
	towerSvc := tower.NewService(tower.NewMemoryRepository(), assetSvc)
	for i := 0; i < 5; i++ {
		if _, err := towerSvc.StartChallenge(ctx, guest.PlayerID, "pagoda"); err != nil {
			t.Fatalf("expected pagoda start success, got %v", err)
		}
		if _, err := towerSvc.StartChallenge(ctx, guest.PlayerID, "spirit"); err != nil {
			t.Fatalf("expected spirit start success, got %v", err)
		}
	}

	arenaSvc := arena.NewService(arena.NewMemoryRepository(), assetSvc)
	if _, err := arenaSvc.RecordBattleResult(ctx, guest.PlayerID, true); err != nil {
		t.Fatalf("expected arena win success, got %v", err)
	}
	rankingSvc := ranking.NewService(accountRepo, growthRepo, dungeonSvc, arenaSvc)
	petSvc := pet.NewService(pet.NewMemoryRepository())
	svc := NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc, arenaSvc, rankingSvc)

	overview, err := svc.GetOverview(ctx, guest.PlayerID, guest.Token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if overview.NextAction.Route != "/arena" {
		t.Fatalf("expected next action /arena after towers exhausted, got %s", overview.NextAction.Route)
	}

	payload, err := json.Marshal(overview)
	if err != nil {
		t.Fatalf("expected overview marshal success, got %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(payload, &raw); err != nil {
		t.Fatalf("expected overview json decode success, got %v", err)
	}

	modules, ok := raw["modules"].(map[string]any)
	if !ok {
		t.Fatal("expected modules object")
	}
	if _, ok := modules["arena"]; !ok {
		t.Fatal("expected arena summary in overview modules")
	}
	if _, ok := modules["ranking"]; !ok {
		t.Fatal("expected ranking summary in overview modules")
	}
}
