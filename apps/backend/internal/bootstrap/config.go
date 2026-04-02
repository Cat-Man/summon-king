package bootstrap

import (
	"os"
	"strconv"
)

type Config struct {
	AppName  string
	HTTPPort int
}

func LoadConfig() Config {
	httpPort := 8080
	if value := os.Getenv("HTTP_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			httpPort = parsed
		}
	}

	return Config{
		AppName:  "zhzw-api",
		HTTPPort: httpPort,
	}
}
