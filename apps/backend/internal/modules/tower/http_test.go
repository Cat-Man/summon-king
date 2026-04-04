package tower

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestHandler_PagodaStatusReturnsUnifiedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository()))
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
