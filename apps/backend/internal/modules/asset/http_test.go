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

func TestHandler_LogsReturnsPagedDataInDescOrderWithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	service := NewService(repo)
	playerID := int64(1001)
	nowSeq := []time.Time{
		time.Date(2026, 3, 23, 12, 0, 1, 0, time.UTC),
		time.Date(2026, 3, 23, 12, 0, 2, 0, time.UTC),
		time.Date(2026, 3, 23, 12, 0, 3, 0, time.UTC),
		time.Date(2026, 3, 23, 12, 0, 4, 0, time.UTC),
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

	if _, err := service.GrantReward(context.Background(), RewardGrant{
		PlayerID: playerID,
		BizID:    "mail-2",
		Source:   "mail_reward",
		Coins:    30,
	}); err != nil {
		t.Fatalf("expected second grant reward success, got %v", err)
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
	request := httptest.NewRequest(http.MethodGet, "/api/v1/player/assets/logs?player_id=1001&page=1&page_size=2&change_type=grant_reward", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response struct {
		httpx.APIResponse
		Data ResourceChangeLogPage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}

	if response.Data.Total != 2 {
		t.Fatalf("expected total=2, got %d", response.Data.Total)
	}
	if response.Data.Page != 1 || response.Data.PageSize != 2 {
		t.Fatalf("expected page=1,page_size=2 got page=%d,page_size=%d", response.Data.Page, response.Data.PageSize)
	}
	if len(response.Data.Items) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(response.Data.Items))
	}
	if response.Data.Items[0].PlayerID != playerID {
		t.Fatalf("expected player_id=%d, got %d", playerID, response.Data.Items[0].PlayerID)
	}
	if response.Data.Items[0].BizID != "mail-2" {
		t.Fatalf("expected first biz_id=mail-2, got %s", response.Data.Items[0].BizID)
	}
	if response.Data.Items[1].BizID != "signin-1" {
		t.Fatalf("expected second biz_id=signin-1, got %s", response.Data.Items[1].BizID)
	}
}

func TestHandler_LogsTreatsAllChangeTypeAsUnfiltered(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	service := NewService(repo)
	playerID := int64(1002)
	service.now = func() time.Time {
		return time.Date(2026, 3, 23, 12, 30, 0, 0, time.UTC)
	}

	if _, err := service.UseInventory(context.Background(), InventoryOperateRequest{
		PlayerID: playerID,
		ItemID:   "potion_small",
		Count:    1,
	}); err != nil {
		t.Fatalf("expected use inventory success, got %v", err)
	}

	router := gin.New()
	NewHandler(service).RegisterRoutes(router.Group("/api/v1/player/assets"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/player/assets/logs?player_id=1002&page=1&page_size=20&change_type=all",
		nil,
	)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response struct {
		httpx.APIResponse
		Data ResourceChangeLogPage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}

	if response.Data.Total != 1 {
		t.Fatalf("expected total=1, got %d", response.Data.Total)
	}
	if len(response.Data.Items) != 1 {
		t.Fatalf("expected 1 log item, got %d", len(response.Data.Items))
	}
	if response.Data.Items[0].ChangeType != "inventory_use" {
		t.Fatalf("expected change_type=inventory_use, got %s", response.Data.Items[0].ChangeType)
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

func TestHandler_LogsReturnsBadRequestWhenPageInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	service := NewService(repo)

	router := gin.New()
	NewHandler(service).RegisterRoutes(router.Group("/api/v1/player/assets"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/player/assets/logs?player_id=1001&page=0", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestHandler_LogsReturnsBadRequestWhenPageSizeInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMemoryRepository()
	service := NewService(repo)

	router := gin.New()
	NewHandler(service).RegisterRoutes(router.Group("/api/v1/player/assets"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/player/assets/logs?player_id=1001&page_size=abc", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}
