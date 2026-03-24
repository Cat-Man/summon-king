package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Cat-Man/summon-king/apps/backend/internal/bootstrap"
)

func newServer() *http.Server {
	cfg := bootstrap.LoadConfig()
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: bootstrap.SetupRouter(),
	}
}

func run(start func(*http.Server) error) error {
	return start(newServer())
}

func main() {
	err := run(func(server *http.Server) error {
		return server.ListenAndServe()
	})
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
