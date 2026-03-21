package dungeon

import (
	"context"
	"testing"
)

func TestEnterDungeon_GivesInitialDice(t *testing.T) {
	svc := newTestDungeonService(t)
	run, err := svc.EnterDungeon(context.Background(), 1001, 2001)
	if err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	if run.RemainDice != 15 {
		t.Fatalf("expected remain dice 15, got %d", run.RemainDice)
	}
}

func TestCultivationClaim_RequiresCompletedState(t *testing.T) {
	svc := newTestDungeonService(t)
	record, err := svc.StartCultivation(context.Background(), 1001, "fire-mountain", 2)
	if err != nil {
		t.Fatalf("expected start cultivation success, got %v", err)
	}
	if _, err := svc.ClaimCultivation(context.Background(), 1001, record.RecordID); err == nil {
		t.Fatal("expected claim before complete to fail")
	}
}

func newTestDungeonService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryRepository())
}
