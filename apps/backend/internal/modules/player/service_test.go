package player

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/commerce"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
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

func TestGetHomeIndex_AggregatesDashboardBlocksAndDynamicFields(t *testing.T) {
	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository())
	dungeonSvc := dungeon.NewService(dungeon.NewMemoryRepository())
	commerceSvc := commerce.NewService(commerce.NewMemoryRepository())

	accountService := account.NewService(accountRepo)
	playerEntity, err := accountService.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected player init success, got %v", err)
	}

	_ = commerceSvc.ClaimSignin(ctx, playerEntity.PlayerID)
	record, err := dungeonSvc.StartCultivation(ctx, playerEntity.PlayerID, "fire-mountain", 2)
	if err != nil {
		t.Fatalf("expected start cultivation success, got %v", err)
	}

	svc := NewService(
		NewRepository(accountRepo, assetRepo),
		WithPetReader(petSvc),
		WithDungeonReader(dungeonSvc),
		WithCommerceReader(commerceSvc),
		WithNow(func() time.Time {
			return record.StartedAt.Add(30 * time.Minute)
		}),
	)
	home, err := svc.GetHomeIndex(ctx, playerEntity.PlayerID)
	if err != nil {
		t.Fatalf("expected home index success, got %v", err)
	}

	if len(home.DailyTodos) < 3 {
		t.Fatalf("expected at least 3 daily_todos, got %d", len(home.DailyTodos))
	}
	if len(home.Messages) == 0 {
		t.Fatal("expected messages to be aggregated")
	}
	if len(home.Entries) == 0 {
		t.Fatal("expected entries to be aggregated")
	}

	if home.ActivityEntry.Title != "修行收益" {
		t.Fatalf("expected dynamic activity title 修行收益, got %s", home.ActivityEntry.Title)
	}
	if home.ActivityEntry.Description != "火焰山修行进行中" {
		t.Fatalf("expected dynamic activity description 火焰山修行进行中, got %s", home.ActivityEntry.Description)
	}
	if home.ActivityEntry.Note != "剩余 1小时30分" {
		t.Fatalf("expected dynamic activity note 剩余 1小时30分, got %s", home.ActivityEntry.Note)
	}

	if home.CultivationSummary.Status != "running" {
		t.Fatalf("expected cultivation_summary status running, got %s", home.CultivationSummary.Status)
	}
	if home.CultivationSummary.MapName != "火焰山" {
		t.Fatalf("expected cultivation_summary map_name 火焰山, got %s", home.CultivationSummary.MapName)
	}
	if !home.CultivationSummary.StartedAt.Equal(record.StartedAt) {
		t.Fatalf("expected cultivation_summary started_at %v, got %v", record.StartedAt, home.CultivationSummary.StartedAt)
	}
	if !home.CultivationSummary.FinishedAt.Equal(record.FinishedAt) {
		t.Fatalf("expected cultivation_summary finished_at %v, got %v", record.FinishedAt, home.CultivationSummary.FinishedAt)
	}
	if home.CultivationSummary.RewardCoins != 240 {
		t.Fatalf("expected cultivation_summary reward_coins 240, got %d", home.CultivationSummary.RewardCoins)
	}
	if home.CultivationSummary.RewardPetExp != 160 {
		t.Fatalf("expected cultivation_summary reward_pet_exp 160, got %d", home.CultivationSummary.RewardPetExp)
	}
	if home.CultivationSummary.Action != "查看修行" {
		t.Fatalf("expected cultivation_summary action 查看修行, got %s", home.CultivationSummary.Action)
	}
	if home.CultivationSummary.RemainingSeconds != 5400 {
		t.Fatalf("expected cultivation_summary remaining_seconds 5400, got %d", home.CultivationSummary.RemainingSeconds)
	}

	vitalityValue := resourceValueByLabel(home.Resources, "活力")
	if vitalityValue == "" {
		t.Fatal("expected vitality resource")
	}
	if vitalityValue == "120 / 120" {
		t.Fatalf("expected vitality to be dynamic instead of fixed placeholder, got %s", vitalityValue)
	}
	if !strings.Contains(vitalityValue, "/") {
		t.Fatalf("expected vitality formatted as current / max, got %s", vitalityValue)
	}

	reputationValue := resourceValueByLabel(home.Resources, "声望")
	if reputationValue == "" {
		t.Fatal("expected reputation resource")
	}
	if reputationValue == "0" {
		t.Fatalf("expected reputation to be dynamic instead of fixed placeholder, got %s", reputationValue)
	}
}

func TestGetHomeIndex_CultivationStatusUsesTimeWindow(t *testing.T) {
	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()
	accountService := account.NewService(accountRepo)

	playerEntity, err := accountService.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected player init success, got %v", err)
	}

	now := time.Now()
	dungeonReader := &fixedDungeonReader{
		world: dungeon.WorldMap{
			PlayerID: playerEntity.PlayerID,
			Region:   "east",
			Cities: []dungeon.City{
				{CityID: "qingyun", CityName: "青云城", Unlocked: true},
				{CityID: "fire-mountain", CityName: "火焰山", Unlocked: true},
			},
		},
		found: true,
		cultivation: dungeon.CultivationRecord{
			RecordID:     "cult-1",
			PlayerID:     playerEntity.PlayerID,
			MapID:        "fire-mountain",
			Hours:        2,
			Status:       "running",
			RewardCoins:  240,
			RewardPetExp: 160,
			StartedAt:    now.Add(-3 * time.Hour),
			FinishedAt:   now.Add(-3 * time.Minute),
		},
	}

	svc := NewService(
		NewRepository(accountRepo, assetRepo),
		WithDungeonReader(dungeonReader),
		WithNow(func() time.Time {
			return now
		}),
	)

	home, err := svc.GetHomeIndex(ctx, playerEntity.PlayerID)
	if err != nil {
		t.Fatalf("expected home index success, got %v", err)
	}

	if home.CultivationSummary.Status != "claimable" {
		t.Fatalf("expected status claimable when finished_at passed, got %s", home.CultivationSummary.Status)
	}
	if home.CultivationSummary.RemainingSeconds != 0 {
		t.Fatalf("expected remaining_seconds 0, got %d", home.CultivationSummary.RemainingSeconds)
	}
	if home.CultivationSummary.Action != "领取收益" {
		t.Fatalf("expected action 领取收益, got %s", home.CultivationSummary.Action)
	}
	if home.ActivityEntry.Title != "修行收益" {
		t.Fatalf("expected activity to prioritize cultivation, got %s", home.ActivityEntry.Title)
	}
}

func TestGetHomeIndex_ActivityEntryFallsBackToSigninWhenNoCultivation(t *testing.T) {
	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()
	accountService := account.NewService(accountRepo)

	playerEntity, err := accountService.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected player init success, got %v", err)
	}

	svc := NewService(
		NewRepository(accountRepo, assetRepo),
		WithDungeonReader(&fixedDungeonReader{
			world: dungeon.WorldMap{
				PlayerID: playerEntity.PlayerID,
				Region:   "east",
				Cities: []dungeon.City{
					{CityID: "qingyun", CityName: "青云城", Unlocked: true},
				},
			},
			found: false,
		}),
		WithCommerceReader(commerce.NewService(commerce.NewMemoryRepository())),
	)

	home, err := svc.GetHomeIndex(ctx, playerEntity.PlayerID)
	if err != nil {
		t.Fatalf("expected home index success, got %v", err)
	}
	if home.CultivationSummary.Status != "idle" {
		t.Fatalf("expected idle cultivation_summary without record, got %s", home.CultivationSummary.Status)
	}
	if home.ActivityEntry.Title != "签到奖励" {
		t.Fatalf("expected signin fallback activity title 签到奖励, got %s", home.ActivityEntry.Title)
	}
}

func resourceValueByLabel(resources []HomeResource, label string) string {
	for _, item := range resources {
		if item.Label == label {
			return item.Value
		}
	}
	return ""
}

type fixedDungeonReader struct {
	world          dungeon.WorldMap
	worldErr       error
	cultivation    dungeon.CultivationRecord
	found          bool
	cultivationErr error
}

func (r *fixedDungeonReader) GetWorldMap(_ context.Context, _ int64) (dungeon.WorldMap, error) {
	return r.world, r.worldErr
}

func (r *fixedDungeonReader) GetLatestCultivation(_ context.Context, _ int64) (dungeon.CultivationRecord, bool, error) {
	return r.cultivation, r.found, r.cultivationErr
}
