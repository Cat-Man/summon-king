package bootstrap

import (
	"fmt"
	"os"
	"strconv"
)

const defaultAppName = "zhzw-api"

type Config struct {
	AppName  string
	HTTPPort int
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

	return Config{
		AppName:  defaultAppName,
		HTTPPort: httpPort,
	}, nil
}
