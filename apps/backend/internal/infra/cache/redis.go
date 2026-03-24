package cache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewRedis(cfg Config) (*redis.Client, error) {
	if cfg.Addr == "" {
		return nil, errors.New("redis addr is required")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client, nil
}
