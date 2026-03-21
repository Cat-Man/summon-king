package growth

import (
	"context"
	"testing"
)

func TestSpiritWash_UsesFreeCountBeforeSpiritPower(t *testing.T) {
	svc := newTestGrowthService(t)
	ctx := context.Background()
	before := svc.GetWallet(ctx, 1001).SpiritPower
	err := svc.WashSpirit(ctx, 1001, 3001, WashOption{})
	after := svc.GetWallet(ctx, 1001).SpiritPower
	if err != nil {
		t.Fatalf("expected wash spirit success, got %v", err)
	}
	if before != after {
		t.Fatalf("expected spirit power unchanged when free count exists, before=%d after=%d", before, after)
	}
}

func TestHarvestManor_RequiresReadyCrop(t *testing.T) {
	svc := newTestGrowthService(t)
	ctx := context.Background()
	_, err := svc.PlantCrop(ctx, 1001, "lotus-seed")
	if err != nil {
		t.Fatalf("expected plant crop success, got %v", err)
	}
	if _, err := svc.HarvestCrop(ctx, 1001); err == nil {
		t.Fatal("expected harvest before ready to fail")
	}
}

func newTestGrowthService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryRepository())
}
