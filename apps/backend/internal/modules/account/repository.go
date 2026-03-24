package account

import (
	"context"
	"errors"
	"sync"
)

var ErrPlayerNotFound = errors.New("player not found")

type Repository interface {
	NextPlayerID(ctx context.Context) (int64, error)
	Save(ctx context.Context, player Player) error
	GetByPlayerID(ctx context.Context, playerID int64) (Player, error)
	GetByToken(ctx context.Context, token string) (Player, error)
}

type MemoryRepository struct {
	mu           sync.Mutex
	nextPlayerID int64
	playersByID  map[int64]Player
	playerTokens map[string]int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextPlayerID: 1,
		playersByID:  make(map[int64]Player),
		playerTokens: make(map[string]int64),
	}
}

func (r *MemoryRepository) NextPlayerID(context.Context) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	playerID := r.nextPlayerID
	r.nextPlayerID++
	return playerID, nil
}

func (r *MemoryRepository) Save(_ context.Context, player Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.playersByID[player.PlayerID] = player
	r.playerTokens[player.Token] = player.PlayerID
	return nil
}

func (r *MemoryRepository) GetByPlayerID(_ context.Context, playerID int64) (Player, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	player, ok := r.playersByID[playerID]
	if !ok {
		return Player{}, ErrPlayerNotFound
	}
	return player, nil
}

func (r *MemoryRepository) GetByToken(ctx context.Context, token string) (Player, error) {
	r.mu.Lock()
	playerID, ok := r.playerTokens[token]
	r.mu.Unlock()
	if !ok {
		return Player{}, ErrPlayerNotFound
	}
	return r.GetByPlayerID(ctx, playerID)
}
