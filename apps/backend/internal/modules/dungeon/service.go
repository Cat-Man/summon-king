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

type rewardRule struct {
	DungeonID int64
	MinFloor  int
	MaxFloor  int
	Status    string
	Reward    RollReward
}

var rewardRules = []rewardRule{
	{
		Status: "exhausted",
		Reward: RollReward{
			Label:       "无掉落",
			SpiritPower: 0,
			SoulPieces:  0,
		},
	},
	{
		DungeonID: 1,
		MinFloor:  10,
		Status:    "boss",
		Reward: RollReward{
			Label:       "深层妖窟Boss掉落",
			SpiritPower: 10,
			SoulPieces:  2,
		},
	},
	{
		DungeonID: 1,
		MinFloor:  10,
		Status:    "ongoing",
		Reward: RollReward{
			Label:       "深层妖窟掉落",
			SpiritPower: 6,
			SoulPieces:  0,
		},
	},
	{
		DungeonID: 2,
		MinFloor:  10,
		Status:    "boss",
		Reward: RollReward{
			Label:       "寒渊深层Boss掉落",
			SpiritPower: 15,
			SoulPieces:  3,
		},
	},
	{
		DungeonID: 2,
		MaxFloor:  9,
		Status:    "boss",
		Reward: RollReward{
			Label:       "寒渊Boss掉落",
			SpiritPower: 12,
			SoulPieces:  2,
		},
	},
	{
		DungeonID: 2,
		MinFloor:  10,
		Status:    "ongoing",
		Reward: RollReward{
			Label:       "寒渊深层掉落",
			SpiritPower: 9,
			SoulPieces:  0,
		},
	},
	{
		DungeonID: 2,
		MaxFloor:  9,
		Status:    "ongoing",
		Reward: RollReward{
			Label:       "寒渊裂隙掉落",
			SpiritPower: 7,
			SoulPieces:  0,
		},
	},
	{
		Status: "boss",
		Reward: RollReward{
			Label:       "Boss掉落",
			SpiritPower: 8,
			SoulPieces:  1,
		},
	},
	{
		Status: "ongoing",
		Reward: RollReward{
			Label:       "怪物掉落",
			SpiritPower: 5,
			SoulPieces:  0,
		},
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
	if err := s.validateDungeonUnlock(ctx, playerID, dungeonID); err != nil {
		return DungeonRun{}, err
	}
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
	for _, rule := range rewardRules {
		if matchesRewardRule(run, rule) {
			return rule.Reward
		}
	}
	return RollReward{
		Label:       "怪物掉落",
		SpiritPower: 5,
		SoulPieces:  0,
	}
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

func matchesRewardRule(run DungeonRun, rule rewardRule) bool {
	if rule.Status != run.Status {
		return false
	}
	if rule.DungeonID != 0 && run.DungeonID != rule.DungeonID {
		return false
	}
	if rule.MinFloor > 0 && run.CurrentFloor < rule.MinFloor {
		return false
	}
	if rule.MaxFloor > 0 && run.CurrentFloor > rule.MaxFloor {
		return false
	}
	return true
}

func (s *Service) validateDungeonUnlock(ctx context.Context, playerID, dungeonID int64) error {
	if s.wallets == nil {
		return nil
	}

	dungeon, found, err := s.findDungeon(ctx, dungeonID)
	if err != nil || !found || dungeon.UnlockSpiritPower == 0 {
		return err
	}

	wallet, err := s.wallets.GetWallet(ctx, playerID)
	if err != nil {
		return err
	}
	if wallet.SpiritPower < dungeon.UnlockSpiritPower {
		return ErrDungeonLocked
	}
	return nil
}

func (s *Service) findDungeon(ctx context.Context, dungeonID int64) (MapDungeon, bool, error) {
	world, err := s.repo.GetWorldMap(ctx)
	if err != nil {
		return MapDungeon{}, false, err
	}
	for _, city := range world.Cities {
		for _, dungeon := range city.Dungeons {
			if dungeon.DungeonID == dungeonID {
				return dungeon, true, nil
			}
		}
	}
	return MapDungeon{}, false, nil
}
