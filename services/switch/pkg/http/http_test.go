package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		statusClass string
	}{
		{"200", 200, "2xx"},
		{"201", 201, "2xx"},
		{"300", 300, "3xx"},
		{"400", 400, "4xx"},
		{"404", 404, "4xx"},
		{"500", 500, "5xx"},
		{"502", 502, "5xx"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}

			wrapperHandler := func(w http.ResponseWriter, r *http.Request) {
				rw := NewResponseWriter(w)
				handler(rw, r)

				assert.Equal(t, tc.statusCode, rw.StatusCode())
				assert.Equal(t, tc.statusClass, rw.StatusClass())
			}

			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			wrapperHandler(w, r)
		})
	}
}
