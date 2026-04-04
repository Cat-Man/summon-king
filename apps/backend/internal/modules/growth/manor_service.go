package growth

import "context"

type ManorService struct {
	repo Repository
}

func NewManorService(repo Repository) *ManorService {
	return &ManorService{repo: repo}
}

func (s *ManorService) ListPlots(ctx context.Context, playerID int64) []ManorPlot {
	plots, _ := s.repo.GetManorPlots(ctx, playerID)
	return plots
}

func (s *ManorService) Harvest(ctx context.Context, playerID int64) (ManorHarvestResult, error) {
	plots, err := s.repo.HarvestManor(ctx, playerID)
	if err != nil {
		return ManorHarvestResult{}, err
	}
	return ManorHarvestResult{
		Message: "庄园收获完成",
		Plots:   plots,
	}, nil
}

func (s *ManorService) Plant(ctx context.Context, playerID int64) (ManorHarvestResult, error) {
	plots, err := s.repo.PlantManor(ctx, playerID)
	if err != nil {
		return ManorHarvestResult{}, err
	}
	return ManorHarvestResult{
		Message: "庄园种植完成",
		Plots:   plots,
	}, nil
}
