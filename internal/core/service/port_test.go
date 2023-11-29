package service_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/domain/domaintest"
	"github.com/rafaeltg/goports/internal/core/port/porttest"
	"github.com/rafaeltg/goports/internal/core/service"
	"github.com/stretchr/testify/assert"
)

var loggerTest = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func TestPortService_Get(t *testing.T) {
	id := "ABC"

	t.Run("database error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedPortRepo := porttest.NewMockPortRepository(ctrl)
		mockedPortRepo.EXPECT().
			Get(gomock.Any(), id).
			Return(nil, errors.New("get err"))

		svc := service.NewPortService(mockedPortRepo, loggerTest)

		port, err := svc.Get(context.Background(), id)
		assert.EqualError(t, err, "get err")
		assert.Nil(t, port)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedPortRepo := porttest.NewMockPortRepository(ctrl)
		mockedPortRepo.EXPECT().
			Get(gomock.Any(), id).
			Return(
				&domain.Port{
					ID: id,
				},
				nil,
			)

		svc := service.NewPortService(mockedPortRepo, loggerTest)

		port, err := svc.Get(context.Background(), id)
		assert.NoError(t, err)
		assert.NotNil(t, port)
	})
}

func TestPortService_BulkUpsert(t *testing.T) {
	ports := domain.Ports{
		{
			ID:   "ABC",
			Name: "Test",
		},
		{
			ID:   "DEF",
			Name: "Test 1",
		},
	}

	t.Run("database error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedPortRepo := porttest.NewMockPortRepository(ctrl)
		mockedPortRepo.EXPECT().
			BulkUpsert(
				gomock.Any(),
				domaintest.PortsMatcher(ports),
			).
			Return(errors.New("upsert err"))

		svc := service.NewPortService(mockedPortRepo, loggerTest)

		err := svc.BulkUpsert(context.Background(), ports)
		assert.EqualError(t, err, "upsert err")
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedPortRepo := porttest.NewMockPortRepository(ctrl)
		mockedPortRepo.EXPECT().
			BulkUpsert(
				gomock.Any(),
				domaintest.PortsMatcher(ports),
			).
			Return(nil)

		svc := service.NewPortService(mockedPortRepo, loggerTest)

		err := svc.BulkUpsert(context.Background(), ports)
		assert.NoError(t, err)
	})
}
