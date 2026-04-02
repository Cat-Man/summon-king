package tower

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) StartChallenge(ctx context.Context, playerID int64, tower string) (TowerResult, error) {
	return s.repo.StartChallenge(ctx, playerID, tower)
}

func (s *Service) GetInfo(ctx context.Context, tower string) TowerInfo {
	return s.repo.GetInfo(ctx, tower)
}
