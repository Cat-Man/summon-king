package bootstrap

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_ModuleRoots(t *testing.T) {
	r := NewRouter()

	tests := []struct {
		name string
		path string
	}{
		{name: "healthz", path: "/healthz"},
		{name: "arena root", path: "/api/v1/arena"},
		{name: "dungeon root", path: "/api/v1/dungeon"},
		{name: "growth root", path: "/api/v1/growth"},
		{name: "pet root", path: "/api/v1/pet"},
		{name: "tower root", path: "/api/v1/tower"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Fatalf("expected 200 for %s, got %d", tt.path, resp.Code)
			}

			var payload struct {
				Code    int                    `json:"code"`
				TraceID string                 `json:"trace_id"`
				Data    map[string]interface{} `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				t.Fatalf("failed to decode payload for %s: %v", tt.path, err)
			}
			if payload.Code != 0 {
				t.Fatalf("expected code 0 for %s, got %d", tt.path, payload.Code)
			}
			if payload.TraceID == "" {
				t.Fatalf("expected trace_id for %s", tt.path)
			}
			if tt.path == "/healthz" {
				if payload.Data["status"] != "ok" {
					t.Fatalf("expected healthz status ok, got %s", payload.Data["status"])
				}
				if payload.Data["app"] != defaultAppName {
					t.Fatalf("expected healthz app %s, got %s", defaultAppName, payload.Data["app"])
				}
			} else {
				if payload.Data["status"] != "ok" {
					t.Fatalf("expected module %s status ok, got %s", tt.name, payload.Data["status"])
				}
				if payload.Data["module"] != tt.name[:len(tt.name)-5] { // e.g. "arena root"
					module := tt.path[len("/api/v1/"):]
					if payload.Data["module"] != module {
						t.Fatalf("expected module %s, got %s", module, payload.Data["module"])
					}
				}
			}
		})
	}
}

func TestRouter_PetMainSwitchAffectsArenaBattlePower(t *testing.T) {
	r := NewRouter()

	switchReq := httptest.NewRequest(http.MethodPost, "/api/v1/pet/team/main", bytes.NewBufferString(`{"player_id":1001,"pet_id":10012}`))
	switchReq.Header.Set("Content-Type", "application/json")
	switchResp := httptest.NewRecorder()
	r.ServeHTTP(switchResp, switchReq)

	if switchResp.Code != http.StatusOK {
		t.Fatalf("expected pet switch 200, got %d", switchResp.Code)
	}

	arenaReq := httptest.NewRequest(http.MethodPost, "/api/v1/arena/challenge", bytes.NewBufferString(`{"player_id":1001,"won":true}`))
	arenaReq.Header.Set("Content-Type", "application/json")
	arenaResp := httptest.NewRecorder()
	r.ServeHTTP(arenaResp, arenaReq)

	if arenaResp.Code != http.StatusOK {
		t.Fatalf("expected arena challenge 200, got %d", arenaResp.Code)
	}

	var payload struct {
		Code int `json:"code"`
		Data struct {
			Battle struct {
				AttackerPower int64 `json:"attacker_power"`
			} `json:"battle"`
		} `json:"data"`
	}
	if err := json.NewDecoder(arenaResp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if payload.Code != 0 {
		t.Fatalf("expected business code 0, got %d", payload.Code)
	}
	if payload.Data.Battle.AttackerPower != 156 {
		t.Fatalf("expected attacker power 156 after switching main pet, got %d", payload.Data.Battle.AttackerPower)
	}
}

func TestRouter_SavedPetTeamAffectsArenaBattlePower(t *testing.T) {
	r := NewRouter()

	saveReq := httptest.NewRequest(http.MethodPost, "/api/v1/pet/team/save", bytes.NewBufferString(`{"player_id":1001,"pet_ids":[10011,10012]}`))
	saveReq.Header.Set("Content-Type", "application/json")
	saveResp := httptest.NewRecorder()
	r.ServeHTTP(saveResp, saveReq)

	if saveResp.Code != http.StatusOK {
		t.Fatalf("expected team save 200, got %d", saveResp.Code)
	}

	arenaReq := httptest.NewRequest(http.MethodPost, "/api/v1/arena/challenge", bytes.NewBufferString(`{"player_id":1001,"won":true}`))
	arenaReq.Header.Set("Content-Type", "application/json")
	arenaResp := httptest.NewRecorder()
	r.ServeHTTP(arenaResp, arenaReq)

	if arenaResp.Code != http.StatusOK {
		t.Fatalf("expected arena challenge 200, got %d", arenaResp.Code)
	}

	var payload struct {
		Code int `json:"code"`
		Data struct {
			Battle struct {
				AttackerPower int64 `json:"attacker_power"`
			} `json:"battle"`
		} `json:"data"`
	}
	if err := json.NewDecoder(arenaResp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected JSON payload, got %v", err)
	}
	if payload.Code != 0 {
		t.Fatalf("expected business code 0, got %d", payload.Code)
	}
	if payload.Data.Battle.AttackerPower != 276 {
		t.Fatalf("expected attacker power 276 after saving team, got %d", payload.Data.Battle.AttackerPower)
	}
}
