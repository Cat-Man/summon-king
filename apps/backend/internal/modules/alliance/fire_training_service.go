package alliance

import "strconv"

type FireTrainingService struct{}

type FireTrainingStatus struct {
	FurnaceLevel  int    `json:"furnace_level"`
	CanClaimStone bool   `json:"can_claim_stone"`
	CurrentRoom   string `json:"current_room"`
	Claimable     bool   `json:"claimable"`
	RewardPreview string `json:"reward_preview"`
}

func NewFireTrainingService() *FireTrainingService {
	return &FireTrainingService{}
}

func (s *FireTrainingService) GetStatus(state allianceState) FireTrainingStatus {
	furnaceLevel := 1
	for _, building := range state.Buildings {
		if building.BuildingType == "fire_forge" {
			furnaceLevel = max(building.Level, 1)
			break
		}
	}

	return FireTrainingStatus{
		FurnaceLevel:  furnaceLevel,
		CanClaimStone: true,
		CurrentRoom:   "未入房",
		Claimable:     false,
		RewardPreview: fireTrainingRewardPreview(furnaceLevel),
	}
}

func fireTrainingRewardPreview(furnaceLevel int) string {
	return "焚火晶 x" + strconv.Itoa(4+furnaceLevel*2) + " / 2小时"
}
