package battle

type StatusEffect struct {
	Code            string `json:"code"`
	RemainingRounds int    `json:"remaining_rounds"`
}

func tickStatuses(unit *Unit) {
	if unit == nil || len(unit.Statuses) == 0 {
		return
	}

	next := unit.Statuses[:0]
	for _, status := range unit.Statuses {
		if status.RemainingRounds > 0 {
			status.RemainingRounds--
		}
		if status.RemainingRounds > 0 {
			next = append(next, status)
		}
	}
	unit.Statuses = next
}
