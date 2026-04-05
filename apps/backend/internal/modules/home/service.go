package home

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/ranking"
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

type petReader interface {
	GetBattleTeam(ctx context.Context, playerID int64) (pet.TeamSnapshot, error)
}

type arenaReader interface {
	GetDailyRecord(ctx context.Context, playerID int64) (arena.DailyRecord, error)
}

type rankingReader interface {
	GetLeaderboard(ctx context.Context, playerID int64, limit int) ([]ranking.LeaderboardEntry, error)
}

type Service struct {
	accounts accountReader
	dungeons dungeonReader
	growth   growth.Repository
	pets     petReader
	towers   towerReader
	arena    arenaReader
	ranking  rankingReader
}

func NewService(
	accounts accountReader,
	dungeons dungeonReader,
	growthRepo growth.Repository,
	pets petReader,
	towers towerReader,
	arena arenaReader,
	ranking rankingReader,
) *Service {
	return &Service{
		accounts: accounts,
		dungeons: dungeons,
		growth:   growthRepo,
		pets:     pets,
		towers:   towers,
		arena:    arena,
		ranking:  ranking,
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
			Pet: PetSummary{},
			Tower: TowerOverview{
				Pagoda: s.towers.GetStatus(ctx, playerID, "pagoda"),
				Spirit: s.towers.GetStatus(ctx, playerID, "spirit"),
			},
			Arena:   ArenaSummary{},
			Ranking: RankingSummary{},
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

	if s.pets != nil {
		if team, err := s.pets.GetBattleTeam(ctx, playerID); err == nil {
			overview.Modules.Pet = summarizePet(team)
		}
	}
	if s.arena != nil {
		if record, err := s.arena.GetDailyRecord(ctx, playerID); err == nil {
			overview.Modules.Arena = ArenaSummary{
				CurrentStreak: record.CurrentStreak,
				LastWin:       record.LastWin,
			}
		}
	}
	if s.ranking != nil {
		if board, err := s.ranking.GetLeaderboard(ctx, playerID, 10); err == nil {
			overview.Modules.Ranking = summarizeRanking(board)
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
	if s.arena != nil {
		return NextAction{
			Title:       "继续竞技场斗法",
			Description: "塔挑战已经打完，去竞技场冲连胜并抬升榜单。",
			Route:       "/arena",
			CTA:         "前往竞技场",
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

func summarizePet(team pet.TeamSnapshot) PetSummary {
	summary := PetSummary{
		TotalPower: team.TotalPower,
	}

	for _, battlePet := range team.Pets {
		if battlePet.IsActive {
			summary.ActiveCount++
		}
		if battlePet.Slot == 1 && battlePet.Name != "" {
			summary.StarterPetName = battlePet.Name
		}
	}

	if summary.StarterPetName == "" && len(team.Pets) > 0 {
		summary.StarterPetName = team.Pets[0].Name
	}

	return summary
}

func summarizeRanking(entries []ranking.LeaderboardEntry) RankingSummary {
	summary := RankingSummary{}
	if len(entries) > 0 {
		summary.TopName = entries[0].Name
		summary.TopScore = entries[0].Score
	}

	for _, entry := range entries {
		if entry.IsSelf {
			summary.SelfRank = entry.Rank
			summary.SelfScore = entry.Score
			break
		}
	}

	return summary
}
