package alliance

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
)

var (
	ErrAllianceNotFound           = errors.New("alliance not found")
	ErrAllianceNameRequired       = errors.New("alliance name is required")
	ErrAllianceNameTaken          = errors.New("alliance name already exists")
	ErrAlreadyInAlliance          = errors.New("player already in alliance")
	ErrAlreadyApplied             = errors.New("alliance application already exists")
	ErrAlliancePermissionDenied   = errors.New("alliance permission denied")
	ErrAllianceApplicationMissing = errors.New("alliance application not found")
	ErrAllianceFull               = errors.New("alliance is full")
)

const (
	defaultAllianceMemberLimit = 20
	roleLeader                 = "leader"
	roleMember                 = "member"
)

type Repository interface {
	CreateAlliance(ctx context.Context, playerID int64, name string) (allianceState, string, error)
	GetAllianceByPlayerID(ctx context.Context, playerID int64) (allianceState, string, bool, error)
	ListHall(ctx context.Context, playerID int64) ([]hallState, error)
	Apply(ctx context.Context, playerID, allianceID int64) error
	ListApplications(ctx context.Context, playerID int64) ([]applicationState, error)
	ApproveApplication(ctx context.Context, approverID, applicantID int64) (allianceState, string, error)
}

type allianceState struct {
	AllianceID  int64
	Name        string
	Level       int
	Notice      string
	MemberLimit int
	Members     []memberState
	Buildings   []buildingState
}

type memberState struct {
	PlayerID int64
	Role     string
}

type buildingState struct {
	BuildingType string
	Level        int
}

type applicationState struct {
	PlayerID int64
}

type hallState struct {
	allianceState
	HasApplied bool
}

type MemoryRepository struct {
	mu                sync.Mutex
	nextAllianceID    int64
	alliances         map[int64]allianceState
	playerAllianceIDs map[int64]int64
	applications      map[int64]map[int64]applicationState
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextAllianceID:    1,
		alliances:         make(map[int64]allianceState),
		playerAllianceIDs: make(map[int64]int64),
		applications:      make(map[int64]map[int64]applicationState),
	}
}

func (r *MemoryRepository) CreateAlliance(_ context.Context, playerID int64, name string) (allianceState, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return allianceState{}, "", ErrAllianceNameRequired
	}
	if _, exists := r.playerAllianceIDs[playerID]; exists {
		return allianceState{}, "", ErrAlreadyInAlliance
	}
	for _, existing := range r.alliances {
		if existing.Name == name {
			return allianceState{}, "", ErrAllianceNameTaken
		}
	}

	allianceID := r.nextAllianceID
	r.nextAllianceID++

	state := allianceState{
		AllianceID:  allianceID,
		Name:        name,
		Level:       1,
		Notice:      "欢迎加入，共筑仙盟。",
		MemberLimit: defaultAllianceMemberLimit,
		Members: []memberState{
			{PlayerID: playerID, Role: roleLeader},
		},
		Buildings: []buildingState{
			{BuildingType: "hall", Level: 1},
			{BuildingType: "fire_forge", Level: 1},
		},
	}
	r.alliances[allianceID] = state
	r.playerAllianceIDs[playerID] = allianceID

	return cloneAllianceState(state), roleLeader, nil
}

func (r *MemoryRepository) GetAllianceByPlayerID(_ context.Context, playerID int64) (allianceState, string, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	allianceID, exists := r.playerAllianceIDs[playerID]
	if !exists {
		return allianceState{}, "", false, nil
	}

	state, exists := r.alliances[allianceID]
	if !exists {
		return allianceState{}, "", false, ErrAllianceNotFound
	}

	for _, member := range state.Members {
		if member.PlayerID == playerID {
			return cloneAllianceState(state), member.Role, true, nil
		}
	}

	return allianceState{}, "", false, ErrAllianceNotFound
}

func (r *MemoryRepository) ListHall(_ context.Context, playerID int64) ([]hallState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entries := make([]hallState, 0, len(r.alliances))
	applied := make(map[int64]struct{})
	for allianceID, applications := range r.applications {
		if _, exists := applications[playerID]; exists {
			applied[allianceID] = struct{}{}
		}
	}

	ids := make([]int64, 0, len(r.alliances))
	for allianceID := range r.alliances {
		ids = append(ids, allianceID)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	for _, allianceID := range ids {
		state := r.alliances[allianceID]
		_, hasApplied := applied[allianceID]
		entries = append(entries, hallState{
			allianceState: cloneAllianceState(state),
			HasApplied:    hasApplied,
		})
	}

	return entries, nil
}

func (r *MemoryRepository) Apply(_ context.Context, playerID, allianceID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.playerAllianceIDs[playerID]; exists {
		return ErrAlreadyInAlliance
	}
	if _, exists := r.alliances[allianceID]; !exists {
		return ErrAllianceNotFound
	}
	for _, applications := range r.applications {
		if _, exists := applications[playerID]; exists {
			return ErrAlreadyApplied
		}
	}
	applications := r.ensureApplicationsLocked(allianceID)
	if _, exists := applications[playerID]; exists {
		return ErrAlreadyApplied
	}

	applications[playerID] = applicationState{PlayerID: playerID}
	return nil
}

func (r *MemoryRepository) ListApplications(_ context.Context, playerID int64) ([]applicationState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	allianceID, role, err := r.findAllianceAndRoleLocked(playerID)
	if err != nil {
		return nil, err
	}
	if role != roleLeader {
		return nil, ErrAlliancePermissionDenied
	}

	applications := r.ensureApplicationsLocked(allianceID)
	ids := make([]int64, 0, len(applications))
	for applicantID := range applications {
		ids = append(ids, applicantID)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	result := make([]applicationState, 0, len(ids))
	for _, applicantID := range ids {
		result = append(result, applications[applicantID])
	}
	return result, nil
}

func (r *MemoryRepository) ApproveApplication(_ context.Context, approverID, applicantID int64) (allianceState, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	allianceID, role, err := r.findAllianceAndRoleLocked(approverID)
	if err != nil {
		return allianceState{}, "", err
	}
	if role != roleLeader {
		return allianceState{}, "", ErrAlliancePermissionDenied
	}
	if _, exists := r.playerAllianceIDs[applicantID]; exists {
		return allianceState{}, "", ErrAlreadyInAlliance
	}

	applications := r.ensureApplicationsLocked(allianceID)
	if _, exists := applications[applicantID]; !exists {
		return allianceState{}, "", ErrAllianceApplicationMissing
	}

	state := r.alliances[allianceID]
	if len(state.Members) >= state.MemberLimit {
		return allianceState{}, "", ErrAllianceFull
	}

	state.Members = append(state.Members, memberState{
		PlayerID: applicantID,
		Role:     roleMember,
	})
	sort.Slice(state.Members, func(i, j int) bool {
		return state.Members[i].PlayerID < state.Members[j].PlayerID
	})
	r.alliances[allianceID] = state
	r.playerAllianceIDs[applicantID] = allianceID
	delete(applications, applicantID)

	return cloneAllianceState(state), roleMember, nil
}

func (r *MemoryRepository) ensureApplicationsLocked(allianceID int64) map[int64]applicationState {
	if applications, exists := r.applications[allianceID]; exists {
		return applications
	}

	applications := make(map[int64]applicationState)
	r.applications[allianceID] = applications
	return applications
}

func (r *MemoryRepository) findAllianceAndRoleLocked(playerID int64) (int64, string, error) {
	allianceID, exists := r.playerAllianceIDs[playerID]
	if !exists {
		return 0, "", ErrAllianceNotFound
	}

	state, exists := r.alliances[allianceID]
	if !exists {
		return 0, "", ErrAllianceNotFound
	}

	for _, member := range state.Members {
		if member.PlayerID == playerID {
			return allianceID, member.Role, nil
		}
	}

	return 0, "", ErrAllianceNotFound
}

func cloneAllianceState(state allianceState) allianceState {
	cloned := state
	cloned.Members = append([]memberState(nil), state.Members...)
	cloned.Buildings = append([]buildingState(nil), state.Buildings...)
	return cloned
}
