package growth

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	GetWallet(ctx context.Context, playerID int64) (Wallet, error)
	BeginFreeWash(ctx context.Context, playerID int64) (bool, error)
	CommitWash(ctx context.Context, playerID int64, delta int64) error
	UpdateSpiritPower(ctx context.Context, playerID int64, delta int64) error
}

type MemoryRepository struct {
	mu      sync.Mutex
	wallets map[int64]Wallet
	updated map[int64]time.Time
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		wallets: make(map[int64]Wallet),
		updated: make(map[int64]time.Time),
	}
}

func (r *MemoryRepository) ensureWallet(playerID int64) Wallet {
	if wallet, ok := r.wallets[playerID]; ok {
		return wallet
	}
	wallet := Wallet{
		PlayerID:       playerID,
		SpiritPower:    100,
		SpiritFreeWash: 3,
		BoneLevel:      1,
		SoulPieces:     0,
		ManorPlots:     2,
	}
	r.wallets[playerID] = wallet
	return wallet
}

func (r *MemoryRepository) GetWallet(ctx context.Context, playerID int64) (Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ensureWallet(playerID), nil
}

func (r *MemoryRepository) BeginFreeWash(ctx context.Context, playerID int64) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	wallet := r.ensureWallet(playerID)
	if wallet.SpiritFreeWash > 0 {
		wallet.SpiritFreeWash--
		r.wallets[playerID] = wallet
		return true, nil
	}
	return false, nil
}

func (r *MemoryRepository) CommitWash(ctx context.Context, playerID int64, delta int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	wallet := r.ensureWallet(playerID)
	wallet.SpiritPower += delta
	if wallet.SpiritPower < 0 {
		wallet.SpiritPower = 0
	}
	r.wallets[playerID] = wallet
	r.updated[playerID] = time.Now()
	return nil
}

func (r *MemoryRepository) UpdateSpiritPower(ctx context.Context, playerID int64, delta int64) error {
	return r.CommitWash(ctx, playerID, delta)
}
