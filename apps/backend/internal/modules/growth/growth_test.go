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
