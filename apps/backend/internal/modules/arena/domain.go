package arena

import (
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type DailyRecord struct {
	PlayerID      int64 `json:"player_id"`
	CurrentStreak int   `json:"current_streak"`
	LastWin       bool  `json:"last_win"`
}

type ArenaRewardDelta struct {
	SpiritPower int64 `json:"spirit_power"`
	SoulPieces  int   `json:"soul_pieces"`
}

type BattleResult struct {
	Record         DailyRecord      `json:"record"`
	RewardDelta    ArenaRewardDelta `json:"reward_delta"`
	BattleResult   battle.Summary   `json:"battle"`
	WalletSnapshot growth.Wallet    `json:"wallet_snapshot"`
}
