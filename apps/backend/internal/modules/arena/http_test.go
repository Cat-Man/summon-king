package arena

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestHandler_StatusReturnsDefaultRecordForNewPlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository()))
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
