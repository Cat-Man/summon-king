package growth

import (
	"context"
	"sync"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
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

type MySQLRepository struct {
	store mysqlstore.ModuleStateStore
}

type mysqlState struct {
	Wallet           Wallet      `json:"wallet"`
	SpiritBonusPower int64       `json:"spirit_bonus_power"`
	Plots            []ManorPlot `json:"plots"`
}

const mysqlModuleName = "growth"

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		wallets: make(map[int64]Wallet),
		updated: make(map[int64]time.Time),
		plots:   make(map[int64][]ManorPlot),
	}
}

func NewMySQLRepository(store mysqlstore.ModuleStateStore) *MySQLRepository {
	return &MySQLRepository{store: store}
}

func (r *MemoryRepository) ensureWallet(playerID int64) Wallet {
	if wallet, ok := r.wallets[playerID]; ok {
		return wallet
	}
	wallet := Wallet{
		PlayerID:         playerID,
		SpiritPower:      100,
		SpiritBonusPower: 100,
		SpiritFreeWash:   3,
		BoneLevel:        1,
		SoulPieces:       0,
		ManorPlots:       2,
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
	r.mu.Lock()
	defer r.mu.Unlock()

	wallet := r.ensureWallet(playerID)
	wallet.SpiritPower += delta
	if wallet.SpiritPower < 0 {
		wallet.SpiritPower = 0
	}
	if delta > 0 {
		wallet.SpiritBonusPower += delta
	}
	r.wallets[playerID] = wallet
	r.updated[playerID] = time.Now()
	return nil
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

func (r *MySQLRepository) GetWallet(ctx context.Context, playerID int64) (Wallet, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return Wallet{}, err
	}
	return state.Wallet, nil
}

func (r *MySQLRepository) BeginFreeWash(ctx context.Context, playerID int64) (bool, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return false, err
	}
	if state.Wallet.SpiritFreeWash <= 0 {
		return false, nil
	}
	state.Wallet.SpiritFreeWash--
	return true, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) CommitWash(ctx context.Context, playerID int64, delta int64) error {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return err
	}
	state.Wallet.SpiritPower += delta
	if state.Wallet.SpiritPower < 0 {
		state.Wallet.SpiritPower = 0
	}
	return r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) UpdateSpiritPower(ctx context.Context, playerID int64, delta int64) error {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return err
	}
	state.Wallet.SpiritPower += delta
	if state.Wallet.SpiritPower < 0 {
		state.Wallet.SpiritPower = 0
	}
	if delta > 0 {
		state.Wallet.SpiritBonusPower += delta
	}
	return r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) UpgradeBoneLevel(ctx context.Context, playerID int64, delta int) (Wallet, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return Wallet{}, err
	}
	state.Wallet.BoneLevel += delta
	if state.Wallet.BoneLevel < 1 {
		state.Wallet.BoneLevel = 1
	}
	return state.Wallet, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) UpgradeSoulPower(ctx context.Context, playerID int64, delta int) (Wallet, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return Wallet{}, err
	}
	state.Wallet.SoulPieces += delta
	if state.Wallet.SoulPieces < 0 {
		state.Wallet.SoulPieces = 0
	}
	return state.Wallet, r.saveState(ctx, playerID, state)
}

func (r *MySQLRepository) GetManorPlots(ctx context.Context, playerID int64) ([]ManorPlot, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return nil, err
	}
	return append([]ManorPlot(nil), state.Plots...), nil
}

func (r *MySQLRepository) HarvestManor(ctx context.Context, playerID int64) ([]ManorPlot, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return nil, err
	}
	next := make([]ManorPlot, len(state.Plots))
	for i, plot := range state.Plots {
		plot.State = "冷却中"
		next[i] = plot
	}
	state.Plots = next
	if err := r.saveState(ctx, playerID, state); err != nil {
		return nil, err
	}
	return append([]ManorPlot(nil), next...), nil
}

func (r *MySQLRepository) PlantManor(ctx context.Context, playerID int64) ([]ManorPlot, error) {
	state, err := r.loadState(ctx, playerID)
	if err != nil {
		return nil, err
	}
	next := make([]ManorPlot, len(state.Plots))
	for i, plot := range state.Plots {
		plot.State = "成长中"
		next[i] = plot
	}
	state.Plots = next
	if err := r.saveState(ctx, playerID, state); err != nil {
		return nil, err
	}
	return append([]ManorPlot(nil), next...), nil
}

func (r *MySQLRepository) loadState(ctx context.Context, playerID int64) (mysqlState, error) {
	var state mysqlState
	ok, err := r.store.LoadModuleState(ctx, playerID, mysqlModuleName, &state)
	if err != nil {
		return mysqlState{}, err
	}
	if !ok {
		state = defaultMySQLState(playerID)
		return state, nil
	}
	if state.SpiritBonusPower > 0 {
		state.Wallet.SpiritBonusPower = state.SpiritBonusPower
	}
	if state.Wallet.PlayerID == 0 {
		state.Wallet.PlayerID = playerID
	}
	if state.Wallet.BoneLevel < 1 {
		state.Wallet.BoneLevel = 1
	}
	if state.Wallet.ManorPlots == 0 {
		state.Wallet.ManorPlots = 2
	}
	if len(state.Plots) == 0 {
		state.Plots = defaultPlots()
	}
	return state, nil
}

func (r *MySQLRepository) saveState(ctx context.Context, playerID int64, state mysqlState) error {
	if state.Wallet.PlayerID == 0 {
		state.Wallet.PlayerID = playerID
	}
	if state.Wallet.BoneLevel < 1 {
		state.Wallet.BoneLevel = 1
	}
	if state.Wallet.ManorPlots == 0 {
		state.Wallet.ManorPlots = 2
	}
	if len(state.Plots) == 0 {
		state.Plots = defaultPlots()
	}
	state.SpiritBonusPower = state.Wallet.SpiritBonusPower
	return r.store.SaveModuleState(ctx, playerID, mysqlModuleName, state)
}

func defaultMySQLState(playerID int64) mysqlState {
	return mysqlState{
		Wallet: normalizeWallet(playerID, Wallet{}),
		Plots:  defaultPlots(),
	}
}

func normalizeWallet(playerID int64, wallet Wallet) Wallet {
	if wallet.PlayerID == 0 {
		wallet.PlayerID = playerID
	}
	if wallet.SpiritPower == 0 {
		wallet.SpiritPower = 100
	}
	if wallet.SpiritBonusPower == 0 {
		wallet.SpiritBonusPower = 100
	}
	if wallet.SpiritFreeWash == 0 {
		wallet.SpiritFreeWash = 3
	}
	if wallet.BoneLevel <= 0 {
		wallet.BoneLevel = 1
	}
	if wallet.ManorPlots == 0 {
		wallet.ManorPlots = 2
	}
	return wallet
}

func defaultPlots() []ManorPlot {
	return []ManorPlot{
		{PlotID: 1, State: "空闲"},
		{PlotID: 2, State: "成长中"},
	}
}
