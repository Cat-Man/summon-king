package asset

import (
	"context"
	"fmt"
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

func TestGetResourceChangeLogs_PaginatedDescAndFiltered(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()
	playerID := int64(2001)

	nowSeq := []time.Time{
		time.Date(2026, 3, 23, 10, 0, 1, 0, time.UTC),
		time.Date(2026, 3, 23, 10, 0, 2, 0, time.UTC),
		time.Date(2026, 3, 23, 10, 0, 3, 0, time.UTC),
		time.Date(2026, 3, 23, 10, 0, 4, 0, time.UTC),
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

	if _, err := svc.GrantReward(ctx, RewardGrant{
		PlayerID: playerID,
		BizID:    "reward-2",
		Source:   "mail_reward",
		Coins:    30,
	}); err != nil {
		t.Fatalf("expected second grant reward success, got %v", err)
	}

	if _, err := svc.SellInventory(ctx, InventoryOperateRequest{
		PlayerID: playerID,
		ItemID:   "summon_scroll",
		Count:    1,
	}); err != nil {
		t.Fatalf("expected sell inventory success, got %v", err)
	}

	page1, err := svc.GetResourceChangeLogs(ctx, playerID, ResourceChangeLogQuery{
		Page:     1,
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("expected get logs success, got %v", err)
	}
	if page1.Total != 4 {
		t.Fatalf("expected total=4, got %d", page1.Total)
	}
	if page1.Page != 1 || page1.PageSize != 2 {
		t.Fatalf("expected page=1,page_size=2 got page=%d,page_size=%d", page1.Page, page1.PageSize)
	}
	if len(page1.Items) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(page1.Items))
	}
	if page1.Items[0].ChangeType != "inventory_sell" {
		t.Fatalf("expected first page first log inventory_sell, got %s", page1.Items[0].ChangeType)
	}
	if page1.Items[1].ChangeType != "grant_reward" {
		t.Fatalf("expected first page second log grant_reward, got %s", page1.Items[1].ChangeType)
	}

	page2, err := svc.GetResourceChangeLogs(ctx, playerID, ResourceChangeLogQuery{
		Page:     2,
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("expected get logs success, got %v", err)
	}
	if len(page2.Items) != 2 {
		t.Fatalf("expected 2 logs on second page, got %d", len(page2.Items))
	}
	if page2.Items[0].ChangeType != "inventory_use" {
		t.Fatalf("expected second page first log inventory_use, got %s", page2.Items[0].ChangeType)
	}
	if page2.Items[1].ChangeType != "grant_reward" {
		t.Fatalf("expected second page second log grant_reward, got %s", page2.Items[1].ChangeType)
	}

	filtered, err := svc.GetResourceChangeLogs(ctx, playerID, ResourceChangeLogQuery{
		Page:       1,
		PageSize:   20,
		ChangeType: "grant_reward",
	})
	if err != nil {
		t.Fatalf("expected get filtered logs success, got %v", err)
	}
	if filtered.Total != 2 {
		t.Fatalf("expected filtered total=2, got %d", filtered.Total)
	}
	if len(filtered.Items) != 2 {
		t.Fatalf("expected filtered items=2, got %d", len(filtered.Items))
	}
	if filtered.Items[0].BizID != "reward-2" {
		t.Fatalf("expected filtered first biz_id=reward-2, got %s", filtered.Items[0].BizID)
	}
	if filtered.Items[1].BizID != "reward-1" {
		t.Fatalf("expected filtered second biz_id=reward-1, got %s", filtered.Items[1].BizID)
	}
}

func TestGetResourceChangeLogs_DefaultPagingWhenNonPositive(t *testing.T) {
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
			BizID:    fmt.Sprintf("reward-%d", i),
			Source:   "test_source",
			Coins:    int64(i),
		})
		if err != nil {
			t.Fatalf("expected grant reward success at %d, got %v", i, err)
		}
	}

	logs, err := svc.GetResourceChangeLogs(ctx, playerID, ResourceChangeLogQuery{
		Page:     0,
		PageSize: -1,
	})
	if err != nil {
		t.Fatalf("expected get logs success, got %v", err)
	}
	if logs.Page != 1 || logs.PageSize != 20 {
		t.Fatalf("expected defaults page=1,page_size=20 got page=%d,page_size=%d", logs.Page, logs.PageSize)
	}
	if logs.Total != 25 {
		t.Fatalf("expected total=25, got %d", logs.Total)
	}
	if len(logs.Items) != 20 {
		t.Fatalf("expected default 20 logs, got %d", len(logs.Items))
	}
	if logs.Items[0].CoinsDelta != 25 {
		t.Fatalf("expected newest log coins delta 25, got %d", logs.Items[0].CoinsDelta)
	}
	if logs.Items[19].CoinsDelta != 6 {
		t.Fatalf("expected 20th log coins delta 6, got %d", logs.Items[19].CoinsDelta)
	}
}
