package bootstrap

import (
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
