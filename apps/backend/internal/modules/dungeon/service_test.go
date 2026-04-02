package dungeon

import (
	"context"
	"testing"
)

func TestEnterDungeon_GivesInitialDice(t *testing.T) {
	svc := newTestDungeonService(t)
	ctx := context.Background()
	playerID := int64(1002)
	dungeonID := int64(1)

	run, err := svc.EnterDungeon(ctx, playerID, dungeonID)
	if err != nil {
		t.Fatalf("expected enter dungeon to succeed, got %v", err)
	}
	if run.RemainDice != 15 {
		t.Fatalf("expected remain dice 15, got %d", run.RemainDice)
	}
}

func newTestDungeonService(t *testing.T) *Service {
	t.Helper()
	repo := NewMemoryRepository()
	return NewService(repo)
}
