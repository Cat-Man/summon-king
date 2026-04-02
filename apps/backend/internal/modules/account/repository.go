package account

import (
	"context"
	"sync"
)

type Repository interface {
	CreateGuest(ctx context.Context, nickname, token string) (GuestLoginResponse, error)
}

type MemoryRepository struct {
	mu           sync.Mutex
	nextPlayerID int64
	accounts     map[string]GuestLoginResponse
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextPlayerID: 1000,
		accounts:     make(map[string]GuestLoginResponse),
	}
}

func (r *MemoryRepository) CreateGuest(_ context.Context, nickname, token string) (GuestLoginResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextPlayerID++
	account := GuestLoginResponse{
		PlayerID: r.nextPlayerID,
		Token:    token,
		Nickname: nickname,
	}
	r.accounts[token] = account

	return account, nil
}
