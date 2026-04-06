package dungeon

import (
	"context"
)

func (s *Service) CultivationSnapshot(ctx context.Context, playerID int64) (CultivationStatus, bool, error) {
	state, err := s.repo.GetCultivation(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, false, err
	}
	claimable := !state.ClaimableAt.IsZero() && !currentTime().Before(state.ClaimableAt)
	return state, claimable, nil
}
