package ranking

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

func TestGetLeaderboard_IncludesCurrentPlayerAndSortsByScore(t *testing.T) {
	ctx := context.Background()

	accountRepo := account.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)
	guest, err := accountSvc.GuestLogin(ctx, "榜首旅人")
	if err != nil {
		t.Fatalf("expected guest login success, got %v", err)
	}

	growthRepo := growth.NewMemoryRepository()
	if err := growthRepo.UpdateSpiritPower(ctx, guest.PlayerID, 800); err != nil {
		t.Fatalf("expected update spirit power success, got %v", err)
	}

	dungeonSvc := dungeon.NewService(dungeon.NewMemoryRepository(), growthRepo)
	if _, err := dungeonSvc.EnterDungeon(ctx, guest.PlayerID, 1); err != nil {
		t.Fatalf("expected enter dungeon success, got %v", err)
	}
	if _, err := dungeonSvc.RollDice(ctx, guest.PlayerID); err != nil {
		t.Fatalf("expected roll dice success, got %v", err)
	}

	svc := NewService(accountRepo, growthRepo, dungeonSvc)
	board, err := svc.GetLeaderboard(ctx, guest.PlayerID, 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(board) < 3 {
		t.Fatalf("expected at least 3 leaderboard entries, got %d", len(board))
	}
	if board[0].Score < board[1].Score {
		t.Fatalf("expected descending score order, got %d then %d", board[0].Score, board[1].Score)
	}

	foundSelf := false
	for _, entry := range board {
		if entry.PlayerID == guest.PlayerID {
			foundSelf = true
			if entry.Name != "榜首旅人" {
				t.Fatalf("expected current player name 榜首旅人, got %s", entry.Name)
			}
			if !entry.IsSelf {
				t.Fatal("expected current player entry to be marked as self")
			}
			if entry.Rank == 0 {
				t.Fatal("expected rank to be assigned")
			}
		}
	}
	if !foundSelf {
		t.Fatal("expected leaderboard to include current player entry")
	}
}
