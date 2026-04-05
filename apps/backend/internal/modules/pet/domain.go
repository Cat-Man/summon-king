package pet

type BattlePet struct {
	PetID    int64  `json:"pet_id"`
	Slot     int    `json:"slot"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Power    int64  `json:"power"`
	IsActive bool   `json:"is_active"`
}

type TeamSnapshot struct {
	PlayerID   int64       `json:"player_id"`
	TotalPower int64       `json:"total_power"`
	Pets       []BattlePet `json:"pets"`
}

type CollectionView struct {
	PlayerID   int64       `json:"player_id"`
	TotalPower int64       `json:"total_power"`
	TeamSize   int         `json:"team_size"`
	ActiveTeam []BattlePet `json:"active_team"`
	Roster     []BattlePet `json:"roster"`
}
