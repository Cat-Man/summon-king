package wxmini

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
)

type Service struct {
	accounts account.Repository
	now      func() time.Time
}

type Option func(*Service)

const defaultGameBaseURL = "https://game.xxx.com"

func WithNow(now func() time.Time) Option {
	return func(s *Service) {
		s.now = now
	}
}

func NewService(accounts account.Repository, options ...Option) *Service {
	svc := &Service{
		accounts: accounts,
		now:      time.Now,
	}
	for _, option := range options {
		if option != nil {
			option(svc)
		}
	}
	return svc
}

func (s *Service) ExchangeLogin(ctx context.Context, code string) (LoginExchangeResponse, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return LoginExchangeResponse{}, ErrCodeRequired
	}

	token := buildUnifiedToken(code, s.now())
	nickname := buildNickname(code)
	guest, err := s.accounts.CreateGuest(ctx, nickname, token)
	if err != nil {
		return LoginExchangeResponse{}, err
	}

	return LoginExchangeResponse{
		UnifiedToken: token,
		PlayerID:     guest.PlayerID,
		Nickname:     guest.Nickname,
		Channel:      Channel,
	}, nil
}

func (s *Service) BootstrapSession(ctx context.Context, unifiedToken, gameBaseURL string) (SessionBootstrapResponse, error) {
	unifiedToken = strings.TrimSpace(unifiedToken)
	if unifiedToken == "" {
		return SessionBootstrapResponse{}, ErrUnifiedTokenRequired
	}
	gameBaseURL = strings.TrimSpace(gameBaseURL)
	if gameBaseURL == "" {
		gameBaseURL = defaultGameBaseURL
	}

	guest, err := s.accounts.GetByToken(ctx, unifiedToken)
	if err != nil {
		if errors.Is(err, account.ErrGuestAccountNotFound) {
			return SessionBootstrapResponse{}, ErrSessionNotFound
		}
		return SessionBootstrapResponse{}, err
	}

	return SessionBootstrapResponse{
		Token:    guest.Token,
		PlayerID: guest.PlayerID,
		Nickname: guest.Nickname,
		Channel:  Channel,
		GameURL:  buildGameURL(gameBaseURL, guest.Token),
	}, nil
}

func buildGameURL(baseURL, token string) string {
	baseWithQuery, hash, _ := strings.Cut(baseURL, "#")
	separator := "?"
	if strings.Contains(baseWithQuery, "?") {
		separator = "&"
	}
	gameURL := fmt.Sprintf("%s%schannel=%s&token=%s", baseWithQuery, separator, Channel, token)
	if hash == "" {
		return gameURL
	}
	return gameURL + "#" + hash
}

func buildUnifiedToken(code string, now time.Time) string {
	digest := sha1.Sum([]byte(code))
	hexDigest := hex.EncodeToString(digest[:])
	return fmt.Sprintf("%s-%d-%s", Channel, now.UnixNano(), hexDigest[:10])
}

func buildNickname(code string) string {
	digest := sha1.Sum([]byte(code))
	hexDigest := hex.EncodeToString(digest[:])
	return "wx_" + hexDigest[:8]
}
