package wxmini

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
)

func TestService_ExchangeLoginCreatesUnifiedSession(t *testing.T) {
	ctx := context.Background()
	repo := account.NewMemoryRepository()
	svc := NewService(repo, WithNow(func() time.Time {
		return time.Date(2026, time.April, 7, 3, 30, 0, 0, time.UTC)
	}))

	resp, err := svc.ExchangeLogin(ctx, "wx-login-code")
	if err != nil {
		t.Fatalf("expected exchange success, got %v", err)
	}
	if resp.UnifiedToken == "" {
		t.Fatal("expected unified token")
	}
	if resp.Channel != Channel {
		t.Fatalf("expected channel %s, got %s", Channel, resp.Channel)
	}

	session, err := repo.GetByToken(ctx, resp.UnifiedToken)
	if err != nil {
		t.Fatalf("expected account persisted by token, got %v", err)
	}
	if session.PlayerID != resp.PlayerID {
		t.Fatalf("expected player_id %d, got %d", resp.PlayerID, session.PlayerID)
	}
	if session.Nickname != resp.Nickname {
		t.Fatalf("expected nickname %s, got %s", resp.Nickname, session.Nickname)
	}
}

func TestService_ExchangeLoginRequiresCode(t *testing.T) {
	svc := NewService(account.NewMemoryRepository())

	_, err := svc.ExchangeLogin(context.Background(), "   ")
	if !errors.Is(err, ErrCodeRequired) {
		t.Fatalf("expected ErrCodeRequired, got %v", err)
	}
}

func TestService_BootstrapSessionReturnsExistingSession(t *testing.T) {
	ctx := context.Background()
	repo := account.NewMemoryRepository()
	svc := NewService(repo)

	exchanged, err := svc.ExchangeLogin(ctx, "wx-code-boot")
	if err != nil {
		t.Fatalf("expected exchange success, got %v", err)
	}

	bootstrapped, err := svc.BootstrapSession(ctx, exchanged.UnifiedToken, "https://game.xxx.com/play?foo=1")
	if err != nil {
		t.Fatalf("expected bootstrap success, got %v", err)
	}
	if bootstrapped.Token != exchanged.UnifiedToken {
		t.Fatalf("expected token %s, got %s", exchanged.UnifiedToken, bootstrapped.Token)
	}
	if bootstrapped.PlayerID != exchanged.PlayerID {
		t.Fatalf("expected player_id %d, got %d", exchanged.PlayerID, bootstrapped.PlayerID)
	}
	if bootstrapped.Channel != Channel {
		t.Fatalf("expected channel %s, got %s", Channel, bootstrapped.Channel)
	}
	expectedGameURL := "https://game.xxx.com/play?foo=1&channel=wxmini&token=" + exchanged.UnifiedToken
	if bootstrapped.GameURL != expectedGameURL {
		t.Fatalf("expected game_url %s, got %s", expectedGameURL, bootstrapped.GameURL)
	}
}

func TestService_BootstrapSessionNotFound(t *testing.T) {
	svc := NewService(account.NewMemoryRepository())

	_, err := svc.BootstrapSession(context.Background(), "wxmini-not-found", "https://game.xxx.com")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}
