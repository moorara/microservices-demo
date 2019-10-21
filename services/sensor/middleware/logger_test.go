package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestWrapWithLogger(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		url                 string
		statusCode          int
		expectedEndpoint    string
		expectedStatusClass string
	}{
		{"101", "GET", "http://service/resource", 101, "/resource", "1xx"},
		{"200", "GET", "http://service/resource", 200, "/resource", "2xx"},
		{"302", "GET", "http://service/resource", 302, "/resource", "3xx"},
		{"404", "GET", "http://service/resource", 404, "/resource", "4xx"},
		{"500", "GET", "http://service/resource", 500, "/resource", "5xx"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader, writer, err := os.Pipe()
			assert.NoError(t, err)

			dec := json.NewDecoder(reader)
			logger := log.NewJSONLogger(writer)
			loggerMiddleware := NewLoggerMiddleware(logger)

			r := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			handler := loggerMiddleware.Wrap(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			})
			handler(w, r)

			var log map[string]interface{}
			res := w.Result()
			err = dec.Decode(&log)
			assert.NoError(t, err)

			assert.Equal(t, tc.statusCode, res.StatusCode)
			assert.Equal(t, tc.method, log["req.method"])
			assert.Equal(t, tc.expectedEndpoint, log["req.endpoint"])
			assert.Equal(t, float64(tc.statusCode), log["res.statusCode"])
			assert.Equal(t, tc.expectedStatusClass, log["res.statusClass"])
		})
	}
}
