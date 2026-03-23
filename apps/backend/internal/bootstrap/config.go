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
	if rawPort := os.Getenv("HTTP_PORT"); rawPort != "" {
		if parsedPort, err := strconv.Atoi(rawPort); err == nil && parsedPort > 0 {
			httpPort = parsedPort
		}
	}

	return Config{
		AppName:  "zhzw-api",
		HTTPPort: httpPort,
	}
}
