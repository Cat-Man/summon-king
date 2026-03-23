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
	now            func() time.Time
}

type HomeIndex struct {
	PlayerID           int64                  `json:"player_id"`
	Nickname           string                 `json:"nickname"`
	Level              int                    `json:"level"`
	Coin               int64                  `json:"coin"`
	Diamond            int64                  `json:"diamond"`
	LastLoginAt        time.Time              `json:"last_login_at"`
	DailyTodos         []HomeTodo             `json:"daily_todos"`
	Resources          []HomeResource         `json:"resources"`
	ResourceIcons      []HomeResourceIcon     `json:"resource_icons"`
	Messages           []string               `json:"messages"`
	Entries            []string               `json:"entries"`
	ActivityEntry      HomeActivityEntry      `json:"activity_entry"`
	CultivationSummary HomeCultivationSummary `json:"cultivation_summary"`
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
	Action      string `json:"action"`
}

type HomeCultivationSummary struct {
	Status           string    `json:"status"`
	MapName          string    `json:"map_name"`
	StartedAt        time.Time `json:"started_at"`
	FinishedAt       time.Time `json:"finished_at"`
	RemainingSeconds int64     `json:"remaining_seconds"`
	RewardCoins      int64     `json:"reward_coins"`
	RewardPetExp     int64     `json:"reward_pet_exp"`
	Action           string    `json:"action"`
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

type cultivationRuntimeState struct {
	record    dungeon.CultivationRecord
	found     bool
	status    string
	mapName   string
	action    string
	remaining int64
	todoValue string
	summary   HomeCultivationSummary
}

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

func WithNow(now func() time.Time) ServiceOption {
	return func(s *Service) {
		s.now = now
	}
}

func NewService(repo Repository, options ...ServiceOption) *Service {
	svc := &Service{
		repo: repo,
		now:  time.Now,
	}
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

	giftState := commerce.GiftState{}
	vipState := commerce.VIPState{}
	if s.commerceReader != nil {
		giftState = s.commerceReader.GetGiftState(ctx, playerID)
		vipState = s.commerceReader.GetVIPState(ctx, playerID)
	}

	world := dungeon.WorldMap{}
	recommended := recommendedDungeon(player.Profile.Level)
	if s.dungeonReader != nil {
		if loadedWorld, worldErr := s.dungeonReader.GetWorldMap(ctx, playerID); worldErr == nil {
			world = loadedWorld
			recommended = recommendedDungeonFromWorld(loadedWorld, player.Profile.Level)
		}
	}

	cultivationState := buildIdleCultivationState(player.Profile.Level, world)
	if s.dungeonReader != nil {
		if record, found, cultivationErr := s.dungeonReader.GetLatestCultivation(ctx, playerID); cultivationErr == nil && found {
			cultivationState = buildCultivationState(record, world, s.now())
		}
	}

	dailySigninClaimed := giftState.SigninDays > 0
	unlockedCities := countUnlockedCities(world)
	vitalityValue := buildVitalityValue(player.Wallet.Vitality, cultivationState)
	reputationValue := buildReputationValue(player.Wallet.Reputation, player.Profile.Level, unlockedCities, giftState.SigninDays)

	return HomeIndex{
		PlayerID:           player.PlayerID,
		Nickname:           player.Profile.Nickname,
		Level:              player.Profile.Level,
		Coin:               player.Wallet.Coin,
		Diamond:            player.Wallet.Diamond,
		LastLoginAt:        player.DailyState.LastLoginAt,
		DailyTodos:         buildDailyTodos(dailySigninClaimed, cultivationState.todoValue, recommended),
		Resources:          buildResources(player.Profile.Level, power, player.Wallet.Coin, player.Wallet.Diamond, vitalityValue, reputationValue),
		ResourceIcons:      buildResourceIcons(player.Wallet.Coin, player.Wallet.Diamond),
		Messages:           buildMessages(recommended, vipState.DailyChestClaimed, cultivationState),
		Entries:            []string{"世界地图", "联盟", "幻兽", "背包", "竞技场", "庄园", "修行", "排行"},
		ActivityEntry:      buildActivityEntry(cultivationState, dailySigninClaimed, giftState.SigninDays, vipState.DailyChestClaimed, recommended),
		CultivationSummary: cultivationState.summary,
	}, nil
}

func buildDailyTodos(dailySigninClaimed bool, cultivationValue string, recommended string) []HomeTodo {
	signinValue := "今日未签"
	signinAction := "前往签到"
	if dailySigninClaimed {
		signinValue = "今日已签"
		signinAction = "明日再来"
	}

	return []HomeTodo{
		{Title: "签到状态", Value: signinValue, Action: signinAction},
		{Title: "当前修行状态", Value: cultivationValue, Action: "查看修行"},
		{Title: "当前推荐副本", Value: recommended, Action: "进入副本"},
	}
}

func buildResources(level int, power int64, coin int64, diamond int64, vitality string, reputation int64) []HomeResource {
	return []HomeResource{
		{Label: "等级", Value: fmt.Sprintf("Lv.%d", level)},
		{Label: "战力", Value: formatInt64(power)},
		{Label: "活力", Value: vitality},
		{Label: "铜钱", Value: formatInt64(coin)},
		{Label: "元宝", Value: formatInt64(diamond)},
		{Label: "声望", Value: formatInt64(reputation)},
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

func recommendedDungeonFromWorld(world dungeon.WorldMap, level int) string {
	for _, city := range world.Cities {
		if city.Unlocked {
			return fmt.Sprintf("%s · Boss 可挑战", city.CityName)
		}
	}
	return recommendedDungeon(level)
}

func buildMessages(recommended string, vipDailyChestClaimed bool, cultivation cultivationRuntimeState) []string {
	systemMsg := "系统消息：VIP 每日宝箱可领取"
	if vipDailyChestClaimed {
		systemMsg = "系统消息：VIP 每日宝箱已领取"
	}

	cultivationMsg := "修行消息：当前暂无进行中的修行"
	switch cultivation.status {
	case "running":
		cultivationMsg = fmt.Sprintf("修行消息：%s剩余 %s", cultivation.mapName, formatDurationCN(cultivation.remaining))
	case "claimable":
		cultivationMsg = fmt.Sprintf("修行消息：%s收益可领取", cultivation.mapName)
	case "claimed":
		cultivationMsg = fmt.Sprintf("修行消息：%s收益已领取", cultivation.mapName)
	}

	return []string{
		"世界消息：" + recommended,
		cultivationMsg,
		systemMsg,
	}
}

func buildActivityEntry(cultivation cultivationRuntimeState, dailySigninClaimed bool, signinDays int, vipDailyChestClaimed bool, recommended string) HomeActivityEntry {
	switch cultivation.status {
	case "running":
		return HomeActivityEntry{
			Title:       "修行收益",
			Description: fmt.Sprintf("%s修行进行中", cultivation.mapName),
			Note:        "剩余 " + formatDurationCN(cultivation.remaining),
			Action:      cultivation.action,
		}
	case "claimable":
		return HomeActivityEntry{
			Title:       "修行收益",
			Description: fmt.Sprintf("%s修行可领取", cultivation.mapName),
			Note:        "可立即领取",
			Action:      cultivation.action,
		}
	case "claimed":
		return HomeActivityEntry{
			Title:       "修行收益",
			Description: fmt.Sprintf("%s收益已领取", cultivation.mapName),
			Note:        "等待下一次修行",
			Action:      "前往修行",
		}
	}

	if !dailySigninClaimed {
		description := "今日签到待领取"
		if signinDays > 0 {
			description = fmt.Sprintf("今日签到待领取，已累计 %d 天", signinDays)
		}
		return HomeActivityEntry{
			Title:       "签到奖励",
			Description: description,
			Note:        "前往签到",
			Action:      "前往签到",
		}
	}

	if !vipDailyChestClaimed {
		return HomeActivityEntry{
			Title:       "VIP宝箱",
			Description: "每日宝箱待领取",
			Note:        "前往领取",
			Action:      "前往领取",
		}
	}

	return HomeActivityEntry{
		Title:       "推荐副本",
		Description: recommended,
		Note:        "进入副本",
		Action:      "进入副本",
	}
}

func buildIdleCultivationState(level int, world dungeon.WorldMap) cultivationRuntimeState {
	mapName := defaultCultivationMapName(level)
	for _, city := range world.Cities {
		if city.Unlocked {
			mapName = city.CityName
			break
		}
	}

	summary := HomeCultivationSummary{
		Status:           "idle",
		MapName:          mapName,
		RemainingSeconds: 0,
		RewardCoins:      0,
		RewardPetExp:     0,
		Action:           "前往修行",
	}

	return cultivationRuntimeState{
		status:    "idle",
		mapName:   mapName,
		action:    summary.Action,
		remaining: 0,
		todoValue: "尚未开始修行",
		summary:   summary,
	}
}

func buildCultivationState(record dungeon.CultivationRecord, world dungeon.WorldMap, now time.Time) cultivationRuntimeState {
	status := resolveCultivationStatus(record, now)
	mapName := resolveCultivationMapName(record.MapID, world)
	remaining := int64(0)
	action := "查看修行"
	todoValue := "尚未开始修行"

	switch status {
	case "running":
		remaining = maxInt64(int64(record.FinishedAt.Sub(now).Seconds()), 0)
		action = "查看修行"
		todoValue = "剩余 " + formatDurationCN(remaining)
	case "claimable":
		action = "领取收益"
		todoValue = "收益可领取"
	case "claimed":
		action = "已领取"
		todoValue = "收益已领取"
	}

	summary := HomeCultivationSummary{
		Status:           status,
		MapName:          mapName,
		StartedAt:        record.StartedAt,
		FinishedAt:       record.FinishedAt,
		RemainingSeconds: remaining,
		RewardCoins:      record.RewardCoins,
		RewardPetExp:     record.RewardPetExp,
		Action:           action,
	}

	return cultivationRuntimeState{
		record:    record,
		found:     true,
		status:    status,
		mapName:   mapName,
		action:    action,
		remaining: remaining,
		todoValue: todoValue,
		summary:   summary,
	}
}

func resolveCultivationStatus(record dungeon.CultivationRecord, now time.Time) string {
	if record.Status == "claimed" {
		return "claimed"
	}
	if !now.Before(record.FinishedAt) {
		return "claimable"
	}
	return "running"
}

func resolveCultivationMapName(mapID string, world dungeon.WorldMap) string {
	for _, city := range world.Cities {
		if city.CityID == mapID {
			return city.CityName
		}
	}
	return mapID
}

func buildVitalityValue(base int, cultivation cultivationRuntimeState) string {
	maxVitality := base
	if maxVitality <= 0 {
		maxVitality = 120
	}

	current := maxVitality
	if cultivation.status == "running" || cultivation.status == "claimable" {
		cost := cultivation.record.Hours * 10
		maxCost := maxVitality / 2
		if cost > maxCost {
			cost = maxCost
		}
		current -= cost
		if current < 0 {
			current = 0
		}
	}

	return fmt.Sprintf("%d / %d", current, maxVitality)
}

func buildReputationValue(base int64, level int, unlockedCities int, signinDays int) int64 {
	derived := int64(level*20 + unlockedCities*15 + signinDays*10)
	if derived > base {
		return derived
	}
	return base
}

func countUnlockedCities(world dungeon.WorldMap) int {
	total := 0
	for _, city := range world.Cities {
		if city.Unlocked {
			total++
		}
	}
	return total
}

func defaultCultivationMapName(level int) string {
	if level >= 10 {
		return "黑石山谷"
	}
	return "青木林地"
}

func formatDurationCN(seconds int64) string {
	if seconds <= 0 {
		return "0分"
	}

	totalMinutes := int((seconds + 59) / 60)
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	if hours == 0 {
		return fmt.Sprintf("%d分", totalMinutes)
	}
	if minutes == 0 {
		return fmt.Sprintf("%d小时", hours)
	}
	return fmt.Sprintf("%d小时%d分", hours, minutes)
}

func maxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
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
