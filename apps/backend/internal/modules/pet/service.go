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

	var totalPower int64
	for _, battlePet := range team.Pets {
		totalPower += battlePet.Power
	}
	team.TotalPower = totalPower

	return team, nil
}
