package main

import (
	"fmt"
	"log"

	"github.com/Cat-Man/summon-king/apps/backend/internal/bootstrap"
)

func main() {
	cfg, err := bootstrap.LoadConfig()
	if err != nil {
		log.Fatalf("config validation failed: %v", err)
	}
	router, err := bootstrap.NewRouterWithConfig(cfg)
	if err != nil {
		log.Fatalf("router setup failed: %v", err)
	}
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
