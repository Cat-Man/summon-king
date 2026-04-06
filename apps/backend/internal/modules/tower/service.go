package tower

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
)

type assetWriter interface {
	Apply(ctx context.Context, playerID int64, req asset.ApplyRequest) (asset.ApplyResult, error)
	Snapshot(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type teamReader interface {
	GetBattleTeam(ctx context.Context, playerID int64) (pet.TeamSnapshot, error)
}

type battleResolver interface {
	Resolve(ctx context.Context, req battle.Request) (battle.Summary, error)
}

type Option func(*Service)

type Service struct {
	repo    Repository
	asset   assetWriter
	teams   teamReader
	battles battleResolver
}

func WithBattleTeamReader(reader teamReader) Option {
	return func(s *Service) {
		s.teams = reader
	}
}

func NewService(repo Repository, assetWriter assetWriter, options ...Option) *Service {
	svc := &Service{
		repo:    repo,
		asset:   assetWriter,
		teams:   pet.NewService(pet.NewMemoryRepository()),
		battles: battle.NewService(),
	}
	for _, option := range options {
		if option != nil {
			option(svc)
		}
	}
	return svc
}

func (s *Service) StartChallenge(ctx context.Context, playerID int64, tower string) (TowerResult, error) {
	result, err := s.repo.StartChallenge(ctx, playerID, tower)
	if err != nil {
		return TowerResult{}, err
	}

	var delta TowerRewardDelta
	switch tower {
	case "pagoda":
		delta.BoneLevel = 1
	case "spirit":
		delta.SpiritPower = 12
		delta.SoulPieces = 1
	}

	result.RewardDelta = delta
	result.BattleResult, err = s.battleSummary(ctx, playerID, result)
	if err != nil {
		return TowerResult{}, err
	}
	wallet, err := s.walletSnapshot(ctx, playerID, tower, delta)
	if err != nil {
		return TowerResult{}, err
	}
	result.WalletSnapshot = wallet
	if err := s.repo.SaveLastReward(ctx, playerID, tower, result.Reward, result.RewardDelta); err != nil {
		return TowerResult{}, err
	}
	return result, nil
}

func (s *Service) GetStatus(ctx context.Context, playerID int64, tower string) TowerStatus {
	return s.repo.GetStatus(ctx, playerID, tower)
}

func (s *Service) battleSummary(ctx context.Context, playerID int64, result TowerResult) (battle.Summary, error) {
	if s.teams == nil || s.battles == nil {
		return battle.Summary{}, nil
	}

	attacker, err := s.teams.GetBattleTeam(ctx, playerID)
	if err != nil {
		return battle.Summary{}, err
	}

	return s.battles.Resolve(ctx, battle.Request{
		Type:         "tower",
		PresetResult: "success",
		Attacker:     attacker,
		Defender:     towerEnemyTeam(result),
	})
}

func (s *Service) walletSnapshot(ctx context.Context, playerID int64, tower string, delta TowerRewardDelta) (growth.Wallet, error) {
	if s.asset == nil {
		return growth.Wallet{}, nil
	}
	if delta.BoneLevel != 0 || delta.SpiritPower != 0 || delta.SoulPieces != 0 {
		result, err := s.asset.Apply(ctx, playerID, asset.ApplyRequest{
			Delta: asset.Delta{
				SpiritPower: delta.SpiritPower,
				BoneLevel:   delta.BoneLevel,
				SoulPieces:  delta.SoulPieces,
			},
			Metadata: asset.ApplyMetadata{
				Source:         "tower",
				Reason:         tower,
				IdempotencyKey: "tower-" + tower,
			},
		})
		if err != nil {
			return growth.Wallet{}, err
		}
		return result.Wallet, nil
	}
	return s.asset.Snapshot(ctx, playerID)
}

func towerEnemyTeam(result TowerResult) pet.TeamSnapshot {
	power := int64(100 + result.Floor*16)
	if result.Tower == "spirit" {
		power += 24
	}

	return pet.TeamSnapshot{
		PlayerID:   0,
		TotalPower: power,
		Pets: []pet.BattlePet{
			{
				PetID:    int64(result.Floor),
				Slot:     1,
				Name:     result.Reward,
				Level:    max(result.Floor, 1),
				Power:    power,
				IsActive: true,
			},
		},
	}
}
