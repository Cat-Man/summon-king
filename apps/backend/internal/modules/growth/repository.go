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
	UpgradeBoneLevel(ctx context.Context, playerID int64, delta int) (Wallet, error)
	UpgradeSoulPower(ctx context.Context, playerID int64, delta int) (Wallet, error)
	GetManorPlots(ctx context.Context, playerID int64) ([]ManorPlot, error)
	HarvestManor(ctx context.Context, playerID int64) ([]ManorPlot, error)
	PlantManor(ctx context.Context, playerID int64) ([]ManorPlot, error)
}

type MemoryRepository struct {
	mu      sync.Mutex
	wallets map[int64]Wallet
	updated map[int64]time.Time
	plots   map[int64][]ManorPlot
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		wallets: make(map[int64]Wallet),
		updated: make(map[int64]time.Time),
		plots:   make(map[int64][]ManorPlot),
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

func (r *MemoryRepository) UpgradeBoneLevel(_ context.Context, playerID int64, delta int) (Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	wallet := r.ensureWallet(playerID)
	wallet.BoneLevel += delta
	if wallet.BoneLevel < 1 {
		wallet.BoneLevel = 1
	}
	r.wallets[playerID] = wallet
	r.updated[playerID] = time.Now()
	return wallet, nil
}

func (r *MemoryRepository) UpgradeSoulPower(_ context.Context, playerID int64, delta int) (Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	wallet := r.ensureWallet(playerID)
	wallet.SoulPieces += delta
	if wallet.SoulPieces < 0 {
		wallet.SoulPieces = 0
	}
	r.wallets[playerID] = wallet
	r.updated[playerID] = time.Now()
	return wallet, nil
}

func (r *MemoryRepository) GetManorPlots(_ context.Context, playerID int64) ([]ManorPlot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]ManorPlot(nil), r.ensurePlots(playerID)...), nil
}

func (r *MemoryRepository) HarvestManor(_ context.Context, playerID int64) ([]ManorPlot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	plots := r.ensurePlots(playerID)
	next := make([]ManorPlot, len(plots))
	for i, plot := range plots {
		plot.State = "冷却中"
		next[i] = plot
	}
	r.plots[playerID] = next
	return append([]ManorPlot(nil), next...), nil
}

func (r *MemoryRepository) PlantManor(_ context.Context, playerID int64) ([]ManorPlot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	plots := r.ensurePlots(playerID)
	next := make([]ManorPlot, len(plots))
	for i, plot := range plots {
		plot.State = "成长中"
		next[i] = plot
	}
	r.plots[playerID] = next
	return append([]ManorPlot(nil), next...), nil
}

func (r *MemoryRepository) ensurePlots(playerID int64) []ManorPlot {
	if plots, ok := r.plots[playerID]; ok {
		return plots
	}
	plots := []ManorPlot{
		{PlotID: 1, State: "空闲"},
		{PlotID: 2, State: "成长中"},
	}
	r.plots[playerID] = plots
	return plots
}
