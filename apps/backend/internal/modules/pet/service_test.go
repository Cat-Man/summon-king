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

func TestService_GetCollectionReturnsBattleTeamAndRoster(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	data, err := svc.GetCollection(ctx, 1001)

	if err != nil {
		t.Fatalf("expected collection success, got %v", err)
	}
	if data.PlayerID != 1001 {
		t.Fatalf("expected player 1001, got %d", data.PlayerID)
	}
	if data.TeamSize != 1 {
		t.Fatalf("expected team size 1, got %d", data.TeamSize)
	}
	if data.TotalPower != 120 {
		t.Fatalf("expected total power 120, got %d", data.TotalPower)
	}
	if len(data.ActiveTeam) != 1 {
		t.Fatalf("expected active team 1, got %d", len(data.ActiveTeam))
	}
	if len(data.Roster) != 2 {
		t.Fatalf("expected roster 2, got %d", len(data.Roster))
	}
	if data.Roster[0].Name != "初始灵狐" {
		t.Fatalf("expected starter pet 初始灵狐, got %s", data.Roster[0].Name)
	}
}

func TestService_SetMainPetSwitchesActiveTeam(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	data, err := svc.SetMainPet(ctx, 1001, 10012)

	if err != nil {
		t.Fatalf("expected set main pet success, got %v", err)
	}
	if len(data.ActiveTeam) != 1 {
		t.Fatalf("expected active team 1, got %d", len(data.ActiveTeam))
	}
	if data.ActiveTeam[0].PetID != 10012 {
		t.Fatalf("expected main pet 10012, got %d", data.ActiveTeam[0].PetID)
	}
	if data.ActiveTeam[0].Name != "玄甲龟" {
		t.Fatalf("expected main pet 玄甲龟, got %s", data.ActiveTeam[0].Name)
	}
	if data.TotalPower != 156 {
		t.Fatalf("expected total power 156, got %d", data.TotalPower)
	}
}

func TestService_SetMainPetDoesNotClearActiveTeamWhenPetMissing(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	_, err := svc.SetMainPet(ctx, 1001, 99999)
	if err == nil {
		t.Fatal("expected missing pet error")
	}
	if err != ErrPetNotFound {
		t.Fatalf("expected ErrPetNotFound, got %v", err)
	}

	data, err := svc.GetCollection(ctx, 1001)
	if err != nil {
		t.Fatalf("expected collection success after failed switch, got %v", err)
	}
	if len(data.ActiveTeam) != 1 {
		t.Fatalf("expected active team to remain 1, got %d", len(data.ActiveTeam))
	}
	if data.ActiveTeam[0].PetID != 10011 {
		t.Fatalf("expected original main pet 10011, got %d", data.ActiveTeam[0].PetID)
	}
	if data.TotalPower != 120 {
		t.Fatalf("expected total power 120, got %d", data.TotalPower)
	}
}
