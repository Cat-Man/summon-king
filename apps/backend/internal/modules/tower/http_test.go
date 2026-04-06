package tower

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/gin-gonic/gin"
)

func TestHandler_PagodaStatusReturnsUnifiedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), asset.NewService(growth.NewMemoryRepository())))
	g := r.Group("/tower")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/tower/pagoda/status?player_id=1001", nil)
	req.Header.Set("X-Trace-ID", "trace-tower-status")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code    int         `json:"code"`
		TraceID string      `json:"trace_id"`
		Data    TowerStatus `json:"data"`
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
	if payload.TraceID != "trace-tower-status" {
		t.Fatalf("expected trace id trace-tower-status, got %s", payload.TraceID)
	}
	if payload.Data.Tower != "pagoda" {
		t.Fatalf("expected tower pagoda, got %s", payload.Data.Tower)
	}
	if payload.Data.CurrentFloor != 0 {
		t.Fatalf("expected current floor 0, got %d", payload.Data.CurrentFloor)
	}
	if payload.Data.RemainingChallenges != 5 {
		t.Fatalf("expected 5 remaining challenges, got %d", payload.Data.RemainingChallenges)
	}
}

func TestHandler_PagodaStartReturnsBattlePayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), asset.NewService(growth.NewMemoryRepository())))
	g := r.Group("/tower")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/tower/pagoda/start", bytes.NewBufferString(`{"player_id":1001}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Trace-ID", "trace-tower-start")
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
	if _, ok := raw.Data["battle"]; !ok {
		t.Fatal("expected battle key in tower start payload")
	}
	if _, ok := raw.Data["battle_result"]; ok {
		t.Fatal("expected tower start payload to stop exposing battle_result")
	}

	var summary struct {
		BattleNo   string `json:"battle_no"`
		WinnerSide string `json:"winner_side"`
	}
	if err := json.Unmarshal(raw.Data["battle"], &summary); err != nil {
		t.Fatalf("expected battle summary JSON, got %v", err)
	}
	if summary.BattleNo == "" {
		t.Fatal("expected battle_no in tower start payload")
	}
	if summary.WinnerSide != "attacker" {
		t.Fatalf("expected winner_side attacker, got %s", summary.WinnerSide)
	}
}
