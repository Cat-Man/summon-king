package arena

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetCurrentStreak(ctx context.Context, playerID int64, streak int) error {
	return s.repo.SetCurrentStreak(ctx, playerID, streak)
}

func (s *Service) RecordBattleResult(ctx context.Context, playerID int64, won bool) error {
	return s.repo.RecordBattleResult(ctx, playerID, won)
}

func (s *Service) GetDailyRecord(ctx context.Context, playerID int64) (DailyRecord, error) {
	return s.repo.GetDailyRecord(ctx, playerID)
}
