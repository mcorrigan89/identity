package main

import (
	"os"
	"sync"

	"github.com/mcorrigan89/identity/internal/config"
	"github.com/mcorrigan89/identity/internal/repositories"
	"github.com/mcorrigan89/identity/internal/services"
	"github.com/rs/zerolog"
)

type application struct {
	config   config.Config
	wg       sync.WaitGroup
	logger   *zerolog.Logger
	services *services.Services
}

func main() {

	cfg := config.Config{}
	config.LoadConfig(&cfg)

	logger := getLogger()

	db, err := openDBPool(&cfg)
	if err != nil {
		logger.Err(err).Msg("Failed to open database connection")
		os.Exit(1)
	}
	defer db.Close()

	wg := sync.WaitGroup{}

	repositories := repositories.NewRepositories(db, &logger, &wg)
	services := services.NewServices(&repositories, &cfg, &logger, &wg)

	app := &application{
		config:   cfg,
		logger:   &logger,
		services: &services,
	}

	err = app.serve()
	if err != nil {
		logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
}
