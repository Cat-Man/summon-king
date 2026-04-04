package battle

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Resolve(_ context.Context, req Request) (Summary, error) {
	result := req.PresetResult
	if result == "" {
		if req.Attacker.TotalPower >= req.Defender.TotalPower {
			result = "success"
		} else {
			result = "fail"
		}
	}

	return Summary{
		BattleType:    req.Type,
		Result:        result,
		Rounds:        1,
		AttackerPower: req.Attacker.TotalPower,
		DefenderPower: req.Defender.TotalPower,
	}, nil
}
