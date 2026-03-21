package battle

import (
	"fmt"
	"sort"
	"time"
)

const (
	defaultBattleHP      = 100
	defaultBattleAttack  = 50
	defaultBattleDefense = 10
	maxRounds            = 20
)

type skillResolution struct {
	name       string
	multiplier float64
}

type turnActor struct {
	side string
	unit *Unit
}

func RunBattle(input BattleInput) BattleResult {
	left := normalizeUnits(input.Left)
	right := normalizeUnits(input.Right)
	result := BattleResult{
		BattleNo:  battleNoOrDefault(input.BattleNo),
		CreatedAt: time.Now(),
	}

	for roundNumber := 1; roundNumber <= maxRounds; roundNumber++ {
		if nextAlive(left) == nil || nextAlive(right) == nil {
			break
		}

		turns := []turnActor{
			{side: "left", unit: nextAlive(left)},
			{side: "right", unit: nextAlive(right)},
		}
		sort.SliceStable(turns, func(i, j int) bool {
			if turns[i].unit.Speed == turns[j].unit.Speed {
				return turns[i].side < turns[j].side
			}
			return turns[i].unit.Speed > turns[j].unit.Speed
		})

		round := BattleRound{Number: roundNumber}
		for _, turn := range turns {
			if turn.unit == nil || turn.unit.HP <= 0 {
				continue
			}

			targets := opposingTeam(turn.side, left, right)
			target := nextAlive(targets)
			if target == nil {
				break
			}

			skill := resolveSkill(turn.unit)
			damage := CalculateNormalAttackDamage(turn.unit.Attack, target.Defense, skill.multiplier)
			target.HP -= damage
			if target.HP < 0 {
				target.HP = 0
			}

			action := BattleAction{
				Actor:          turn.unit.Name,
				Target:         target.Name,
				Damage:         damage,
				TargetHP:       target.HP,
				TargetDefeated: target.HP == 0,
			}
			if skill.name != "" {
				action.Skill = skill.name
			}
			round.Actions = append(round.Actions, action)

			tickStatuses(turn.unit)
			tickStatuses(target)

			if nextAlive(targets) == nil {
				break
			}
		}

		if len(round.Actions) > 0 {
			result.Rounds = append(result.Rounds, round)
		}
	}

	result.Winner = decideWinner(left, right)
	result.Left = left
	result.Right = right
	return result
}

func normalizeUnits(units []Unit) []Unit {
	if len(units) == 0 {
		return nil
	}

	normalized := make([]Unit, 0, len(units))
	for index, unit := range units {
		if unit.Name == "" {
			unit.Name = fmt.Sprintf("unit-%d", index+1)
		}
		if unit.MaxHP <= 0 {
			if unit.HP > 0 {
				unit.MaxHP = unit.HP
			} else {
				unit.MaxHP = defaultBattleHP
			}
		}
		if unit.HP <= 0 {
			unit.HP = unit.MaxHP
		}
		if unit.Attack <= 0 {
			unit.Attack = defaultBattleAttack
		}
		if unit.Defense < 0 {
			unit.Defense = 0
		}
		if unit.Defense == 0 {
			unit.Defense = defaultBattleDefense
		}
		normalized = append(normalized, unit)
	}

	return normalized
}

func nextAlive(units []Unit) *Unit {
	for index := range units {
		if units[index].HP > 0 {
			return &units[index]
		}
	}
	return nil
}

func opposingTeam(side string, left, right []Unit) []Unit {
	if side == "left" {
		return right
	}
	return left
}

func resolveSkill(unit *Unit) skillResolution {
	if unit == nil || len(unit.Skills) == 0 {
		return skillResolution{multiplier: 1}
	}

	for _, skill := range unit.Skills {
		if !skill.Enabled {
			continue
		}
		multiplier := skill.Multiplier
		if multiplier <= 0 {
			multiplier = 1
		}
		if skill.TriggerRate >= 1 {
			return skillResolution{name: skill.Name, multiplier: multiplier}
		}
	}

	return skillResolution{multiplier: 1}
}

func decideWinner(left, right []Unit) string {
	leftAlive := nextAlive(left) != nil
	rightAlive := nextAlive(right) != nil

	switch {
	case leftAlive && !rightAlive:
		return "left"
	case rightAlive && !leftAlive:
		return "right"
	default:
		return "draw"
	}
}

func battleNoOrDefault(battleNo string) string {
	if battleNo != "" {
		return battleNo
	}
	return fmt.Sprintf("battle-%d", time.Now().UnixNano())
}
