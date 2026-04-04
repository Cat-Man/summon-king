package pet

import (
	"context"
	"testing"
)

func TestService_GetBattleTeamReturnsStarterPet(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	team, err := svc.GetBattleTeam(ctx, 1001)

	if err != nil {
		t.Fatalf("expected battle team success, got %v", err)
	}
	if team.PlayerID != 1001 {
		t.Fatalf("expected player 1001, got %d", team.PlayerID)
	}
	if len(team.Pets) != 1 {
		t.Fatalf("expected 1 starter pet, got %d", len(team.Pets))
	}
	if team.Pets[0].Slot != 1 {
		t.Fatalf("expected starter slot 1, got %d", team.Pets[0].Slot)
	}
	if team.Pets[0].Power <= 0 {
		t.Fatalf("expected starter power > 0, got %d", team.Pets[0].Power)
	}
}

func TestService_GetBattleTeamAggregatesTeamPower(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	team, err := svc.GetBattleTeam(ctx, 1002)

	if err != nil {
		t.Fatalf("expected battle team success, got %v", err)
	}
	if team.TotalPower != team.Pets[0].Power {
		t.Fatalf("expected total power %d, got %d", team.Pets[0].Power, team.TotalPower)
	}
}
