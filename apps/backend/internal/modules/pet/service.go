package pet

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type growthReader interface {
	GetWallet(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type Option func(*Service)

type Service struct {
	repo   Repository
	growth growthReader
}

func WithGrowthReader(reader growthReader) Option {
	return func(s *Service) {
		s.growth = reader
	}
}

func NewService(repo Repository, options ...Option) *Service {
	svc := &Service{repo: repo}
	for _, option := range options {
		if option != nil {
			option(svc)
		}
	}
	return svc
}

func (s *Service) GetBattleTeam(ctx context.Context, playerID int64) (TeamSnapshot, error) {
	team, err := s.repo.GetBattleTeam(ctx, playerID)
	if err != nil {
		return TeamSnapshot{}, err
	}

	team = s.applyGrowthBonuses(ctx, playerID, team)
	team.TotalPower = totalPower(team.Pets)

	return team, nil
}

func (s *Service) GetCollection(ctx context.Context, playerID int64) (CollectionView, error) {
	team, err := s.GetBattleTeam(ctx, playerID)
	if err != nil {
		return CollectionView{}, err
	}

	roster, err := s.repo.ListPets(ctx, playerID)
	if err != nil {
		return CollectionView{}, err
	}

	return CollectionView{
		PlayerID:   playerID,
		TotalPower: team.TotalPower,
		TeamSize:   len(team.Pets),
		ActiveTeam: team.Pets,
		Roster:     s.applyRosterPowerOverrides(roster, team.Pets),
	}, nil
}

func (s *Service) SetMainPet(ctx context.Context, playerID, petID int64) (CollectionView, error) {
	if err := s.repo.SetMainPet(ctx, playerID, petID); err != nil {
		return CollectionView{}, err
	}
	return s.GetCollection(ctx, playerID)
}

func (s *Service) SaveTeam(ctx context.Context, playerID int64, petIDs []int64) (CollectionView, error) {
	if err := s.repo.SaveTeam(ctx, playerID, petIDs); err != nil {
		return CollectionView{}, err
	}
	return s.GetCollection(ctx, playerID)
}

func (s *Service) GrantActiveTeamExperience(ctx context.Context, playerID int64, exp int64) (TeamSnapshot, error) {
	team, err := s.repo.GrantActiveTeamExperience(ctx, playerID, exp)
	if err != nil {
		return TeamSnapshot{}, err
	}

	team.TotalPower = totalPower(team.Pets)
	return team, nil
}

func totalPower(pets []BattlePet) int64 {
	var total int64
	for _, battlePet := range pets {
		total += battlePet.Power
	}
	return total
}

const boneBonusPerLevel int64 = 24
const soulBonusPerPiece int64 = 8

func (s *Service) applyGrowthBonuses(ctx context.Context, playerID int64, team TeamSnapshot) TeamSnapshot {
	if s.growth == nil || len(team.Pets) == 0 {
		return team
	}

	wallet, err := s.growth.GetWallet(ctx, playerID)
	if err != nil {
		return team
	}

	teamBonus := growthTeamBonus(wallet)
	if teamBonus == 0 {
		return team
	}

	pets := clonePets(team.Pets)
	share := teamBonus / int64(len(pets))
	remainder := teamBonus % int64(len(pets))
	for idx := range pets {
		pets[idx].Power += share
		if int64(idx) < remainder {
			pets[idx].Power++
		}
	}
	team.Pets = pets
	return team
}

func applyPowerOverrides(pets []BattlePet, overrides map[int64]BattlePet) []BattlePet {
	next := clonePets(pets)
	for idx := range next {
		if override, ok := overrides[next[idx].PetID]; ok {
			next[idx].Power = override.Power
		}
	}
	return next
}

func (s *Service) applyRosterPowerOverrides(roster []BattlePet, active []BattlePet) []BattlePet {
	overrides := make(map[int64]BattlePet, len(active))
	for _, battlePet := range active {
		overrides[battlePet.PetID] = battlePet
	}
	return applyPowerOverrides(roster, overrides)
}

func growthTeamBonus(wallet growth.Wallet) int64 {
	boneLevels := max(wallet.BoneLevel-1, 0)
	soulPieces := max(wallet.SoulPieces, 0)
	return int64(boneLevels)*boneBonusPerLevel + int64(soulPieces)*soulBonusPerPiece
}
