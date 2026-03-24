package alliance

import (
	"context"
	"time"
)

type Alliance struct {
	AllianceID string         `json:"alliance_id"`
	Name       string         `json:"name"`
	OwnerID    int64          `json:"owner_id"`
	Notice     string         `json:"notice"`
	Funds      int64          `json:"funds"`
	Buildings  map[string]int `json:"buildings"`
	Members    []int64        `json:"members"`
	CreatedAt  time.Time      `json:"created_at"`
}

type FireTrainingRecord struct {
	RecordID      string    `json:"record_id"`
	AllianceID    string    `json:"alliance_id"`
	PlayerID      int64     `json:"player_id"`
	RoomID        string    `json:"room_id"`
	Status        string    `json:"status"`
	RewardCrystal int64     `json:"reward_crystal"`
	StartedAt     time.Time `json:"started_at"`
	FinishedAt    time.Time `json:"finished_at"`
}

type AllianceWarRecord struct {
	AllianceID string  `json:"alliance_id"`
	Signed     bool    `json:"signed"`
	CheckedIn  []int64 `json:"checked_in"`
	OpponentID string  `json:"opponent_id"`
	Status     string  `json:"status"`
	WarPoints  int64   `json:"war_points"`
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) CreateAlliance(ctx context.Context, playerID int64, name string) error {
	if s.repo.GetPlayerLevel(ctx, playerID) < 30 {
		return ErrLevelTooLow
	}
	_, err := s.repo.CreateAlliance(ctx, playerID, name, s.now())
	return err
}

func (s *Service) CreateAllianceAsGM(ctx context.Context, playerID int64, name string) (string, error) {
	alliance, err := s.repo.CreateAlliance(ctx, playerID, name, s.now())
	if err != nil {
		return "", err
	}
	return alliance.AllianceID, nil
}

func (s *Service) ApplyJoin(ctx context.Context, allianceID string, playerID int64) error {
	alliance, err := s.repo.GetAlliance(ctx, allianceID)
	if err != nil {
		return err
	}
	alliance.Members = append(alliance.Members, playerID)
	return s.repo.SaveAlliance(ctx, alliance)
}

func (s *Service) ApproveApplication(ctx context.Context, allianceID string, playerID int64) error {
	return s.ApplyJoin(ctx, allianceID, playerID)
}

func (s *Service) UpdateNotice(ctx context.Context, allianceID, notice string) error {
	alliance, err := s.repo.GetAlliance(ctx, allianceID)
	if err != nil {
		return err
	}
	alliance.Notice = notice
	return s.repo.SaveAlliance(ctx, alliance)
}

func (s *Service) Donate(ctx context.Context, allianceID string, coins int64) (Alliance, error) {
	alliance, err := s.repo.GetAlliance(ctx, allianceID)
	if err != nil {
		return Alliance{}, err
	}
	alliance.Funds += coins
	return alliance, s.repo.SaveAlliance(ctx, alliance)
}

func (s *Service) UpgradeBuilding(ctx context.Context, allianceID, building string) (Alliance, error) {
	alliance, err := s.repo.GetAlliance(ctx, allianceID)
	if err != nil {
		return Alliance{}, err
	}
	alliance.Buildings[building]++
	return alliance, s.repo.SaveAlliance(ctx, alliance)
}
