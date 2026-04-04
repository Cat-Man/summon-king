package tower

import (
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type TowerStatus struct {
	Tower               string `json:"tower"`
	Label               string `json:"label"`
	CurrentFloor        int    `json:"current_floor"`
	MaxFloor            int    `json:"max_floor"`
	RemainingChallenges int    `json:"remaining_challenges"`
	RewardPreview       string `json:"reward_preview"`
}

type TowerRewardDelta struct {
	BoneLevel   int   `json:"bone_level"`
	SpiritPower int64 `json:"spirit_power"`
	SoulPieces  int   `json:"soul_pieces"`
}

type TowerResult struct {
	PlayerID            int64            `json:"player_id"`
	Tower               string           `json:"tower"`
	Floor               int              `json:"floor"`
	Reward              string           `json:"reward"`
	RemainingChallenges int              `json:"remaining_challenges"`
	RewardDelta         TowerRewardDelta `json:"reward_delta"`
	BattleResult        battle.Summary   `json:"battle"`
	WalletSnapshot      growth.Wallet    `json:"wallet_snapshot"`
}
