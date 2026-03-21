package alliance

import "context"

func (s *Service) StartFireTraining(ctx context.Context, allianceID string, playerID int64, roomID string) (FireTrainingRecord, error) {
	return s.repo.StartFireTraining(ctx, allianceID, playerID, roomID, s.now())
}

func (s *Service) ClaimFireTraining(ctx context.Context, allianceID string, playerID int64, recordID string) (FireTrainingRecord, error) {
	record, err := s.repo.GetFireTraining(ctx, allianceID, recordID)
	if err != nil {
		return FireTrainingRecord{}, err
	}
	if record.Status == "claimed" {
		return FireTrainingRecord{}, ErrFireTrainingClaimed
	}
	if s.now().Before(record.FinishedAt) {
		return FireTrainingRecord{}, ErrFireTrainingNotReady
	}
	record.Status = "claimed"
	if err := s.repo.SaveFireTraining(ctx, record); err != nil {
		return FireTrainingRecord{}, err
	}
	return record, nil
}
