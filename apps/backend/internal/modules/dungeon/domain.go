package dungeon

import (
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type DungeonRun struct {
	PlayerID       int64          `json:"player_id"`
	DungeonID      int64          `json:"dungeon_id"`
	RemainDice     int            `json:"remain_dice"`
	CurrentFloor   int            `json:"current_floor"`
	Status         string         `json:"status"`
	StartedAt      time.Time      `json:"started_at"`
	LastReward     RollReward     `json:"last_reward"`
	PetGrowth      PetGrowth      `json:"pet_growth"`
	BattleResult   battle.Summary `json:"last_battle"`
	WalletSnapshot growth.Wallet  `json:"wallet_snapshot"`
}

type RollReward struct {
	Label       string `json:"label"`
	SpiritPower int64  `json:"spirit_power"`
	SoulPieces  int    `json:"soul_pieces"`
}

type WorldMap struct {
	Name   string    `json:"name"`
	Cities []MapCity `json:"cities"`
}

type MapDungeon struct {
	DungeonID         int64  `json:"dungeon_id"`
	DungeonName       string `json:"dungeon_name"`
	UnlockSpiritPower int64  `json:"unlock_spirit_power"`
	UnlockBoneLevel   int    `json:"unlock_bone_level"`
	UnlockSoulPieces  int    `json:"unlock_soul_pieces"`
}

type MapCity struct {
	CityID   int64        `json:"city_id"`
	Name     string       `json:"name"`
	Region   string       `json:"region"`
	LocX     float64      `json:"loc_x"`
	LocY     float64      `json:"loc_y"`
	Dungeons []MapDungeon `json:"dungeons"`
}

type CultivationStatus struct {
	PlayerID    int64     `json:"player_id"`
	SpiritPower int       `json:"spirit_power"`
	State       string    `json:"state"`
	StartAt     time.Time `json:"start_at,omitempty"`
	ClaimableAt time.Time `json:"claimable_at,omitempty"`
	PetGrowth   PetGrowth `json:"pet_growth"`
}

type PetGrowth struct {
	Exp            int64 `json:"exp"`
	TeamTotalPower int64 `json:"team_total_power"`
}
