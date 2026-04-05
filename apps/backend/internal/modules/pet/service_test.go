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

func TestService_SaveTeamAddsSecondActivePet(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	data, err := svc.SaveTeam(ctx, 1001, []int64{10011, 10012})

	if err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}
	if len(data.ActiveTeam) != 2 {
		t.Fatalf("expected active team 2, got %d", len(data.ActiveTeam))
	}
	if data.ActiveTeam[0].PetID != 10011 {
		t.Fatalf("expected slot 1 pet 10011, got %d", data.ActiveTeam[0].PetID)
	}
	if data.ActiveTeam[1].PetID != 10012 {
		t.Fatalf("expected slot 2 pet 10012, got %d", data.ActiveTeam[1].PetID)
	}
	if data.TotalPower != 276 {
		t.Fatalf("expected total power 276, got %d", data.TotalPower)
	}
}

func TestService_SetMainPetReordersExistingActiveTeam(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())
	if _, err := svc.SaveTeam(ctx, 1001, []int64{10011, 10012}); err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}

	data, err := svc.SetMainPet(ctx, 1001, 10012)

	if err != nil {
		t.Fatalf("expected set main pet success, got %v", err)
	}
	if len(data.ActiveTeam) != 2 {
		t.Fatalf("expected active team 2, got %d", len(data.ActiveTeam))
	}
	if data.ActiveTeam[0].PetID != 10012 {
		t.Fatalf("expected slot 1 pet 10012, got %d", data.ActiveTeam[0].PetID)
	}
	if data.ActiveTeam[1].PetID != 10011 {
		t.Fatalf("expected slot 2 pet 10011, got %d", data.ActiveTeam[1].PetID)
	}
	if data.TotalPower != 276 {
		t.Fatalf("expected total power 276, got %d", data.TotalPower)
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

func TestService_SaveTeamDoesNotClearActiveTeamWhenPetMissing(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	_, err := svc.SaveTeam(ctx, 1001, []int64{10011, 99999})
	if err == nil {
		t.Fatal("expected missing pet error")
	}
	if err != ErrPetNotFound {
		t.Fatalf("expected ErrPetNotFound, got %v", err)
	}

	data, err := svc.GetCollection(ctx, 1001)
	if err != nil {
		t.Fatalf("expected collection success after failed save, got %v", err)
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

func TestService_GrantActiveTeamExperienceLevelsMainPet(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	team, err := svc.GrantActiveTeamExperience(ctx, 1001, 100)

	if err != nil {
		t.Fatalf("expected grant pet exp success, got %v", err)
	}
	if len(team.Pets) != 1 {
		t.Fatalf("expected active team 1, got %d", len(team.Pets))
	}
	if team.Pets[0].Level != 2 {
		t.Fatalf("expected level 2 after 100 exp, got %d", team.Pets[0].Level)
	}
	if team.Pets[0].Exp != 0 {
		t.Fatalf("expected remaining exp 0 after level up, got %d", team.Pets[0].Exp)
	}
	if team.Pets[0].NextLevelExp != 200 {
		t.Fatalf("expected next level exp 200, got %d", team.Pets[0].NextLevelExp)
	}
	if team.Pets[0].Power != 144 {
		t.Fatalf("expected level 2 power 144, got %d", team.Pets[0].Power)
	}
	if team.TotalPower != 144 {
		t.Fatalf("expected total power 144, got %d", team.TotalPower)
	}
}

func TestService_GrantActiveTeamExperienceDoesNotAffectInactivePets(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	if _, err := svc.GrantActiveTeamExperience(ctx, 1001, 50); err != nil {
		t.Fatalf("expected grant pet exp success, got %v", err)
	}

	data, err := svc.GetCollection(ctx, 1001)
	if err != nil {
		t.Fatalf("expected collection success, got %v", err)
	}
	if data.ActiveTeam[0].Exp != 50 {
		t.Fatalf("expected main pet exp 50, got %d", data.ActiveTeam[0].Exp)
	}
	if data.ActiveTeam[0].NextLevelExp != 100 {
		t.Fatalf("expected next level exp 100, got %d", data.ActiveTeam[0].NextLevelExp)
	}
	if data.Roster[1].PetID != 10012 {
		t.Fatalf("expected inactive pet 10012, got %d", data.Roster[1].PetID)
	}
	if data.Roster[1].Level != 1 {
		t.Fatalf("expected inactive pet level 1, got %d", data.Roster[1].Level)
	}
	if data.Roster[1].Exp != 0 {
		t.Fatalf("expected inactive pet exp 0, got %d", data.Roster[1].Exp)
	}
	if data.Roster[1].Power != 156 {
		t.Fatalf("expected inactive pet power 156, got %d", data.Roster[1].Power)
	}
}

func TestService_GrantActiveTeamExperienceAppliesToWholeActiveTeam(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())
	if _, err := svc.SaveTeam(ctx, 1001, []int64{10011, 10012}); err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}

	team, err := svc.GrantActiveTeamExperience(ctx, 1001, 100)

	if err != nil {
		t.Fatalf("expected grant pet exp success, got %v", err)
	}
	if len(team.Pets) != 2 {
		t.Fatalf("expected active team 2, got %d", len(team.Pets))
	}
	if team.Pets[0].Level != 2 {
		t.Fatalf("expected pet 1 level 2, got %d", team.Pets[0].Level)
	}
	if team.Pets[1].Level != 2 {
		t.Fatalf("expected pet 2 level 2, got %d", team.Pets[1].Level)
	}
	if team.Pets[0].Power != 144 {
		t.Fatalf("expected pet 1 power 144, got %d", team.Pets[0].Power)
	}
	if team.Pets[1].Power != 180 {
		t.Fatalf("expected pet 2 power 180, got %d", team.Pets[1].Power)
	}
	if team.TotalPower != 324 {
		t.Fatalf("expected total power 324, got %d", team.TotalPower)
	}
}
