package asset

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

func TestService_ApplyAccumulatesWalletDelta(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(growthRepo)

	result, err := svc.Apply(ctx, 1001, Delta{
		SpiritPower: 12,
		BoneLevel:   1,
		SoulPieces:  2,
	})
	if err != nil {
		t.Fatalf("expected apply success, got %v", err)
	}
	if result.Wallet.SpiritPower != 112 {
		t.Fatalf("expected spirit power 112, got %d", result.Wallet.SpiritPower)
	}
	if result.Wallet.BoneLevel != 2 {
		t.Fatalf("expected bone level 2, got %d", result.Wallet.BoneLevel)
	}
	if result.Wallet.SoulPieces != 2 {
		t.Fatalf("expected soul pieces 2, got %d", result.Wallet.SoulPieces)
	}
}

func TestService_SnapshotReturnsCurrentWallet(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(growthRepo)

	snapshot, err := svc.Snapshot(ctx, 1002)
	if err != nil {
		t.Fatalf("expected snapshot success, got %v", err)
	}
	if snapshot.PlayerID != 1002 {
		t.Fatalf("expected player id 1002, got %d", snapshot.PlayerID)
	}
}
