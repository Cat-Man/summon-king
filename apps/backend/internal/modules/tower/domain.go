package tower

type TowerInfo struct {
	Level    string `json:"level"`
	MaxFloor int    `json:"max_floor"`
}

type TowerResult struct {
	PlayerID int64  `json:"player_id"`
	Floor    int    `json:"floor"`
	Reward   string `json:"reward"`
}
