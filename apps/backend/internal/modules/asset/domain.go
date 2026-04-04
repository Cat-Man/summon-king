package asset

import "github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"

type Delta struct {
	SpiritPower int64 `json:"spirit_power"`
	BoneLevel   int   `json:"bone_level"`
	SoulPieces  int   `json:"soul_pieces"`
}

type ApplyResult struct {
	Wallet growth.Wallet `json:"wallet"`
}
