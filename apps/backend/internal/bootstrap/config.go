package bootstrap

type Config struct {
	AppName  string
	HTTPPort int
}

func LoadConfig() Config {
	return Config{
		AppName:  "zhzw-api",
		HTTPPort: 8080,
	}
}
