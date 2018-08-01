package arango

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPService(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{"WithoutTimeout", "http://localhost:9999"},
		{"WithTimeout", "http://localhost:9999"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewHTTPService(tc.address)

			assert.NotNil(t, service)
		})
	}
}

func TestNotifyReady(t *testing.T) {
	tests := []struct {
		name             string
		mockServer       bool
		serverStatusCode int
		contextTimeout   time.Duration
		expectError      bool
	}{
		{
			"ServerUnavailable",
			false,
			0,
			100 * time.Millisecond,
			true,
		},
		{
			"ServerNotReady",
			true,
			http.StatusServiceUnavailable,
			100 * time.Millisecond,
			true,
		},
		{
			"ServerReady",
			true,
			http.StatusOK,
			100 * time.Millisecond,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &httpService{
				client: &http.Client{},
			}

			if tc.mockServer {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.serverStatusCode)
				}))
				defer ts.Close()
				service.address = ts.URL
			}

			ctx, cancel := context.WithTimeout(context.Background(), tc.contextTimeout)
			defer cancel()

			ch := service.NotifyReady(ctx)
			err := <-ch

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name             string
		mockServer       bool
		serverStatusCode int
		serverResponse   string
		user, password   string
		expectError      bool
	}{
		{
			"ServerUnavailable",
			false,
			0,
			``,
			"root", "pass",
			true,
		},
		{
			"RequestFailed",
			true,
			http.StatusUnauthorized,
			`{}`,
			"root", "pass",
			true,
		},
		{
			"InvalidResponse",
			true,
			http.StatusOK,
			`{"jwt"}`,
			"root", "pass",
			true,
		},
		{
			"Successful",
			true,
			http.StatusOK,
			`{"jwt":"token"}`,
			"root", "pass",
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &httpService{
				client: &http.Client{},
			}

			if tc.mockServer {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.serverStatusCode)
					w.Write([]byte(tc.serverResponse))
				}))
				defer ts.Close()
				service.address = ts.URL
			}

			ctx := context.Background()
			err := service.Login(ctx, tc.user, tc.password)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCall(t *testing.T) {
	tests := []struct {
		name                   string
		mockServer             bool
		serverStatusCode       int
		serverResponse         string
		method, endpoint, body string
		expectError            bool
	}{
		{
			"InvalidRequest",
			false,
			0,
			``,
			"MOCK", ":", ``,
			true,
		},
		{
			"ServerUnavailable",
			false,
			0,
			``,
			"", "", ``,
			true,
		},
		{
			"InvalidResponse",
			true,
			http.StatusOK,
			`{"result"}`,
			"GET", "/_db/example", ``,
			true,
		},
		{
			"Successful",
			true,
			http.StatusOK,
			`{"result":1}`,
			"POST", "/_api/database", `{"name":"example"}`,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &httpService{
				client: &http.Client{},
			}

			if tc.mockServer {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.serverStatusCode)
					w.Write([]byte(tc.serverResponse))
				}))
				defer ts.Close()
				service.address = ts.URL
			}

			ctx := context.Background()
			_, statusCode, err := service.Call(ctx, tc.method, tc.endpoint, tc.body)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.serverStatusCode, statusCode)
			}
		})
	}
}
