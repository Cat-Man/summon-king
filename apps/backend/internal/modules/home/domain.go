package home

import (
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
)

type Overview struct {
	PlayerID   int64           `json:"player_id"`
	Nickname   string          `json:"nickname"`
	Wallet     growth.Wallet   `json:"wallet"`
	Modules    ModulesOverview `json:"modules"`
	NextAction NextAction      `json:"next_action"`
}

type ModulesOverview struct {
	MapLabel     string             `json:"map_label"`
	MapCityCount int                `json:"map_city_count"`
	Dungeon      DungeonSummary     `json:"dungeon"`
	Cultivation  CultivationSummary `json:"cultivation"`
	Tower        TowerOverview      `json:"tower"`
}

type DungeonSummary struct {
	Status       string `json:"status"`
	CurrentFloor int    `json:"current_floor"`
	RemainDice   int    `json:"remain_dice"`
	DungeonID    int64  `json:"dungeon_id"`
}

type CultivationSummary struct {
	State       string    `json:"state"`
	SpiritPower int       `json:"spirit_power"`
	Claimable   bool      `json:"claimable"`
	ClaimableAt time.Time `json:"claimable_at,omitempty"`
}

type TowerOverview struct {
	Pagoda tower.TowerStatus `json:"pagoda"`
	Spirit tower.TowerStatus `json:"spirit"`
}

type NextAction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Route       string `json:"route"`
	CTA         string `json:"cta"`
}
