package asset

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) GetWallet(ctx context.Context, playerID int64) (Wallet, error) {
	return s.repo.GetWallet(ctx, playerID)
}

func (s *Service) GetInventory(ctx context.Context, playerID int64) ([]InventoryItem, error) {
	return s.repo.ListInventory(ctx, playerID)
}

func (s *Service) GetResourceChangeLogs(ctx context.Context, playerID int64, limit int) ([]ResourceChangeLog, error) {
	return s.repo.ListResourceChangeLogs(ctx, playerID, limit)
}

func (s *Service) GrantReward(ctx context.Context, grant RewardGrant) (GrantRewardResult, error) {
	return s.repo.GrantRewardIdempotent(ctx, grant, s.now())
}

func (s *Service) UseInventory(ctx context.Context, req InventoryOperateRequest) (InventoryOperateResult, error) {
	return s.repo.UseInventoryItem(ctx, req, s.now())
}

func (s *Service) SellInventory(ctx context.Context, req InventoryOperateRequest) (InventoryOperateResult, error) {
	return s.repo.SellInventoryItem(ctx, req, s.now())
}
