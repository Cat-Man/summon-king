package bootstrap

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultAppName = "zhzw-api"
const defaultStorageDriver = "memory"

type Config struct {
	AppName       string
	HTTPPort      int
	StorageDriver string
	MySQLDSN      string
	RedisAddr     string
}

func LoadConfig() (Config, error) {
	httpPort := 8080
	if value := os.Getenv("HTTP_PORT"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 || parsed > 65535 {
			return Config{}, fmt.Errorf("invalid HTTP_PORT: %q", value)
		}
		httpPort = parsed
	}

	storageDriver := strings.ToLower(strings.TrimSpace(os.Getenv("STORAGE_DRIVER")))
	if storageDriver == "" {
		storageDriver = defaultStorageDriver
	}

	mysqlDSN := strings.TrimSpace(os.Getenv("MYSQL_DSN"))
	redisAddr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))

	switch storageDriver {
	case defaultStorageDriver:
	case "mysql":
		if mysqlDSN == "" {
			return Config{}, fmt.Errorf("MYSQL_DSN is required when STORAGE_DRIVER=mysql")
		}
	default:
		return Config{}, fmt.Errorf("invalid STORAGE_DRIVER: %q", storageDriver)
	}

	return Config{
		AppName:       defaultAppName,
		HTTPPort:      httpPort,
		StorageDriver: storageDriver,
		MySQLDSN:      mysqlDSN,
		RedisAddr:     redisAddr,
	}, nil
}
