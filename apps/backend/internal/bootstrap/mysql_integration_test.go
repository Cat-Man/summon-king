package bootstrap

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/home"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

func TestNewRouterWithConfig_UsesRealMySQLStorageAndPersistsAcrossRebuilds(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("MYSQL_SMOKE_DSN"))
	if dsn == "" {
		t.Skip("MYSQL_SMOKE_DSN is not set")
	}

	resetMySQLSmokeDatabase(t, dsn)
	t.Cleanup(func() {
		resetMySQLSmokeDatabase(t, dsn)
	})

	cfg := Config{
		AppName:       defaultAppName,
		HTTPPort:      8080,
		StorageDriver: "mysql",
		MySQLDSN:      dsn,
	}

	router, err := NewRouterWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected mysql router setup success, got %v", err)
	}

	guest := guestLogin(t, router, "云上旅人")
	postJSON(t, router, "/api/v1/pet/team/main", fmt.Sprintf(`{"player_id":%d,"pet_id":%d}`, guest.PlayerID, guest.PlayerID*10+2), nil)
	postJSON(t, router, "/api/v1/arena/challenge", fmt.Sprintf(`{"player_id":%d,"won":true}`, guest.PlayerID), nil)
	postForm(t, router, fmt.Sprintf("/api/v1/dungeon/enter?player_id=%d", guest.PlayerID), "dungeon_id=1")
	postJSON(t, router, "/api/v1/tower/pagoda/start", fmt.Sprintf(`{"player_id":%d}`, guest.PlayerID), nil)

	restarted, err := NewRouterWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected mysql router rebuild success, got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/home/overview?player_id=%d", guest.PlayerID), nil)
	req.Header.Set("Authorization", "Bearer "+guest.Token)
	resp := httptest.NewRecorder()
	restarted.ServeHTTP(resp, req)

	var homePayload struct {
		Code int           `json:"code"`
		Data home.Overview `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&homePayload); err != nil {
		t.Fatalf("expected overview JSON payload, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("expected home overview 200, got %d", resp.Code)
	}
	if homePayload.Code != 0 {
		t.Fatalf("expected home overview business code 0, got %d", homePayload.Code)
	}
	if homePayload.Data.Nickname != "云上旅人" {
		t.Fatalf("expected persisted nickname 云上旅人, got %s", homePayload.Data.Nickname)
	}
	if homePayload.Data.Modules.Pet.StarterPetName != "玄甲龟" {
		t.Fatalf("expected persisted starter pet 玄甲龟, got %s", homePayload.Data.Modules.Pet.StarterPetName)
	}
	if homePayload.Data.Modules.Arena.CurrentStreak != 1 {
		t.Fatalf("expected persisted arena streak 1, got %d", homePayload.Data.Modules.Arena.CurrentStreak)
	}
	if homePayload.Data.Wallet.SpiritPower != 118 {
		t.Fatalf("expected persisted spirit power 118, got %d", homePayload.Data.Wallet.SpiritPower)
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
	if dungeonPayload.Code != 0 {
		t.Fatalf("expected dungeon status business code 0, got %d", dungeonPayload.Code)
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
	if towerPayload.Code != 0 {
		t.Fatalf("expected tower status business code 0, got %d", towerPayload.Code)
	}
	if towerPayload.Data.CurrentFloor != 1 {
		t.Fatalf("expected persisted tower floor 1, got %d", towerPayload.Data.CurrentFloor)
	}
	if towerPayload.Data.RemainingChallenges != 4 {
		t.Fatalf("expected persisted tower remaining challenges 4, got %d", towerPayload.Data.RemainingChallenges)
	}
}

func resetMySQLSmokeDatabase(t *testing.T, dsn string) {
	t.Helper()

	store, err := mysqlstore.Open(dsn)
	if err != nil {
		t.Fatalf("expected mysql smoke store open success, got %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("expected mysql smoke store close success, got %v", err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("expected mysql smoke db open success, got %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("expected mysql smoke db close success, got %v", err)
		}
	}()

	statements := []string{
		`TRUNCATE TABLE module_states`,
		`TRUNCATE TABLE guest_accounts`,
		`ALTER TABLE guest_accounts AUTO_INCREMENT = 1001`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("expected mysql smoke reset statement success for %q, got %v", statement, err)
		}
	}
}
