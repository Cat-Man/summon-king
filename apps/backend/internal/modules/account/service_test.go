package account

import (
	"context"
	"testing"
)

func TestCreateGuestPlayer_InitializesProfileAndWallet(t *testing.T) {
	svc := newTestAccountService(t)
	player, err := svc.CreateGuestPlayer(context.Background(), "web")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if player.PlayerID == 0 {
		t.Fatal("expected non-zero player id")
	}
	if player.Profile.PlayerID != player.PlayerID {
		t.Fatal("expected profile to be initialized")
	}
	if player.Wallet.PlayerID != player.PlayerID {
		t.Fatal("expected wallet to be initialized")
	}
}

func newTestAccountService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryRepository())
}
