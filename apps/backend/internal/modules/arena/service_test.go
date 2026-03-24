package arena

import (
	"context"
	"testing"
)

func TestArena_FailResetsCurrentStreak(t *testing.T) {
	svc := newTestArenaService(t)
	_ = svc.SetCurrentStreak(context.Background(), 1001, 10)
	_ = svc.RecordBattleResult(context.Background(), 1001, false)
	rec := svc.GetDailyRecord(context.Background(), 1001)
	if rec.CurrentStreak != 0 {
		t.Fatalf("expected current streak 0, got %d", rec.CurrentStreak)
	}
}

func newTestArenaService(t *testing.T) *Service {
	t.Helper()
	return NewService()
}
