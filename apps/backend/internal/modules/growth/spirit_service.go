package growth

import "context"

type SpiritService struct {
	repo Repository
}

func NewSpiritService(repo Repository) *SpiritService {
	return &SpiritService{repo: repo}
}

func (s *SpiritService) GetWallet(ctx context.Context, playerID int64) Wallet {
	wallet, _ := s.repo.GetWallet(ctx, playerID)
	return wallet
}

func (s *SpiritService) WashSpirit(ctx context.Context, playerID int64, opt WashOption) error {
	free, _ := s.repo.BeginFreeWash(ctx, playerID)
	if free {
		return nil
	}
	return s.repo.CommitWash(ctx, playerID, -opt.SpiritIDOrOne())
}

func (opt WashOption) SpiritIDOrOne() int64 {
	if opt.SpiritID > 0 {
		return opt.SpiritID
	}
	return 1
}
