package pet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestHandler_TeamReturnsPetCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository()))
	g := r.Group("/pet")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/pet/team?player_id=1001", nil)
	req.Header.Set("X-Trace-ID", "trace-pet-team")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code int            `json:"code"`
		Data CollectionView `json:"data"`
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
	if payload.Data.TeamSize != 1 {
		t.Fatalf("expected team size 1, got %d", payload.Data.TeamSize)
	}
	if len(payload.Data.ActiveTeam) != 1 {
		t.Fatalf("expected active team 1, got %d", len(payload.Data.ActiveTeam))
	}
}
