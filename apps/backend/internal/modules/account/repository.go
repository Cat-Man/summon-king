package account

import (
	"context"
	"errors"
	"sync"
)

type Repository interface {
	CreateGuest(ctx context.Context, nickname, token string) (GuestLoginResponse, error)
	GetByToken(ctx context.Context, token string) (GuestLoginResponse, error)
}

var ErrGuestAccountNotFound = errors.New("guest account not found")

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

func (r *MemoryRepository) GetByToken(_ context.Context, token string) (GuestLoginResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, ok := r.accounts[token]
	if !ok {
		return GuestLoginResponse{}, ErrGuestAccountNotFound
	}
	return account, nil
}
