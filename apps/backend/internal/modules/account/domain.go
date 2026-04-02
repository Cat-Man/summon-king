package account

type GuestLoginResponse struct {
	PlayerID int64  `json:"player_id"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
}
