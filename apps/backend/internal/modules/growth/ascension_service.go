package growth

import "context"

func (s *Service) StartAscension(ctx context.Context, playerID, petID int64) (AscensionRecord, error) {
	record := s.repo.GetAscension(ctx, playerID)
	record.PetID = petID
	record.State = "running"
	s.repo.SaveAscension(ctx, record)
	return record, nil
}
