package dungeon

import (
	"context"
)

type spiritWalletUpdater interface {
	UpdateSpiritPower(ctx context.Context, playerID int64, delta int64) error
}

type Service struct {
	repo    Repository
	wallets spiritWalletUpdater
}

func NewService(repo Repository, wallets ...spiritWalletUpdater) *Service {
	service := &Service{repo: repo}
	if len(wallets) > 0 {
		service.wallets = wallets[0]
	}
	return service
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
	status, err := s.repo.ClaimCultivation(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, err
	}
	if s.wallets != nil && status.SpiritPower > 0 {
		if err := s.wallets.UpdateSpiritPower(ctx, playerID, int64(status.SpiritPower)); err != nil {
			return CultivationStatus{}, err
		}
	}
	return status, nil
}

func (s *Service) GetCultivationStatus(ctx context.Context, playerID int64) (CultivationStatus, error) {
	return s.repo.GetCultivation(ctx, playerID)
}
