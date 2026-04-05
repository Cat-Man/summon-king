package arena

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
)

type assetWriter interface {
	Apply(ctx context.Context, playerID int64, delta asset.Delta) (asset.ApplyResult, error)
	Snapshot(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type teamReader interface {
	GetBattleTeam(ctx context.Context, playerID int64) (pet.TeamSnapshot, error)
}

type battleResolver interface {
	Resolve(ctx context.Context, req battle.Request) (battle.Summary, error)
}

type Option func(*Service)

type Service struct {
	repo    Repository
	asset   assetWriter
	teams   teamReader
	battles battleResolver
}

func WithBattleTeamReader(reader teamReader) Option {
	return func(s *Service) {
		s.teams = reader
	}
}

func NewService(repo Repository, assetWriter assetWriter, options ...Option) *Service {
	svc := &Service{
		repo:    repo,
		asset:   assetWriter,
		teams:   pet.NewService(pet.NewMemoryRepository()),
		battles: battle.NewService(),
	}
	for _, option := range options {
		if option != nil {
			option(svc)
		}
	}
	return svc
}

func (s *Service) SetCurrentStreak(ctx context.Context, playerID int64, streak int) error {
	return s.repo.SetCurrentStreak(ctx, playerID, streak)
}

func (s *Service) RecordBattleResult(ctx context.Context, playerID int64, won bool) (BattleResult, error) {
	attacker, err := s.currentAttackerTeam(ctx, playerID)
	if err != nil {
		return BattleResult{}, err
	}

	return s.finishBattle(ctx, playerID, won, attacker, defaultArenaDefender(attacker.TotalPower, won))
}

func (s *Service) GetIndex(ctx context.Context, playerID int64) (IndexView, error) {
	version, err := s.repo.GetRefreshVersion(ctx, playerID)
	if err != nil {
		return IndexView{}, err
	}

	return IndexView{
		PlayerID:  playerID,
		Opponents: buildOpponents(playerID, version),
	}, nil
}

func (s *Service) RefreshOpponents(ctx context.Context, playerID int64) (IndexView, error) {
	version, err := s.repo.IncrementRefreshVersion(ctx, playerID)
	if err != nil {
		return IndexView{}, err
	}

	return IndexView{
		PlayerID:  playerID,
		Opponents: buildOpponents(playerID, version),
	}, nil
}

func (s *Service) ChallengeOpponent(ctx context.Context, playerID, opponentID int64) (BattleResult, error) {
	index, err := s.GetIndex(ctx, playerID)
	if err != nil {
		return BattleResult{}, err
	}

	var target Opponent
	found := false
	for _, opponent := range index.Opponents {
		if opponent.OpponentID == opponentID {
			target = opponent
			found = true
			break
		}
	}
	if !found {
		return BattleResult{}, ErrArenaRecordNotFound
	}

	attacker, err := s.currentAttackerTeam(ctx, playerID)
	if err != nil {
		return BattleResult{}, err
	}

	won := attacker.TotalPower >= target.Power
	return s.finishBattle(ctx, playerID, won, attacker, opponentTeam(target))
}

func (s *Service) GetDailyRecord(ctx context.Context, playerID int64) (DailyRecord, error) {
	return s.repo.GetDailyRecord(ctx, playerID)
}

func (s *Service) finishBattle(
	ctx context.Context,
	playerID int64,
	won bool,
	attacker pet.TeamSnapshot,
	defender pet.TeamSnapshot,
) (BattleResult, error) {
	if err := s.repo.RecordBattleResult(ctx, playerID, won); err != nil {
		return BattleResult{}, err
	}

	record, err := s.repo.GetDailyRecord(ctx, playerID)
	if err != nil {
		return BattleResult{}, err
	}

	delta := rewardDeltaForResult(won)
	summary, err := s.resolveBattle(ctx, attacker, defender, won)
	if err != nil {
		return BattleResult{}, err
	}

	wallet, err := s.walletSnapshot(ctx, playerID, delta)
	if err != nil {
		return BattleResult{}, err
	}

	return BattleResult{
		Record:         record,
		RewardDelta:    delta,
		BattleResult:   summary,
		WalletSnapshot: wallet,
	}, nil
}

func rewardDeltaForResult(won bool) ArenaRewardDelta {
	if won {
		return ArenaRewardDelta{
			SpiritPower: 18,
			SoulPieces:  1,
		}
	}
	return ArenaRewardDelta{SpiritPower: 6}
}

func (s *Service) walletSnapshot(ctx context.Context, playerID int64, delta ArenaRewardDelta) (growth.Wallet, error) {
	if s.asset == nil {
		return growth.Wallet{}, nil
	}
	result, err := s.asset.Apply(ctx, playerID, asset.Delta{
		SpiritPower: delta.SpiritPower,
		SoulPieces:  delta.SoulPieces,
	})
	if err != nil {
		return growth.Wallet{}, err
	}
	return result.Wallet, nil
}

func (s *Service) currentAttackerTeam(ctx context.Context, playerID int64) (pet.TeamSnapshot, error) {
	if s.teams == nil || s.battles == nil {
		return pet.TeamSnapshot{}, nil
	}

	return s.teams.GetBattleTeam(ctx, playerID)
}

func (s *Service) resolveBattle(
	ctx context.Context,
	attacker pet.TeamSnapshot,
	defender pet.TeamSnapshot,
	won bool,
) (battle.Summary, error) {
	if s.battles == nil {
		return battle.Summary{}, nil
	}

	presetResult := "success"
	if !won {
		presetResult = "fail"
	}

	return s.battles.Resolve(ctx, battle.Request{
		Type:         "arena",
		PresetResult: presetResult,
		Attacker:     attacker,
		Defender:     defender,
	})
}

func buildOpponents(playerID int64, version int) []Opponent {
	baseID := playerID*100 + int64(version*10)
	return []Opponent{
		{
			OpponentID: baseID + 1,
			Name:       "流火镜像",
			Power:      108 + int64(version*6),
		},
		{
			OpponentID: baseID + 2,
			Name:       "寒锋镜像",
			Power:      148 + int64(version*8),
		},
	}
}

func defaultArenaDefender(attackerPower int64, won bool) pet.TeamSnapshot {
	defenderPower := attackerPower - 20
	if !won {
		defenderPower = attackerPower + 20
	}

	return opponentTeam(Opponent{
		OpponentID: 1,
		Name:       "竞技镜像",
		Power:      defenderPower,
	})
}

func opponentTeam(opponent Opponent) pet.TeamSnapshot {
	return pet.TeamSnapshot{
		PlayerID:   0,
		TotalPower: opponent.Power,
		Pets: []pet.BattlePet{
			{
				PetID:    opponent.OpponentID,
				Slot:     1,
				Name:     opponent.Name,
				Level:    1,
				Power:    opponent.Power,
				IsActive: true,
			},
		},
	}
}
