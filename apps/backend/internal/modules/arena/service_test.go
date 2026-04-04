package arena

import (
	"context"
	"testing"
)

func TestArena_FailResetsCurrentStreak(t *testing.T) {
	ctx := context.Background()
	svc := newTestArenaService(t)
	playerID := int64(1001)

	if err := svc.SetCurrentStreak(ctx, playerID, 10); err != nil {
		t.Fatalf("expected SetCurrentStreak success, got %v", err)
	}

	if err := svc.RecordBattleResult(ctx, playerID, false); err != nil {
		t.Fatalf("expected RecordBattleResult success, got %v", err)
	}

	rec, err := svc.GetDailyRecord(ctx, playerID)
	if err != nil {
		t.Fatalf("expected GetDailyRecord success, got %v", err)
	}
	if rec.CurrentStreak != 0 {
		t.Fatalf("expected streak reset to 0, got %d", rec.CurrentStreak)
	}
}

func newTestArenaService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryRepository())
}

func TestArena_NewPlayerStartsWithEmptyRecord(t *testing.T) {
	ctx := context.Background()
	svc := newTestArenaService(t)
	playerID := int64(2002)

	rec, err := svc.GetDailyRecord(ctx, playerID)
	if err != nil {
		t.Fatalf("expected empty record success, got %v", err)
	}
	if rec.PlayerID != playerID {
		t.Fatalf("expected player id %d, got %d", playerID, rec.PlayerID)
	}
	if rec.CurrentStreak != 0 {
		t.Fatalf("expected current streak 0, got %d", rec.CurrentStreak)
	}
}
