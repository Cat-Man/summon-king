package alliance

import "context"

type FireTrainingService struct{}

type FireTrainingStatus struct {
	CanClaimStone bool   `json:"can_claim_stone"`
	CurrentRoom   string `json:"current_room"`
	Claimable     bool   `json:"claimable"`
}

func NewFireTrainingService() *FireTrainingService {
	return &FireTrainingService{}
}

func (s *FireTrainingService) GetStatus(_ context.Context, _ int64) FireTrainingStatus {
	return FireTrainingStatus{
		CanClaimStone: true,
		CurrentRoom:   "未入房",
		Claimable:     false,
	}
}
