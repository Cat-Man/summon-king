package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunUsesConfiguredHTTPServer(t *testing.T) {
	var got *http.Server

	err := run(func(server *http.Server) error {
		got = server
		return http.ErrServerClosed
	})
	if !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("expected http.ErrServerClosed, got %v", err)
	}
	if got == nil {
		t.Fatal("expected server to be created")
	}
	if got.Addr != ":8080" {
		t.Fatalf("expected addr :8080, got %s", got.Addr)
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	got.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected healthz 200, got %d", rec.Code)
	}
}
