package growth

import "time"

type Wallet struct {
	PlayerID            int64 `json:"player_id"`
	SpiritPower         int64 `json:"spirit_power"`
	FreeSpiritWashCount int   `json:"free_spirit_wash_count"`
	BonePower           int64 `json:"bone_power"`
	SoulPower           int64 `json:"soul_power"`
}

type WashOption struct{}

type Bone struct {
	BoneID string `json:"bone_id"`
	Level  int    `json:"level"`
}

type Spirit struct {
	SpiritID  int64  `json:"spirit_id"`
	AttrLine  string `json:"attr_line"`
	WashCount int    `json:"wash_count"`
	SellPrice int64  `json:"sell_price"`
}

type Soul struct {
	SoulID string `json:"soul_id"`
	Level  int    `json:"level"`
}

type CauldronRecord struct {
	PlayerID       int64 `json:"player_id"`
	RefineCount    int   `json:"refine_count"`
	CurrentQuality int   `json:"current_quality"`
}

type AscensionRecord struct {
	PlayerID int64  `json:"player_id"`
	PetID    int64  `json:"pet_id"`
	State    string `json:"state"`
}

type Crop struct {
	PlayerID   int64     `json:"player_id"`
	SeedID     string    `json:"seed_id"`
	Status     string    `json:"status"`
	RewardCoin int64     `json:"reward_coin"`
	ReadyAt    time.Time `json:"ready_at"`
}
