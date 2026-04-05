package arena

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

func TestHandler_StatusReturnsDefaultRecordForNewPlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), asset.NewService(growth.NewMemoryRepository())))
	g := r.Group("/arena")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/arena/status?player_id=1001", nil)
	req.Header.Set("X-Trace-ID", "trace-arena-status")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code    int         `json:"code"`
		TraceID string      `json:"trace_id"`
		Data    DailyRecord `json:"data"`
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
	if payload.TraceID != "trace-arena-status" {
		t.Fatalf("expected trace id trace-arena-status, got %s", payload.TraceID)
	}
	if payload.Data.PlayerID != 1001 {
		t.Fatalf("expected player id 1001, got %d", payload.Data.PlayerID)
	}
	if payload.Data.CurrentStreak != 0 {
		t.Fatalf("expected current streak 0, got %d", payload.Data.CurrentStreak)
	}
}

func TestHandler_ChallengeReturnsRewardPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), asset.NewService(growth.NewMemoryRepository())))
	g := r.Group("/arena")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/arena/challenge", bytes.NewBufferString(`{"player_id":1001,"won":true}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Trace-ID", "trace-arena-challenge")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	body := resp.Body.Bytes()

	var payload struct {
		Code    int          `json:"code"`
		TraceID string       `json:"trace_id"`
		Data    BattleResult `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if payload.Data.RewardDelta.SpiritPower == 0 {
		t.Fatal("expected spirit reward in challenge response")
	}
	if payload.Data.Record.CurrentStreak != 1 {
		t.Fatalf("expected current streak 1, got %d", payload.Data.Record.CurrentStreak)
	}

	var raw struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("expected raw JSON payload, got %v", err)
	}
	if _, ok := raw.Data["battle"]; !ok {
		t.Fatal("expected battle key in arena challenge payload")
	}
	if _, ok := raw.Data["battle_result"]; ok {
		t.Fatal("expected arena challenge payload to stop exposing battle_result")
	}

	var summary struct {
		BattleNo   string `json:"battle_no"`
		WinnerSide string `json:"winner_side"`
	}
	if err := json.Unmarshal(raw.Data["battle"], &summary); err != nil {
		t.Fatalf("expected battle summary JSON, got %v", err)
	}
	if summary.BattleNo == "" {
		t.Fatal("expected battle_no in arena challenge payload")
	}
	if summary.WinnerSide != "attacker" {
		t.Fatalf("expected winner_side attacker, got %s", summary.WinnerSide)
	}
}

func TestHandler_IndexReturnsOpponentPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), asset.NewService(growth.NewMemoryRepository())))
	g := r.Group("/arena")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/arena/index?player_id=1001", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code int `json:"code"`
		Data struct {
			Opponents []struct {
				OpponentID int64  `json:"opponent_id"`
				Name       string `json:"name"`
			} `json:"opponents"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if len(payload.Data.Opponents) != 2 {
		t.Fatalf("expected 2 opponents, got %d", len(payload.Data.Opponents))
	}
	if payload.Data.Opponents[0].OpponentID == 0 {
		t.Fatal("expected non-zero opponent id")
	}
}

func TestHandler_RefreshChangesOpponentPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), asset.NewService(growth.NewMemoryRepository())))
	g := r.Group("/arena")
	h.RegisterRoutes(g)

	indexReq := httptest.NewRequest(http.MethodGet, "/arena/index?player_id=1002", nil)
	indexResp := httptest.NewRecorder()
	r.ServeHTTP(indexResp, indexReq)

	var indexPayload struct {
		Data struct {
			Opponents []struct {
				OpponentID int64 `json:"opponent_id"`
			} `json:"opponents"`
		} `json:"data"`
	}
	if err := json.NewDecoder(indexResp.Body).Decode(&indexPayload); err != nil {
		t.Fatalf("expected index JSON payload, got %v", err)
	}

	refreshReq := httptest.NewRequest(http.MethodPost, "/arena/refresh", bytes.NewBufferString(`{"player_id":1002}`))
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshResp := httptest.NewRecorder()
	r.ServeHTTP(refreshResp, refreshReq)

	var refreshPayload struct {
		Data struct {
			Opponents []struct {
				OpponentID int64 `json:"opponent_id"`
			} `json:"opponents"`
		} `json:"data"`
	}
	if err := json.NewDecoder(refreshResp.Body).Decode(&refreshPayload); err != nil {
		t.Fatalf("expected refresh JSON payload, got %v", err)
	}
	if refreshResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", refreshResp.Code)
	}
	if len(refreshPayload.Data.Opponents) != 2 {
		t.Fatalf("expected 2 refreshed opponents, got %d", len(refreshPayload.Data.Opponents))
	}
	if refreshPayload.Data.Opponents[0].OpponentID == indexPayload.Data.Opponents[0].OpponentID {
		t.Fatalf("expected refreshed opponent id to change, still got %d", refreshPayload.Data.Opponents[0].OpponentID)
	}
}
