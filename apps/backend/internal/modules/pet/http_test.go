package pet

import (
	"bytes"
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
	if len(payload.Data.Roster) != 2 {
		t.Fatalf("expected roster 2, got %d", len(payload.Data.Roster))
	}
}

func TestHandler_SetMainPetReturnsUpdatedCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository()))
	g := r.Group("/pet")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/pet/team/main", bytes.NewBufferString(`{"player_id":1001,"pet_id":10012}`))
	req.Header.Set("Content-Type", "application/json")
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
	if payload.Data.ActiveTeam[0].PetID != 10012 {
		t.Fatalf("expected main pet 10012, got %d", payload.Data.ActiveTeam[0].PetID)
	}
	if payload.Data.ActiveTeam[0].Name != "玄甲龟" {
		t.Fatalf("expected main pet 玄甲龟, got %s", payload.Data.ActiveTeam[0].Name)
	}
}
