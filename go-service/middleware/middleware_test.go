package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	mockMiddleware struct {
		called bool
	}
)

func (m *mockMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.called = true
		next(w, r)
	}
}

func TestWrapAll(t *testing.T) {
	tests := []struct {
		name       string
		middleware []Middleware
	}{
		{
			"Simple",
			[]Middleware{
				&mockMiddleware{},
				&mockMiddleware{},
				&mockMiddleware{},
				&mockMiddleware{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://service/votes", nil)
			w := httptest.NewRecorder()

			handler := WrapAll(http.NotFound, tc.middleware...)
			handler(w, r)
			res := w.Result()

			assert.Equal(t, http.StatusNotFound, res.StatusCode)
			for _, m := range tc.middleware {
				mock, ok := m.(*mockMiddleware)
				assert.True(t, ok)
				assert.True(t, mock.called)
			}
		})
	}
}
