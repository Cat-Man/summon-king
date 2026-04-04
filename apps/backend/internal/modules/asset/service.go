package asset

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type Service struct {
	growth growth.Repository
}

func NewService(growthRepo growth.Repository) *Service {
	return &Service{growth: growthRepo}
}

func (s *Service) Apply(ctx context.Context, playerID int64, delta Delta) (ApplyResult, error) {
	if delta.SpiritPower != 0 {
		if err := s.growth.UpdateSpiritPower(ctx, playerID, delta.SpiritPower); err != nil {
			return ApplyResult{}, err
		}
	}
	if delta.BoneLevel != 0 {
		if _, err := s.growth.UpgradeBoneLevel(ctx, playerID, delta.BoneLevel); err != nil {
			return ApplyResult{}, err
		}
	}
	if delta.SoulPieces != 0 {
		if _, err := s.growth.UpgradeSoulPower(ctx, playerID, delta.SoulPieces); err != nil {
			return ApplyResult{}, err
		}
	}
	wallet, err := s.growth.GetWallet(ctx, playerID)
	if err != nil {
		return ApplyResult{}, err
	}
	return ApplyResult{Wallet: wallet}, nil
}

func (s *Service) Snapshot(ctx context.Context, playerID int64) (growth.Wallet, error) {
	return s.growth.GetWallet(ctx, playerID)
}
