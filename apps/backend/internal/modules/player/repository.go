package player

import (
	"context"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
)

type Repository interface {
	GetByPlayerID(ctx context.Context, playerID int64) (account.Player, error)
}

type AssetWalletReader interface {
	GetWallet(ctx context.Context, playerID int64) (asset.Wallet, error)
}

type repository struct {
	accountRepo account.Repository
	assetRepo   AssetWalletReader
}

func NewRepository(accountRepo account.Repository, assetRepo ...AssetWalletReader) Repository {
	repo := &repository{accountRepo: accountRepo}
	if len(assetRepo) > 0 {
		repo.assetRepo = assetRepo[0]
	}
	return repo
}

func (r *repository) GetByPlayerID(ctx context.Context, playerID int64) (account.Player, error) {
	player, err := r.accountRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return account.Player{}, err
	}

	if r.assetRepo != nil {
		if wallet, err := r.assetRepo.GetWallet(ctx, playerID); err == nil {
			player.Wallet.Coin = wallet.Coins
			player.Wallet.Diamond = wallet.Diamonds
		}
	}

	return player, nil
}
