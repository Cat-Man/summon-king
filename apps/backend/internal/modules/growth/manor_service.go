package growth

import "context"

func (s *Service) PlantCrop(ctx context.Context, playerID int64, seedID string) (Crop, error) {
	crop := Crop{PlayerID: playerID, SeedID: seedID, Status: "growing", RewardCoin: 180, ReadyAt: cropReadyAt(s.now())}
	s.repo.SaveCrop(ctx, crop)
	return crop, nil
}

func (s *Service) HarvestCrop(ctx context.Context, playerID int64) (Crop, error) {
	crop, ok := s.repo.GetCrop(ctx, playerID)
	if !ok || s.now().Before(crop.ReadyAt) {
		return Crop{}, ErrCropNotReady
	}
	crop.Status = "harvested"
	s.repo.SaveCrop(ctx, crop)
	return crop, nil
}
