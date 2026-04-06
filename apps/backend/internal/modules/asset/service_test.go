package asset

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

func TestService_ApplyWithMetadataPopulatesResult(t *testing.T) {
	ctx := context.Background()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(growthRepo)

	result, err := svc.ApplyWithMetadata(ctx, 1001, ApplyRequest{
		Delta: Delta{
			SpiritPower: 12,
			BoneLevel:   1,
			SoulPieces:  2,
		},
		Metadata: ApplyMetadata{
			Source:         "test",
			Reason:         "asset_delta",
			IdempotencyKey: "asset-1001",
		},
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
	if result.Metadata.Source != "test" || result.Metadata.Reason != "asset_delta" || result.Metadata.IdempotencyKey != "asset-1001" {
		t.Fatalf("expected metadata round trip, got %+v", result.Metadata)
	}
}

func TestService_SnapshotReturnsWallet(t *testing.T) {
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
