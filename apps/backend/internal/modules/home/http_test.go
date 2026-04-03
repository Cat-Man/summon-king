package home

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
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
	dungeonSvc := dungeon.NewService(dungeonRepo)
	if _, err := dungeonSvc.EnterDungeon(ctx, guest.PlayerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	growthRepo := growth.NewMemoryRepository()

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	r.Use(middleware.InjectAuthToken())

	h := NewHandler(NewService(accountRepo, dungeonSvc, growthRepo))
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
}

func TestHandler_OverviewRequiresPlayerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(account.NewMemoryRepository(), dungeon.NewService(dungeon.NewMemoryRepository()), growth.NewMemoryRepository()))
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
