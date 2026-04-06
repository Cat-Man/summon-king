package pet

type BattlePet struct {
	PetID          int64          `json:"pet_id"`
	Slot           int            `json:"slot"`
	Name           string         `json:"name"`
	Level          int            `json:"level"`
	Exp            int64          `json:"exp"`
	NextLevelExp   int64          `json:"next_level_exp"`
	Power          int64          `json:"power"`
	PowerBreakdown PowerBreakdown `json:"power_breakdown"`
	IsActive       bool           `json:"is_active"`
	BasePower      int64          `json:"-"`
}

type PowerBreakdown struct {
	Base   int64 `json:"base"`
	Level  int64 `json:"level"`
	Bone   int64 `json:"bone"`
	Spirit int64 `json:"spirit"`
	Soul   int64 `json:"soul"`
	Total  int64 `json:"total"`
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
