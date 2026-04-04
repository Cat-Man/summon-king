package tower

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
	wallet, err := s.walletSnapshot(ctx, playerID, delta)
	if err != nil {
		return TowerResult{}, err
	}
	result.WalletSnapshot = wallet
	return result, nil
}

func (s *Service) GetStatus(ctx context.Context, playerID int64, tower string) TowerStatus {
	return s.repo.GetStatus(ctx, playerID, tower)
}

func (s *Service) walletSnapshot(ctx context.Context, playerID int64, delta TowerRewardDelta) (growth.Wallet, error) {
	if s.asset == nil {
		return growth.Wallet{}, nil
	}
	if delta.BoneLevel != 0 || delta.SpiritPower != 0 || delta.SoulPieces != 0 {
		result, err := s.asset.Apply(ctx, playerID, asset.Delta{
			SpiritPower: delta.SpiritPower,
			BoneLevel:   delta.BoneLevel,
			SoulPieces:  delta.SoulPieces,
		})
		if err != nil {
			return growth.Wallet{}, err
		}
		return result.Wallet, nil
	}
	return s.asset.Snapshot(ctx, playerID)
}
