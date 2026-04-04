package battle

import (
	"context"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
)

func TestService_ResolveChoosesHigherPowerWithoutPreset(t *testing.T) {
	ctx := context.Background()
	svc := NewService()

	result, err := svc.Resolve(ctx, Request{
		Type: "tower",
		Attacker: pet.TeamSnapshot{
			PlayerID:   1001,
			TotalPower: 180,
		},
		Defender: pet.TeamSnapshot{
			PlayerID:   0,
			TotalPower: 120,
		},
	})

	if err != nil {
		t.Fatalf("expected resolve success, got %v", err)
	}
	if result.Result != "success" {
		t.Fatalf("expected success, got %s", result.Result)
	}
	if result.BattleNo != "tower-1001-180-120" {
		t.Fatalf("expected battle no tower-1001-180-120, got %s", result.BattleNo)
	}
	if result.WinnerSide != "attacker" {
		t.Fatalf("expected winner side attacker, got %s", result.WinnerSide)
	}
}

func TestService_ResolveUsesPresetOutcomeWhenProvided(t *testing.T) {
	ctx := context.Background()
	svc := NewService()

	result, err := svc.Resolve(ctx, Request{
		Type:         "arena",
		PresetResult: "fail",
		Attacker: pet.TeamSnapshot{
			PlayerID:   1001,
			TotalPower: 180,
		},
		Defender: pet.TeamSnapshot{
			PlayerID:   0,
			TotalPower: 120,
		},
	})

	if err != nil {
		t.Fatalf("expected resolve success, got %v", err)
	}
	if result.Result != "fail" {
		t.Fatalf("expected fail, got %s", result.Result)
	}
	if result.BattleType != "arena" {
		t.Fatalf("expected battle type arena, got %s", result.BattleType)
	}
	if result.BattleNo != "arena-1001-180-120" {
		t.Fatalf("expected battle no arena-1001-180-120, got %s", result.BattleNo)
	}
	if result.WinnerSide != "defender" {
		t.Fatalf("expected winner side defender, got %s", result.WinnerSide)
	}
}
