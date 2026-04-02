package growth

import "context"

type ManorService struct {
	repo Repository
}

func NewManorService(repo Repository) *ManorService {
	return &ManorService{repo: repo}
}

func (s *ManorService) ListPlots(ctx context.Context, playerID int64) []ManorPlot {
	return []ManorPlot{
		{PlotID: 1, State: "空闲"},
		{PlotID: 2, State: "成长中"},
	}
}
