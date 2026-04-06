package commerce

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

func TestService_ClaimSigninAppliesReward(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 4, 8, 9, 0, 0, 0, time.UTC)
	repo := NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(repo, asset.NewService(growthRepo), WithNow(func() time.Time {
		return now
	}))

	result, err := svc.ClaimSignin(ctx, 1001)
	if err != nil {
		t.Fatalf("expected signin claim success, got %v", err)
	}
	if result.RewardDelta.SpiritPower != 8 {
		t.Fatalf("expected reward 8, got %d", result.RewardDelta.SpiritPower)
	}
	if !result.TodayClaimed {
		t.Fatal("expected today claimed")
	}
	if result.StreakDays != 1 {
		t.Fatalf("expected streak 1, got %d", result.StreakDays)
	}
	if result.WalletSnapshot.SpiritPower != 108 {
		t.Fatalf("expected wallet spirit 108, got %d", result.WalletSnapshot.SpiritPower)
	}
	if result.RewardDelta.SpiritPower != 8 {
		t.Fatalf("expected reward 8, got %d", result.RewardDelta.SpiritPower)
	}
	if result.RewardDelta.BoneLevel != 0 {
		t.Fatalf("expected bone delta 0, got %d", result.RewardDelta.BoneLevel)
	}
	if result.RewardDelta.SoulPieces != 0 {
		t.Fatalf("expected soul delta 0, got %d", result.RewardDelta.SoulPieces)
	}

	_, err = svc.ClaimSignin(ctx, 1001)
	if !errors.Is(err, ErrAlreadyClaimedToday) {
		t.Fatalf("expected ErrAlreadyClaimedToday, got %v", err)
	}
}

func TestService_ClaimSigninDoublesAtFifthDay(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	base := time.Date(2026, 4, 8, 9, 0, 0, 0, time.UTC)
	svc := NewService(repo, asset.NewService(growthRepo), WithNow(func() time.Time {
		return base
	}))

	for i := 0; i < 4; i++ {
		if _, err := svc.ClaimSignin(ctx, 1100); err != nil {
			t.Fatalf("expected day %d success, got %v", i+1, err)
		}
		base = base.Add(24 * time.Hour)
		svc = NewService(repo, asset.NewService(growthRepo), WithNow(func() time.Time {
			return base
		}))
	}

	result, err := svc.ClaimSignin(ctx, 1100)
	if err != nil {
		t.Fatalf("expected fifth day claim success, got %v", err)
	}
	if result.StreakDays != 5 {
		t.Fatalf("expected streak 5, got %d", result.StreakDays)
	}
	if result.RewardDelta.SpiritPower != 16 {
		t.Fatalf("expected reward 16, got %d", result.RewardDelta.SpiritPower)
	}
}

func TestService_VIPDailyRequiresSingleClaim(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 4, 8, 10, 0, 0, 0, time.UTC)
	repo := NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(repo, asset.NewService(growthRepo), WithNow(func() time.Time {
		return now
	}))

	result, err := svc.ClaimVIPDaily(ctx, 1200)
	if err != nil {
		t.Fatalf("expected vip claim success, got %v", err)
	}
	if !result.DailyClaimed {
		t.Fatal("expected daily claimed")
	}
	if result.RewardDelta.SpiritPower != 6 {
		t.Fatalf("expected reward 6, got %d", result.RewardDelta.SpiritPower)
	}

	_, err = svc.ClaimVIPDaily(ctx, 1200)
	if !errors.Is(err, ErrAlreadyClaimedToday) {
		t.Fatalf("expected ErrAlreadyClaimedToday, got %v", err)
	}
}
