package dungeon

import (
	"context"
	"time"
)

func (s *Service) CultivationSnapshot(ctx context.Context, playerID int64) (CultivationStatus, bool, error) {
	state, err := s.repo.StartCultivation(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, false, err
	}
	now := time.Now()
	return state, now.After(state.ClaimableAt), nil
}
