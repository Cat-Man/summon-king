package account

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) CreateGuestPlayer(ctx context.Context, channel string) (Player, error) {
	if channel == "" {
		channel = "web"
	}

	playerID, err := s.repo.NextPlayerID(ctx)
	if err != nil {
		return Player{}, err
	}

	now := s.now()
	player := Player{
		PlayerID: playerID,
		Channel:  channel,
		Token:    fmt.Sprintf("guest-token-%d", playerID),
		Profile: PlayerProfile{
			PlayerID: playerID,
			Nickname: fmt.Sprintf("游客%d", playerID),
			Level:    1,
		},
		Wallet: PlayerWallet{
			PlayerID: playerID,
			Coin:     0,
			Diamond:  0,
		},
		DailyState: PlayerDailyState{
			PlayerID:    playerID,
			LastLoginAt: now,
		},
	}

	if err := s.repo.Save(ctx, player); err != nil {
		return Player{}, err
	}

	return player, nil
}

func (s *Service) GetPlayerByID(ctx context.Context, playerID int64) (Player, error) {
	return s.repo.GetByPlayerID(ctx, playerID)
}
