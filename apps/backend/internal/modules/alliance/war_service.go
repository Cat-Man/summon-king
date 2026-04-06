package alliance

import "context"

type WarService struct{}

type WarIndex struct {
	Phase       string `json:"phase"`
	TargetLabel string `json:"target_label"`
	CanRegister bool   `json:"can_register"`
}

func NewWarService() *WarService {
	return &WarService{}
}

func (s *WarService) GetIndex(_ context.Context, _ int64) WarIndex {
	return WarIndex{
		Phase:       "preparing",
		TargetLabel: "待开放目标",
		CanRegister: false,
	}
}
