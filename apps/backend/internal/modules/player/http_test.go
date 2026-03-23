package player

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/commerce"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/gin-gonic/gin"
)

type homeIndexHTTPEnvelope struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    HomeIndex `json:"data"`
	TraceID string    `json:"trace_id"`
}

func TestHandler_IndexReturnsHomeDashboardBlocks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	playerEntity, err := accountSvc.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected guest init success, got %v", err)
	}

	petSvc := pet.NewService(pet.NewMemoryRepository())
	dungeonSvc := dungeon.NewService(dungeon.NewMemoryRepository())
	commerceSvc := commerce.NewService(commerce.NewMemoryRepository())
	_ = commerceSvc.ClaimSignin(ctx, playerEntity.PlayerID)
	record, _ := dungeonSvc.StartCultivation(ctx, playerEntity.PlayerID, "fire-mountain", 2)

	playerSvc := NewService(
		NewRepository(accountRepo, assetRepo),
		WithPetReader(petSvc),
		WithDungeonReader(dungeonSvc),
		WithCommerceReader(commerceSvc),
		WithNow(func() time.Time {
			return record.StartedAt.Add(30 * time.Minute)
		}),
	)

	router := gin.New()
	router.Use(middleware.TraceID())
	NewHandler(playerSvc).RegisterRoutes(router.Group("/api/v1/player/home"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/player/home/index?player_id=1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected http 200, got %d", rec.Code)
	}

	var envelope homeIndexHTTPEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if envelope.Code != 0 {
		t.Fatalf("expected business code 0, got %d", envelope.Code)
	}
	if len(envelope.Data.DailyTodos) == 0 {
		t.Fatal("expected daily_todos in response")
	}
	if len(envelope.Data.Resources) == 0 {
		t.Fatal("expected resources in response")
	}
	if len(envelope.Data.ResourceIcons) == 0 {
		t.Fatal("expected resource_icons in response")
	}
	if len(envelope.Data.Messages) == 0 {
		t.Fatal("expected messages in response")
	}
	if len(envelope.Data.Entries) == 0 {
		t.Fatal("expected entries in response")
	}
	if envelope.Data.ActivityEntry.Title != "修行收益" {
		t.Fatalf("expected dynamic activity_entry title 修行收益, got %s", envelope.Data.ActivityEntry.Title)
	}
	if envelope.Data.CultivationSummary.MapName == "" {
		t.Fatal("expected cultivation_summary in response")
	}
	if envelope.Data.CultivationSummary.Status != "running" {
		t.Fatalf("expected cultivation_summary status running, got %s", envelope.Data.CultivationSummary.Status)
	}
	if envelope.Data.CultivationSummary.RemainingSeconds != 5400 {
		t.Fatalf("expected cultivation_summary remaining_seconds 5400, got %d", envelope.Data.CultivationSummary.RemainingSeconds)
	}
	if envelope.Data.CultivationSummary.RewardCoins <= 0 {
		t.Fatalf("expected cultivation_summary.reward_coins > 0, got %d", envelope.Data.CultivationSummary.RewardCoins)
	}
	if envelope.Data.CultivationSummary.RewardPetExp <= 0 {
		t.Fatalf("expected cultivation_summary.reward_pet_exp > 0, got %d", envelope.Data.CultivationSummary.RewardPetExp)
	}

	vitalityValue := resourceValueByLabel(envelope.Data.Resources, "活力")
	reputationValue := resourceValueByLabel(envelope.Data.Resources, "声望")
	if vitalityValue == "120 / 120" {
		t.Fatalf("expected vitality to be dynamic instead of fixed placeholder, got %s", vitalityValue)
	}
	if reputationValue == "0" {
		t.Fatalf("expected reputation to be dynamic instead of fixed placeholder, got %s", reputationValue)
	}
}

func TestHandler_IndexRequiresPlayerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.TraceID())
	NewHandler(NewService(NewRepository(account.NewMemoryRepository(), asset.NewMemoryRepository()))).
		RegisterRoutes(router.Group("/api/v1/player/home"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/player/home/index", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected http 400, got %d", rec.Code)
	}

	var envelope struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
		TraceID string          `json:"trace_id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if envelope.Code != 4002 {
		t.Fatalf("expected code 4002, got %d", envelope.Code)
	}
}
