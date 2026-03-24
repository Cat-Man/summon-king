package dungeon

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) GetWorldMap(ctx context.Context, playerID int64) (WorldMap, error) {
	return s.repo.GetWorldMap(ctx, playerID)
}

func (s *Service) Teleport(ctx context.Context, playerID int64, cityID string) (City, error) {
	return s.repo.Teleport(ctx, playerID, cityID)
}

func (s *Service) EnterDungeon(ctx context.Context, playerID, dungeonID int64) (DungeonRun, error) {
	return s.repo.CreateDungeonRun(ctx, playerID, dungeonID, s.now())
}

func (s *Service) RollDice(ctx context.Context, runID string) (RollResult, error) {
	run, err := s.repo.GetDungeonRun(ctx, runID)
	if err != nil {
		return RollResult{}, err
	}
	if run.RemainDice <= 0 {
		return RollResult{}, ErrNoDiceRemaining
	}

	run.RemainDice--
	run.CurrentFloor += 1
	if err := s.repo.SaveDungeonRun(ctx, run); err != nil {
		return RollResult{}, err
	}

	return RollResult{Run: run, DicePoint: 1, FloorResolved: run.CurrentFloor}, nil
}

func (s *Service) ClaimBossReward(ctx context.Context, runID string) (DungeonRun, error) {
	run, err := s.repo.GetDungeonRun(ctx, runID)
	if err != nil {
		return DungeonRun{}, err
	}
	run.BossRewarded = true
	if err := s.repo.SaveDungeonRun(ctx, run); err != nil {
		return DungeonRun{}, err
	}
	return run, nil
}

func (s *Service) GetLatestCultivation(ctx context.Context, playerID int64) (CultivationRecord, bool, error) {
	return s.repo.GetLatestCultivation(ctx, playerID)
}
