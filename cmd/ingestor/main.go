package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rafaeltg/goports/internal/adapters/client/http"
	"github.com/rafaeltg/goports/internal/adapters/handler/ingest"
	"github.com/rafaeltg/goports/internal/core/config"
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
		syscall.SIGKILL,
	)
	defer cancel()

	// Dependency injection
	httpClient := http.NewCient(cfg.Server.Host())
	portClient := http.NewPortClient(httpClient, logger)
	portIngestor := ingest.NewPortIngestor(
		portClient,
		logger,
		ingest.WithBatchSize(cfg.Ingestor.BatchSize),
	)

	logger.Info("running ingestor",
		slog.Any("config", cfg),
	)

	err = portIngestor.Process(ctx, cfg.Ingestor.Filepath)
	if err != nil {
		logger.Error(
			"error importing ports data",
			logging.Error(err),
		)
	} else {
		logger.Info("done importing ports data")
	}
}
