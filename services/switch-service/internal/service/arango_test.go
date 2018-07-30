package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockResp struct {
	StatusCode int
	Body       string
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name                 string
		contextTimeout       time.Duration
		user, password       string
		database, collection string
		router               map[string]mockResp
		expectError          bool
	}{
		{
			"Unavailable",
			50 * time.Millisecond,
			"user", "pass",
			"animals", "mammals",
			map[string]mockResp{
				"/_admin/server/availability": mockResp{http.StatusServiceUnavailable, ``},
			},
			true,
		},
		{
			"Unauthorized",
			50 * time.Millisecond,
			"user", "pass",
			"animals", "mammals",
			map[string]mockResp{
				"/_admin/server/availability":          mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/database/current":   mockResp{http.StatusUnauthorized, `{"code":401, "error":true, "errorMessage":"not authorized"}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusUnauthorized, `{"code":401, "error":true, "errorMessage":"not authorized"}`},
			},
			true,
		},
		{
			"Unauthorized",
			50 * time.Millisecond,
			"user", "pass",
			"animals", "mammals",
			map[string]mockResp{
				"/_admin/server/availability":          mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/database/current":   mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusUnauthorized, `{"code":401, "error":true, "errorMessage":"not authorized"}`},
			},
			true,
		},
		{
			"DatabaseAndCollectionExist",
			50 * time.Millisecond,
			"user", "pass",
			"animals", "mammals",
			map[string]mockResp{
				"/_admin/server/availability":          mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/database/current":   mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusOK, `{"code":200, "error":false}`},
			},
			false,
		},
		{
			"CreateDatabaseError",
			50 * time.Millisecond,
			"user", "pass",
			"animals", "mammals",
			map[string]mockResp{
				"/_admin/server/availability":        mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/database/current": mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"database not found"}`},
				"/_db/_system/_api/database":         mockResp{http.StatusBadRequest, `{"code":400, "error":false}`},
			},
			true,
		},
		{
			"CreateCollectionError",
			50 * time.Millisecond,
			"user", "pass",
			"animals", "mammals",
			map[string]mockResp{
				"/_admin/server/availability":          mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/database/current":   mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"database not found"}`},
				"/_db/_system/_api/database":           mockResp{http.StatusCreated, `{"code":201, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"collection not found"}`},
				"/_db/animals/_api/collection":         mockResp{http.StatusBadRequest, `{"code":400, "error":false}`},
			},
			true,
		},
		{
			"CreateDatabaseAndCollection",
			50 * time.Millisecond,
			"user", "pass",
			"animals", "mammals",
			map[string]mockResp{
				"/_admin/server/availability":          mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/database/current":   mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"database not found"}`},
				"/_db/_system/_api/database":           mockResp{http.StatusCreated, `{"code":201, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"collection not found"}`},
				"/_db/animals/_api/collection":         mockResp{http.StatusOK, `{"code":200, "error":false}`},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock a back-end
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.router[path].StatusCode)
				w.Write([]byte(tc.router[path].Body))
			}))
			defer ts.Close()

			ctx, cancel := context.WithTimeout(context.Background(), tc.contextTimeout)
			defer cancel()

			service := NewArangoService()
			err := service.Connect(ctx, []string{ts.URL}, tc.user, tc.password, tc.database, tc.collection)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
