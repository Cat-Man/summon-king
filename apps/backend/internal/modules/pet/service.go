package pet

import (
	"context"
	"errors"
)

var ErrDuplicatePetInTeam = errors.New("duplicate pet in team")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCatalog(ctx context.Context, playerID int64) ([]PetCatalogItem, error) {
	return s.repo.GetCatalog(ctx, playerID)
}

func (s *Service) GetPlayerPets(ctx context.Context, playerID int64) ([]PlayerPet, error) {
	return s.repo.GetPlayerPets(ctx, playerID)
}

func (s *Service) GetDetail(ctx context.Context, playerID, petID int64) (PetCatalogItem, error) {
	return s.repo.GetPetDetail(ctx, playerID, petID)
}

func (s *Service) SaveTeam(ctx context.Context, playerID int64, petIDs []int64) error {
	seen := make(map[int64]struct{}, len(petIDs))
	for _, petID := range petIDs {
		if petID == 0 {
			continue
		}
		if _, ok := seen[petID]; ok {
			return ErrDuplicatePetInTeam
		}
		seen[petID] = struct{}{}
	}

	return s.repo.SaveTeam(ctx, PetTeam{
		PlayerID: playerID,
		PetIDs:   petIDs,
	})
}

func (s *Service) GetTeam(ctx context.Context, playerID int64) (PetTeam, error) {
	return s.repo.GetTeam(ctx, playerID)
}
