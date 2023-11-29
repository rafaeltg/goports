package service

import (
	"context"
	"log/slog"

	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/port"
)

type PortService struct {
	productRepo port.PortRepository
	logger      *slog.Logger
}

func NewPortService(repo port.PortRepository, logger *slog.Logger) *PortService {
	return &PortService{
		productRepo: repo,
		logger:      logger,
	}
}

func (svc *PortService) Get(ctx context.Context, id string) (*domain.Port, error) {
	svc.logger.DebugContext(ctx,
		"[PortService.Get] executing",
		slog.String("id", id),
	)

	return svc.productRepo.Get(ctx, id)
}

func (svc *PortService) BulkUpsert(ctx context.Context, ports domain.Ports) error {
	svc.logger.DebugContext(ctx,
		"[PortService.BulkUpsert] executing",
		slog.Any("ports", ports),
	)

	return svc.productRepo.BulkUpsert(ctx, ports)
}
