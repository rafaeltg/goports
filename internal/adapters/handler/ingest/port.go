package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/port"
	"github.com/rafaeltg/goports/pkg/logging"
)

const batchSize int = 10

type PortsIngestor struct {
	portSvc port.PortService
	logger  *slog.Logger
}

func NewPortsIngestor(portSvc port.PortService, logger *slog.Logger) *PortsIngestor {
	return &PortsIngestor{
		portSvc: portSvc,
		logger:  logger,
	}
}

func (i *PortsIngestor) Process(ctx context.Context, filepath string) error {
	l := i.logger.With(
		slog.String("filepath", filepath),
	)

	l.InfoContext(ctx, "processing file")

	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			l.Error(
				"error on closing file",
				logging.Error(err),
			)
		}
	}(f)

	if _, err = f.Stat(); err != nil {
		return fmt.Errorf("finvalid file: %w", err)
	}

	dec := json.NewDecoder(f)

	// read opening JSON delimiter
	token, err := dec.Token()
	if err != nil {
		return fmt.Errorf("failed to read opening delimiter: %w", err)
	}

	if token != json.Delim('{') {
		return fmt.Errorf("unexpected token encountered on reading opening delimiterr: %s", token)
	}

	wg := sync.WaitGroup{}
	errCh := make(chan error)
	batch := make(domain.Ports, 0, batchSize)

	done := false
	for !done && dec.More() {
		select {
		case <-ctx.Done():
			done = true
		case err = <-errCh:
			done = true
		default:
			// get key
			portKey, err := dec.Token()
			if err != nil {
				return fmt.Errorf("failed to read port key: %w", err)
			}

			// check key is a string
			key, ok := portKey.(string)
			if !ok {
				return fmt.Errorf("unexpected type for port key: '%T'", token)
			}

			// read the rest of the port JSON
			var port domain.Port
			err = dec.Decode(&port)
			if err != nil {
				return fmt.Errorf("error on decoding port with key '%s': %w", key, err)
			}

			port.ID = key

			batch = append(batch, port)

			if len(batch) == batchSize {
				wg.Add(1)
				go func(ports domain.Ports) {
					defer wg.Done()

					err := i.portSvc.BulkUpsert(ctx, batch)
					if err != nil {
						errCh <- err
					}
				}(batch)

				batch = make(domain.Ports, 0, batchSize)
			}
		}
	}

	if !done && len(batch) > 0 {
		err = i.portSvc.BulkUpsert(ctx, batch)
	}

	wg.Wait()

	return err
}
