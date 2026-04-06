package account

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GuestLogin(ctx context.Context, nickname string) (GuestLoginResponse, error) {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return GuestLoginResponse{}, errors.New("nickname is required")
	}

	token := fmt.Sprintf("guest-%d", time.Now().UnixNano())
	return s.repo.CreateGuest(ctx, nickname, token)
}
