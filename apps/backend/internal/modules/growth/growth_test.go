package growth

import (
	"context"
	"testing"
)

func TestSpiritWash_UsesFreeCountBeforeSpiritPower(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewSpiritService(repo)
	ctx := context.Background()
	playerID := int64(1001)

	before := svc.GetWallet(ctx, playerID).SpiritPower
	err := svc.WashSpirit(ctx, playerID, WashOption{})
	after := svc.GetWallet(ctx, playerID).SpiritPower

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if before != after {
		t.Fatalf("expected spirit power unchanged when free wash exists, before=%d after=%d", before, after)
	}
}

func TestBoneUpgrade_IncreasesBoneLevel(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewBoneService(repo)
	ctx := context.Background()
	playerID := int64(1002)

	before := svc.GetBoneState(ctx, playerID)
	after, err := svc.Upgrade(ctx, playerID)

	if err != nil {
		t.Fatalf("expected bone upgrade success, got %v", err)
	}
	if after.Level != before.Level+1 {
		t.Fatalf("expected level %d, got %d", before.Level+1, after.Level)
	}
}

func TestManorHarvest_UpdatesPlotState(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewManorService(repo)
	ctx := context.Background()
	playerID := int64(1003)

	result, err := svc.Harvest(ctx, playerID)
	if err != nil {
		t.Fatalf("expected manor harvest success, got %v", err)
	}
	if result.Message == "" {
		t.Fatal("expected harvest message")
	}
	if len(result.Plots) == 0 {
		t.Fatal("expected manor plots in result")
	}
	if result.Plots[0].State != "冷却中" {
		t.Fatalf("expected first plot cooling down, got %s", result.Plots[0].State)
	}
}
