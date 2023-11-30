package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/pkg/cid"
	"github.com/rafaeltg/goports/pkg/logging"
)

const bulkUpsertPortsPath = "/ports/bulk-upsert"

type (
	HttpClient interface {
		Do(*Request, *Response) error
	}

	PortClient struct {
		client HttpClient
		logger *slog.Logger
	}
)

func NewPortClient(client HttpClient, logger *slog.Logger) *PortClient {
	return &PortClient{
		client: client,
		logger: logger,
	}
}

func (p *PortClient) BulkUpsert(ctx context.Context, ports domain.Ports) error {
	p.logger.DebugContext(ctx,
		"[PortClient.BulkUpsert] executing",
		slog.Int("ports.length", len(ports)),
	)

	req := &Request{
		Path:   bulkUpsertPortsPath,
		Method: http.MethodPost,
		Body:   ports,
	}

	corrId, ok := cid.FromContext(ctx)
	if !ok {
		id, _ := uuid.NewV4()
		corrId = id.String()
	}

	req.Headers = map[string]string{
		"Content-Type": "application/json",
		"X-Request-Id": corrId,
	}

	res := &Response{
		StatusCode: http.StatusCreated,
		OutError:   &ApiErrorResponse{},
	}

	err := p.client.Do(req, res)
	if err != nil {
		p.logger.ErrorContext(ctx,
			"[PortClient.BulkUpsert] failed to execute request",
			logging.Error(err),
		)
	}

	return err
}
