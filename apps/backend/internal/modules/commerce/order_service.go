package commerce

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrOrderNotFound = errors.New("order not found")

type Repository interface {
	GetVIPState(ctx context.Context, playerID int64) VIPState
	SaveVIPState(ctx context.Context, state VIPState)
	GetGiftState(ctx context.Context, playerID int64) GiftState
	SaveGiftState(ctx context.Context, state GiftState)
	CreateOrder(ctx context.Context, playerID int64, idemKey string, productID int64, now time.Time) (Order, error)
	GetOrder(ctx context.Context, orderNo string) (Order, error)
	SaveOrder(ctx context.Context, order Order) error
}

type MemoryRepository struct {
	mu          sync.Mutex
	nextOrder   int64
	vip         map[int64]VIPState
	gifts       map[int64]GiftState
	orders      map[string]Order
	idemToOrder map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{nextOrder: 1, vip: make(map[int64]VIPState), gifts: make(map[int64]GiftState), orders: make(map[string]Order), idemToOrder: make(map[string]string)}
}

func (r *MemoryRepository) GetVIPState(_ context.Context, playerID int64) VIPState {
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.vip[playerID]
	if !ok {
		state = VIPState{PlayerID: playerID, Level: 0}
		r.vip[playerID] = state
	}
	return state
}

func (r *MemoryRepository) SaveVIPState(_ context.Context, state VIPState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vip[state.PlayerID] = state
}

func (r *MemoryRepository) GetGiftState(_ context.Context, playerID int64) GiftState {
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.gifts[playerID]
	if !ok {
		state = GiftState{PlayerID: playerID, GiftClaims: map[string]bool{}, RedeemRecords: map[string]bool{}}
		r.gifts[playerID] = state
	}
	return state
}

func (r *MemoryRepository) SaveGiftState(_ context.Context, state GiftState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gifts[state.PlayerID] = state
}

func (r *MemoryRepository) CreateOrder(_ context.Context, playerID int64, idemKey string, productID int64, now time.Time) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%d:%s:%d", playerID, idemKey, productID)
	if orderNo, ok := r.idemToOrder[key]; ok {
		return r.orders[orderNo], nil
	}
	orderNo := fmt.Sprintf("ord-%d", r.nextOrder)
	r.nextOrder++
	order := Order{OrderNo: orderNo, PlayerID: playerID, IdempotencyKey: idemKey, ProductID: productID, Status: "created", CreatedAt: now}
	r.orders[orderNo] = order
	r.idemToOrder[key] = orderNo
	return order, nil
}

func (r *MemoryRepository) GetOrder(_ context.Context, orderNo string) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[orderNo]
	if !ok {
		return Order{}, ErrOrderNotFound
	}
	return order, nil
}

func (r *MemoryRepository) SaveOrder(_ context.Context, order Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.OrderNo] = order
	return nil
}
