package main

import (
	"context"
	"fmt"
	"log"
	gohttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rafaeltg/goports/internal/adapters/handler/http"
	"github.com/rafaeltg/goports/internal/adapters/repository/memory"
	"github.com/rafaeltg/goports/internal/core/config"
	"github.com/rafaeltg/goports/internal/core/service"
	"github.com/rafaeltg/goports/pkg/logging"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Load configuration from env vars.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
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

	router := mux.NewRouter()
	srv := &gohttp.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           handlers.RecoveryHandler()(router),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-gCtx.Done()

		logger.InfoContext(ctx, "shutting down http server...")

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			logger.ErrorContext(ctx,
				"failed to shutdown http server",
				logging.Error(err),
			)
		}

		return err
	})

	g.Go(func() error {
		http.WithPortHandlers(
			router,
			portSvc,
			logger,
		)

		err := srv.ListenAndServe()
		if err != gohttp.ErrServerClosed {
			logger.ErrorContext(ctx,
				"failed to start http server",
				logging.Error(err),
			)
		}

		return err
	})

	logger.InfoContext(ctx, "application is up")

	_ = g.Wait()

	logger.InfoContext(ctx, "application gracefully stopped")
}
