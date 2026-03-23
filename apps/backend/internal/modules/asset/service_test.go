package asset

import (
	"context"
	"testing"
	"time"
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

func TestGetResourceChangeLogs_OrderedDescAndLimited(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()
	playerID := int64(2001)

	nowSeq := []time.Time{
		time.Date(2026, 3, 23, 10, 0, 1, 0, time.UTC),
		time.Date(2026, 3, 23, 10, 0, 2, 0, time.UTC),
		time.Date(2026, 3, 23, 10, 0, 3, 0, time.UTC),
	}
	idx := 0
	svc.now = func() time.Time {
		v := nowSeq[idx]
		idx++
		return v
	}

	if _, err := svc.GrantReward(ctx, RewardGrant{
		PlayerID: playerID,
		BizID:    "reward-1",
		Source:   "daily_signin",
		Coins:    50,
	}); err != nil {
		t.Fatalf("expected grant reward success, got %v", err)
	}

	if _, err := svc.UseInventory(ctx, InventoryOperateRequest{
		PlayerID: playerID,
		ItemID:   "potion_small",
		Count:    1,
	}); err != nil {
		t.Fatalf("expected use inventory success, got %v", err)
	}

	if _, err := svc.SellInventory(ctx, InventoryOperateRequest{
		PlayerID: playerID,
		ItemID:   "summon_scroll",
		Count:    1,
	}); err != nil {
		t.Fatalf("expected sell inventory success, got %v", err)
	}

	logs, err := svc.GetResourceChangeLogs(ctx, playerID, 2)
	if err != nil {
		t.Fatalf("expected get logs success, got %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}
	if logs[0].ChangeType != "inventory_sell" {
		t.Fatalf("expected first log inventory_sell, got %s", logs[0].ChangeType)
	}
	if logs[1].ChangeType != "inventory_use" {
		t.Fatalf("expected second log inventory_use, got %s", logs[1].ChangeType)
	}
}

func TestGetResourceChangeLogs_DefaultLimitWhenNonPositive(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()
	playerID := int64(2002)

	base := time.Date(2026, 3, 23, 11, 0, 0, 0, time.UTC)
	idx := 0
	svc.now = func() time.Time {
		v := base.Add(time.Duration(idx) * time.Second)
		idx++
		return v
	}

	for i := 1; i <= 25; i++ {
		_, err := svc.GrantReward(ctx, RewardGrant{
			PlayerID: playerID,
			BizID:    "reward-" + time.Duration(i).String(),
			Source:   "test_source",
			Coins:    int64(i),
		})
		if err != nil {
			t.Fatalf("expected grant reward success at %d, got %v", i, err)
		}
	}

	logs, err := svc.GetResourceChangeLogs(ctx, playerID, 0)
	if err != nil {
		t.Fatalf("expected get logs success, got %v", err)
	}
	if len(logs) != 20 {
		t.Fatalf("expected default 20 logs, got %d", len(logs))
	}
	if logs[0].CoinsDelta != 25 {
		t.Fatalf("expected newest log coins delta 25, got %d", logs[0].CoinsDelta)
	}
	if logs[19].CoinsDelta != 6 {
		t.Fatalf("expected 20th log coins delta 6, got %d", logs[19].CoinsDelta)
	}
}
