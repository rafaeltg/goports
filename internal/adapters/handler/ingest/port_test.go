package ingest_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaeltg/goports/internal/adapters/handler/ingest"
	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/domain/domaintest"
	"github.com/rafaeltg/goports/internal/core/port/porttest"
	"github.com/stretchr/testify/assert"
)

var loggerTest = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func TestProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("file not found", func(t *testing.T) {
		ingestor := ingest.NewPortIngestor(nil, loggerTest)

		err := ingestor.Process(context.Background(), "abc.json")
		assert.EqualError(t, err, "failed to read file: open abc.json: no such file or directory")
	})

	t.Run("invalid file", func(t *testing.T) {
		ingestor := ingest.NewPortIngestor(nil, loggerTest)

		err := ingestor.Process(context.Background(), "testdata/ports_invalid.json")
		assert.EqualError(t, err, "unexpected token encountered on reading opening delimiterr: abc")
	})

	t.Run("invalid key", func(t *testing.T) {
		ingestor := ingest.NewPortIngestor(nil, loggerTest)

		err := ingestor.Process(context.Background(), "testdata/ports_invalid_key.json")
		assert.EqualError(t, err, "failed to read port key: invalid character '1'")
	})

	t.Run("invalid port", func(t *testing.T) {
		ingestor := ingest.NewPortIngestor(nil, loggerTest)

		err := ingestor.Process(context.Background(), "testdata/ports_invalid_port.json")
		assert.EqualError(t, err, "error on decoding port with id 'ABC': json: cannot unmarshal number into Go struct field Port.id of type string")
	})

	t.Run("failed to insert ports", func(t *testing.T) {
		mockedPortSvc := porttest.NewMockPortService(ctrl)
		mockedPortSvc.EXPECT().
			BulkUpsert(gomock.Any(), gomock.Any()).
			Return(errors.New("bulk err"))

		ingestor := ingest.NewPortIngestor(mockedPortSvc, loggerTest)

		err := ingestor.Process(context.Background(), "testdata/ports_valid.json")
		assert.EqualError(t, err, "bulk err")
	})

	t.Run("success", func(t *testing.T) {
		mockedPortSvc := porttest.NewMockPortService(ctrl)
		mockedPortSvc.EXPECT().
			BulkUpsert(
				gomock.Any(),
				domaintest.PortsMatcher(
					domain.Ports{
						{
							ID:   "AEAJM",
							Name: "Ajman",
						},
						{
							ID:   "AEAUH",
							Name: "Abu Dhabi",
						},
					},
				),
			).
			Return(nil)

		mockedPortSvc.EXPECT().
			BulkUpsert(
				gomock.Any(),
				domaintest.PortsMatcher(
					domain.Ports{
						{
							ID:   "AEDXB",
							Name: "Dubai",
						},
						{
							ID:   "AEFJR",
							Name: "Al Fujayrah",
						},
					},
				),
			).
			Return(nil)

		ingestor := ingest.NewPortIngestor(
			mockedPortSvc,
			loggerTest,
			ingest.WithBatchSize(2),
		)

		err := ingestor.Process(context.Background(), "testdata/ports_valid.json")
		assert.NoError(t, err)
	})
}
