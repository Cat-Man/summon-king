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

func (s *Service) Apply(ctx context.Context, playerID int64, req ApplyRequest) (ApplyResult, error) {
	return s.apply(ctx, playerID, req)
}

func (s *Service) ApplyWithMetadata(ctx context.Context, playerID int64, req ApplyRequest) (ApplyResult, error) {
	return s.Apply(ctx, playerID, req)
}

func (s *Service) apply(ctx context.Context, playerID int64, req ApplyRequest) (ApplyResult, error) {
	if req.Delta.SpiritPower != 0 {
		if err := s.growth.UpdateSpiritPower(ctx, playerID, req.Delta.SpiritPower); err != nil {
			return ApplyResult{}, err
		}
	}
	if req.Delta.BoneLevel != 0 {
		if _, err := s.growth.UpgradeBoneLevel(ctx, playerID, req.Delta.BoneLevel); err != nil {
			return ApplyResult{}, err
		}
	}
	if req.Delta.SoulPieces != 0 {
		if _, err := s.growth.UpgradeSoulPower(ctx, playerID, req.Delta.SoulPieces); err != nil {
			return ApplyResult{}, err
		}
	}
	wallet, err := s.growth.GetWallet(ctx, playerID)
	if err != nil {
		return ApplyResult{}, err
	}
	return ApplyResult{
		Wallet:   wallet,
		Metadata: req.Metadata,
	}, nil
}

func (s *Service) Snapshot(ctx context.Context, playerID int64) (growth.Wallet, error) {
	return s.growth.GetWallet(ctx, playerID)
}
