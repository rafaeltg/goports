package cid_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/rafaeltg/goports/pkg/cid"
	"github.com/stretchr/testify/assert"
)

func TestFromRequest(t *testing.T) {
	uuidVal := "0768b925-5aca-4f86-983b-8331c263d2ee"

	tcs := []struct {
		name          string
		req           *http.Request
		expectedErr   error
		expectedCid   string
		mockedUUIDGen *generatorMock
	}{
		{
			name: "request with header",
			req: &http.Request{
				Header: http.Header{
					cid.HeaderKey: []string{uuidVal},
				},
			},
			expectedCid: uuidVal,
		},
		{
			name: "request without header",
			req: &http.Request{
				Header: http.Header{},
			},
			expectedCid: uuidVal,
			mockedUUIDGen: &generatorMock{
				newV4Fn: func() (uuid.UUID, error) {
					return uuid.FromStringOrNil(uuidVal), nil
				},
			},
		},
		{
			name: "failed to generate uuid",
			req: &http.Request{
				Header: http.Header{},
			},
			expectedCid: "",
			expectedErr: errors.New("failed to generate correlation id: err"),
			mockedUUIDGen: &generatorMock{
				newV4Fn: func() (uuid.UUID, error) {
					return uuid.UUID{}, errors.New("err")
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockedUUIDGen != nil {
				uuid.DefaultGenerator = tc.mockedUUIDGen
			}

			actualCid, actualErr := cid.FromRequest(tc.req)

			assert.Equal(t, tc.expectedCid, actualCid)
			if tc.expectedErr != nil {
				assert.EqualError(t, actualErr, tc.expectedErr.Error())
			} else {
				assert.NoError(t, actualErr)
			}
		})
	}
}

type generatorMock struct {
	uuid.Generator
	newV4Fn func() (uuid.UUID, error)
}

func (g generatorMock) NewV4() (uuid.UUID, error) {
	return g.newV4Fn()
}
