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
