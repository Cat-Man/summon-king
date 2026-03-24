package bootstrap

import "testing"

func TestLoadConfig_DefaultEnv(t *testing.T) {
	cfg := LoadConfig()
	if cfg.AppName == "" {
		t.Fatal("app name should not be empty")
	}
	if cfg.HTTPPort == 0 {
		t.Fatal("http port should not be zero")
	}
}

func TestLoadConfig_UsesHTTPPortEnv(t *testing.T) {
	t.Setenv("HTTP_PORT", "18080")

	cfg := LoadConfig()
	if cfg.HTTPPort != 18080 {
		t.Fatalf("expected http port 18080, got %d", cfg.HTTPPort)
	}
}
