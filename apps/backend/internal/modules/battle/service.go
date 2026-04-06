package battle

import (
	"context"
	"fmt"
)

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
		BattleNo:      buildBattleNo(req),
		BattleType:    req.Type,
		Result:        result,
		WinnerSide:    winnerSide(result),
		Rounds:        1,
		AttackerPower: req.Attacker.TotalPower,
		DefenderPower: req.Defender.TotalPower,
	}, nil
}

func buildBattleNo(req Request) string {
	return fmt.Sprintf("%s-%d-%d-%d", req.Type, req.Attacker.PlayerID, req.Attacker.TotalPower, req.Defender.TotalPower)
}

func winnerSide(result string) string {
	switch result {
	case "success":
		return "attacker"
	case "fail":
		return "defender"
	default:
		return ""
	}
}
