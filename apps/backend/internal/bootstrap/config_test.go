package bootstrap

import "testing"

func TestLoadConfig_DefaultEnv(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AppName == "" {
		t.Fatal("app name should not be empty")
	}
	if cfg.HTTPPort != 8080 {
		t.Fatalf("default HTTP port should be 8080, got %d", cfg.HTTPPort)
	}
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "-1")
	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected error for invalid port")
	}
}
