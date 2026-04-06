package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/home"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

func TestNewRouterWithConfig_UsesMySQLStorageAndPersistsHomeArenaOverviewAcrossRebuilds(t *testing.T) {
	store := newFakeMySQLStore()
	restore := stubMySQLStoreOpener(t, store)
	defer restore()

	cfg := Config{
		AppName:       defaultAppName,
		HTTPPort:      8080,
		StorageDriver: "mysql",
		MySQLDSN:      "fake",
	}

	router, err := NewRouterWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected mysql router setup success, got %v", err)
	}

	guest := guestLogin(t, router, "云上旅人")
	postJSON(t, router, "/api/v1/pet/team/main", fmt.Sprintf(`{"player_id":%d,"pet_id":%d}`, guest.PlayerID, guest.PlayerID*10+2), nil)
	postJSON(t, router, "/api/v1/arena/challenge", fmt.Sprintf(`{"player_id":%d,"won":true}`, guest.PlayerID), nil)

	restarted, err := NewRouterWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected mysql router rebuild success, got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/home/overview?player_id=%d", guest.PlayerID), nil)
	req.Header.Set("Authorization", "Bearer "+guest.Token)
	resp := httptest.NewRecorder()
	restarted.ServeHTTP(resp, req)

	var payload struct {
		Code int           `json:"code"`
		Data home.Overview `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected overview JSON payload, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if payload.Code != 0 {
		t.Fatalf("expected business code 0, got %d", payload.Code)
	}
	if payload.Data.Nickname != "云上旅人" {
		t.Fatalf("expected persisted nickname 云上旅人, got %s", payload.Data.Nickname)
	}
	if payload.Data.Modules.Pet.StarterPetName != "玄甲龟" {
		t.Fatalf("expected persisted starter pet 玄甲龟, got %s", payload.Data.Modules.Pet.StarterPetName)
	}
	if payload.Data.Modules.Arena.CurrentStreak != 1 {
		t.Fatalf("expected persisted arena streak 1, got %d", payload.Data.Modules.Arena.CurrentStreak)
	}
	if payload.Data.Wallet.SpiritPower != 118 {
		t.Fatalf("expected persisted spirit power 118, got %d", payload.Data.Wallet.SpiritPower)
	}
}

func TestNewRouterWithConfig_UsesMySQLStorageAndPersistsDungeonTowerStateAcrossRebuilds(t *testing.T) {
	store := newFakeMySQLStore()
	restore := stubMySQLStoreOpener(t, store)
	defer restore()

	cfg := Config{
		AppName:       defaultAppName,
		HTTPPort:      8080,
		StorageDriver: "mysql",
		MySQLDSN:      "fake",
	}

	router, err := NewRouterWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected mysql router setup success, got %v", err)
	}

	guest := guestLogin(t, router, "洞府守望")
	postForm(t, router, fmt.Sprintf("/api/v1/dungeon/enter?player_id=%d", guest.PlayerID), "dungeon_id=1")
	postJSON(t, router, "/api/v1/tower/pagoda/start", fmt.Sprintf(`{"player_id":%d}`, guest.PlayerID), nil)

	restarted, err := NewRouterWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected mysql router rebuild success, got %v", err)
	}

	dungeonReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/dungeon/status?player_id=%d", guest.PlayerID), nil)
	dungeonResp := httptest.NewRecorder()
	restarted.ServeHTTP(dungeonResp, dungeonReq)

	var dungeonPayload struct {
		Code int `json:"code"`
		Data struct {
			DungeonID    int64 `json:"dungeon_id"`
			CurrentFloor int   `json:"current_floor"`
			RemainDice   int   `json:"remain_dice"`
		} `json:"data"`
	}
	if err := json.NewDecoder(dungeonResp.Body).Decode(&dungeonPayload); err != nil {
		t.Fatalf("expected dungeon JSON payload, got %v", err)
	}
	if dungeonResp.Code != http.StatusOK {
		t.Fatalf("expected dungeon status 200, got %d", dungeonResp.Code)
	}
	if dungeonPayload.Data.DungeonID != 1 || dungeonPayload.Data.CurrentFloor != 1 || dungeonPayload.Data.RemainDice != 15 {
		t.Fatalf("expected persisted dungeon run, got %+v", dungeonPayload.Data)
	}

	towerReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tower/pagoda/status?player_id=%d", guest.PlayerID), nil)
	towerResp := httptest.NewRecorder()
	restarted.ServeHTTP(towerResp, towerReq)

	var towerPayload struct {
		Code int               `json:"code"`
		Data tower.TowerStatus `json:"data"`
	}
	if err := json.NewDecoder(towerResp.Body).Decode(&towerPayload); err != nil {
		t.Fatalf("expected tower JSON payload, got %v", err)
	}
	if towerResp.Code != http.StatusOK {
		t.Fatalf("expected tower status 200, got %d", towerResp.Code)
	}
	if towerPayload.Data.CurrentFloor != 1 {
		t.Fatalf("expected persisted tower floor 1, got %d", towerPayload.Data.CurrentFloor)
	}
	if towerPayload.Data.RemainingChallenges != 4 {
		t.Fatalf("expected persisted tower remaining challenges 4, got %d", towerPayload.Data.RemainingChallenges)
	}
}

type fakeMySQLStore struct {
	mu              sync.Mutex
	nextPlayerID    int64
	accountsByToken map[string]mysqlstore.GuestAccountRecord
	accountsByID    map[int64]mysqlstore.GuestAccountRecord
	moduleStates    map[string][]byte
}

func newFakeMySQLStore() *fakeMySQLStore {
	return &fakeMySQLStore{
		nextPlayerID:    1000,
		accountsByToken: make(map[string]mysqlstore.GuestAccountRecord),
		accountsByID:    make(map[int64]mysqlstore.GuestAccountRecord),
		moduleStates:    make(map[string][]byte),
	}
}

func (s *fakeMySQLStore) CreateGuestAccount(_ context.Context, nickname, token string) (mysqlstore.GuestAccountRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextPlayerID++
	record := mysqlstore.GuestAccountRecord{
		PlayerID: s.nextPlayerID,
		Token:    token,
		Nickname: nickname,
	}
	s.accountsByToken[token] = record
	s.accountsByID[record.PlayerID] = record
	return record, nil
}

func (s *fakeMySQLStore) GetGuestAccountByToken(_ context.Context, token string) (mysqlstore.GuestAccountRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.accountsByToken[token]
	if !ok {
		return mysqlstore.GuestAccountRecord{}, errors.New("guest account not found")
	}
	return record, nil
}

func (s *fakeMySQLStore) GetGuestAccountByPlayerID(_ context.Context, playerID int64) (mysqlstore.GuestAccountRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.accountsByID[playerID]
	if !ok {
		return mysqlstore.GuestAccountRecord{}, errors.New("guest account not found")
	}
	return record, nil
}

func (s *fakeMySQLStore) LoadModuleState(_ context.Context, playerID int64, module string, target any) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload, ok := s.moduleStates[moduleStateKey(playerID, module)]
	if !ok {
		return false, nil
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return false, err
	}
	return true, nil
}

func (s *fakeMySQLStore) SaveModuleState(_ context.Context, playerID int64, module string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.moduleStates[moduleStateKey(playerID, module)] = payload
	return nil
}

func stubMySQLStoreOpener(t *testing.T, store mysqlStore) func() {
	t.Helper()

	previous := openMySQLStore
	openMySQLStore = func(_ Config) (mysqlStore, error) {
		return store, nil
	}
	return func() {
		openMySQLStore = previous
	}
}

func guestLogin(t *testing.T, router http.Handler, nickname string) mysqlstore.GuestAccountRecord {
	t.Helper()

	resp := postJSON(t, router, "/api/v1/auth/guest-login", fmt.Sprintf(`{"nickname":%q}`, nickname), nil)

	var payload struct {
		Code int                           `json:"code"`
		Data mysqlstore.GuestAccountRecord `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected guest login JSON payload, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("expected guest login 200, got %d", resp.Code)
	}
	if payload.Code != 0 {
		t.Fatalf("expected guest login business code 0, got %d", payload.Code)
	}
	return payload.Data
}

func postJSON(t *testing.T, router http.Handler, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func postForm(t *testing.T, router http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func moduleStateKey(playerID int64, module string) string {
	return fmt.Sprintf("%d:%s", playerID, module)
}
