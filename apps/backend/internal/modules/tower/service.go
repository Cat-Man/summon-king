package tower

import "context"

type Service struct {
	pagodaFloor      map[int64]int
	spiritTowerFloor map[int64]int
}

type ChallengeResult struct {
	PlayerID int64  `json:"player_id"`
	Tower    string `json:"tower"`
	Floor    int    `json:"floor"`
	Victory  bool   `json:"victory"`
}

func NewService() *Service {
	return &Service{pagodaFloor: make(map[int64]int), spiritTowerFloor: make(map[int64]int)}
}

func (s *Service) ChallengePagoda(_ context.Context, playerID int64, victory bool) ChallengeResult {
	if victory {
		s.pagodaFloor[playerID]++
	}
	if s.pagodaFloor[playerID] == 0 {
		s.pagodaFloor[playerID] = 1
	}
	return ChallengeResult{PlayerID: playerID, Tower: "pagoda", Floor: s.pagodaFloor[playerID], Victory: victory}
}

func (s *Service) ChallengeSpiritTower(_ context.Context, playerID int64, victory bool) ChallengeResult {
	if victory {
		s.spiritTowerFloor[playerID]++
	}
	if s.spiritTowerFloor[playerID] == 0 {
		s.spiritTowerFloor[playerID] = 1
	}
	return ChallengeResult{PlayerID: playerID, Tower: "spirit", Floor: s.spiritTowerFloor[playerID], Victory: victory}
}
