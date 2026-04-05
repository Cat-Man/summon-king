package pet

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBattleTeam(ctx context.Context, playerID int64) (TeamSnapshot, error) {
	team, err := s.repo.GetBattleTeam(ctx, playerID)
	if err != nil {
		return TeamSnapshot{}, err
	}

	team.TotalPower = totalPower(team.Pets)

	return team, nil
}

func (s *Service) GetCollection(ctx context.Context, playerID int64) (CollectionView, error) {
	team, err := s.GetBattleTeam(ctx, playerID)
	if err != nil {
		return CollectionView{}, err
	}

	roster, err := s.repo.ListPets(ctx, playerID)
	if err != nil {
		return CollectionView{}, err
	}

	return CollectionView{
		PlayerID:   playerID,
		TotalPower: team.TotalPower,
		TeamSize:   len(team.Pets),
		ActiveTeam: team.Pets,
		Roster:     roster,
	}, nil
}

func totalPower(pets []BattlePet) int64 {
	var total int64
	for _, battlePet := range pets {
		total += battlePet.Power
	}
	return total
}
