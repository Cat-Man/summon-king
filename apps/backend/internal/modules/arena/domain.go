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

type Opponent struct {
	OpponentID int64  `json:"opponent_id"`
	Name       string `json:"name"`
	Power      int64  `json:"power"`
}

type IndexView struct {
	PlayerID  int64      `json:"player_id"`
	Opponents []Opponent `json:"opponents"`
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
