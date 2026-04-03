package home

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type accountReader interface {
	GetByToken(ctx context.Context, token string) (account.GuestLoginResponse, error)
}

type dungeonReader interface {
	GetWorldMap(ctx context.Context) (dungeon.WorldMap, error)
	GetRun(ctx context.Context, playerID int64) (dungeon.DungeonRun, error)
	GetCultivationStatus(ctx context.Context, playerID int64) (dungeon.CultivationStatus, error)
}

type Service struct {
	accounts accountReader
	dungeons dungeonReader
	growth   growth.Repository
}

func NewService(accounts accountReader, dungeons dungeonReader, growthRepo growth.Repository) *Service {
	return &Service{
		accounts: accounts,
		dungeons: dungeons,
		growth:   growthRepo,
	}
}

func (s *Service) GetOverview(ctx context.Context, playerID int64, token string) (Overview, error) {
	wallet, err := s.growth.GetWallet(ctx, playerID)
	if err != nil {
		return Overview{}, err
	}

	world, err := s.dungeons.GetWorldMap(ctx)
	if err != nil {
		return Overview{}, err
	}

	overview := Overview{
		PlayerID: playerID,
		Nickname: s.resolveNickname(ctx, playerID, token),
		Wallet:   wallet,
		Modules: ModulesOverview{
			MapLabel:     s.buildMapLabel(world),
			MapCityCount: len(world.Cities),
			Dungeon: DungeonSummary{
				Status: "idle",
			},
			Cultivation: CultivationSummary{
				State: "idle",
			},
		},
	}

	if run, err := s.dungeons.GetRun(ctx, playerID); err == nil {
		overview.Modules.Dungeon = DungeonSummary{
			Status:       run.Status,
			CurrentFloor: run.CurrentFloor,
			RemainDice:   run.RemainDice,
			DungeonID:    run.DungeonID,
		}
	}

	if status, err := s.dungeons.GetCultivationStatus(ctx, playerID); err == nil {
		overview.Modules.Cultivation = CultivationSummary{
			State:       status.State,
			SpiritPower: status.SpiritPower,
			Claimable:   !status.ClaimableAt.IsZero() && time.Now().After(status.ClaimableAt),
			ClaimableAt: status.ClaimableAt,
		}
	}

	return overview, nil
}

func (s *Service) resolveNickname(ctx context.Context, playerID int64, token string) string {
	token = strings.TrimSpace(token)
	if token != "" {
		if guest, err := s.accounts.GetByToken(ctx, token); err == nil && guest.Nickname != "" {
			return guest.Nickname
		}
	}
	return fmt.Sprintf("访客修士#%d", playerID)
}

func (s *Service) buildMapLabel(world dungeon.WorldMap) string {
	if len(world.Cities) == 0 {
		return world.Name
	}
	return fmt.Sprintf("%s · 已开放 %d 城", world.Name, len(world.Cities))
}
