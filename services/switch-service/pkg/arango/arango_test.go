package arango

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPService(t *testing.T) {
	tests := []struct {
		name    string
		address string
		timeout time.Duration
	}{
		{"WithoutTimeout", "http://localhost:9999", 0},
		{"WithTimeout", "http://localhost:9999", 2 * time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewHTTPService(tc.address, tc.timeout)

			assert.NotNil(t, service)
		})
	}
}

func TestNotifyReady(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		contextTimeout time.Duration
		expectError    bool
	}{
		{
			"ContextTimeout",
			"http://localhost:9999",
			100 * time.Millisecond,
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &httpService{
				client:  &http.Client{},
				address: tc.address,
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
		name           string
		address        string
		user, password string
		expectError    bool
	}{
		{
			"Failure",
			"http://localhost:9999",
			"root", "pass",
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &httpService{
				client:  &http.Client{},
				address: tc.address,
			}

			err := service.Login(tc.user, tc.password)

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
		address                string
		method, endpoint, body string
		expectError            bool
	}{
		{
			"InvalidRequest",
			"",
			"", "", ``,
			true,
		},
		{
			"CreateDatabase",
			"http://localhost:9999",
			"POST", "/_api/database", `{"name":"example"}`,
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &httpService{
				client:  &http.Client{},
				address: tc.address,
			}

			_, _, err := service.Call(context.Background(), tc.method, tc.endpoint, tc.body)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
