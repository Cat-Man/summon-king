package asset

import (
	"context"
	"testing"
)

func TestGrantReward_IdempotentByBizID(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()
	playerID := int64(1001)

	reward := RewardGrant{
		PlayerID: playerID,
		BizID:    "daily-signin-2026-03-21",
		Source:   "daily_signin",
		Coins:    100,
		Diamonds: 5,
	}

	first, err := svc.GrantReward(ctx, reward)
	if err != nil {
		t.Fatalf("expected first grant success, got error: %v", err)
	}
	if first.Idempotent {
		t.Fatal("expected first grant to be non-idempotent")
	}

	second, err := svc.GrantReward(ctx, reward)
	if err != nil {
		t.Fatalf("expected second grant success, got error: %v", err)
	}
	if !second.Idempotent {
		t.Fatal("expected second grant to be idempotent hit")
	}

	wallet, err := svc.GetWallet(ctx, playerID)
	if err != nil {
		t.Fatalf("expected wallet query success, got error: %v", err)
	}
	if wallet.Coins != 100 {
		t.Fatalf("expected coins=100, got %d", wallet.Coins)
	}
	if wallet.Diamonds != 5 {
		t.Fatalf("expected diamonds=5, got %d", wallet.Diamonds)
	}
}
