package commerce

import (
	"context"
	"time"
)

type VIPState struct {
	PlayerID          int64 `json:"player_id"`
	Level             int   `json:"level"`
	DailyChestClaimed bool  `json:"daily_chest_claimed"`
}

type Order struct {
	OrderNo        string    `json:"order_no"`
	PlayerID       int64     `json:"player_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	ProductID      int64     `json:"product_id"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	Delivered      bool      `json:"delivered"`
}

type GiftState struct {
	PlayerID      int64           `json:"player_id"`
	GiftClaims    map[string]bool `json:"gift_claims"`
	SigninDays    int             `json:"signin_days"`
	RedeemRecords map[string]bool `json:"redeem_records"`
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) GetVIPState(ctx context.Context, playerID int64) VIPState {
	return s.repo.GetVIPState(ctx, playerID)
}

func (s *Service) CreateOrder(ctx context.Context, playerID int64, idemKey string, productID int64) (Order, error) {
	return s.repo.CreateOrder(ctx, playerID, idemKey, productID, s.now())
}

func (s *Service) MarkOrderPaid(ctx context.Context, orderNo string) (Order, error) {
	order, err := s.repo.GetOrder(ctx, orderNo)
	if err != nil {
		return Order{}, err
	}
	order.Status = "paid"
	return order, s.repo.SaveOrder(ctx, order)
}

func (s *Service) DeliverOrder(ctx context.Context, orderNo string) (Order, error) {
	order, err := s.repo.GetOrder(ctx, orderNo)
	if err != nil {
		return Order{}, err
	}
	order.Delivered = true
	order.Status = "delivered"
	return order, s.repo.SaveOrder(ctx, order)
}

func (s *Service) ClaimVIPDailyChest(ctx context.Context, playerID int64) VIPState {
	state := s.repo.GetVIPState(ctx, playerID)
	state.DailyChestClaimed = true
	s.repo.SaveVIPState(ctx, state)
	return state
}

func (s *Service) ClaimSignin(ctx context.Context, playerID int64) GiftState {
	state := s.repo.GetGiftState(ctx, playerID)
	state.SigninDays++
	s.repo.SaveGiftState(ctx, state)
	return state
}

func (s *Service) RedeemCode(ctx context.Context, playerID int64, code string) GiftState {
	state := s.repo.GetGiftState(ctx, playerID)
	state.RedeemRecords[code] = true
	s.repo.SaveGiftState(ctx, state)
	return state
}
