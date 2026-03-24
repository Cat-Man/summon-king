package battle

const damageRandomFactor = 0.07

func CalculateNormalAttackDamage(attack, defense int, multiplier float64) int {
	if multiplier <= 0 {
		multiplier = 1
	}

	gap := attack - defense
	if gap < 0 {
		gap = 1
	}

	damage := float64(gap) * damageRandomFactor * defenseBandMultiplier(defense) * multiplier
	if damage < 1 {
		return 1
	}

	return int(damage + 0.5)
}

func defenseBandMultiplier(defense int) float64 {
	switch {
	case defense < 1000:
		return 3.8
	case defense < 2000:
		return 2.0
	case defense < 3000:
		return 1.6
	case defense < 4000:
		return 1.3
	case defense < 5000:
		return 1.1
	default:
		return 1.0
	}
}
