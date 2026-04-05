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

func TestLoadConfig_DefaultStorageDriver(t *testing.T) {
	t.Setenv("STORAGE_DRIVER", "")
	t.Setenv("MYSQL_DSN", "")
	t.Setenv("REDIS_ADDR", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StorageDriver != "memory" {
		t.Fatalf("expected default storage driver memory, got %s", cfg.StorageDriver)
	}
	if cfg.MySQLDSN != "" {
		t.Fatalf("expected empty mysql dsn, got %s", cfg.MySQLDSN)
	}
	if cfg.RedisAddr != "" {
		t.Fatalf("expected empty redis addr, got %s", cfg.RedisAddr)
	}
}

func TestLoadConfig_InvalidStorageDriver(t *testing.T) {
	t.Setenv("STORAGE_DRIVER", "mongo")
	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected error for invalid storage driver")
	}
}

func TestLoadConfig_MySQLRequiresDSN(t *testing.T) {
	t.Setenv("STORAGE_DRIVER", "mysql")
	t.Setenv("MYSQL_DSN", "")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected error when mysql storage is configured without MYSQL_DSN")
	}
}
