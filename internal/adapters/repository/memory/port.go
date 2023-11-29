package memory

import (
	"context"
	"log/slog"

	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/port"
)

// PortRepository implements PortRepository interface.
type PortRepository struct {
	db     *Database
	logger *slog.Logger
}

// NewPortRepository creates a new port repository instance.
func NewPortRepository(db *Database, logger *slog.Logger) *PortRepository {
	return &PortRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PortRepository) Get(ctx context.Context, id string) (*domain.Port, error) {
	r.logger.DebugContext(ctx,
		"[PortRepository.Get] executing",
		slog.String("id", id),
	)

	result, ok := r.db.Get(ctx, id)
	if !ok {
		return nil, port.ErrPortNotFound
	}

	return result.(*domain.Port), nil
}

func (r *PortRepository) BulkUpsert(ctx context.Context, ports domain.Ports) error {
	r.logger.DebugContext(ctx,
		"[PortRepository.BulkUpsert] executing",
		slog.Int("ports.length", len(ports)),
	)

	for _, p := range ports {
		select {
		case <-ctx.Done():
			return nil
		default:
			r.db.Set(ctx, p.ID, p)
		}

	}

	return nil
}
