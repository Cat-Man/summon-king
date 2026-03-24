package dungeon

import "time"

type City struct {
	CityID   string `json:"city_id"`
	CityName string `json:"city_name"`
	Unlocked bool   `json:"unlocked"`
}

type WorldMap struct {
	PlayerID int64  `json:"player_id"`
	Region   string `json:"region"`
	Cities   []City `json:"cities"`
}

type DungeonRun struct {
	RunID        string    `json:"run_id"`
	PlayerID     int64     `json:"player_id"`
	DungeonID    int64     `json:"dungeon_id"`
	CurrentFloor int       `json:"current_floor"`
	RemainDice   int       `json:"remain_dice"`
	BossRewarded bool      `json:"boss_rewarded"`
	CreatedAt    time.Time `json:"created_at"`
}

type RollResult struct {
	Run           DungeonRun `json:"run"`
	DicePoint     int        `json:"dice_point"`
	FloorResolved int        `json:"floor_resolved"`
}

type CultivationRecord struct {
	RecordID     string    `json:"record_id"`
	PlayerID     int64     `json:"player_id"`
	MapID        string    `json:"map_id"`
	Hours        int       `json:"hours"`
	Status       string    `json:"status"`
	RewardCoins  int64     `json:"reward_coins"`
	RewardPetExp int64     `json:"reward_pet_exp"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
}
