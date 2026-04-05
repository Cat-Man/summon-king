package home

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/ranking"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
	"github.com/gin-gonic/gin"
)

func TestHandler_OverviewReturnsAggregatedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "云上旅人")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	assetSvc := asset.NewService(growthRepo)
	dungeonSvc := dungeon.NewService(dungeonRepo, assetSvc)
	if _, err := dungeonSvc.EnterDungeon(ctx, guest.PlayerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	towerSvc := tower.NewService(tower.NewMemoryRepository(), assetSvc)
	petSvc := pet.NewService(pet.NewMemoryRepository())
	arenaSvc := arena.NewService(arena.NewMemoryRepository(), assetSvc)
	if _, err := arenaSvc.RecordBattleResult(ctx, guest.PlayerID, true); err != nil {
		t.Fatalf("expected arena win success, got %v", err)
	}
	rankingSvc := ranking.NewService(accountRepo, growthRepo, dungeonSvc, arenaSvc)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	r.Use(middleware.InjectAuthToken())

	h := NewHandler(NewService(accountRepo, dungeonSvc, growthRepo, petSvc, towerSvc, arenaSvc, rankingSvc))
	g := r.Group("/home")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/home/overview?player_id=1001", nil)
	req.Header.Set("Authorization", "Bearer "+guest.Token)
	req.Header.Set("X-Trace-ID", "trace-home-overview")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code    int      `json:"code"`
		TraceID string   `json:"trace_id"`
		Data    Overview `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if payload.Code != 0 {
		t.Fatalf("expected business code 0, got %d", payload.Code)
	}
	if payload.TraceID != "trace-home-overview" {
		t.Fatalf("expected trace id trace-home-overview, got %s", payload.TraceID)
	}
	if payload.Data.Nickname != "云上旅人" {
		t.Fatalf("expected nickname 云上旅人, got %s", payload.Data.Nickname)
	}
	if payload.Data.Wallet.PlayerID != guest.PlayerID {
		t.Fatalf("expected wallet player id %d, got %d", guest.PlayerID, payload.Data.Wallet.PlayerID)
	}
	if payload.Data.Modules.Dungeon.CurrentFloor != 1 {
		t.Fatalf("expected current floor 1, got %d", payload.Data.Modules.Dungeon.CurrentFloor)
	}
	if payload.Data.Modules.Pet.TotalPower <= 0 {
		t.Fatalf("expected pet total power > 0, got %d", payload.Data.Modules.Pet.TotalPower)
	}
	if payload.Data.Modules.Pet.ActiveCount != 1 {
		t.Fatalf("expected active pet count 1, got %d", payload.Data.Modules.Pet.ActiveCount)
	}
	if payload.Data.Modules.Pet.StarterPetName != "初始灵狐" {
		t.Fatalf("expected starter pet 初始灵狐, got %s", payload.Data.Modules.Pet.StarterPetName)
	}

	raw, err := json.Marshal(payload.Data)
	if err != nil {
		t.Fatalf("expected marshal success, got %v", err)
	}
	var view map[string]any
	if err := json.Unmarshal(raw, &view); err != nil {
		t.Fatalf("expected decode success, got %v", err)
	}
	modules, ok := view["modules"].(map[string]any)
	if !ok {
		t.Fatal("expected modules object")
	}
	if _, ok := modules["arena"]; !ok {
		t.Fatal("expected arena module in home overview payload")
	}
	if _, ok := modules["ranking"]; !ok {
		t.Fatal("expected ranking module in home overview payload")
	}
}

func TestHandler_OverviewRequiresPlayerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	growthRepo := growth.NewMemoryRepository()
	assetSvc := asset.NewService(growthRepo)
	towerSvc2 := tower.NewService(tower.NewMemoryRepository(), assetSvc)
	petSvc := pet.NewService(pet.NewMemoryRepository())
	arenaSvc := arena.NewService(arena.NewMemoryRepository(), assetSvc)
	rankingSvc := ranking.NewService(account.NewMemoryRepository(), growthRepo, dungeon.NewService(dungeon.NewMemoryRepository(), assetSvc), arenaSvc)
	h := NewHandler(NewService(account.NewMemoryRepository(), dungeon.NewService(dungeon.NewMemoryRepository(), assetSvc), growthRepo, petSvc, towerSvc2, arenaSvc, rankingSvc))
	g := r.Group("/home")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/home/overview", nil)
	req.Header.Set("X-Trace-ID", "trace-home-missing-player")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code    int    `json:"code"`
		TraceID string `json:"trace_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
	if payload.Code == 0 {
		t.Fatal("expected non-zero business code")
	}
	if payload.TraceID != "trace-home-missing-player" {
		t.Fatalf("expected trace id trace-home-missing-player, got %s", payload.TraceID)
	}
}
