package battle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBattle_SpeedDecidesTurnOrder(t *testing.T) {
	result := RunBattle(BattleInput{
		BattleNo: "battle-speed-order",
		Left:     []Unit{{Name: "A", Speed: 20}},
		Right:    []Unit{{Name: "B", Speed: 10}},
	})

	if len(result.Rounds) == 0 {
		t.Fatal("expected at least one round")
	}
	if len(result.Rounds[0].Actions) == 0 {
		t.Fatal("expected at least one action in first round")
	}
	if result.Rounds[0].Actions[0].Actor != "A" {
		t.Fatalf("expected first actor A, got %s", result.Rounds[0].Actions[0].Actor)
	}
}

func TestMemoryStore_SavesBattleDetailAndReplay(t *testing.T) {
	store := NewMemoryStore()
	result := RunBattle(BattleInput{
		BattleNo: "battle-store-1",
		Left:     []Unit{{Name: "A", Attack: 120, Speed: 20}},
		Right:    []Unit{{Name: "B", HP: 60, Defense: 0, Speed: 10}},
	})

	if err := store.Save(context.Background(), result); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	detail, err := store.GetDetail(context.Background(), "battle-store-1")
	if err != nil {
		t.Fatalf("expected detail success, got %v", err)
	}
	if detail.BattleNo != "battle-store-1" {
		t.Fatalf("expected battle no battle-store-1, got %s", detail.BattleNo)
	}

	replay, err := store.GetReplay(context.Background(), "battle-store-1")
	if err != nil {
		t.Fatalf("expected replay success, got %v", err)
	}
	if len(replay.Rounds) == 0 {
		t.Fatal("expected replay rounds")
	}
}

func TestHandler_ProvidesDetailAndReplayEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := NewMemoryStore()
	result := RunBattle(BattleInput{
		BattleNo: "battle-http-1",
		Left:     []Unit{{Name: "A", Attack: 120, Speed: 20}},
		Right:    []Unit{{Name: "B", HP: 60, Defense: 0, Speed: 10}},
	})
	if err := store.Save(context.Background(), result); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	router := gin.New()
	NewHandler(store).RegisterRoutes(router.Group("/api/v1/player/battles"))

	for _, path := range []string{
		"/api/v1/player/battles/detail?battle_no=battle-http-1",
		"/api/v1/player/battles/replay?battle_no=battle-http-1",
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, path, nil)
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200 for %s, got %d", path, recorder.Code)
		}
	}
}
