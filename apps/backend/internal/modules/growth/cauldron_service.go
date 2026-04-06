package growth

import "context"

type CauldronService struct {
	repo Repository
}

func NewCauldronService(repo Repository) *CauldronService {
	return &CauldronService{repo: repo}
}

func (s *CauldronService) Stir(ctx context.Context, playerID int64) (string, error) {
	_, err := s.repo.GetWallet(ctx, playerID)
	if err != nil {
		return "", err
	}
	return "炼妖成功", nil
}
