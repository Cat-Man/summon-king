package pet

import (
	"context"
	"errors"
	"sync"
)

var ErrPetNotFound = errors.New("pet not found")

type Repository interface {
	GetBattleTeam(ctx context.Context, playerID int64) (TeamSnapshot, error)
	ListPets(ctx context.Context, playerID int64) ([]BattlePet, error)
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

	pets := clonePets(r.ensurePetsLocked(playerID))
	active := make([]BattlePet, 0, len(pets))
	for _, battlePet := range pets {
		if battlePet.IsActive {
			active = append(active, battlePet)
		}
	}

	return TeamSnapshot{
		PlayerID: playerID,
		Pets:     active,
	}, nil
}

func (r *MemoryRepository) ListPets(_ context.Context, playerID int64) ([]BattlePet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return clonePets(r.ensurePetsLocked(playerID)), nil
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

	for idx := range pets {
		pets[idx].IsActive = idx == targetIndex
		if pets[idx].IsActive {
			pets[idx].Slot = 1
			continue
		}
		pets[idx].Slot = 0
	}

	r.petsByPlayer[playerID] = pets
	return nil
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
