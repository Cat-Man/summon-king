package growth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestWallet_UsesPlayerIDAndTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewMemoryRepository())
	g := r.Group("/growth")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/growth/wallet?player_id=2002", nil)
	req.Header.Set("X-Trace-ID", "trace-growth")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	body := resp.Body.String()
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if !strings.Contains(body, "\"player_id\":2002") {
		t.Fatalf("expected response to contain player_id 2002, got %s", body)
	}
	if !strings.Contains(body, "\"trace_id\":\"trace-growth\"") {
		t.Fatalf("expected response to contain trace-growth, got %s", body)
	}
}

func TestBoneUpgrade_ReturnsUpdatedBoneState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewMemoryRepository())
	g := r.Group("/growth")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/growth/bone/upgrade?player_id=2002", nil)
	req.Header.Set("X-Trace-ID", "trace-bone-upgrade")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	body := resp.Body.String()
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if !strings.Contains(body, "\"level\":2") {
		t.Fatalf("expected response to contain level 2, got %s", body)
	}
	if !strings.Contains(body, "\"trace_id\":\"trace-bone-upgrade\"") {
		t.Fatalf("expected response to contain trace-bone-upgrade, got %s", body)
	}
}

func TestManorHarvest_ReturnsUpdatedPlots(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewMemoryRepository())
	g := r.Group("/growth")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/growth/manor/harvest?player_id=2003", nil)
	req.Header.Set("X-Trace-ID", "trace-manor-harvest")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	body := resp.Body.String()
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if !strings.Contains(body, "冷却中") {
		t.Fatalf("expected response to contain cooling plot state, got %s", body)
	}
	if !strings.Contains(body, "\"trace_id\":\"trace-manor-harvest\"") {
		t.Fatalf("expected response to contain trace-manor-harvest, got %s", body)
	}
}
