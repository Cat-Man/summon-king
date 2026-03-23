package player

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/commerce"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
)

type Service struct {
	repo           Repository
	petReader      PetReader
	commerceReader CommerceReader
	dungeonReader  DungeonReader
}

type HomeIndex struct {
	PlayerID      int64              `json:"player_id"`
	Nickname      string             `json:"nickname"`
	Level         int                `json:"level"`
	Coin          int64              `json:"coin"`
	Diamond       int64              `json:"diamond"`
	LastLoginAt   time.Time          `json:"last_login_at"`
	DailyTodos    []HomeTodo         `json:"daily_todos"`
	Resources     []HomeResource     `json:"resources"`
	ResourceIcons []HomeResourceIcon `json:"resource_icons"`
	Messages      []string           `json:"messages"`
	Entries       []string           `json:"entries"`
	ActivityEntry HomeActivityEntry  `json:"activity_entry"`
}

type HomeTodo struct {
	Title  string `json:"title"`
	Value  string `json:"value"`
	Action string `json:"action"`
}

type HomeResource struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type HomeResourceIcon struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type HomeActivityEntry struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Note        string `json:"note"`
}

type PetReader interface {
	GetPlayerPets(ctx context.Context, playerID int64) ([]pet.PlayerPet, error)
}

type CommerceReader interface {
	GetVIPState(ctx context.Context, playerID int64) commerce.VIPState
	GetGiftState(ctx context.Context, playerID int64) commerce.GiftState
}

type DungeonReader interface {
	GetWorldMap(ctx context.Context, playerID int64) (dungeon.WorldMap, error)
	GetLatestCultivation(ctx context.Context, playerID int64) (dungeon.CultivationRecord, bool, error)
}

type ServiceOption func(*Service)

func WithPetReader(reader PetReader) ServiceOption {
	return func(s *Service) {
		s.petReader = reader
	}
}

func WithCommerceReader(reader CommerceReader) ServiceOption {
	return func(s *Service) {
		s.commerceReader = reader
	}
}

func WithDungeonReader(reader DungeonReader) ServiceOption {
	return func(s *Service) {
		s.dungeonReader = reader
	}
}

func NewService(repo Repository, options ...ServiceOption) *Service {
	svc := &Service{repo: repo}
	for _, option := range options {
		option(svc)
	}
	return svc
}

func (s *Service) GetHomeIndex(ctx context.Context, playerID int64) (HomeIndex, error) {
	player, err := s.repo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return HomeIndex{}, err
	}

	power := int64(1000 + player.Profile.Level*80)
	if s.petReader != nil {
		if pets, petErr := s.petReader.GetPlayerPets(ctx, playerID); petErr == nil {
			power = 0
			for _, item := range pets {
				power += int64(item.Power)
			}
		}
	}

	dailySigninClaimed := false
	vipDailyChestClaimed := false
	if s.commerceReader != nil {
		giftState := s.commerceReader.GetGiftState(ctx, playerID)
		dailySigninClaimed = giftState.SigninDays > 0
		vipState := s.commerceReader.GetVIPState(ctx, playerID)
		vipDailyChestClaimed = vipState.DailyChestClaimed
	}

	recommended := recommendedDungeon(player.Profile.Level)
	currentCultivationText := "尚未开始修行"
	if s.dungeonReader != nil {
		if world, worldErr := s.dungeonReader.GetWorldMap(ctx, playerID); worldErr == nil {
			for _, city := range world.Cities {
				if city.Unlocked {
					recommended = fmt.Sprintf("%s · Boss 可挑战", city.CityName)
					break
				}
			}
		}
		if cultivation, found, cultivationErr := s.dungeonReader.GetLatestCultivation(ctx, playerID); cultivationErr == nil && found {
			if cultivation.Status == "running" {
				currentCultivationText = "修行进行中"
			} else {
				currentCultivationText = "修行可领取"
			}
		}
	}

	return HomeIndex{
		PlayerID:      player.PlayerID,
		Nickname:      player.Profile.Nickname,
		Level:         player.Profile.Level,
		Coin:          player.Wallet.Coin,
		Diamond:       player.Wallet.Diamond,
		LastLoginAt:   player.DailyState.LastLoginAt,
		DailyTodos:    buildDailyTodos(dailySigninClaimed, currentCultivationText, recommended),
		Resources:     buildResources(player.Profile.Level, power, player.Wallet.Coin, player.Wallet.Diamond),
		ResourceIcons: buildResourceIcons(player.Wallet.Coin, player.Wallet.Diamond),
		Messages:      buildMessages(recommended, vipDailyChestClaimed),
		Entries:       []string{"世界地图", "联盟", "幻兽", "背包", "竞技场", "庄园", "修行", "排行"},
		ActivityEntry: HomeActivityEntry{
			Title:       "今日活动",
			Description: "夺宝双倍",
			Note:        "21:00 开始",
		},
	}, nil
}

func buildDailyTodos(dailySigninClaimed bool, currentCultivationText string, recommended string) []HomeTodo {
	signinValue := "今日未签"
	signinAction := "前往签到"
	if dailySigninClaimed {
		signinValue = "今日已签"
		signinAction = "明日再来"
	}

	return []HomeTodo{
		{Title: "签到状态", Value: signinValue, Action: signinAction},
		{Title: "当前修行状态", Value: currentCultivationText, Action: "查看修行"},
		{Title: "当前推荐副本", Value: recommended, Action: "进入副本"},
	}
}

func buildResources(level int, power int64, coin int64, diamond int64) []HomeResource {
	return []HomeResource{
		{Label: "等级", Value: fmt.Sprintf("Lv.%d", level)},
		{Label: "战力", Value: formatInt64(power)},
		{Label: "活力", Value: "120 / 120"},
		{Label: "铜钱", Value: formatInt64(coin)},
		{Label: "元宝", Value: formatInt64(diamond)},
		{Label: "声望", Value: "0"},
	}
}

func buildResourceIcons(coin int64, diamond int64) []HomeResourceIcon {
	return []HomeResourceIcon{
		{Label: "铜钱", Value: formatInt64(coin)},
		{Label: "元宝", Value: formatInt64(diamond)},
	}
}

func recommendedDungeon(level int) string {
	if level >= 10 {
		return "黑石山谷 · Boss 可挑战"
	}
	return "青木林地 · 可挑战"
}

func buildMessages(recommended string, vipDailyChestClaimed bool) []string {
	systemMsg := "系统消息：VIP 每日宝箱可领取"
	if vipDailyChestClaimed {
		systemMsg = "系统消息：VIP 每日宝箱已领取"
	}

	return []string{
		"世界消息：" + recommended,
		"联盟消息：今晚 20:00 盟战锁定名单",
		systemMsg,
	}
}

func formatInt64(value int64) string {
	negative := value < 0
	if negative {
		value = -value
	}

	raw := strconv.FormatInt(value, 10)
	if len(raw) <= 3 {
		if negative {
			return "-" + raw
		}
		return raw
	}

	formatted := raw[:len(raw)%3]
	if formatted == "" {
		formatted = raw[:3]
		raw = raw[3:]
	} else {
		raw = raw[len(formatted):]
	}

	for len(raw) > 0 {
		formatted += "," + raw[:3]
		raw = raw[3:]
	}

	if negative {
		return "-" + formatted
	}
	return formatted
}
