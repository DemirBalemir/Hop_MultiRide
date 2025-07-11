// service/osrm_test.go
package service

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeRoundTripper struct {
	resp *http.Response
	err  error
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.resp, f.err
}

func TestGetOSRMDistance(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse string
		mockStatus   int
		mockError    error
		expectError  bool
		wantDist     float64
		wantDur      float64
	}{
		{
			name: "valid response",
			mockResponse: `{
				"routes": [{"distance": 1500.5, "duration": 300.3}]
			}`,
			mockStatus:  200,
			expectError: false,
			wantDist:    1500.5,
			wantDur:     300.3,
		},
		{
			name:         "no routes returned",
			mockResponse: `{ "routes": [] }`,
			mockStatus:   200,
			expectError:  true,
		},
		{
			name:         "malformed JSON",
			mockResponse: `{ bad json }`,
			mockStatus:   200,
			expectError:  true,
		},
		{
			name:         "http 500 error",
			mockResponse: `internal server error`,
			mockStatus:   500,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{
				Transport: &fakeRoundTripper{
					resp: &http.Response{
						StatusCode: tt.mockStatus,
						Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					},
					err: tt.mockError,
				},
			}

			dist, dur, err := GetOSRMDistance(client, 32.8, 39.9, 33.0, 40.0)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantDist, dist)
				assert.Equal(t, tt.wantDur, dur)
			}
		})
	}
}
