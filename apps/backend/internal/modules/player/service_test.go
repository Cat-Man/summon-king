package player

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
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
