package main

import (
	"fmt"
	"log"

	"github.com/Cat-Man/summon-king/apps/backend/internal/bootstrap"
)

func main() {
	cfg := bootstrap.LoadConfig()
	router := bootstrap.NewRouter()
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
