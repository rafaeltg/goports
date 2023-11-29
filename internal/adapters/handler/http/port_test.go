package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	gohttp "net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rafaeltg/goports/internal/adapters/handler/http"
	"github.com/rafaeltg/goports/internal/core/domain"
	"github.com/rafaeltg/goports/internal/core/port"
	"github.com/rafaeltg/goports/internal/core/port/porttest"
	"github.com/stretchr/testify/assert"
)

var loggerTest = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func TestGetPort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "ABC"

	tcs := []struct {
		name               string
		existingPort       *domain.Port
		svcError           error
		expectedStatusCode int
		expectedResponse   any
	}{
		{
			name:               "internal server error",
			svcError:           errors.New("internal"),
			expectedStatusCode: gohttp.StatusInternalServerError,
			expectedResponse: http.ErrorResponse{
				Error: http.ErrorData{
					Message: "internal",
				},
			},
		},
		{
			name:               "not found",
			svcError:           port.ErrPortNotFound,
			expectedStatusCode: gohttp.StatusNotFound,
			expectedResponse: http.ErrorResponse{
				Error: http.ErrorData{
					Message: port.ErrPortNotFound.Error(),
				},
			},
		},
		{
			name: "success",
			existingPort: &domain.Port{
				ID:     id,
				Unlocs: []string{"ABC"},
			},
			expectedStatusCode: gohttp.StatusOK,
			expectedResponse: domain.Port{
				ID:     id,
				Unlocs: []string{"ABC"},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockedPortSvc := porttest.NewMockPortService(ctrl)
			mockedPortSvc.EXPECT().
				Get(gomock.Any(), id).
				Return(tc.existingPort, tc.svcError)

			router := mux.NewRouter()
			http.WithPortHandlers(
				router,
				mockedPortSvc,
				loggerTest,
			)

			srv := httptest.NewServer(router)
			defer srv.Close()

			client := &gohttp.Client{}
			req, err := gohttp.NewRequestWithContext(
				context.Background(),
				gohttp.MethodGet,
				fmt.Sprintf("%s/ports/%s", srv.URL, id),
				nil,
			)
			assert.NoError(t, err)

			resp, err := client.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			expectedResp, err := json.Marshal(tc.expectedResponse)
			assert.NoError(t, err)

			actualResp, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			// Read all adds an exta \n at the end
			assert.Equal(t, string(expectedResp), string(actualResp)[:len(actualResp)-1])
		})
	}
}
