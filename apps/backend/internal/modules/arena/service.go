package arena

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
)

type assetWriter interface {
	Apply(ctx context.Context, playerID int64, delta asset.Delta) (asset.ApplyResult, error)
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

func (s *Service) SetCurrentStreak(ctx context.Context, playerID int64, streak int) error {
	return s.repo.SetCurrentStreak(ctx, playerID, streak)
}

func (s *Service) RecordBattleResult(ctx context.Context, playerID int64, won bool) (BattleResult, error) {
	if err := s.repo.RecordBattleResult(ctx, playerID, won); err != nil {
		return BattleResult{}, err
	}

	record, err := s.repo.GetDailyRecord(ctx, playerID)
	if err != nil {
		return BattleResult{}, err
	}

	delta := ArenaRewardDelta{}
	if won {
		delta.SpiritPower = 18
		delta.SoulPieces = 1
	} else {
		delta.SpiritPower = 6
	}

	wallet, err := s.walletSnapshot(ctx, playerID, delta)
	if err != nil {
		return BattleResult{}, err
	}

	summary, err := s.battleSummary(ctx, playerID, won)
	if err != nil {
		return BattleResult{}, err
	}

	return BattleResult{
		Record:         record,
		RewardDelta:    delta,
		BattleResult:   summary,
		WalletSnapshot: wallet,
	}, nil
}

func (s *Service) GetDailyRecord(ctx context.Context, playerID int64) (DailyRecord, error) {
	return s.repo.GetDailyRecord(ctx, playerID)
}

func (s *Service) walletSnapshot(ctx context.Context, playerID int64, delta ArenaRewardDelta) (growth.Wallet, error) {
	if s.asset == nil {
		return growth.Wallet{}, nil
	}
	result, err := s.asset.Apply(ctx, playerID, asset.Delta{
		SpiritPower: delta.SpiritPower,
		SoulPieces:  delta.SoulPieces,
	})
	if err != nil {
		return growth.Wallet{}, err
	}
	return result.Wallet, nil
}

func (s *Service) battleSummary(ctx context.Context, playerID int64, won bool) (battle.Summary, error) {
	if s.teams == nil || s.battles == nil {
		return battle.Summary{}, nil
	}

	attacker, err := s.teams.GetBattleTeam(ctx, playerID)
	if err != nil {
		return battle.Summary{}, err
	}

	defenderPower := attacker.TotalPower - 20
	presetResult := "success"
	if !won {
		defenderPower = attacker.TotalPower + 20
		presetResult = "fail"
	}

	return s.battles.Resolve(ctx, battle.Request{
		Type:         "arena",
		PresetResult: presetResult,
		Attacker:     attacker,
		Defender: pet.TeamSnapshot{
			PlayerID:   0,
			TotalPower: defenderPower,
			Pets: []pet.BattlePet{
				{
					PetID:    1,
					Slot:     1,
					Name:     "竞技镜像",
					Level:    1,
					Power:    defenderPower,
					IsActive: true,
				},
			},
		},
	})
}
