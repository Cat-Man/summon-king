package player

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
)

type Repository interface {
	GetByPlayerID(ctx context.Context, playerID int64) (account.Player, error)
}

type repository struct {
	accountRepo account.Repository
}

func NewRepository(accountRepo account.Repository) Repository {
	return &repository{accountRepo: accountRepo}
}

func (r *repository) GetByPlayerID(ctx context.Context, playerID int64) (account.Player, error) {
	return r.accountRepo.GetByPlayerID(ctx, playerID)
}
