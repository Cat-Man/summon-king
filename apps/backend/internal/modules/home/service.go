package home

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
)

type accountReader interface {
	GetByToken(ctx context.Context, token string) (account.GuestLoginResponse, error)
}

type dungeonReader interface {
	GetWorldMap(ctx context.Context) (dungeon.WorldMap, error)
	GetRun(ctx context.Context, playerID int64) (dungeon.DungeonRun, error)
	GetCultivationStatus(ctx context.Context, playerID int64) (dungeon.CultivationStatus, error)
}

type towerReader interface {
	GetStatus(ctx context.Context, playerID int64, tower string) tower.TowerStatus
}

type Service struct {
	accounts accountReader
	dungeons dungeonReader
	growth   growth.Repository
	towers   towerReader
}

func NewService(accounts accountReader, dungeons dungeonReader, growthRepo growth.Repository, towers towerReader) *Service {
	return &Service{
		accounts: accounts,
		dungeons: dungeons,
		growth:   growthRepo,
		towers:   towers,
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
			Tower: TowerOverview{
				Pagoda: s.towers.GetStatus(ctx, playerID, "pagoda"),
				Spirit: s.towers.GetStatus(ctx, playerID, "spirit"),
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

	overview.NextAction = s.determineNextAction(overview.Modules)
	return overview, nil
}

func (s *Service) determineNextAction(modules ModulesOverview) NextAction {
	if modules.Cultivation.Claimable {
		return NextAction{
			Title:       "领取修行收益",
			Description: "修行已完成，先把灵力写回钱包。",
			Route:       "/cultivation",
			CTA:         "立即领取",
		}
	}

	if modules.Dungeon.Status == "ongoing" || modules.Dungeon.Status == "boss" {
		return NextAction{
			Title:       "继续地下城挑战",
			Description: "还有未完成的副本进度，先把它推到下一层。",
			Route:       "/dungeon",
			CTA:         "继续挑战",
		}
	}

	if modules.Tower.Pagoda.RemainingChallenges > 0 {
		return NextAction{
			Title:       "通天塔挑战",
			Description: "今日还有通天塔挑战次数，用战骨提升成长。",
			Route:       "/tower/pagoda",
			CTA:         "前往塔楼",
		}
	}

	if modules.Tower.Spirit.RemainingChallenges > 0 {
		return NextAction{
			Title:       "战灵塔试炼",
			Description: "灵力与魔魂可在战灵塔获得，先去试炼。",
			Route:       "/tower/spirit",
			CTA:         "前往战灵塔",
		}
	}

	return NextAction{
		Title:       "探索世界地图",
		Description: "地图已开放新城，继续迈向下一片大地。",
		Route:       "/map",
		CTA:         "前往地图",
	}
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
