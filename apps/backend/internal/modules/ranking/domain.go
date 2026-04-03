package ranking

type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	PlayerID int64  `json:"player_id"`
	Name     string `json:"name"`
	Score    int64  `json:"score"`
	Updated  int64  `json:"updated"`
	IsSelf   bool   `json:"is_self"`
}
