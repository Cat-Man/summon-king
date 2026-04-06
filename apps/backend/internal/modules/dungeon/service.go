package dungeon

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
)

type assetWriter interface {
	Apply(ctx context.Context, playerID int64, req asset.ApplyRequest) (asset.ApplyResult, error)
	Snapshot(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type teamReader interface {
	GetBattleTeam(ctx context.Context, playerID int64) (pet.TeamSnapshot, error)
}

type petProgressor interface {
	GrantActiveTeamExperience(ctx context.Context, playerID int64, exp int64) (pet.TeamSnapshot, error)
}

type battleResolver interface {
	Resolve(ctx context.Context, req battle.Request) (battle.Summary, error)
}

type Option func(*Service)

type Service struct {
	repo       Repository
	asset      assetWriter
	teams      teamReader
	progressor petProgressor
	battles    battleResolver
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

func WithBattleTeamReader(reader teamReader) Option {
	return func(s *Service) {
		s.teams = reader
	}
}

func WithPetProgressor(progressor petProgressor) Option {
	return func(s *Service) {
		s.progressor = progressor
	}
}

func NewService(repo Repository, assetWriter assetWriter, options ...Option) *Service {
	svc := &Service{
		repo:       repo,
		asset:      assetWriter,
		teams:      pet.NewService(pet.NewMemoryRepository()),
		progressor: nil,
		battles:    battle.NewService(),
	}
	for _, option := range options {
		if option != nil {
			option(svc)
		}
	}
	return svc
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
	run, err = s.attachBattleSummary(ctx, playerID, run)
	if err != nil {
		return DungeonRun{}, err
	}
	if s.asset != nil {
		if err := s.applyReward(ctx, playerID, run.LastReward); err != nil {
			return DungeonRun{}, err
		}
	}
	run, err = s.attachRollPetGrowth(ctx, playerID, run)
	if err != nil {
		return DungeonRun{}, err
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
	current, err := s.repo.GetCultivation(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, err
	}
	if current.State != "cultivating" || current.ClaimableAt.IsZero() || currentTime().Before(current.ClaimableAt) {
		return CultivationStatus{}, ErrCultivationNotReady
	}

	status, err := s.repo.ClaimCultivation(ctx, playerID)
	if err != nil {
		return CultivationStatus{}, err
	}
	if s.asset != nil && status.SpiritPower > 0 {
		if _, err := s.asset.Apply(ctx, playerID, asset.ApplyRequest{
			Delta: asset.Delta{SpiritPower: int64(status.SpiritPower)},
			Metadata: asset.ApplyMetadata{
				Source:         "dungeon",
				Reason:         "cultivation_claim",
				IdempotencyKey: "dungeon-cultivation",
			},
		}); err != nil {
			return CultivationStatus{}, err
		}
	}
	status, err = s.attachCultivationPetGrowth(ctx, playerID, status)
	if err != nil {
		return CultivationStatus{}, err
	}
	return status, nil
}

func (s *Service) GetCultivationStatus(ctx context.Context, playerID int64) (CultivationStatus, error) {
	return s.repo.GetCultivation(ctx, playerID)
}

func (s *Service) attachWalletSnapshot(ctx context.Context, playerID int64, run DungeonRun) (DungeonRun, error) {
	if s.asset == nil {
		return run, nil
	}
	wallet, err := s.asset.Snapshot(ctx, playerID)
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
	if s.asset == nil {
		return nil
	}
	_, err := s.asset.Apply(ctx, playerID, asset.ApplyRequest{
		Delta: asset.Delta{
			SpiritPower: reward.SpiritPower,
			SoulPieces:  reward.SoulPieces,
		},
		Metadata: asset.ApplyMetadata{
			Source:         "dungeon",
			Reason:         reward.Label,
			IdempotencyKey: "dungeon-roll",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) attachBattleSummary(ctx context.Context, playerID int64, run DungeonRun) (DungeonRun, error) {
	if run.Status == "exhausted" || s.teams == nil || s.battles == nil {
		return run, nil
	}

	attacker, err := s.teams.GetBattleTeam(ctx, playerID)
	if err != nil {
		return DungeonRun{}, err
	}

	summary, err := s.battles.Resolve(ctx, battle.Request{
		Type:         "dungeon",
		PresetResult: "success",
		Attacker:     attacker,
		Defender:     dungeonEnemyTeam(run),
	})
	if err != nil {
		return DungeonRun{}, err
	}

	run.BattleResult = summary
	return run, nil
}

func dungeonEnemyTeam(run DungeonRun) pet.TeamSnapshot {
	power := int64(90 + run.CurrentFloor*12)
	if run.DungeonID == 2 {
		power += 24
	}
	if run.Status == "boss" {
		power += 36
	}

	level := run.CurrentFloor
	if level < 1 {
		level = 1
	}

	return pet.TeamSnapshot{
		PlayerID:   0,
		TotalPower: power,
		Pets: []pet.BattlePet{
			{
				PetID:    run.DungeonID*1000 + int64(level),
				Slot:     1,
				Name:     run.LastReward.Label,
				Level:    level,
				Power:    power,
				IsActive: true,
			},
		},
	}
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

func (s *Service) attachRollPetGrowth(ctx context.Context, playerID int64, run DungeonRun) (DungeonRun, error) {
	exp := resolveRollPetExp(run)
	if exp == 0 || s.progressor == nil {
		return run, nil
	}

	team, err := s.progressor.GrantActiveTeamExperience(ctx, playerID, exp)
	if err != nil {
		return DungeonRun{}, err
	}
	run.PetGrowth = PetGrowth{
		Exp:            exp,
		TeamTotalPower: team.TotalPower,
	}
	return run, nil
}

func (s *Service) attachCultivationPetGrowth(ctx context.Context, playerID int64, status CultivationStatus) (CultivationStatus, error) {
	exp := resolveCultivationPetExp(status)
	if exp == 0 || s.progressor == nil {
		return status, nil
	}

	team, err := s.progressor.GrantActiveTeamExperience(ctx, playerID, exp)
	if err != nil {
		return CultivationStatus{}, err
	}
	status.PetGrowth = PetGrowth{
		Exp:            exp,
		TeamTotalPower: team.TotalPower,
	}
	return status, nil
}

func resolveRollPetExp(run DungeonRun) int64 {
	switch run.Status {
	case "exhausted":
		return 0
	case "boss":
		return 100
	default:
		return 50
	}
}

func resolveCultivationPetExp(status CultivationStatus) int64 {
	if status.SpiritPower <= 0 {
		return 0
	}
	return 100
}

func (s *Service) validateDungeonUnlock(ctx context.Context, playerID, dungeonID int64) error {
	if s.asset == nil {
		return nil
	}

	dungeon, found, err := s.findDungeon(ctx, dungeonID)
	if err != nil || !found || dungeon.UnlockSpiritPower == 0 {
		return err
	}

	wallet, err := s.asset.Snapshot(ctx, playerID)
	if err != nil {
		return err
	}
	if wallet.SpiritPower < dungeon.UnlockSpiritPower {
		return ErrDungeonLocked
	}
	if dungeon.UnlockBoneLevel > 0 && wallet.BoneLevel < dungeon.UnlockBoneLevel {
		return ErrDungeonLocked
	}
	if dungeon.UnlockSoulPieces > 0 && wallet.SoulPieces < dungeon.UnlockSoulPieces {
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
