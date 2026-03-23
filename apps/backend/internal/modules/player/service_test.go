package player

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/commerce"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
)

func TestGetHomeIndex_UsesAssetWalletSnapshot(t *testing.T) {
	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()

	accountService := account.NewService(accountRepo)
	playerEntity, err := accountService.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected player init success, got %v", err)
	}

	assetService := asset.NewService(assetRepo)
	_, err = assetService.GrantReward(ctx, asset.RewardGrant{
		PlayerID: playerEntity.PlayerID,
		BizID:    "home-wallet-sync",
		Source:   "test",
		Coins:    88,
		Diamonds: 6,
	})
	if err != nil {
		t.Fatalf("expected asset grant success, got %v", err)
	}

	svc := NewService(NewRepository(accountRepo, assetRepo))
	home, err := svc.GetHomeIndex(ctx, playerEntity.PlayerID)
	if err != nil {
		t.Fatalf("expected home index success, got %v", err)
	}
	if home.Coin != 88 {
		t.Fatalf("expected coins from asset wallet, got %d", home.Coin)
	}
	if home.Diamond != 6 {
		t.Fatalf("expected diamonds from asset wallet, got %d", home.Diamond)
	}
}

func TestGetHomeIndex_AggregatesDashboardBlocks(t *testing.T) {
	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository())
	dungeonSvc := dungeon.NewService(dungeon.NewMemoryRepository())
	commerceSvc := commerce.NewService(commerce.NewMemoryRepository())

	accountService := account.NewService(accountRepo)
	playerEntity, err := accountService.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected player init success, got %v", err)
	}

	_ = commerceSvc.ClaimSignin(ctx, playerEntity.PlayerID)
	if _, err := dungeonSvc.StartCultivation(ctx, playerEntity.PlayerID, "fire-mountain", 2); err != nil {
		t.Fatalf("expected start cultivation success, got %v", err)
	}

	svc := NewService(
		NewRepository(accountRepo, assetRepo),
		WithPetReader(petSvc),
		WithDungeonReader(dungeonSvc),
		WithCommerceReader(commerceSvc),
	)
	home, err := svc.GetHomeIndex(ctx, playerEntity.PlayerID)
	if err != nil {
		t.Fatalf("expected home index success, got %v", err)
	}

	if len(home.DailyTodos) == 0 {
		t.Fatal("expected daily_todos to be aggregated")
	}
	if home.DailyTodos[0].Title != "签到状态" {
		t.Fatalf("expected first todo title 签到状态, got %s", home.DailyTodos[0].Title)
	}
	if len(home.Resources) == 0 {
		t.Fatal("expected resources to be aggregated")
	}
	if len(home.ResourceIcons) < 2 {
		t.Fatalf("expected at least 2 resource_icons, got %d", len(home.ResourceIcons))
	}
	if len(home.Messages) == 0 {
		t.Fatal("expected messages to be aggregated")
	}
	if len(home.Entries) == 0 {
		t.Fatal("expected entries to be aggregated")
	}
	if home.ActivityEntry.Title == "" {
		t.Fatal("expected activity_entry to be aggregated")
	}
	if len(home.DailyTodos) < 3 {
		t.Fatalf("expected at least 3 daily_todos, got %d", len(home.DailyTodos))
	}
	if home.Resources[1].Label != "战力" {
		t.Fatalf("expected second resource label 战力, got %s", home.Resources[1].Label)
	}
	if home.Resources[1].Value != "450" {
		t.Fatalf("expected power resource value 450, got %s", home.Resources[1].Value)
	}
}
