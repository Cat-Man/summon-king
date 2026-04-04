package battle

import "github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"

type Request struct {
	Type         string           `json:"type"`
	PresetResult string           `json:"preset_result,omitempty"`
	Attacker     pet.TeamSnapshot `json:"attacker"`
	Defender     pet.TeamSnapshot `json:"defender"`
}

type Summary struct {
	BattleType    string `json:"battle_type"`
	Result        string `json:"result"`
	Rounds        int    `json:"rounds"`
	AttackerPower int64  `json:"attacker_power"`
	DefenderPower int64  `json:"defender_power"`
}
