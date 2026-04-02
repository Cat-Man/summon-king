package arena

type DailyRecord struct {
	PlayerID      int64 `json:"player_id"`
	CurrentStreak int   `json:"current_streak"`
	LastWin       bool  `json:"last_win"`
}
