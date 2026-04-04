package tower

type TowerStatus struct {
	Tower               string `json:"tower"`
	Label               string `json:"label"`
	CurrentFloor        int    `json:"current_floor"`
	MaxFloor            int    `json:"max_floor"`
	RemainingChallenges int    `json:"remaining_challenges"`
	RewardPreview       string `json:"reward_preview"`
}

type TowerResult struct {
	PlayerID            int64  `json:"player_id"`
	Tower               string `json:"tower"`
	Floor               int    `json:"floor"`
	Reward              string `json:"reward"`
	RemainingChallenges int    `json:"remaining_challenges"`
}
