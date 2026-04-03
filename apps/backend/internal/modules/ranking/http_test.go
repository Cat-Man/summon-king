package ranking

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/gin-gonic/gin"
)

func TestHandler_LeaderboardReturnsUnifiedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "榜单玩家")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	growthRepo := growth.NewMemoryRepository()
	dungeonSvc := dungeon.NewService(dungeon.NewMemoryRepository(), growthRepo)
	svc := NewService(accountRepo, growthRepo, dungeonSvc)

	r := gin.New()
	r.Use(middleware.InjectTraceID())
	h := NewHandler(svc)
	g := r.Group("/ranking")
	h.RegisterRoutes(g)

	req := httptest.NewRequest(http.MethodGet, "/ranking/leaderboard?player_id=1001", nil)
	req.Header.Set("X-Trace-ID", "trace-ranking")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var payload struct {
		Code    int                `json:"code"`
		TraceID string             `json:"trace_id"`
		Data    []LeaderboardEntry `json:"data"`
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
	if payload.TraceID != "trace-ranking" {
		t.Fatalf("expected trace id trace-ranking, got %s", payload.TraceID)
	}
	if len(payload.Data) == 0 {
		t.Fatal("expected leaderboard entries")
	}
	foundSelf := false
	for _, entry := range payload.Data {
		if entry.PlayerID == guest.PlayerID {
			foundSelf = true
		}
	}
	if !foundSelf {
		t.Fatal("expected current player entry in leaderboard response")
	}
}
