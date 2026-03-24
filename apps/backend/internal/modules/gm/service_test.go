package gm

import (
	"context"
	"testing"
)

func TestGrantCurrency_AppendsAuditLog(t *testing.T) {
	svc := NewService()
	ctx := context.Background()

	err := svc.GrantCurrency(ctx, GrantCurrencyRequest{
		OperatorID: 9001,
		PlayerID:   1001,
		Currency:   "gold",
		Amount:     888,
		Reason:     "compensate",
		TicketNo:   "TK-001",
	})
	if err != nil {
		t.Fatalf("expected grant currency success, got %v", err)
	}

	state := svc.GetPlayerOpsState(ctx, 1001)
	if state.Wallet["gold"] != 888 {
		t.Fatalf("expected gold=888, got %d", state.Wallet["gold"])
	}

	logs := svc.ListAuditLogs(ctx)
	if len(logs) != 1 {
		t.Fatalf("expected 1 audit log, got %d", len(logs))
	}
	if logs[0].Action != "grant_currency" {
		t.Fatalf("expected action grant_currency, got %s", logs[0].Action)
	}
}
