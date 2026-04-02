package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Healthz(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
