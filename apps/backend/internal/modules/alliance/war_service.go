package alliance

import "context"

func (s *Service) SignAllianceWar(ctx context.Context, allianceID string) (AllianceWarRecord, error) {
	if _, err := s.repo.GetAlliance(ctx, allianceID); err != nil {
		return AllianceWarRecord{}, err
	}
	record := AllianceWarRecord{AllianceID: allianceID, Signed: true, Status: "signed", CheckedIn: []int64{}}
	return record, s.repo.SaveWarSignUp(ctx, record)
}

func (s *Service) CheckInAllianceWar(ctx context.Context, allianceID string, playerID int64) (AllianceWarRecord, error) {
	record, _ := s.repo.GetWarRecord(ctx, allianceID)
	record.AllianceID = allianceID
	record.Signed = true
	record.Status = "checkin"
	record.CheckedIn = append(record.CheckedIn, playerID)
	return record, s.repo.SaveWarSignUp(ctx, record)
}

func (s *Service) MatchAllianceWar(ctx context.Context, allianceID, opponentID string) (AllianceWarRecord, error) {
	record, _ := s.repo.GetWarRecord(ctx, allianceID)
	record.AllianceID = allianceID
	record.OpponentID = opponentID
	record.Status = "matched"
	return record, s.repo.SaveWarSignUp(ctx, record)
}

func (s *Service) SettleAllianceWar(ctx context.Context, allianceID string, won bool) (AllianceWarRecord, error) {
	record, _ := s.repo.GetWarRecord(ctx, allianceID)
	record.AllianceID = allianceID
	record.Status = "settled"
	if won {
		record.WarPoints += 100
	}
	return record, s.repo.SaveWarSignUp(ctx, record)
}

func (s *Service) RedeemWarPoints(ctx context.Context, allianceID string, points int64) (AllianceWarRecord, error) {
	record, _ := s.repo.GetWarRecord(ctx, allianceID)
	record.AllianceID = allianceID
	if record.WarPoints >= points {
		record.WarPoints -= points
	}
	return record, s.repo.SaveWarSignUp(ctx, record)
}
