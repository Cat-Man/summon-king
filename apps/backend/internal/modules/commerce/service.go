package commerce

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

var ErrAlreadyClaimedToday = errors.New("already claimed today")

type assetWriter interface {
	Apply(ctx context.Context, playerID int64, req asset.ApplyRequest) (asset.ApplyResult, error)
	Snapshot(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type Option func(*Service)

type Service struct {
	repo  Repository
	asset assetWriter
	now   func() time.Time
}

type SigninIndex struct {
	PlayerID       int64         `json:"player_id"`
	TodayClaimed   bool          `json:"today_claimed"`
	StreakDays     int           `json:"streak_days"`
	TodayReward    asset.Delta   `json:"today_reward"`
	WalletSnapshot growth.Wallet `json:"wallet_snapshot"`
}

type SigninClaimResult struct {
	PlayerID       int64         `json:"player_id"`
	TodayClaimed   bool          `json:"today_claimed"`
	StreakDays     int           `json:"streak_days"`
	RewardDelta    asset.Delta   `json:"reward_delta"`
	WalletSnapshot growth.Wallet `json:"wallet_snapshot"`
}

type VIPIndex struct {
	PlayerID        int64         `json:"player_id"`
	VIPLevel        int           `json:"vip_level"`
	NextVIPLevel    int           `json:"next_vip_level"`
	DailyClaimed    bool          `json:"daily_claimed"`
	DailyReward     asset.Delta   `json:"daily_reward"`
	CurrentBenefits []string      `json:"current_benefits"`
	NextBenefits    []string      `json:"next_benefits"`
	WalletSnapshot  growth.Wallet `json:"wallet_snapshot"`
}

type VIPClaimResult struct {
	PlayerID        int64         `json:"player_id"`
	VIPLevel        int           `json:"vip_level"`
	NextVIPLevel    int           `json:"next_vip_level"`
	DailyClaimed    bool          `json:"daily_claimed"`
	DailyReward     asset.Delta   `json:"daily_reward"`
	RewardDelta     asset.Delta   `json:"reward_delta"`
	CurrentBenefits []string      `json:"current_benefits"`
	NextBenefits    []string      `json:"next_benefits"`
	WalletSnapshot  growth.Wallet `json:"wallet_snapshot"`
}

const dateLayout = "2006-01-02"

func WithNow(now func() time.Time) Option {
	return func(s *Service) {
		s.now = now
	}
}

func NewService(repo Repository, assetWriter assetWriter, options ...Option) *Service {
	svc := &Service{
		repo:  repo,
		asset: assetWriter,
		now:   time.Now,
	}
	for _, option := range options {
		if option != nil {
			option(svc)
		}
	}
	return svc
}

func (s *Service) GetSigninIndex(ctx context.Context, playerID int64) (SigninIndex, error) {
	state, err := s.repo.GetState(ctx, playerID)
	if err != nil {
		return SigninIndex{}, err
	}
	wallet, err := s.asset.Snapshot(ctx, playerID)
	if err != nil {
		return SigninIndex{}, err
	}

	today := formatDate(s.now())
	todayClaimed := state.SigninLastDate == today
	streakDays := state.SigninStreak
	if !todayClaimed {
		streakDays = nextSigninStreak(state, s.now())
	}

	return SigninIndex{
		PlayerID:       playerID,
		TodayClaimed:   todayClaimed,
		StreakDays:     streakDays,
		TodayReward:    signinReward(streakDays),
		WalletSnapshot: wallet,
	}, nil
}

func (s *Service) ClaimSignin(ctx context.Context, playerID int64) (SigninClaimResult, error) {
	state, err := s.repo.GetState(ctx, playerID)
	if err != nil {
		return SigninClaimResult{}, err
	}

	now := s.now()
	today := formatDate(now)
	if state.SigninLastDate == today {
		return SigninClaimResult{}, ErrAlreadyClaimedToday
	}

	nextStreak := nextSigninStreak(state, now)
	reward := signinReward(nextStreak)
	applyResult, err := s.asset.Apply(ctx, playerID, asset.ApplyRequest{
		Delta: reward,
		Metadata: asset.ApplyMetadata{
			Source:         "signin",
			Reason:         "daily_signin",
			IdempotencyKey: buildDailyKey("signin", playerID, today),
		},
	})
	if err != nil {
		return SigninClaimResult{}, err
	}

	state.SigninLastDate = today
	state.SigninStreak = nextStreak
	if err := s.repo.SaveState(ctx, playerID, state); err != nil {
		return SigninClaimResult{}, err
	}

	return SigninClaimResult{
		PlayerID:       playerID,
		TodayClaimed:   true,
		StreakDays:     nextStreak,
		RewardDelta:    reward,
		WalletSnapshot: applyResult.Wallet,
	}, nil
}

func (s *Service) GetVIPIndex(ctx context.Context, playerID int64) (VIPIndex, error) {
	state, err := s.repo.GetState(ctx, playerID)
	if err != nil {
		return VIPIndex{}, err
	}
	wallet, err := s.asset.Snapshot(ctx, playerID)
	if err != nil {
		return VIPIndex{}, err
	}

	return VIPIndex{
		PlayerID:        playerID,
		VIPLevel:        state.VIPLevel,
		NextVIPLevel:    state.VIPLevel + 1,
		DailyClaimed:    state.VIPDailyDate == formatDate(s.now()),
		DailyReward:     vipDailyReward(state.VIPLevel),
		CurrentBenefits: vipBenefits(state.VIPLevel),
		NextBenefits:    vipBenefits(state.VIPLevel + 1),
		WalletSnapshot:  wallet,
	}, nil
}

func (s *Service) ClaimVIPDaily(ctx context.Context, playerID int64) (VIPClaimResult, error) {
	state, err := s.repo.GetState(ctx, playerID)
	if err != nil {
		return VIPClaimResult{}, err
	}

	today := formatDate(s.now())
	if state.VIPDailyDate == today {
		return VIPClaimResult{}, ErrAlreadyClaimedToday
	}

	reward := vipDailyReward(state.VIPLevel)
	applyResult, err := s.asset.Apply(ctx, playerID, asset.ApplyRequest{
		Delta: reward,
		Metadata: asset.ApplyMetadata{
			Source:         "vip",
			Reason:         "vip_daily_chest",
			IdempotencyKey: buildDailyKey("vip_daily", playerID, today),
		},
	})
	if err != nil {
		return VIPClaimResult{}, err
	}

	state.VIPDailyDate = today
	if err := s.repo.SaveState(ctx, playerID, state); err != nil {
		return VIPClaimResult{}, err
	}

	return VIPClaimResult{
		PlayerID:        playerID,
		VIPLevel:        state.VIPLevel,
		NextVIPLevel:    state.VIPLevel + 1,
		DailyClaimed:    true,
		DailyReward:     reward,
		RewardDelta:     reward,
		CurrentBenefits: vipBenefits(state.VIPLevel),
		NextBenefits:    vipBenefits(state.VIPLevel + 1),
		WalletSnapshot:  applyResult.Wallet,
	}, nil
}

func nextSigninStreak(state State, now time.Time) int {
	if state.SigninLastDate == "" {
		return 1
	}

	lastDate, err := time.Parse(dateLayout, state.SigninLastDate)
	if err != nil {
		return 1
	}
	if formatDate(lastDate.AddDate(0, 0, 1)) == formatDate(now) {
		return state.SigninStreak + 1
	}
	return 1
}

func signinReward(streakDays int) asset.Delta {
	reward := asset.Delta{SpiritPower: 8}
	if streakDays >= 5 {
		reward.SpiritPower *= 2
	}
	return reward
}

func vipDailyReward(vipLevel int) asset.Delta {
	reward := asset.Delta{SpiritPower: int64(6 + vipLevel*4)}
	if vipLevel >= 1 {
		reward.SoulPieces = 1
	}
	return reward
}

func vipBenefits(vipLevel int) []string {
	switch {
	case vipLevel <= 0:
		return []string{
			"每日宝箱：灵力 +6",
			"副本次数：1 次/天",
		}
	case vipLevel == 1:
		return []string{
			"每日宝箱：灵力 +10，灵魂碎片 +1",
			"战灵免费洗炼次数 +1",
		}
	default:
		return []string{
			"每日宝箱持续增强",
			"副本次数与修行便利性提升",
		}
	}
}

func formatDate(now time.Time) string {
	return now.UTC().Format(dateLayout)
}

func buildDailyKey(prefix string, playerID int64, date string) string {
	return prefix + ":" + date + ":" + strconv.FormatInt(playerID, 10)
}
