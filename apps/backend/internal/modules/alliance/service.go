package alliance

import "context"

type Service struct {
	repo Repository
}

type Index struct {
	PlayerID            int64                 `json:"player_id"`
	HasAlliance         bool                  `json:"has_alliance"`
	CurrentRole         string                `json:"current_role"`
	Alliance            *AllianceSummary      `json:"alliance,omitempty"`
	FireTraining        *FireTrainingStatus   `json:"fire_training,omitempty"`
	War                 *WarIndex             `json:"war,omitempty"`
	PendingApplications []AllianceApplication `json:"pending_applications,omitempty"`
}

type AllianceSummary struct {
	AllianceID  int64              `json:"alliance_id"`
	Name        string             `json:"name"`
	Level       int                `json:"level"`
	Notice      string             `json:"notice"`
	MemberCount int                `json:"member_count"`
	MemberLimit int                `json:"member_limit"`
	Members     []AllianceMember   `json:"members"`
	Buildings   []AllianceBuilding `json:"buildings"`
}

type AllianceMember struct {
	PlayerID int64  `json:"player_id"`
	Role     string `json:"role"`
}

type AllianceBuilding struct {
	BuildingType string `json:"building_type"`
	Level        int    `json:"level"`
}

type AllianceApplication struct {
	PlayerID int64 `json:"player_id"`
}

type HallEntry struct {
	AllianceID  int64  `json:"alliance_id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	MemberCount int    `json:"member_count"`
	MemberLimit int    `json:"member_limit"`
	Notice      string `json:"notice"`
	HasApplied  bool   `json:"has_applied"`
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetIndex(ctx context.Context, playerID int64) (Index, error) {
	state, role, found, err := s.repo.GetAllianceByPlayerID(ctx, playerID)
	if err != nil {
		return Index{}, err
	}
	if !found {
		return Index{
			PlayerID:    playerID,
			HasAlliance: false,
		}, nil
	}

	index := buildIndex(playerID, role, state)
	if role == roleLeader {
		applications, err := s.repo.ListApplications(ctx, playerID)
		if err != nil && err != ErrAlliancePermissionDenied {
			return Index{}, err
		}
		index.PendingApplications = buildApplications(applications)
	}

	return index, nil
}

func (s *Service) ListHall(ctx context.Context, playerID int64) ([]HallEntry, error) {
	entries, err := s.repo.ListHall(ctx, playerID)
	if err != nil {
		return nil, err
	}

	result := make([]HallEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, HallEntry{
			AllianceID:  entry.AllianceID,
			Name:        entry.Name,
			Level:       entry.Level,
			MemberCount: len(entry.Members),
			MemberLimit: entry.MemberLimit,
			Notice:      entry.Notice,
			HasApplied:  entry.HasApplied,
		})
	}
	return result, nil
}

func (s *Service) CreateAlliance(ctx context.Context, playerID int64, name string) (Index, error) {
	state, role, err := s.repo.CreateAlliance(ctx, playerID, name)
	if err != nil {
		return Index{}, err
	}
	return buildIndex(playerID, role, state), nil
}

func (s *Service) Apply(ctx context.Context, playerID, allianceID int64) error {
	return s.repo.Apply(ctx, playerID, allianceID)
}

func (s *Service) ListApplications(ctx context.Context, playerID int64) ([]AllianceApplication, error) {
	applications, err := s.repo.ListApplications(ctx, playerID)
	if err != nil {
		return nil, err
	}
	return buildApplications(applications), nil
}

func (s *Service) ApproveApplication(ctx context.Context, approverID, applicantID int64) (Index, error) {
	state, role, err := s.repo.ApproveApplication(ctx, approverID, applicantID)
	if err != nil {
		return Index{}, err
	}
	return buildIndex(applicantID, role, state), nil
}

func buildIndex(playerID int64, role string, state allianceState) Index {
	fireTraining := NewFireTrainingService().GetStatus(state)
	war := buildWarIndex(role, state)

	return Index{
		PlayerID:    playerID,
		HasAlliance: true,
		CurrentRole: role,
		Alliance: &AllianceSummary{
			AllianceID:  state.AllianceID,
			Name:        state.Name,
			Level:       state.Level,
			Notice:      state.Notice,
			MemberCount: len(state.Members),
			MemberLimit: state.MemberLimit,
			Members:     buildMembers(state.Members),
			Buildings:   buildBuildings(state.Buildings),
		},
		FireTraining: &fireTraining,
		War:          &war,
	}
}

func buildMembers(members []memberState) []AllianceMember {
	result := make([]AllianceMember, 0, len(members))
	for _, member := range members {
		result = append(result, AllianceMember{
			PlayerID: member.PlayerID,
			Role:     member.Role,
		})
	}
	return result
}

func buildBuildings(buildings []buildingState) []AllianceBuilding {
	result := make([]AllianceBuilding, 0, len(buildings))
	for _, building := range buildings {
		result = append(result, AllianceBuilding{
			BuildingType: building.BuildingType,
			Level:        building.Level,
		})
	}
	return result
}

func buildApplications(applications []applicationState) []AllianceApplication {
	result := make([]AllianceApplication, 0, len(applications))
	for _, application := range applications {
		result = append(result, AllianceApplication{
			PlayerID: application.PlayerID,
		})
	}
	return result
}
