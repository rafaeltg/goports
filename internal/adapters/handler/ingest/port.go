package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/port"
	"github.com/rafaeltg/goports/pkg/logging"
)

const batchSizeDefault int = 20

type (
	PortIngestor struct {
		portSvc   port.PortService
		batchSize int
		logger    *slog.Logger
	}

	PortIngestorOption func(*PortIngestor)
)

func NewPortIngestor(svc port.PortService, logger *slog.Logger, opts ...PortIngestorOption) *PortIngestor {
	i := &PortIngestor{
		portSvc:   svc,
		logger:    logger,
		batchSize: batchSizeDefault,
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *PortIngestor) Process(ctx context.Context, filename string) error {
	l := i.logger.With(
		slog.String("filepath", filename),
	)

	l.InfoContext(ctx, "[PortIngestor.Process] processing")

	f, err := os.Open(filepath.Clean(filename))
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
	batch := make(domain.Ports, 0, i.batchSize)

	done := false
	for !done && dec.More() {
		select {
		case <-ctx.Done():
			done = true
		case err = <-errCh:
			done = true
		default:
			var id json.Token

			id, err = dec.Token()
			if err != nil {
				err = fmt.Errorf("failed to read port key: %w", err)
				done = true

				continue
			}

			key, ok := id.(string)
			if !ok {
				err = fmt.Errorf("unexpected type for port key: '%T'", token)
				done = true

				continue
			}

			// read the rest of the port JSON
			var port domain.Port

			err = dec.Decode(&port)
			if err != nil {
				err = fmt.Errorf("error on decoding port with id '%s': %w", key, err)
				done = true

				continue
			}

			port.ID = key

			batch = append(batch, port)

			if len(batch) == i.batchSize {
				wg.Add(1)

				go func(ports domain.Ports) {
					defer wg.Done()

					err := i.portSvc.BulkUpsert(ctx, ports)
					if err != nil {
						errCh <- err
					}
				}(batch)

				batch = make(domain.Ports, 0, i.batchSize)
			}
		}
	}

	if !done && len(batch) > 0 {
		err = i.portSvc.BulkUpsert(ctx, batch)
	}

	wg.Wait()

	if err != nil {
		i.logger.ErrorContext(ctx,
			"[PortIngestor.Process] failed to process file",
			logging.Error(err),
		)
	}

	return err
}

func WithBatchSize(v int) PortIngestorOption {
	return func(pi *PortIngestor) {
		pi.batchSize = v
	}
}
