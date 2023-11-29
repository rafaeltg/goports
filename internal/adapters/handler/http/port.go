package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/port"
	"github.com/rafaeltg/goports/pkg/cid"
	"github.com/rafaeltg/goports/pkg/logging"
)

func getPortHandler(
	portSvc port.PortService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := getContext(r)

		id := mux.Vars(r)["id"]

		p, err := portSvc.Get(ctx, id)
		if err != nil {
			switch err {
			case port.ErrPortNotFound:
				writeResponse(
					w,
					withStatusCode(http.StatusNotFound),
					withError(err),
				)
			default:
				logger.ErrorContext(ctx,
					"failed to get port",
					logging.Error(err),
				)

				writeResponse(
					w,
					withStatusCode(http.StatusInternalServerError),
					withError(err),
				)
			}

			return
		}

		writeResponse(
			w,
			withStatusCode(http.StatusOK),
			withBody(p),
		)
	})
}

func bulkUpsertHandler(
	portSvc port.PortService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := getContext(r)

		var ports domain.Ports

		err := json.NewDecoder(r.Body).Decode(&ports)
		if err != nil {
			logger.ErrorContext(ctx,
				"failed to decode request body",
				logging.Error(err),
			)

			writeResponse(
				w,
				withStatusCode(http.StatusBadRequest),
				withError(errBadRequest),
			)

			return
		}

		err = portSvc.BulkUpsert(ctx, ports)
		if err != nil {
			writeResponse(
				w,
				withStatusCode(http.StatusInternalServerError),
				withError(err),
			)

			return
		}

		writeResponse(
			w,
			withStatusCode(http.StatusCreated),
		)
	})
}

func getContext(r *http.Request) context.Context {
	ctx := context.Background()

	corrId, err := cid.FromRequest(r)
	if err == nil {
		ctx = cid.NewContext(ctx, corrId)
	}

	return ctx
}

// WithPortHandlers setup port API handlers.
func WithPortHandlers(
	router *mux.Router,
	portSvc port.PortService,
	logger *slog.Logger,
) {
	router.Handle("/ports/{id}", getPortHandler(portSvc, logger)).
		Methods(http.MethodGet).
		Name("getPort")

	router.Handle("/ports/bulk-upsert", bulkUpsertHandler(portSvc, logger)).
		Methods(http.MethodPost).
		Name("bulkUpsertPorts")
}
