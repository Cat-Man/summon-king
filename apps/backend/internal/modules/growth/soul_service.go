package growth

import "context"

func (s *Service) HuntSoul(ctx context.Context, playerID int64, soulID string) (Soul, error) {
	soul := s.repo.GetSoul(ctx, playerID, soulID)
	soul.Level++
	s.repo.SaveSoul(ctx, playerID, soul)
	return soul, nil
}
