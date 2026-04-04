package dungeon

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type spiritWalletUpdater interface {
	UpdateSpiritPower(ctx context.Context, playerID int64, delta int64) error
	UpgradeSoulPower(ctx context.Context, playerID int64, delta int) (growth.Wallet, error)
	GetWallet(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type Service struct {
	repo    Repository
	wallets spiritWalletUpdater
}

var rewardRules = map[string]RollReward{
	"ongoing": {
		Label:       "怪物掉落",
		SpiritPower: 5,
		SoulPieces:  0,
	},
	"boss": {
		Label:       "Boss掉落",
		SpiritPower: 8,
		SoulPieces:  1,
	},
	"exhausted": {
		Label:       "无掉落",
		SpiritPower: 0,
		SoulPieces:  0,
	},
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
	run, err := s.repo.EnterDungeon(ctx, playerID, dungeonID)
	if err != nil {
		return DungeonRun{}, err
	}
	return s.attachWalletSnapshot(ctx, playerID, run)
}

func (s *Service) RollDice(ctx context.Context, playerID int64) (DungeonRun, error) {
	run, err := s.repo.RollDungeonDice(ctx, playerID)
	if err != nil {
		return DungeonRun{}, err
	}
	run.LastReward = resolveRollReward(run)
	if s.wallets != nil {
		if err := s.applyReward(ctx, playerID, run.LastReward); err != nil {
			return DungeonRun{}, err
		}
	}
	return s.attachWalletSnapshot(ctx, playerID, run)
}

func (s *Service) GetRun(ctx context.Context, playerID int64) (DungeonRun, error) {
	run, err := s.repo.GetDungeonRun(ctx, playerID)
	if err != nil {
		return DungeonRun{}, err
	}
	return s.attachWalletSnapshot(ctx, playerID, run)
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

func (s *Service) attachWalletSnapshot(ctx context.Context, playerID int64, run DungeonRun) (DungeonRun, error) {
	if s.wallets == nil {
		return run, nil
	}
	wallet, err := s.wallets.GetWallet(ctx, playerID)
	if err != nil {
		return DungeonRun{}, err
	}
	run.WalletSnapshot = wallet
	return run, nil
}

func resolveRollReward(run DungeonRun) RollReward {
	if reward, ok := rewardRules[run.Status]; ok {
		return reward
	}
	return rewardRules["ongoing"]
}

func (s *Service) applyReward(ctx context.Context, playerID int64, reward RollReward) error {
	if reward.SpiritPower > 0 {
		if err := s.wallets.UpdateSpiritPower(ctx, playerID, reward.SpiritPower); err != nil {
			return err
		}
	}
	if reward.SoulPieces > 0 {
		if _, err := s.wallets.UpgradeSoulPower(ctx, playerID, reward.SoulPieces); err != nil {
			return err
		}
	}
	return nil
}
