package asset

import "github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"

type Delta struct {
	SpiritPower int64 `json:"spirit_power"`
	BoneLevel   int   `json:"bone_level"`
	SoulPieces  int   `json:"soul_pieces"`
}

type ApplyMetadata struct {
	Source         string `json:"source"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type ApplyRequest struct {
	Delta    Delta         `json:"delta"`
	Metadata ApplyMetadata `json:"metadata"`
}

type ApplyResult struct {
	Wallet   growth.Wallet `json:"wallet"`
	Metadata ApplyMetadata `json:"metadata"`
}
