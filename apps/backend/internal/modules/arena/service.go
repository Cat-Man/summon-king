package arena

import "context"

type DailyRecord struct {
	PlayerID      int64 `json:"player_id"`
	CurrentStreak int   `json:"current_streak"`
	BestStreak    int   `json:"best_streak"`
	RewardClaimed bool  `json:"reward_claimed"`
}

type Opponent struct {
	PlayerID int64  `json:"player_id"`
	Name     string `json:"name"`
	Power    int    `json:"power"`
}

type Service struct {
	records   map[int64]DailyRecord
	opponents map[int64][]Opponent
}

func NewService() *Service {
	return &Service{records: make(map[int64]DailyRecord), opponents: make(map[int64][]Opponent)}
}

func (s *Service) SetCurrentStreak(_ context.Context, playerID int64, streak int) error {
	rec := s.records[playerID]
	rec.PlayerID = playerID
	rec.CurrentStreak = streak
	if streak > rec.BestStreak {
		rec.BestStreak = streak
	}
	s.records[playerID] = rec
	return nil
}

func (s *Service) RecordBattleResult(_ context.Context, playerID int64, won bool) error {
	rec := s.records[playerID]
	rec.PlayerID = playerID
	if won {
		rec.CurrentStreak++
		if rec.CurrentStreak > rec.BestStreak {
			rec.BestStreak = rec.CurrentStreak
		}
	} else {
		rec.CurrentStreak = 0
	}
	s.records[playerID] = rec
	return nil
}

func (s *Service) GetDailyRecord(_ context.Context, playerID int64) DailyRecord {
	rec := s.records[playerID]
	rec.PlayerID = playerID
	return rec
}

func (s *Service) GetOpponents(_ context.Context, playerID int64) []Opponent {
	if list, ok := s.opponents[playerID]; ok {
		return list
	}
	list := []Opponent{{PlayerID: 2001, Name: "青云弟子", Power: 180}, {PlayerID: 2002, Name: "火焰守卫", Power: 210}, {PlayerID: 2003, Name: "霜港行者", Power: 230}}
	s.opponents[playerID] = list
	return list
}

func (s *Service) RefreshOpponents(ctx context.Context, playerID int64) []Opponent {
	delete(s.opponents, playerID)
	return s.GetOpponents(ctx, playerID)
}

func (s *Service) ClaimReward(_ context.Context, playerID int64) DailyRecord {
	rec := s.records[playerID]
	rec.PlayerID = playerID
	rec.RewardClaimed = true
	s.records[playerID] = rec
	return rec
}
