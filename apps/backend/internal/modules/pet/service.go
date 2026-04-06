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

type growthBonusBreakdown struct {
	Bone   int64
	Spirit int64
	Soul   int64
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

	team.Pets = withBasePowerBreakdowns(team.Pets)
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
	roster = withBasePowerBreakdowns(roster)

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
const spiritBonusBaseline int64 = 100
const spiritBonusPerPoint int64 = 1
const soulBonusPerPiece int64 = 8

func (s *Service) applyGrowthBonuses(ctx context.Context, playerID int64, team TeamSnapshot) TeamSnapshot {
	if s.growth == nil || len(team.Pets) == 0 {
		return team
	}

	wallet, err := s.growth.GetWallet(ctx, playerID)
	if err != nil {
		return team
	}

	bonuses := growthBonusParts(wallet)
	if bonuses.total() == 0 {
		return team
	}

	pets := clonePets(team.Pets)
	applyDistributedBonus(pets, bonuses.Bone, func(breakdown *PowerBreakdown, bonus int64) {
		breakdown.Bone += bonus
	})
	applyDistributedBonus(pets, bonuses.Spirit, func(breakdown *PowerBreakdown, bonus int64) {
		breakdown.Spirit += bonus
	})
	applyDistributedBonus(pets, bonuses.Soul, func(breakdown *PowerBreakdown, bonus int64) {
		breakdown.Soul += bonus
	})
	team.Pets = pets
	return team
}

func applyPowerOverrides(pets []BattlePet, overrides map[int64]BattlePet) []BattlePet {
	next := clonePets(pets)
	for idx := range next {
		if override, ok := overrides[next[idx].PetID]; ok {
			next[idx].Power = override.Power
			next[idx].PowerBreakdown = override.PowerBreakdown
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
	return growthBonusParts(wallet).total()
}

func growthBonusParts(wallet growth.Wallet) growthBonusBreakdown {
	boneLevels := max(wallet.BoneLevel-1, 0)
	spiritPower := wallet.SpiritBonusPower
	if spiritPower == 0 {
		spiritPower = wallet.SpiritPower
	}
	spiritPower = max(spiritPower-spiritBonusBaseline, 0)
	soulPieces := max(wallet.SoulPieces, 0)
	return growthBonusBreakdown{
		Bone:   int64(boneLevels) * boneBonusPerLevel,
		Spirit: spiritPower * spiritBonusPerPoint,
		Soul:   int64(soulPieces) * soulBonusPerPiece,
	}
}

func (b growthBonusBreakdown) total() int64 {
	return b.Bone + b.Spirit + b.Soul
}

func withBasePowerBreakdowns(pets []BattlePet) []BattlePet {
	next := clonePets(pets)
	for idx := range next {
		next[idx].PowerBreakdown = basePowerBreakdown(next[idx])
		next[idx].Power = next[idx].PowerBreakdown.Total
	}
	return next
}

func basePowerBreakdown(pet BattlePet) PowerBreakdown {
	base := pet.BasePower
	if base <= 0 {
		base = pet.Power
	}
	level := pet.Power - base
	if level < 0 {
		level = 0
	}
	return PowerBreakdown{
		Base:  base,
		Level: level,
		Total: base + level,
	}
}

func applyDistributedBonus(pets []BattlePet, total int64, assign func(*PowerBreakdown, int64)) {
	if total == 0 || len(pets) == 0 {
		return
	}

	share := total / int64(len(pets))
	remainder := total % int64(len(pets))
	for idx := range pets {
		bonus := share
		if int64(idx) < remainder {
			bonus++
		}
		if bonus == 0 {
			continue
		}
		assign(&pets[idx].PowerBreakdown, bonus)
		pets[idx].PowerBreakdown.Total += bonus
		pets[idx].Power += bonus
	}
}
