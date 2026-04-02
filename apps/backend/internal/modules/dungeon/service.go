package dungeon

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetWorldMap(ctx context.Context) (WorldMap, error) {
	return s.repo.GetWorldMap(ctx)
}

func (s *Service) Teleport(ctx context.Context, cityID int64) (MapCity, error) {
	return s.repo.TeleportToCity(ctx, cityID)
}

func (s *Service) EnterDungeon(ctx context.Context, playerID, dungeonID int64) (DungeonRun, error) {
	return s.repo.EnterDungeon(ctx, playerID, dungeonID)
}

func (s *Service) RollDice(ctx context.Context, playerID int64) (DungeonRun, error) {
	return s.repo.RollDungeonDice(ctx, playerID)
}

func (s *Service) GetRun(ctx context.Context, playerID int64) (DungeonRun, error) {
	return s.repo.GetDungeonRun(ctx, playerID)
}

func (s *Service) StartCultivation(ctx context.Context, playerID int64) (CultivationStatus, error) {
	return s.repo.StartCultivation(ctx, playerID)
}

func (s *Service) ClaimCultivation(ctx context.Context, playerID int64) (CultivationStatus, error) {
	return s.repo.ClaimCultivation(ctx, playerID)
}
