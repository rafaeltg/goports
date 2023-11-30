package redis

import (
	"context"
	"log/slog"

	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/port"
	"github.com/rafaeltg/goports/pkg/logging"
)

type (
	// PortRepository implements PortRepository interface.
	PortRepository struct {
		db     *Database
		logger *slog.Logger
	}

	portDTO struct {
		ID          string    `redis:"id"`
		Name        string    `redis:"name"`
		City        string    `redis:"city"`
		Country     string    `redis:"country"`
		Alias       []string  `redis:"alias,omitempty"`
		Regions     []string  `redis:"regions,omitempty"`
		Coordinates []float64 `redis:"coordinates,omitempty"`
		Province    string    `redis:"province"`
		Timezone    string    `redis:"timezone"`
		Unlocs      []string  `redis:"unlocs,omitempty"`
		Code        string    `redis:"code"`
	}
)

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

	var p portDTO

	err := r.db.MGet(ctx, id).Scan(&p)
	if err != nil {
		r.logger.ErrorContext(ctx,
			"[PortRepository.Get] error getting port",
			logging.Error(err),
		)

		return nil, err
	}

	if p.ID == "" {
		return nil, port.ErrPortNotFound
	}

	return p.ToDomain(), nil
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
			err := r.db.MSet(ctx, p.ID, p).Err()
			if err != nil {
				r.logger.ErrorContext(ctx,
					"[PortRepository.BulkUpsert] error inserting port",
					logging.Error(err),
				)

				return err
			}
		}
	}

	return nil
}

func (p portDTO) ToDomain() *domain.Port {
	return &domain.Port{
		ID:          p.ID,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
}
