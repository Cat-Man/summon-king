package pet

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
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
	if payload.Data.ActiveTeam[0].Exp != 0 {
		t.Fatalf("expected active pet exp 0, got %d", payload.Data.ActiveTeam[0].Exp)
	}
	if payload.Data.ActiveTeam[0].NextLevelExp != 100 {
		t.Fatalf("expected active pet next level exp 100, got %d", payload.Data.ActiveTeam[0].NextLevelExp)
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Total != 120 {
		t.Fatalf("expected active pet breakdown total 120, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Total)
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Base != 120 {
		t.Fatalf("expected active pet breakdown base 120, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Base)
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Level != 0 {
		t.Fatalf("expected active pet breakdown level 0, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Level)
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

func TestHandler_SaveTeamReturnsUpdatedCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository()))
	g := r.Group("/pet")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodPost, "/pet/team/save", bytes.NewBufferString(`{"player_id":1001,"pet_ids":[10011,10012]}`))
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
	if len(payload.Data.ActiveTeam) != 2 {
		t.Fatalf("expected active team 2, got %d", len(payload.Data.ActiveTeam))
	}
	if payload.Data.TotalPower != 276 {
		t.Fatalf("expected total power 276, got %d", payload.Data.TotalPower)
	}
}

func TestHandler_TeamReturnsPowerBreakdown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	growthRepo := growth.NewMemoryRepository()
	ctx := context.Background()
	if _, err := growthRepo.UpgradeBoneLevel(ctx, 1001, 1); err != nil {
		t.Fatalf("expected bone upgrade success, got %v", err)
	}
	if err := growthRepo.UpdateSpiritPower(ctx, 1001, 8); err != nil {
		t.Fatalf("expected spirit update success, got %v", err)
	}
	if _, err := growthRepo.UpgradeSoulPower(ctx, 1001, 2); err != nil {
		t.Fatalf("expected soul upgrade success, got %v", err)
	}

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(NewService(NewMemoryRepository(), WithGrowthReader(growthRepo)))
	g := r.Group("/pet")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/pet/team?player_id=1001", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code int `json:"code"`
		Data struct {
			ActiveTeam []struct {
				PowerBreakdown struct {
					Base   int64 `json:"base"`
					Bone   int64 `json:"bone"`
					Spirit int64 `json:"spirit"`
					Soul   int64 `json:"soul"`
					Total  int64 `json:"total"`
				} `json:"power_breakdown"`
			} `json:"active_team"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if payload.Code != 0 {
		t.Fatalf("expected business code 0, got %d", payload.Code)
	}
	if len(payload.Data.ActiveTeam) != 1 {
		t.Fatalf("expected active team 1, got %d", len(payload.Data.ActiveTeam))
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Base != 120 {
		t.Fatalf("expected breakdown base 120, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Base)
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Bone != 24 {
		t.Fatalf("expected breakdown bone 24, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Bone)
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Spirit != 8 {
		t.Fatalf("expected breakdown spirit 8, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Spirit)
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Soul != 16 {
		t.Fatalf("expected breakdown soul 16, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Soul)
	}
	if payload.Data.ActiveTeam[0].PowerBreakdown.Total != 168 {
		t.Fatalf("expected breakdown total 168, got %d", payload.Data.ActiveTeam[0].PowerBreakdown.Total)
	}
}
