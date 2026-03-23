package asset

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/gin-gonic/gin"
)

func TestHandler_LogsReturnsResourceChangeEntriesInDescOrderWithLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	service := NewService(repo)
	playerID := int64(1001)
	nowSeq := []time.Time{
		time.Date(2026, 3, 23, 12, 0, 1, 0, time.UTC),
		time.Date(2026, 3, 23, 12, 0, 2, 0, time.UTC),
		time.Date(2026, 3, 23, 12, 0, 3, 0, time.UTC),
	}
	idx := 0
	service.now = func() time.Time {
		v := nowSeq[idx]
		idx++
		return v
	}

	if _, err := service.GrantReward(context.Background(), RewardGrant{
		PlayerID: playerID,
		BizID:    "signin-1",
		Source:   "daily_signin",
		Coins:    100,
	}); err != nil {
		t.Fatalf("expected grant reward success, got %v", err)
	}

	if _, err := service.UseInventory(context.Background(), InventoryOperateRequest{
		PlayerID: playerID,
		ItemID:   "potion_small",
		Count:    1,
	}); err != nil {
		t.Fatalf("expected use inventory success, got %v", err)
	}

	if _, err := service.SellInventory(context.Background(), InventoryOperateRequest{
		PlayerID: playerID,
		ItemID:   "summon_scroll",
		Count:    1,
	}); err != nil {
		t.Fatalf("expected sell inventory success, got %v", err)
	}

	router := gin.New()
	NewHandler(service).RegisterRoutes(router.Group("/api/v1/player/assets"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/player/assets/logs?player_id=1001&limit=2", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response struct {
		httpx.APIResponse
		Data []ResourceChangeLog `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}

	if len(response.Data) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(response.Data))
	}
	if response.Data[0].PlayerID != playerID {
		t.Fatalf("expected player_id=%d, got %d", playerID, response.Data[0].PlayerID)
	}
	if response.Data[0].ChangeType != "inventory_sell" {
		t.Fatalf("expected first change type inventory_sell, got %s", response.Data[0].ChangeType)
	}
	if response.Data[1].ChangeType != "inventory_use" {
		t.Fatalf("expected second change type inventory_use, got %s", response.Data[1].ChangeType)
	}
}

func TestHandler_LogsReturnsBadRequestWhenPlayerIDMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	service := NewService(repo)

	router := gin.New()
	NewHandler(service).RegisterRoutes(router.Group("/api/v1/player/assets"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/player/assets/logs", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}
