package growth

import (
	"context"
)

type SoulService struct {
	repo Repository
}

func NewSoulService(repo Repository) *SoulService {
	return &SoulService{repo: repo}
}

func (s *SoulService) GetSoulState(ctx context.Context, playerID int64) Soul {
	wallet, _ := s.repo.GetWallet(ctx, playerID)
	return Soul{
		Name:  "魔魂",
		Power: wallet.SoulPieces,
	}
}

func (s *SoulService) Upgrade(ctx context.Context, playerID int64) (Soul, error) {
	wallet, err := s.repo.UpgradeSoulPower(ctx, playerID, 1)
	if err != nil {
		return Soul{}, err
	}
	return Soul{
		Name:  "魔魂",
		Power: wallet.SoulPieces,
	}, nil
}
