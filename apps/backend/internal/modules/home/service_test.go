package home

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
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
	dungeonSvc := dungeon.NewService(dungeonRepo)
	if _, err := dungeonSvc.EnterDungeon(ctx, guest.PlayerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	if _, err := dungeonSvc.StartCultivation(ctx, guest.PlayerID); err != nil {
		t.Fatalf("expected cultivation start success, got %v", err)
	}

	growthRepo := growth.NewMemoryRepository()
	svc := NewService(accountRepo, dungeonSvc, growthRepo)

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
}
