package asset

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidRequest      = errors.New("invalid request")
	ErrInventoryNotEnough  = errors.New("inventory not enough")
	ErrInventoryItemAbsent = errors.New("inventory item not found")
)

type Repository interface {
	GetWallet(ctx context.Context, playerID int64) (Wallet, error)
	ListInventory(ctx context.Context, playerID int64) ([]InventoryItem, error)
	GrantRewardIdempotent(ctx context.Context, grant RewardGrant, now time.Time) (GrantRewardResult, error)
	UseInventoryItem(ctx context.Context, req InventoryOperateRequest, now time.Time) (InventoryOperateResult, error)
	SellInventoryItem(ctx context.Context, req InventoryOperateRequest, now time.Time) (InventoryOperateResult, error)
}

type MemoryRepository struct {
	mu                sync.Mutex
	wallets           map[int64]Wallet
	inventories       map[int64]map[string]InventoryItem
	resourceChangeLog map[int64][]ResourceChangeLog
	rewardClaimLog    map[int64]map[string]RewardClaimLog
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		wallets:           make(map[int64]Wallet),
		inventories:       make(map[int64]map[string]InventoryItem),
		resourceChangeLog: make(map[int64][]ResourceChangeLog),
		rewardClaimLog:    make(map[int64]map[string]RewardClaimLog),
	}
}

func (r *MemoryRepository) GetWallet(_ context.Context, playerID int64) (Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ensureWalletUnlocked(playerID), nil
}

func (r *MemoryRepository) ListInventory(_ context.Context, playerID int64) ([]InventoryItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := r.ensureInventoryUnlocked(playerID)
	result := make([]InventoryItem, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result, nil
}

func (r *MemoryRepository) GrantRewardIdempotent(_ context.Context, grant RewardGrant, now time.Time) (GrantRewardResult, error) {
	if grant.PlayerID == 0 || grant.BizID == "" {
		return GrantRewardResult{}, ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	claims := r.ensureRewardClaimsUnlocked(grant.PlayerID)
	if _, exists := claims[grant.BizID]; exists {
		return GrantRewardResult{
			Wallet:     r.ensureWalletUnlocked(grant.PlayerID),
			Idempotent: true,
		}, nil
	}

	wallet := r.ensureWalletUnlocked(grant.PlayerID)
	wallet.Coins += grant.Coins
	wallet.Diamonds += grant.Diamonds
	wallet.UpdatedAt = now
	r.wallets[grant.PlayerID] = wallet

	claims[grant.BizID] = RewardClaimLog{
		PlayerID:  grant.PlayerID,
		BizID:     grant.BizID,
		Source:    grant.Source,
		Coins:     grant.Coins,
		Diamonds:  grant.Diamonds,
		ClaimedAt: now,
	}

	r.resourceChangeLog[grant.PlayerID] = append(r.resourceChangeLog[grant.PlayerID], ResourceChangeLog{
		PlayerID:      grant.PlayerID,
		ChangeType:    "grant_reward",
		BizID:         grant.BizID,
		CoinsDelta:    grant.Coins,
		DiamondsDelta: grant.Diamonds,
		Reason:        grant.Source,
		CreatedAt:     now,
	})

	return GrantRewardResult{
		Wallet:     wallet,
		Idempotent: false,
	}, nil
}

func (r *MemoryRepository) UseInventoryItem(_ context.Context, req InventoryOperateRequest, now time.Time) (InventoryOperateResult, error) {
	if req.PlayerID == 0 || req.ItemID == "" || req.Count <= 0 {
		return InventoryOperateResult{}, ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	items := r.ensureInventoryUnlocked(req.PlayerID)
	item, ok := items[req.ItemID]
	if !ok {
		return InventoryOperateResult{}, ErrInventoryItemAbsent
	}
	if item.Quantity < req.Count {
		return InventoryOperateResult{}, ErrInventoryNotEnough
	}

	item.Quantity -= req.Count
	item.PlayerID = req.PlayerID
	items[req.ItemID] = item

	r.resourceChangeLog[req.PlayerID] = append(r.resourceChangeLog[req.PlayerID], ResourceChangeLog{
		PlayerID:      req.PlayerID,
		ChangeType:    "inventory_use",
		BizID:         req.ItemID,
		CoinsDelta:    0,
		DiamondsDelta: 0,
		Reason:        "inventory_use",
		CreatedAt:     now,
	})

	return InventoryOperateResult{
		Wallet: r.ensureWalletUnlocked(req.PlayerID),
		Item:   item,
	}, nil
}

func (r *MemoryRepository) SellInventoryItem(_ context.Context, req InventoryOperateRequest, now time.Time) (InventoryOperateResult, error) {
	if req.PlayerID == 0 || req.ItemID == "" || req.Count <= 0 {
		return InventoryOperateResult{}, ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	items := r.ensureInventoryUnlocked(req.PlayerID)
	item, ok := items[req.ItemID]
	if !ok {
		return InventoryOperateResult{}, ErrInventoryItemAbsent
	}
	if item.Quantity < req.Count {
		return InventoryOperateResult{}, ErrInventoryNotEnough
	}

	item.Quantity -= req.Count
	item.PlayerID = req.PlayerID
	items[req.ItemID] = item

	wallet := r.ensureWalletUnlocked(req.PlayerID)
	coinDelta := item.SellPrice * req.Count
	wallet.Coins += coinDelta
	wallet.UpdatedAt = now
	r.wallets[req.PlayerID] = wallet

	r.resourceChangeLog[req.PlayerID] = append(r.resourceChangeLog[req.PlayerID], ResourceChangeLog{
		PlayerID:      req.PlayerID,
		ChangeType:    "inventory_sell",
		BizID:         req.ItemID,
		CoinsDelta:    coinDelta,
		DiamondsDelta: 0,
		Reason:        "inventory_sell",
		CreatedAt:     now,
	})

	return InventoryOperateResult{
		Wallet: wallet,
		Item:   item,
	}, nil
}

func (r *MemoryRepository) ensureWalletUnlocked(playerID int64) Wallet {
	wallet, ok := r.wallets[playerID]
	if !ok {
		wallet = Wallet{PlayerID: playerID}
		r.wallets[playerID] = wallet
	}
	return wallet
}

func (r *MemoryRepository) ensureInventoryUnlocked(playerID int64) map[string]InventoryItem {
	items, ok := r.inventories[playerID]
	if !ok {
		items = make(map[string]InventoryItem, len(defaultInventory))
		for _, template := range defaultInventory {
			template.PlayerID = playerID
			items[template.ItemID] = template
		}
		r.inventories[playerID] = items
	}
	return items
}

func (r *MemoryRepository) ensureRewardClaimsUnlocked(playerID int64) map[string]RewardClaimLog {
	claims, ok := r.rewardClaimLog[playerID]
	if !ok {
		claims = make(map[string]RewardClaimLog)
		r.rewardClaimLog[playerID] = claims
	}
	return claims
}
