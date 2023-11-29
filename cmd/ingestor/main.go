package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rafaeltg/goports/internal/adapters/handler/ingest"
	"github.com/rafaeltg/goports/internal/adapters/repository/memory"
	"github.com/rafaeltg/goports/internal/core/config"
	"github.com/rafaeltg/goports/internal/core/service"
	"github.com/rafaeltg/goports/pkg/logging"
)

type Config struct {
	config.Configuration
}

func main() {
	// Load configuration from env vars.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	if len(cfg.Ingestor.Filepath) == 0 {
		log.Fatalf("missing name of the file to process")
	}

	// Setup logger
	logger := logging.NewLogger(
		logging.WithLevel(cfg.LogLevel),
		logging.WithSource(!cfg.Environment.IsProduction()),
		logging.WithField("service", cfg.Application.Name),
		logging.WithField("version", cfg.Application.Version),
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer cancel()

	// Dependency injection
	memDB := memory.NewDatabase()
	portRepo := memory.NewPortRepository(memDB, logger)
	portSvc := service.NewPortService(portRepo, logger)

	portsIngestor := ingest.NewPortsIngestor(portSvc, logger)

	err = portsIngestor.Process(ctx, cfg.Ingestor.Filepath)
	if err != nil {
		logger.Error(
			"error importing ports data",
			logging.Error(err),
		)
	} else {
		logger.Info("done importing ports data")
	}
}
