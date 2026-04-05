package pet

import (
	"context"
	"errors"
	"sort"
	"sync"
)

var ErrPetNotFound = errors.New("pet not found")
var ErrInvalidTeam = errors.New("invalid team")

const maxActiveTeamSize = 2

type Repository interface {
	GetBattleTeam(ctx context.Context, playerID int64) (TeamSnapshot, error)
	ListPets(ctx context.Context, playerID int64) ([]BattlePet, error)
	SaveTeam(ctx context.Context, playerID int64, petIDs []int64) error
	SetMainPet(ctx context.Context, playerID, petID int64) error
}

type MemoryRepository struct {
	mu           sync.Mutex
	petsByPlayer map[int64][]BattlePet
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		petsByPlayer: make(map[int64][]BattlePet),
	}
}

func (r *MemoryRepository) GetBattleTeam(_ context.Context, playerID int64) (TeamSnapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return TeamSnapshot{
		PlayerID: playerID,
		Pets:     activePets(r.ensurePetsLocked(playerID)),
	}, nil
}

func (r *MemoryRepository) ListPets(_ context.Context, playerID int64) ([]BattlePet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return rosterPets(r.ensurePetsLocked(playerID)), nil
}

func (r *MemoryRepository) SaveTeam(_ context.Context, playerID int64, petIDs []int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pets := r.ensurePetsLocked(playerID)
	return r.saveTeamLocked(playerID, pets, petIDs)
}

func (r *MemoryRepository) SetMainPet(_ context.Context, playerID, petID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pets := r.ensurePetsLocked(playerID)
	targetIndex := -1
	for idx := range pets {
		if pets[idx].PetID == petID {
			targetIndex = idx
			break
		}
	}
	if targetIndex == -1 {
		return ErrPetNotFound
	}

	teamIDs := []int64{petID}
	if pets[targetIndex].IsActive {
		teamIDs = append(teamIDs, removePetID(activePetIDs(pets), petID)...)
	}

	return r.saveTeamLocked(playerID, pets, teamIDs)
}

func (r *MemoryRepository) ensurePetsLocked(playerID int64) []BattlePet {
	if pets, ok := r.petsByPlayer[playerID]; ok {
		return pets
	}

	pets := []BattlePet{
		{
			PetID:    playerID*10 + 1,
			Slot:     1,
			Name:     "初始灵狐",
			Level:    1,
			Power:    120,
			IsActive: true,
		},
		{
			PetID:    playerID*10 + 2,
			Slot:     0,
			Name:     "玄甲龟",
			Level:    1,
			Power:    156,
			IsActive: false,
		},
	}
	r.petsByPlayer[playerID] = pets
	return pets
}

func clonePets(pets []BattlePet) []BattlePet {
	cloned := make([]BattlePet, len(pets))
	copy(cloned, pets)
	return cloned
}

func (r *MemoryRepository) saveTeamLocked(playerID int64, pets []BattlePet, petIDs []int64) error {
	if len(petIDs) == 0 || len(petIDs) > maxActiveTeamSize {
		return ErrInvalidTeam
	}

	slots := make(map[int64]int, len(petIDs))
	existing := make(map[int64]struct{}, len(pets))
	for _, battlePet := range pets {
		existing[battlePet.PetID] = struct{}{}
	}
	for idx, petID := range petIDs {
		if _, exists := slots[petID]; exists {
			return ErrInvalidTeam
		}
		if _, exists := existing[petID]; !exists {
			return ErrPetNotFound
		}
		slots[petID] = idx + 1
	}

	for idx := range pets {
		slot, exists := slots[pets[idx].PetID]
		if exists {
			pets[idx].IsActive = true
			pets[idx].Slot = slot
			continue
		}
		pets[idx].IsActive = false
		pets[idx].Slot = 0
	}

	r.petsByPlayer[playerID] = pets
	return nil
}

func activePets(pets []BattlePet) []BattlePet {
	active := make([]BattlePet, 0, len(pets))
	for _, battlePet := range pets {
		if battlePet.IsActive {
			active = append(active, battlePet)
		}
	}
	sort.Slice(active, func(i, j int) bool {
		if active[i].Slot == active[j].Slot {
			return active[i].PetID < active[j].PetID
		}
		return active[i].Slot < active[j].Slot
	})
	return clonePets(active)
}

func rosterPets(pets []BattlePet) []BattlePet {
	roster := clonePets(pets)
	sort.Slice(roster, func(i, j int) bool {
		if roster[i].IsActive != roster[j].IsActive {
			return roster[i].IsActive
		}
		if roster[i].Slot != roster[j].Slot {
			if roster[i].Slot == 0 {
				return false
			}
			if roster[j].Slot == 0 {
				return true
			}
			return roster[i].Slot < roster[j].Slot
		}
		return roster[i].PetID < roster[j].PetID
	})
	return roster
}

func activePetIDs(pets []BattlePet) []int64 {
	active := activePets(pets)
	ids := make([]int64, 0, len(active))
	for _, battlePet := range active {
		ids = append(ids, battlePet.PetID)
	}
	return ids
}

func removePetID(petIDs []int64, target int64) []int64 {
	filtered := make([]int64, 0, len(petIDs))
	for _, petID := range petIDs {
		if petID != target {
			filtered = append(filtered, petID)
		}
	}
	return filtered
}
