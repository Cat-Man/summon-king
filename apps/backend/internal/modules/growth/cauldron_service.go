package growth

import "context"

func (s *Service) RefineCauldron(ctx context.Context, playerID int64) (CauldronRecord, error) {
	record := s.repo.GetCauldron(ctx, playerID)
	record.RefineCount++
	record.CurrentQuality++
	s.repo.SaveCauldron(ctx, record)
	return record, nil
}
