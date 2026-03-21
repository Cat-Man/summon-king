package asset

import "time"

var defaultInventory = []InventoryItem{
	{
		ItemID:    "potion_small",
		ItemName:  "Small Potion",
		Quantity:  5,
		SellPrice: 10,
	},
	{
		ItemID:    "summon_scroll",
		ItemName:  "Summon Scroll",
		Quantity:  1,
		SellPrice: 50,
	},
}

type Wallet struct {
	PlayerID  int64     `json:"player_id"`
	Coins     int64     `json:"coins"`
	Diamonds  int64     `json:"diamonds"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryItem struct {
	PlayerID  int64  `json:"player_id"`
	ItemID    string `json:"item_id"`
	ItemName  string `json:"item_name"`
	Quantity  int64  `json:"quantity"`
	SellPrice int64  `json:"sell_price"`
}

type ResourceChangeLog struct {
	PlayerID      int64     `json:"player_id"`
	ChangeType    string    `json:"change_type"`
	BizID         string    `json:"biz_id"`
	CoinsDelta    int64     `json:"coins_delta"`
	DiamondsDelta int64     `json:"diamonds_delta"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created_at"`
}

type RewardClaimLog struct {
	PlayerID  int64     `json:"player_id"`
	BizID     string    `json:"biz_id"`
	Source    string    `json:"source"`
	Coins     int64     `json:"coins"`
	Diamonds  int64     `json:"diamonds"`
	ClaimedAt time.Time `json:"claimed_at"`
}

type RewardGrant struct {
	PlayerID int64  `json:"player_id"`
	BizID    string `json:"biz_id"`
	Source   string `json:"source"`
	Coins    int64  `json:"coins"`
	Diamonds int64  `json:"diamonds"`
}

type GrantRewardResult struct {
	Wallet     Wallet `json:"wallet"`
	Idempotent bool   `json:"idempotent"`
}

type InventoryOperateRequest struct {
	PlayerID int64  `json:"player_id"`
	ItemID   string `json:"item_id"`
	Count    int64  `json:"count"`
}

type InventoryOperateResult struct {
	Wallet Wallet        `json:"wallet"`
	Item   InventoryItem `json:"item"`
}
