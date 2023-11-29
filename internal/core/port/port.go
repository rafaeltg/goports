package port

import (
	"context"
	"errors"

	"github.com/rafaeltg/goports/internal/core/domain"
)

var ErrPortNotFound = errors.New("port not found")

//go:generate mockgen -source=port.go -destination=porttest/port_mock.go -package=porttest
type (
	// PortRepository is an interface for interacting with port-related data.
	PortRepository interface {
		Get(ctx context.Context, id string) (*domain.Port, error)
		BulkUpsert(ctx context.Context, ports domain.Ports) error
	}

	// PortService is an interface for interacting with port-related business logic.
	PortService interface {
		Get(context.Context, string) (*domain.Port, error)
		BulkUpsert(context.Context, domain.Ports) error
	}
)
