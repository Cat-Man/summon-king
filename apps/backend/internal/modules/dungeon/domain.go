package dungeon

import "time"

type DungeonRun struct {
	PlayerID     int64     `json:"player_id"`
	DungeonID    int64     `json:"dungeon_id"`
	RemainDice   int       `json:"remain_dice"`
	CurrentFloor int       `json:"current_floor"`
	Status       string    `json:"status"`
	StartedAt    time.Time `json:"started_at"`
}

type WorldMap struct {
	Name   string    `json:"name"`
	Cities []MapCity `json:"cities"`
}

type MapCity struct {
	CityID int64   `json:"city_id"`
	Name   string  `json:"name"`
	Region string  `json:"region"`
	LocX   float64 `json:"loc_x"`
	LocY   float64 `json:"loc_y"`
}

type CultivationStatus struct {
	PlayerID    int64     `json:"player_id"`
	SpiritPower int       `json:"spirit_power"`
	State       string    `json:"state"`
	StartAt     time.Time `json:"start_at,omitempty"`
	ClaimableAt time.Time `json:"claimable_at,omitempty"`
}
