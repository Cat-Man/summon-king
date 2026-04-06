package alliance

import "context"

var defaultWarTargets = []string{"赤焰谷", "寒月岭"}

type WarService struct {
	repo Repository
}

type WarIndex struct {
	HasAlliance      bool     `json:"has_alliance"`
	CurrentRole      string   `json:"current_role,omitempty"`
	Phase            string   `json:"phase"`
	TargetLabel      string   `json:"target_label"`
	CanRegister      bool     `json:"can_register"`
	AvailableTargets []string `json:"available_targets,omitempty"`
}

func NewWarService(repo Repository) *WarService {
	return &WarService{repo: repo}
}

func (s *WarService) GetIndex(ctx context.Context, playerID int64) (WarIndex, error) {
	state, role, found, err := s.repo.GetAllianceByPlayerID(ctx, playerID)
	if err != nil {
		return WarIndex{}, err
	}
	if !found {
		return WarIndex{
			HasAlliance: false,
			Phase:       "preparing",
			TargetLabel: "未加入联盟",
		}, nil
	}

	return buildWarIndex(role, state), nil
}

func (s *WarService) RegisterTarget(ctx context.Context, playerID int64, target string) (WarIndex, error) {
	state, role, err := s.repo.RegisterWarTarget(ctx, playerID, target)
	if err != nil {
		return WarIndex{}, err
	}

	return buildWarIndex(role, state), nil
}

func buildWarIndex(role string, state allianceState) WarIndex {
	canRegister := role == roleLeader
	targetLabel := state.WarTarget
	if targetLabel == "" {
		targetLabel = "待选择本轮目标"
	}

	index := WarIndex{
		HasAlliance: true,
		CurrentRole: role,
		Phase:       "preparing",
		TargetLabel: targetLabel,
		CanRegister: canRegister,
	}
	if canRegister {
		index.AvailableTargets = append([]string(nil), defaultWarTargets...)
	}
	return index
}

func isAvailableWarTarget(target string) bool {
	for _, candidate := range defaultWarTargets {
		if candidate == target {
			return true
		}
	}
	return false
}
