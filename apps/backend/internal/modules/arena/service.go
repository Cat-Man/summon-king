package arena

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type assetWriter interface {
	Apply(ctx context.Context, playerID int64, delta asset.Delta) (asset.ApplyResult, error)
	Snapshot(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type Service struct {
	repo  Repository
	asset assetWriter
}

func NewService(repo Repository, assetWriter assetWriter) *Service {
	return &Service{
		repo:  repo,
		asset: assetWriter,
	}
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

	return BattleResult{
		Record:         record,
		RewardDelta:    delta,
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
