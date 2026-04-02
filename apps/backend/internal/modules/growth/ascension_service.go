package growth

import "context"

type AscensionService struct {
	repo Repository
}

func NewAscensionService(repo Repository) *AscensionService {
	return &AscensionService{repo: repo}
}

func (s *AscensionService) Elevate(ctx context.Context, playerID int64) string {
	return "化仙完成"
}
