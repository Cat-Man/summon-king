package battle

import "github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"

type Request struct {
	Type         string           `json:"type"`
	PresetResult string           `json:"preset_result,omitempty"`
	Attacker     pet.TeamSnapshot `json:"attacker"`
	Defender     pet.TeamSnapshot `json:"defender"`
}

type Summary struct {
	BattleNo      string `json:"battle_no"`
	BattleType    string `json:"battle_type"`
	Result        string `json:"result"`
	WinnerSide    string `json:"winner_side"`
	Rounds        int    `json:"rounds"`
	AttackerPower int64  `json:"attacker_power"`
	DefenderPower int64  `json:"defender_power"`
}
