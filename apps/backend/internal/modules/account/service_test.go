package account

import (
	"context"
	"testing"
)

func TestGuestLogin_ReturnsToken(t *testing.T) {
	svc := NewService(NewMemoryRepository())

	resp, err := svc.GuestLogin(context.Background(), "测试玩家")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token not empty")
	}
	if resp.PlayerID == 0 {
		t.Fatal("expected player id not empty")
	}
}
