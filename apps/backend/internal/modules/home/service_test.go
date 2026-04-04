package home

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
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
	svc := NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc)

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
	svc := NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc)

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
