package account

import (
	"context"
	"errors"
	"sync"

	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

type Repository interface {
	CreateGuest(ctx context.Context, nickname, token string) (GuestLoginResponse, error)
	GetByToken(ctx context.Context, token string) (GuestLoginResponse, error)
	GetByPlayerID(ctx context.Context, playerID int64) (GuestLoginResponse, error)
}

var ErrGuestAccountNotFound = errors.New("guest account not found")

type MemoryRepository struct {
	mu           sync.Mutex
	nextPlayerID int64
	accounts     map[string]GuestLoginResponse
}

type MySQLRepository struct {
	store mysqlstore.GuestAccountStore
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextPlayerID: 1000,
		accounts:     make(map[string]GuestLoginResponse),
	}
}

func NewMySQLRepository(store mysqlstore.GuestAccountStore) *MySQLRepository {
	return &MySQLRepository{store: store}
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

func (r *MemoryRepository) GetByPlayerID(_ context.Context, playerID int64) (GuestLoginResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, account := range r.accounts {
		if account.PlayerID == playerID {
			return account, nil
		}
	}
	return GuestLoginResponse{}, ErrGuestAccountNotFound
}

func (r *MySQLRepository) CreateGuest(ctx context.Context, nickname, token string) (GuestLoginResponse, error) {
	record, err := r.store.CreateGuestAccount(ctx, nickname, token)
	if err != nil {
		return GuestLoginResponse{}, err
	}
	return GuestLoginResponse{
		PlayerID: record.PlayerID,
		Token:    record.Token,
		Nickname: record.Nickname,
	}, nil
}

func (r *MySQLRepository) GetByToken(ctx context.Context, token string) (GuestLoginResponse, error) {
	record, err := r.store.GetGuestAccountByToken(ctx, token)
	if err != nil {
		return GuestLoginResponse{}, ErrGuestAccountNotFound
	}
	return GuestLoginResponse{
		PlayerID: record.PlayerID,
		Token:    record.Token,
		Nickname: record.Nickname,
	}, nil
}

func (r *MySQLRepository) GetByPlayerID(ctx context.Context, playerID int64) (GuestLoginResponse, error) {
	record, err := r.store.GetGuestAccountByPlayerID(ctx, playerID)
	if err != nil {
		return GuestLoginResponse{}, ErrGuestAccountNotFound
	}
	return GuestLoginResponse{
		PlayerID: record.PlayerID,
		Token:    record.Token,
		Nickname: record.Nickname,
	}, nil
}
