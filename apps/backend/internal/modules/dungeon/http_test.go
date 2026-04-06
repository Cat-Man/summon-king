package dungeon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/gin-gonic/gin"
)

func TestHandler_StatusRequiresPlayerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), nil))
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
	svc := NewService(repo, nil)
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
	svc := NewService(repo, asset.NewService(growthRepo))

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

func TestHandler_RollReturnsLastBattlePayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository())
	svc := NewService(
		repo,
		asset.NewService(growthRepo),
		WithBattleTeamReader(petSvc),
		WithPetProgressor(petSvc),
	)
	if _, err := svc.EnterDungeon(context.Background(), 1002, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(svc)
	g := r.Group("/dungeon")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/dungeon/roll", strings.NewReader("player_id=1002"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Trace-ID", "trace-dungeon-roll")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var raw struct {
		Code int                        `json:"code"`
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &raw); err != nil {
		t.Fatalf("expected raw JSON payload, got %v", err)
	}
	if raw.Code != 0 {
		t.Fatalf("expected business code 0, got %d", raw.Code)
	}
	if _, ok := raw.Data["last_battle"]; !ok {
		t.Fatal("expected last_battle key in dungeon roll payload")
	}
	if _, ok := raw.Data["pet_growth"]; !ok {
		t.Fatal("expected pet_growth key in dungeon roll payload")
	}
	if _, ok := raw.Data["battle_result"]; ok {
		t.Fatal("expected dungeon roll payload to stop exposing battle_result")
	}

	var summary struct {
		BattleNo   string `json:"battle_no"`
		WinnerSide string `json:"winner_side"`
	}
	if err := json.Unmarshal(raw.Data["last_battle"], &summary); err != nil {
		t.Fatalf("expected battle summary JSON, got %v", err)
	}
	if summary.BattleNo == "" {
		t.Fatal("expected battle_no in dungeon roll payload")
	}
	if summary.WinnerSide != "attacker" {
		t.Fatalf("expected winner_side attacker, got %s", summary.WinnerSide)
	}

	var growth PetGrowth
	if err := json.Unmarshal(raw.Data["pet_growth"], &growth); err != nil {
		t.Fatalf("expected pet growth JSON, got %v", err)
	}
	if growth.Exp != 50 {
		t.Fatalf("expected pet growth exp 50, got %d", growth.Exp)
	}
	if growth.TeamTotalPower != 120 {
		t.Fatalf("expected pet growth team total power 120, got %d", growth.TeamTotalPower)
	}
}

func TestHandler_ClaimCultivationReturnsPetGrowth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	petSvc := pet.NewService(pet.NewMemoryRepository())
	svc := NewService(
		repo,
		asset.NewService(growthRepo),
		WithPetProgressor(petSvc),
	)
	base := time.Date(2026, time.April, 6, 10, 0, 0, 0, time.UTC)
	current := base
	restore := SetNowForTesting(func() time.Time { return current })
	defer restore()
	if _, err := svc.StartCultivation(context.Background(), 1003); err != nil {
		t.Fatalf("expected cultivation start success, got %v", err)
	}
	current = base.Add(time.Hour + time.Second)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(svc)
	g := r.Group("/dungeon")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/dungeon/cultivation/claim", strings.NewReader("player_id=1003"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Trace-ID", "trace-dungeon-cultivation-claim")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code int               `json:"code"`
		Data CultivationStatus `json:"data"`
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
	if payload.Data.PetGrowth.Exp != 100 {
		t.Fatalf("expected cultivation pet growth exp 100, got %d", payload.Data.PetGrowth.Exp)
	}
	if payload.Data.PetGrowth.TeamTotalPower != 144 {
		t.Fatalf("expected cultivation pet growth team total power 144, got %d", payload.Data.PetGrowth.TeamTotalPower)
	}
}

func TestHandler_ClaimCultivationRejectsBeforeClaimableAt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	svc := NewService(repo, asset.NewService(growth.NewMemoryRepository()))
	if _, err := svc.StartCultivation(context.Background(), 1004); err != nil {
		t.Fatalf("expected cultivation start success, got %v", err)
	}

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(svc)
	g := r.Group("/dungeon")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/dungeon/cultivation/claim", strings.NewReader("player_id=1004"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Trace-ID", "trace-dungeon-cultivation-not-ready")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload httpx.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if resp.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp.Code)
	}
	if payload.Code != 4091 {
		t.Fatalf("expected business code 4091, got %d", payload.Code)
	}
	if payload.TraceID != "trace-dungeon-cultivation-not-ready" {
		t.Fatalf("expected propagated trace id, got %s", payload.TraceID)
	}
}
