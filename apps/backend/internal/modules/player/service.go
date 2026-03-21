package player

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

type HomeIndex struct {
	PlayerID    int64     `json:"player_id"`
	Nickname    string    `json:"nickname"`
	Level       int       `json:"level"`
	Coin        int64     `json:"coin"`
	Diamond     int64     `json:"diamond"`
	LastLoginAt time.Time `json:"last_login_at"`
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetHomeIndex(ctx context.Context, playerID int64) (HomeIndex, error) {
	player, err := s.repo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return HomeIndex{}, err
	}

	return HomeIndex{
		PlayerID:    player.PlayerID,
		Nickname:    player.Profile.Nickname,
		Level:       player.Profile.Level,
		Coin:        player.Wallet.Coin,
		Diamond:     player.Wallet.Diamond,
		LastLoginAt: player.DailyState.LastLoginAt,
	}, nil
}
