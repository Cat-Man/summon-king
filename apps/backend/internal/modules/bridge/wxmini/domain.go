package wxmini

import "errors"

const Channel = "wxmini"

var (
	ErrCodeRequired         = errors.New("code is required")
	ErrUnifiedTokenRequired = errors.New("unified token is required")
	ErrSessionNotFound      = errors.New("session not found")
)

type LoginExchangeResponse struct {
	UnifiedToken string `json:"unified_token"`
	PlayerID     int64  `json:"player_id"`
	Nickname     string `json:"nickname"`
	Channel      string `json:"channel"`
}

type SessionBootstrapResponse struct {
	Token    string `json:"token"`
	PlayerID int64  `json:"player_id"`
	Nickname string `json:"nickname"`
	Channel  string `json:"channel"`
	GameURL  string `json:"game_url"`
}
