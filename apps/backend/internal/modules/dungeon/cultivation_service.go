package dungeon

import "context"

func (s *Service) StartCultivation(ctx context.Context, playerID int64, mapID string, hours int) (CultivationRecord, error) {
	if hours <= 0 {
		hours = 2
	}
	return s.repo.StartCultivation(ctx, playerID, mapID, hours, s.now())
}

func (s *Service) ClaimCultivation(ctx context.Context, playerID int64, recordID string) (CultivationRecord, error) {
	record, err := s.repo.GetCultivation(ctx, playerID, recordID)
	if err != nil {
		return CultivationRecord{}, err
	}
	if record.Status == "claimed" {
		return CultivationRecord{}, ErrCultivationClaimed
	}
	if s.now().Before(record.FinishedAt) {
		return CultivationRecord{}, ErrCultivationUnfinished
	}
	record.Status = "claimed"
	if err := s.repo.SaveCultivation(ctx, record); err != nil {
		return CultivationRecord{}, err
	}
	return record, nil
}
