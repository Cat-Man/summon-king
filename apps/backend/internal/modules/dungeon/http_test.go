package dungeon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/gin-gonic/gin"
)

func TestHandler_StatusRequiresPlayerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository()))
	g := r.Group("/dungeon")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/dungeon/status", nil)
	req.Header.Set("X-Trace-ID", "trace-dungeon-missing-player")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload httpx.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
	if payload.Code != 4001 {
		t.Fatalf("expected business code 4001, got %d", payload.Code)
	}
	if payload.TraceID != "trace-dungeon-missing-player" {
		t.Fatalf("expected propagated trace id, got %s", payload.TraceID)
	}
}

func TestHandler_StatusReadsPlayerIDFromQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	svc := NewService(repo)
	if _, err := svc.EnterDungeon(context.Background(), 1002, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(svc)
	g := r.Group("/dungeon")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/dungeon/status?player_id=1002", nil)
	req.Header.Set("X-Trace-ID", "trace-dungeon-status")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code    int        `json:"code"`
		TraceID string     `json:"trace_id"`
		Data    DungeonRun `json:"data"`
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
	if payload.TraceID != "trace-dungeon-status" {
		t.Fatalf("expected propagated trace id, got %s", payload.TraceID)
	}
	if payload.Data.PlayerID != 1002 {
		t.Fatalf("expected response player_id 1002, got %d", payload.Data.PlayerID)
	}
}

func TestHandler_EnterDungeonRejectsLockedDungeon(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	svc := NewService(repo, growthRepo)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(svc)
	g := r.Group("/dungeon")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/dungeon/enter", strings.NewReader("player_id=1002&dungeon_id=2"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Trace-ID", "trace-dungeon-locked")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload httpx.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.Code)
	}
	if payload.Code != 4031 {
		t.Fatalf("expected business code 4031, got %d", payload.Code)
	}
	if payload.TraceID != "trace-dungeon-locked" {
		t.Fatalf("expected propagated trace id, got %s", payload.TraceID)
	}
}
