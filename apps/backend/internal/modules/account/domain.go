package account

import "time"

type Player struct {
	PlayerID   int64            `json:"player_id"`
	Channel    string           `json:"channel"`
	Token      string           `json:"token"`
	Profile    PlayerProfile    `json:"profile"`
	Wallet     PlayerWallet     `json:"wallet"`
	DailyState PlayerDailyState `json:"daily_state"`
}

type PlayerProfile struct {
	PlayerID int64  `json:"player_id"`
	Nickname string `json:"nickname"`
	Level    int    `json:"level"`
}

type PlayerWallet struct {
	PlayerID   int64 `json:"player_id"`
	Coin       int64 `json:"coin"`
	Diamond    int64 `json:"diamond"`
	Vitality   int   `json:"vitality"`
	Reputation int64 `json:"reputation"`
}

type PlayerDailyState struct {
	PlayerID    int64     `json:"player_id"`
	LastLoginAt time.Time `json:"last_login_at"`
}
