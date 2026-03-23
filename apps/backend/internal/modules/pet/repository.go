package pet

import (
	"context"
	"errors"
	"sort"
	"sync"
)

var ErrPetNotFound = errors.New("pet not found")

type Repository interface {
	GetCatalog(ctx context.Context, playerID int64) ([]PetCatalogItem, error)
	GetPlayerPets(ctx context.Context, playerID int64) ([]PlayerPet, error)
	GetPetDetail(ctx context.Context, playerID, petID int64) (PetCatalogItem, error)
	SaveTeam(ctx context.Context, team PetTeam) error
	GetTeam(ctx context.Context, playerID int64) (PetTeam, error)
}

type MemoryRepository struct {
	mu      sync.RWMutex
	catalog []PetCatalogItem
	pets    map[int64][]PlayerPet
	teams   map[int64]PetTeam
}

func NewMemoryRepository() *MemoryRepository {
	catalog := []PetCatalogItem{
		{PetID: 1, Name: "赤焰狼", Rarity: 2, Element: "fire", Locked: false, Owned: true, Power: 120, Portrait: "pet_1.png", StoryText: "火山边缘的狩猎者"},
		{PetID: 2, Name: "寒霜鹤", Rarity: 3, Element: "water", Locked: false, Owned: true, Power: 150, Portrait: "pet_2.png", StoryText: "来自北境冰湖的守望者"},
		{PetID: 3, Name: "雷角牛", Rarity: 4, Element: "thunder", Locked: false, Owned: true, Power: 180, Portrait: "pet_3.png", StoryText: "暴风中的冲锋者"},
	}

	ownedPets := []PlayerPet{
		{PlayerID: 0, PetID: 1, Level: 12, Star: 1, Power: 120},
		{PlayerID: 0, PetID: 2, Level: 16, Star: 2, Power: 150},
		{PlayerID: 0, PetID: 3, Level: 20, Star: 3, Power: 180},
	}

	return &MemoryRepository{
		catalog: catalog,
		pets:    map[int64][]PlayerPet{0: ownedPets},
		teams:   make(map[int64]PetTeam),
	}
}

func (r *MemoryRepository) GetCatalog(_ context.Context, _ int64) ([]PetCatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]PetCatalogItem, len(r.catalog))
	copy(result, r.catalog)
	return result, nil
}

func (r *MemoryRepository) GetPlayerPets(_ context.Context, playerID int64) ([]PlayerPet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pets := r.playerPetsLocked(playerID)
	result := make([]PlayerPet, len(pets))
	copy(result, pets)
	return result, nil
}

func (r *MemoryRepository) GetPetDetail(_ context.Context, _ int64, petID int64) (PetCatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.catalog {
		if item.PetID == petID {
			return item, nil
		}
	}

	return PetCatalogItem{}, ErrPetNotFound
}

func (r *MemoryRepository) SaveTeam(_ context.Context, team PetTeam) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	petIDs := make([]int64, len(team.PetIDs))
	copy(petIDs, team.PetIDs)
	r.teams[team.PlayerID] = PetTeam{
		PlayerID: team.PlayerID,
		PetIDs:   petIDs,
	}
	return nil
}

func (r *MemoryRepository) GetTeam(_ context.Context, playerID int64) (PetTeam, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if team, ok := r.teams[playerID]; ok {
		result := make([]int64, len(team.PetIDs))
		copy(result, team.PetIDs)
		return PetTeam{PlayerID: team.PlayerID, PetIDs: result}, nil
	}

	pets := r.playerPetsLocked(playerID)
	return PetTeam{
		PlayerID: playerID,
		PetIDs:   defaultTeamPetIDs(pets),
	}, nil
}

func defaultTeamPetIDs(pets []PlayerPet) []int64 {
	sorted := make([]PlayerPet, len(pets))
	copy(sorted, pets)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Power == sorted[j].Power {
			return sorted[i].PetID < sorted[j].PetID
		}
		return sorted[i].Power > sorted[j].Power
	})

	teamSize := len(sorted)
	if teamSize > 5 {
		teamSize = 5
	}

	result := make([]int64, 0, teamSize)
	for i := 0; i < teamSize; i++ {
		if sorted[i].PetID == 0 {
			continue
		}
		result = append(result, sorted[i].PetID)
	}
	return result
}

func (r *MemoryRepository) playerPetsLocked(playerID int64) []PlayerPet {
	if pets, ok := r.pets[playerID]; ok {
		return pets
	}

	basePets := r.pets[0]
	cloned := make([]PlayerPet, len(basePets))
	for i, p := range basePets {
		cloned[i] = p
		cloned[i].PlayerID = playerID
	}
	r.pets[playerID] = cloned
	return cloned
}
