package ranking

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
)

type accountReader interface {
	GetByPlayerID(ctx context.Context, playerID int64) (account.GuestLoginResponse, error)
}

type growthReader interface {
	GetWallet(ctx context.Context, playerID int64) (growth.Wallet, error)
}

type dungeonReader interface {
	GetRun(ctx context.Context, playerID int64) (dungeon.DungeonRun, error)
}

type Service struct {
	accounts accountReader
	growth   growthReader
	dungeons dungeonReader
}

func NewService(accounts accountReader, growth growthReader, dungeons dungeonReader) *Service {
	return &Service{
		accounts: accounts,
		growth:   growth,
		dungeons: dungeons,
	}
}

func (s *Service) GetLeaderboard(ctx context.Context, playerID int64, limit int) ([]LeaderboardEntry, error) {
	now := time.Now()
	list := s.seedEntries(now)

	if playerID > 0 {
		list = append(list, s.buildCurrentPlayerEntry(ctx, playerID, now))
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Score == list[j].Score {
			return list[i].Updated > list[j].Updated
		}
		return list[i].Score > list[j].Score
	})

	for i := range list {
		list[i].Rank = i + 1
	}
	if limit > 0 && limit < len(list) {
		list = list[:limit]
	}
	return list, nil
}

func (s *Service) seedEntries(now time.Time) []LeaderboardEntry {
	return []LeaderboardEntry{
		{PlayerID: 9001, Name: "星痕丶苍穹", Score: 1480, Updated: now.Add(-2 * time.Minute).UnixMilli()},
		{PlayerID: 9002, Name: "雨落荒原", Score: 1320, Updated: now.Add(-5 * time.Minute).UnixMilli()},
		{PlayerID: 9003, Name: "剑舞倾城", Score: 1200, Updated: now.Add(-8 * time.Minute).UnixMilli()},
	}
}

func (s *Service) buildCurrentPlayerEntry(ctx context.Context, playerID int64, now time.Time) LeaderboardEntry {
	name := fmt.Sprintf("访客修士#%d", playerID)
	if guest, err := s.accounts.GetByPlayerID(ctx, playerID); err == nil && guest.Nickname != "" {
		name = guest.Nickname
	}

	var score int64 = 200
	if wallet, err := s.growth.GetWallet(ctx, playerID); err == nil {
		score += wallet.SpiritPower
		score += int64(wallet.BoneLevel * 120)
		score += int64(wallet.SoulPieces * 40)
		score += int64(wallet.ManorPlots * 50)
	}
	if run, err := s.dungeons.GetRun(ctx, playerID); err == nil {
		score += int64(run.CurrentFloor * 180)
		score += int64(run.RemainDice * 15)
	}

	return LeaderboardEntry{
		PlayerID: playerID,
		Name:     name,
		Score:    score,
		Updated:  now.UnixMilli(),
		IsSelf:   true,
	}
}
